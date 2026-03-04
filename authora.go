package authora

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	Version        = "0.1.0"
	DefaultBaseURL = "https://api.authora.dev/api/v1"
	DefaultTimeout = 30 * time.Second
)

type Client struct {
	Agents        *AgentsService
	Roles         *RolesService
	Permissions   *PermissionsService
	Delegations   *DelegationsService
	Policies      *PoliciesService
	Mcp           *McpService
	Audit         *AuditService
	Notifications *NotificationsService
	Webhooks      *WebhooksService
	Alerts        *AlertsService
	APIKeys       *APIKeysService
	Organizations *OrganizationsService
	Workspaces    *WorkspacesService

	http *httpClient
}

type Option func(*clientConfig)

type clientConfig struct {
	baseURL    string
	timeout    time.Duration
	httpClient *http.Client
}

func WithBaseURL(url string) Option {
	return func(c *clientConfig) {
		c.baseURL = url
	}
}

func WithTimeout(d time.Duration) Option {
	return func(c *clientConfig) {
		c.timeout = d
	}
}

func WithHTTPClient(hc *http.Client) Option {
	return func(c *clientConfig) {
		c.httpClient = hc
	}
}

func NewClient(apiKey string, opts ...Option) *Client {
	cfg := &clientConfig{
		baseURL: DefaultBaseURL,
		timeout: DefaultTimeout,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	var hc *http.Client
	if cfg.httpClient != nil {
		hc = cfg.httpClient
	} else {
		hc = &http.Client{
			Timeout: cfg.timeout,
		}
	}

	httpC := newHTTPClient(cfg.baseURL, apiKey, hc)

	client := &Client{
		http: httpC,
	}

	client.Agents = &AgentsService{client: httpC}
	client.Roles = &RolesService{client: httpC}
	client.Permissions = &PermissionsService{client: httpC}
	client.Delegations = &DelegationsService{client: httpC}
	client.Policies = &PoliciesService{client: httpC}
	client.Mcp = &McpService{client: httpC}
	client.Audit = &AuditService{client: httpC}
	client.Notifications = &NotificationsService{client: httpC}
	client.Webhooks = &WebhooksService{client: httpC}
	client.Alerts = &AlertsService{client: httpC}
	client.APIKeys = &APIKeysService{client: httpC}
	client.Organizations = &OrganizationsService{client: httpC}
	client.Workspaces = &WorkspacesService{client: httpC}

	return client
}

func (c *Client) CreateAgent(ctx context.Context, input *CreateAgentInput) (*AgentRuntime, *KeyPair, error) {
	agent, err := c.Agents.Create(ctx, input)
	if err != nil {
		return nil, nil, fmt.Errorf("authora: create agent: %w", err)
	}

	kp, err := GenerateKeyPair()
	if err != nil {
		return nil, nil, err
	}

	_, err = c.Agents.Activate(ctx, agent.ID, &ActivateAgentInput{
		PublicKey: kp.PublicKey,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("authora: activate agent: %w", err)
	}

	runtime, err := NewAgent(AgentOptions{
		AgentID:    agent.ID,
		PrivateKey: kp.PrivateKey,
		BaseURL:    c.http.baseURL,
	})
	if err != nil {
		return nil, nil, err
	}

	return runtime, kp, nil
}

func (c *Client) LoadAgent(opts AgentOptions) (*AgentRuntime, error) {
	if opts.BaseURL == "" {
		opts.BaseURL = c.http.baseURL
	}
	return NewAgent(opts)
}

func (c *Client) LoadDelegatedAgent(opts AgentOptions) (*AgentRuntime, error) {
	if opts.DelegationToken == "" {
		return nil, fmt.Errorf("authora: delegationToken required")
	}
	if opts.BaseURL == "" {
		opts.BaseURL = c.http.baseURL
	}
	return NewAgent(opts)
}

func (c *Client) VerifyAgent(ctx context.Context, agentID string) (*VerifyAgentResponse, error) {
	return c.Agents.Verify(ctx, agentID)
}
