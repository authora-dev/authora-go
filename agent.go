package authora

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type AgentOptions struct {
	AgentID             string
	PrivateKey          string
	BaseURL             string
	Timeout             time.Duration
	PermissionsCacheTTL time.Duration
	DelegationToken     string
}

type AgentRuntime struct {
	agentID         string
	privateKey      string
	publicKey       string
	baseURL         string
	httpClient      *http.Client
	cacheTTL        time.Duration
	delegationToken string

	mu          sync.RWMutex
	cachedAllow []string
	cachedDeny  []string
	cacheTime   time.Time
}

func NewAgent(opts AgentOptions) (*AgentRuntime, error) {
	if opts.AgentID == "" {
		return nil, fmt.Errorf("authora: agentId required")
	}
	if opts.PrivateKey == "" {
		return nil, fmt.Errorf("authora: privateKey required")
	}
	pubKey, err := GetPublicKey(opts.PrivateKey)
	if err != nil {
		return nil, err
	}
	baseURL := opts.BaseURL
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	baseURL = strings.TrimRight(baseURL, "/")
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	cacheTTL := opts.PermissionsCacheTTL
	if cacheTTL == 0 {
		cacheTTL = 5 * time.Minute
	}
	return &AgentRuntime{
		agentID:         opts.AgentID,
		privateKey:      opts.PrivateKey,
		publicKey:       pubKey,
		baseURL:         baseURL,
		httpClient:      &http.Client{Timeout: timeout},
		cacheTTL:        cacheTTL,
		delegationToken: opts.DelegationToken,
	}, nil
}

func (a *AgentRuntime) SignedRequest(ctx context.Context, method, path string, body interface{}) (*SignedResponse, error) {
	method = strings.ToUpper(method)
	fullURL := a.baseURL + path

	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("authora: failed to marshal body: %w", err)
		}
	}

	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")

	var bodyPtr *string
	if bodyBytes != nil {
		s := string(bodyBytes)
		bodyPtr = &s
	}

	payload := BuildSignaturePayload(method, path, timestamp, bodyPtr)
	sig, err := Sign(payload, a.privateKey)
	if err != nil {
		return nil, err
	}

	var bodyReader io.Reader
	if bodyBytes != nil {
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("authora: failed to create request: %w", err)
	}

	if bodyBytes != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-authora-agent-id", a.agentID)
	req.Header.Set("x-authora-timestamp", timestamp)
	req.Header.Set("x-authora-signature", sig)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("authora: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("authora: failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		apiErr := &AuthoraError{StatusCode: resp.StatusCode}
		if len(respBody) > 0 {
			var errResp struct {
				Message string `json:"message"`
				Code    string `json:"code"`
				Error   string `json:"error"`
			}
			if jsonErr := json.Unmarshal(respBody, &errResp); jsonErr == nil {
				apiErr.Message = errResp.Message
				apiErr.Code = errResp.Code
				if apiErr.Message == "" {
					apiErr.Message = errResp.Error
				}
			}
			if apiErr.Message == "" {
				apiErr.Message = string(respBody)
			}
		}
		if apiErr.Message == "" {
			apiErr.Message = http.StatusText(resp.StatusCode)
		}
		return nil, apiErr
	}

	unwrapped := unwrapResponse(respBody)
	return &SignedResponse{
		Data:       json.RawMessage(unwrapped),
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
	}, nil
}

func (a *AgentRuntime) CheckPermission(ctx context.Context, resource, action string, context_ map[string]interface{}) (*CheckPermissionResponse, error) {
	body := map[string]interface{}{
		"agentId":  a.agentID,
		"resource": resource,
		"action":   action,
	}
	if context_ != nil {
		body["context"] = context_
	}
	resp, err := a.SignedRequest(ctx, "POST", "/permissions/check", body)
	if err != nil {
		return nil, err
	}
	var result CheckPermissionResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("authora: failed to decode response: %w", err)
	}
	return &result, nil
}

