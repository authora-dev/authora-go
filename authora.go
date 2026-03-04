// Package authora provides a Go client for the Authora API.
//
// Authora is an identity and authorization platform for AI agents.
// This SDK provides typed access to all Authora API endpoints with
// zero external dependencies.
//
// Usage:
//
//	client := authora.NewClient("authora_live_...")
//	agent, err := client.Agents.Create(ctx, &authora.CreateAgentInput{
//	    WorkspaceID: "ws_...",
//	    Name:        "my-agent",
//	    CreatedBy:   "user_...",
//	})
package authora

import (
	"net/http"
	"time"
)

const (
	// Version is the SDK version.
	Version = "0.1.0"

	// DefaultBaseURL is the default Authora API base URL.
	DefaultBaseURL = "https://api.authora.dev/api/v1"

	// DefaultTimeout is the default HTTP client timeout.
	DefaultTimeout = 30 * time.Second
)

// Client is the Authora API client. It exposes service objects for each
// API resource, following the pattern client.Resource.Method(ctx, input).
type Client struct {
	// Service objects for each API resource.
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

// Option configures the Client.
type Option func(*clientConfig)

type clientConfig struct {
	baseURL    string
	timeout    time.Duration
	httpClient *http.Client
}

// WithBaseURL overrides the default API base URL.
func WithBaseURL(url string) Option {
	return func(c *clientConfig) {
		c.baseURL = url
	}
}

// WithTimeout sets the HTTP client timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *clientConfig) {
		c.timeout = d
	}
}

// WithHTTPClient sets a custom *http.Client for all requests.
// When set, the timeout option is ignored (configure it on the client directly).
func WithHTTPClient(hc *http.Client) Option {
	return func(c *clientConfig) {
		c.httpClient = hc
	}
}

// NewClient creates a new Authora API client.
//
//	client := authora.NewClient("authora_live_...",
//	    authora.WithBaseURL("https://custom.api.dev/api/v1"),
//	    authora.WithTimeout(10 * time.Second),
//	)
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
