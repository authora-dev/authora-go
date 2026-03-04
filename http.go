package authora

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// httpClient wraps net/http to provide JSON request/response handling.
type httpClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// newHTTPClient creates a new httpClient.
func newHTTPClient(baseURL, apiKey string, client *http.Client) *httpClient {
	return &httpClient{
		baseURL:    strings.TrimRight(baseURL, "/"),
		apiKey:     apiKey,
		httpClient: client,
	}
}

// request performs an HTTP request and decodes the JSON response into dest.
// If dest is nil, the response body is discarded.
// If body is nil, no request body is sent.
func (h *httpClient) request(ctx context.Context, method, path string, body interface{}, dest interface{}) error {
	return h.requestWithAuth(ctx, method, path, body, dest, true)
}

// requestNoAuth performs an HTTP request without authentication.
func (h *httpClient) requestNoAuth(ctx context.Context, method, path string, body interface{}, dest interface{}) error {
	return h.requestWithAuth(ctx, method, path, body, dest, false)
}

func (h *httpClient) requestWithAuth(ctx context.Context, method, path string, body interface{}, dest interface{}, auth bool) error {
	fullURL := h.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("authora: failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return fmt.Errorf("authora: failed to create request: %w", err)
	}

	// Only set Content-Type when there is a body to send.
	// Sending Content-Type: application/json with an empty body causes Fastify
	// to reject the request with FST_ERR_CTP_EMPTY_JSON_BODY.
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	if auth && h.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+h.apiKey)
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("authora: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("authora: failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		apiErr := &AuthoraError{
			StatusCode: resp.StatusCode,
		}
		if len(respBody) > 0 {
			// Try to parse structured error.
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
		return apiErr
	}

	if dest != nil && len(respBody) > 0 {
		// Unwrap the backend's { "data": T } / { "data": [], "pagination": {} } envelope.
		unwrapped := unwrapResponse(respBody)
		if err := json.Unmarshal(unwrapped, dest); err != nil {
			return fmt.Errorf("authora: failed to decode response: %w", err)
		}
	}

	return nil
}

// unwrapResponse extracts the inner payload from the API's response envelope.
// The Authora API wraps responses in { "data": T } or { "data": [...], "pagination": {...} }.
func unwrapResponse(raw []byte) []byte {
	var envelope struct {
		Data       json.RawMessage `json:"data"`
		Pagination json.RawMessage `json:"pagination"`
		Meta       json.RawMessage `json:"meta"`
	}
	if err := json.Unmarshal(raw, &envelope); err != nil || envelope.Data == nil {
		return raw
	}

	// Determine if data is an array
	trimmed := bytes.TrimSpace(envelope.Data)
	isArray := len(trimmed) > 0 && trimmed[0] == '['

	// Get pagination info from either "pagination" or "meta"
	paginationRaw := envelope.Pagination
	if paginationRaw == nil {
		paginationRaw = envelope.Meta
	}

	if isArray && paginationRaw != nil {
		// Build a paginated list: { "items": [...], "total": N, "page": N, "limit": N }
		var pg struct {
			Total int `json:"total"`
			Page  int `json:"page"`
			Limit int `json:"limit"`
		}
		if err := json.Unmarshal(paginationRaw, &pg); err == nil {
			result := fmt.Sprintf(`{"items":%s,"total":%d,"page":%d,"limit":%d}`,
				string(envelope.Data), pg.Total, pg.Page, pg.Limit)
			return []byte(result)
		}
	}

	if isArray {
		// Array without pagination: return the raw array so that both []T
		// and struct{Items []T} targets can be unmarshalled correctly.
		return envelope.Data
	}

	// Single entity: return unwrapped data
	return envelope.Data
}

// queryString builds a URL query string from a map of parameters.
// Nil values are skipped.
func queryString(params map[string]interface{}) string {
	values := url.Values{}
	for k, v := range params {
		if v == nil {
			continue
		}
		switch val := v.(type) {
		case string:
			if val != "" {
				values.Set(k, val)
			}
		case *string:
			if val != nil {
				values.Set(k, *val)
			}
		case int:
			values.Set(k, fmt.Sprintf("%d", val))
		case *int:
			if val != nil {
				values.Set(k, fmt.Sprintf("%d", *val))
			}
		case bool:
			values.Set(k, fmt.Sprintf("%t", val))
		case *bool:
			if val != nil {
				values.Set(k, fmt.Sprintf("%t", *val))
			}
		default:
			values.Set(k, fmt.Sprintf("%v", val))
		}
	}
	encoded := values.Encode()
	if encoded == "" {
		return ""
	}
	return "?" + encoded
}