func (a *AgentRuntime) CheckPermissions(ctx context.Context, checks []BatchCheckItem) ([]CheckPermissionResponse, error) {
	body := map[string]interface{}{
		"agentId": a.agentID,
		"checks":  checks,
	}
	resp, err := a.SignedRequest(ctx, "POST", "/permissions/check-batch", body)
	if err != nil {
		return nil, err
	}
	var result struct {
		Results []CheckPermissionResponse `json:"results"`
	}
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("authora: failed to decode response: %w", err)
	}
	return result.Results, nil
}

func (a *AgentRuntime) FetchPermissions(ctx context.Context) (*EffectivePermissionsResponse, error) {
	resp, err := a.SignedRequest(ctx, "GET", "/agents/"+a.agentID+"/permissions", nil)
	if err != nil {
		return nil, err
	}
	var result EffectivePermissionsResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("authora: failed to decode response: %w", err)
	}

	allow := make([]string, len(result.Permissions))
	for i, p := range result.Permissions {
		allow[i] = p.Permission
	}

	a.mu.Lock()
	a.cachedAllow = allow
	a.cachedDeny = nil
	a.cacheTime = time.Now()
	a.mu.Unlock()

	return &result, nil
}

func (a *AgentRuntime) HasPermission(ctx context.Context, resource string) (bool, error) {
	a.mu.RLock()
	expired := a.cachedAllow == nil || time.Since(a.cacheTime) > a.cacheTTL
	a.mu.RUnlock()

	if expired {
		if _, err := a.FetchPermissions(ctx); err != nil {
			return false, err
		}
	}

	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.cachedDeny != nil && MatchAnyPermission(a.cachedDeny, resource) {
		return false, nil
	}
	return MatchAnyPermission(a.cachedAllow, resource), nil
}

func (a *AgentRuntime) InvalidatePermissionsCache() {
	a.mu.Lock()
	a.cachedAllow = nil
	a.cachedDeny = nil
	a.cacheTime = time.Time{}
	a.mu.Unlock()
}

func (a *AgentRuntime) Delegate(ctx context.Context, targetAgentID string, permissions []string, constraints *DelegationConstraints) (*Delegation, error) {
	body := map[string]interface{}{
		"issuerAgentId": a.agentID,
		"targetAgentId": targetAgentID,
		"permissions":   permissions,
	}
	if constraints != nil {
		body["constraints"] = constraints
	}
	resp, err := a.SignedRequest(ctx, "POST", "/delegations", body)
	if err != nil {
		return nil, err
	}
	var result Delegation
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("authora: failed to decode response: %w", err)
	}
	return &result, nil
}

func (a *AgentRuntime) CallTool(ctx context.Context, params *ToolCallParams) (*McpProxyResponse, error) {
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	mcpSig, err := Sign(BuildSignaturePayload("POST", "/mcp/proxy", timestamp, nil), a.privateKey)
	if err != nil {
		return nil, err
	}

	id := params.ID
	if id == nil {
		id = fmt.Sprintf("%s-%d", a.agentID, time.Now().UnixMilli())
	}

	mcpMethod := params.Method
	if mcpMethod == "" {
		mcpMethod = "tools/call"
	}

	token := params.DelegationToken
	if token == "" {
		token = a.delegationToken
	}

	authora := map[string]interface{}{
		"agentId":   a.agentID,
		"signature": mcpSig,
		"timestamp": timestamp,
	}
	if token != "" {
		authora["delegationToken"] = token
	}

	body := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"method":  mcpMethod,
		"params": map[string]interface{}{
			"name":      params.ToolName,
			"arguments": params.Arguments,
			"_authora":  authora,
		},
	}

	resp, err := a.SignedRequest(ctx, "POST", "/mcp/proxy", body)
	if err != nil {
		return nil, err
	}
	var result McpProxyResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("authora: failed to decode response: %w", err)
	}
	return &result, nil
}

func (a *AgentRuntime) RotateKey(ctx context.Context) (*Agent, *KeyPair, error) {
	kp, err := GenerateKeyPair()
	if err != nil {
		return nil, nil, err
	}
	resp, err := a.SignedRequest(ctx, "POST", "/agents/"+a.agentID+"/rotate-key", map[string]string{
		"publicKey": kp.PublicKey,
	})
	if err != nil {
		return nil, nil, err
	}
	var agent Agent
	if err := json.Unmarshal(resp.Data, &agent); err != nil {
		return nil, nil, fmt.Errorf("authora: failed to decode response: %w", err)
	}
	return &agent, kp, nil
}

func (a *AgentRuntime) Suspend(ctx context.Context) (*Agent, error) {
	resp, err := a.SignedRequest(ctx, "POST", "/agents/"+a.agentID+"/suspend", nil)
	if err != nil {
		return nil, err
	}
	var agent Agent
	if err := json.Unmarshal(resp.Data, &agent); err != nil {
		return nil, fmt.Errorf("authora: failed to decode response: %w", err)
	}
	return &agent, nil
}

func (a *AgentRuntime) Reactivate(ctx context.Context) (*Agent, *KeyPair, error) {
	kp, err := GenerateKeyPair()
	if err != nil {
		return nil, nil, err
	}
	resp, err := a.SignedRequest(ctx, "POST", "/agents/"+a.agentID+"/activate", map[string]string{
		"publicKey": kp.PublicKey,
	})
	if err != nil {
		return nil, nil, err
	}
	var agent Agent
	if err := json.Unmarshal(resp.Data, &agent); err != nil {
		return nil, nil, fmt.Errorf("authora: failed to decode response: %w", err)
	}
	return &agent, kp, nil
}

func (a *AgentRuntime) Revoke(ctx context.Context) (*Agent, error) {
	resp, err := a.SignedRequest(ctx, "POST", "/agents/"+a.agentID+"/revoke", nil)
	if err != nil {
		return nil, err
	}
	var agent Agent
	if err := json.Unmarshal(resp.Data, &agent); err != nil {
		return nil, fmt.Errorf("authora: failed to decode response: %w", err)
	}
	return &agent, nil
}

func (a *AgentRuntime) GetIdentityDocument(ctx context.Context) (*VerifyAgentResponse, error) {
	fullURL := a.baseURL + "/agents/" + a.agentID + "/verify"
	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("authora: failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("authora: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("authora: failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		apiErr := &AuthoraError{StatusCode: resp.StatusCode}
		if len(respBody) > 0 {
			var errResp struct {
				Message string `json:"message"`
				Code    string `json:"code"`
				Error   string `json:"error"`
			}
			if jsonErr := json.Unmarshal(respBody, &errResp); jsonErr == nil {
				apiErr.Message = errResp.Message
				apiErr.Code = errResp.Code
				if apiErr.Message == "" {
					apiErr.Message = errResp.Error
				}
			}
			if apiErr.Message == "" {
				apiErr.Message = string(respBody)
			}
		}
		if apiErr.Message == "" {
			apiErr.Message = http.StatusText(resp.StatusCode)
		}
		return nil, apiErr
	}

	unwrapped := unwrapResponse(respBody)
	var result VerifyAgentResponse
	if err := json.Unmarshal(unwrapped, &result); err != nil {
		return nil, fmt.Errorf("authora: failed to decode response: %w", err)
	}
	return &result, nil
}

func (a *AgentRuntime) GetProfile(ctx context.Context) (*Agent, error) {
	resp, err := a.SignedRequest(ctx, "GET", "/agents/"+a.agentID, nil)
	if err != nil {
		return nil, err
	}
	var agent Agent
	if err := json.Unmarshal(resp.Data, &agent); err != nil {
		return nil, fmt.Errorf("authora: failed to decode response: %w", err)
	}
	return &agent, nil
}

func (a *AgentRuntime) GetPublicKey() string {
	return a.publicKey
}
