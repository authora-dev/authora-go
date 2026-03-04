package authora

import "time"

// ---------------------------------------------------------------------------
// Common / Pagination
// ---------------------------------------------------------------------------

// PaginatedResponse wraps list endpoints that return paginated data.
// After envelope unwrapping, the HTTP layer returns { "items": [...], "total": N, ... }.
type PaginatedResponse[T any] struct {
	Items      []T  `json:"items"`
	Total      int  `json:"total"`
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	TotalPages int  `json:"totalPages"`
	HasMore    bool `json:"hasMore"`
}

// ---------------------------------------------------------------------------
// Agents
// ---------------------------------------------------------------------------

// Agent represents an Authora agent identity.
type Agent struct {
	ID          string                 `json:"id"`
	WorkspaceID string                 `json:"workspaceId"`
	Status      string                 `json:"status"`
	CreatedBy   string                 `json:"createdBy"`
	PublicKey   *string                `json:"publicKey,omitempty"`
	IdentityDoc map[string]interface{} `json:"identityDoc,omitempty"`
	ExpiresAt   *time.Time             `json:"expiresAt,omitempty"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// CreateAgentInput is the request body for POST /agents.
type CreateAgentInput struct {
	WorkspaceID string                 `json:"workspaceId"`
	Name        string                 `json:"name"`
	Description *string                `json:"description,omitempty"`
	CreatedBy   string                 `json:"createdBy"`
	ExpiresIn   *string                `json:"expiresIn,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// CreateAgentResponse is the response from POST /agents.
// The response envelope is unwrapped, so the agent fields are at top level.
type CreateAgentResponse = Agent

// ListAgentsInput contains query parameters for GET /agents.
type ListAgentsInput struct {
	WorkspaceID string  `json:"workspaceId"`
	Status      *string `json:"status,omitempty"`
	Page        *int    `json:"page,omitempty"`
	Limit       *int    `json:"limit,omitempty"`
}

// VerifyAgentResponse is the response from GET /agents/:agentId/verify.
type VerifyAgentResponse struct {
	Valid  bool   `json:"valid"`
	Status string `json:"status"`
	Agent  *Agent `json:"agent,omitempty"`
}

// ActivateAgentInput is the request body for POST /agents/:agentId/activate.
type ActivateAgentInput struct {
	PublicKey string `json:"publicKey"`
}

// RotateKeyInput is the request body for POST /agents/:agentId/rotate-key.
type RotateKeyInput struct {
	PublicKey string `json:"publicKey"`
}

// RotateKeyResponse is the response from POST /agents/:agentId/rotate-key.
type RotateKeyResponse struct {
	Agent     Agent  `json:"agent"`
	NewAPIKey string `json:"newApiKey"`
}

// ---------------------------------------------------------------------------
// Roles
// ---------------------------------------------------------------------------

// Role represents a role definition.
type Role struct {
	ID                 string   `json:"id"`
	WorkspaceID        string   `json:"workspaceId"`
	Name               string   `json:"name"`
	Description        *string  `json:"description,omitempty"`
	Permissions        []string `json:"permissions"`
	DenyPermissions    []string `json:"denyPermissions,omitempty"`
	Stage              *string  `json:"stage,omitempty"`
	MaxSessionDuration *int     `json:"maxSessionDuration,omitempty"`
	CreatedAt          string   `json:"createdAt"`
	UpdatedAt          string   `json:"updatedAt"`
}

// CreateRoleInput is the request body for POST /roles.
type CreateRoleInput struct {
	WorkspaceID        string   `json:"workspaceId"`
	Name               string   `json:"name"`
	Description        *string  `json:"description,omitempty"`
	Permissions        []string `json:"permissions"`
	DenyPermissions    []string `json:"denyPermissions,omitempty"`
	Stage              *string  `json:"stage,omitempty"`
	MaxSessionDuration *int     `json:"maxSessionDuration,omitempty"`
}

// ListRolesInput contains query parameters for GET /roles.
type ListRolesInput struct {
	WorkspaceID string `json:"workspaceId"`
	Page        *int   `json:"page,omitempty"`
	Limit       *int   `json:"limit,omitempty"`
}

// UpdateRoleInput is the request body for PATCH /roles/:roleId.
type UpdateRoleInput struct {
	Name               *string  `json:"name,omitempty"`
	Description        *string  `json:"description,omitempty"`
	Permissions        []string `json:"permissions,omitempty"`
	DenyPermissions    []string `json:"denyPermissions,omitempty"`
	Stage              *string  `json:"stage,omitempty"`
	MaxSessionDuration *int     `json:"maxSessionDuration,omitempty"`
}

// ---------------------------------------------------------------------------
// Agent Role Assignments
// ---------------------------------------------------------------------------

// AgentRole represents an agent-role assignment.
type AgentRole struct {
	ID        string     `json:"id"`
	AgentID   string     `json:"agentId"`
	RoleID    string     `json:"roleId"`
	GrantedBy *string    `json:"grantedBy,omitempty"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
	CreatedAt string     `json:"createdAt"`
	Role      *Role      `json:"role,omitempty"`
}

// AssignRoleInput is the request body for POST /agents/:agentId/roles.
type AssignRoleInput struct {
	RoleID    string     `json:"roleId"`
	GrantedBy *string    `json:"grantedBy,omitempty"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
}

// ---------------------------------------------------------------------------
// Permissions
// ---------------------------------------------------------------------------

// CheckPermissionInput is the request body for POST /permissions/check.
type CheckPermissionInput struct {
	AgentID  string                 `json:"agentId"`
	Resource string                 `json:"resource"`
	Action   string                 `json:"action"`
	Context  map[string]interface{} `json:"context,omitempty"`
}

// CheckPermissionResponse is the response from POST /permissions/check.
type CheckPermissionResponse struct {
	Allowed bool     `json:"allowed"`
	Reason  *string  `json:"reason,omitempty"`
	Matched []string `json:"matched,omitempty"`
	Denied  []string `json:"denied,omitempty"`
}

// BatchCheckItem is a single permission check within a batch request.
type BatchCheckItem struct {
	Resource string                 `json:"resource"`
	Action   string                 `json:"action"`
	Context  map[string]interface{} `json:"context,omitempty"`
}

// BatchCheckInput is the request body for POST /permissions/check-batch.
type BatchCheckInput struct {
	AgentID string           `json:"agentId"`
	Checks  []BatchCheckItem `json:"checks"`
}

// BatchCheckResult is a single result in a batch check response.
type BatchCheckResult struct {
	Resource string   `json:"resource"`
	Action   string   `json:"action"`
	Allowed  bool     `json:"allowed"`
	Reason   *string  `json:"reason,omitempty"`
	Matched  []string `json:"matched,omitempty"`
	Denied   []string `json:"denied,omitempty"`
}

// BatchCheckResponse is the response from POST /permissions/check-batch.
type BatchCheckResponse struct {
	Results []BatchCheckResult `json:"results"`
}

// EffectivePermission represents a resolved permission for an agent.
type EffectivePermission struct {
	Permission string `json:"permission"`
	Source     string `json:"source"`
	RoleID     string `json:"roleId,omitempty"`
	RoleName   string `json:"roleName,omitempty"`
}

// EffectivePermissionsResponse is the response from GET /agents/:agentId/permissions.
type EffectivePermissionsResponse struct {
	AgentID     string                `json:"agentId"`
	Permissions []EffectivePermission `json:"permissions"`
}

// ---------------------------------------------------------------------------
// Delegations
// ---------------------------------------------------------------------------

// Delegation represents a delegation of authority between agents.
type Delegation struct {
	ID            string                   `json:"id"`
	IssuerAgentID string                   `json:"issuerAgentId"`
	TargetAgentID string                   `json:"targetAgentId"`
	Permissions   []string                 `json:"permissions"`
	Constraints   map[string]interface{}   `json:"constraints,omitempty"`
	Chain         []map[string]interface{} `json:"chain,omitempty"`
	Signature     *string                  `json:"signature,omitempty"`
	Status        string                   `json:"status"`
	ExpiresAt     *time.Time               `json:"expiresAt,omitempty"`
	CreatedAt     string                   `json:"createdAt"`
}

// CreateDelegationInput is the request body for POST /delegations.
type CreateDelegationInput struct {
	IssuerAgentID string                 `json:"issuerAgentId"`
	TargetAgentID string                 `json:"targetAgentId"`
	Permissions   []string               `json:"permissions"`
	Constraints   map[string]interface{} `json:"constraints,omitempty"`
	ExpiresIn     *string                `json:"expiresIn,omitempty"`
}

// VerifyDelegationInput is the request body for POST /delegations/verify.
type VerifyDelegationInput struct {
	DelegationID string `json:"delegationId"`
}

// VerifyDelegationResponse is the response from POST /delegations/verify.
type VerifyDelegationResponse struct {
	Valid      bool        `json:"valid"`
	Delegation *Delegation `json:"delegation,omitempty"`
	Reason     *string     `json:"reason,omitempty"`
}

// ListDelegationsInput contains query parameters for listing delegations.
type ListDelegationsInput struct {
	Status *string `json:"status,omitempty"`
	Page   *int    `json:"page,omitempty"`
	Limit  *int    `json:"limit,omitempty"`
}

// ListAgentDelegationsInput contains query parameters for listing an agent's delegations.
type ListAgentDelegationsInput struct {
	Direction *string `json:"direction,omitempty"`
	Status    *string `json:"status,omitempty"`
	Page      *int    `json:"page,omitempty"`
	Limit     *int    `json:"limit,omitempty"`
}

// ---------------------------------------------------------------------------
// Policies
// ---------------------------------------------------------------------------

// Policy represents an authorization policy.
type Policy struct {
	ID          string                 `json:"id"`
	WorkspaceID string                 `json:"workspaceId"`
	Name        string                 `json:"name"`
	Description *string                `json:"description,omitempty"`
	Effect      string                 `json:"effect"`
	Principals  map[string]interface{} `json:"principals,omitempty"`
	Conditions  map[string]interface{} `json:"conditions,omitempty"`
	Resources   []string               `json:"resources,omitempty"`
	Actions     []string               `json:"actions,omitempty"`
	Priority    *int                   `json:"priority,omitempty"`
	Enabled     *bool                  `json:"enabled,omitempty"`
	CreatedAt   string                 `json:"createdAt"`
	UpdatedAt   string                 `json:"updatedAt"`
}

// CreatePolicyInput is the request body for POST /policies.
type CreatePolicyInput struct {
	WorkspaceID string                 `json:"workspaceId"`
	Name        string                 `json:"name"`
	Description *string                `json:"description,omitempty"`
	Effect      string                 `json:"effect"`
	Principals  map[string]interface{} `json:"principals"`
	Conditions  map[string]interface{} `json:"conditions,omitempty"`
	Resources   []string               `json:"resources,omitempty"`
	Actions     []string               `json:"actions,omitempty"`
	Priority    *int                   `json:"priority,omitempty"`
	Enabled     *bool                  `json:"enabled,omitempty"`
}

// ListPoliciesInput contains query parameters for GET /policies.
type ListPoliciesInput struct {
	WorkspaceID string `json:"workspaceId"`
	Page        *int   `json:"page,omitempty"`
	Limit       *int   `json:"limit,omitempty"`
}

// UpdatePolicyInput is the request body for PATCH /policies/:policyId.
type UpdatePolicyInput struct {
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	Effect      *string                `json:"effect,omitempty"`
	Principals  map[string]interface{} `json:"principals,omitempty"`
	Conditions  map[string]interface{} `json:"conditions,omitempty"`
	Resources   []string               `json:"resources,omitempty"`
	Actions     []string               `json:"actions,omitempty"`
	Priority    *int                   `json:"priority,omitempty"`
	Enabled     *bool                  `json:"enabled,omitempty"`
}

// SimulatePolicyInput is the request body for POST /policies/simulate.
type SimulatePolicyInput struct {
	WorkspaceID string                 `json:"workspaceId"`
	AgentID     string                 `json:"agentId"`
	Resource    string                 `json:"resource"`
	Action      string                 `json:"action"`
	Context     map[string]interface{} `json:"context,omitempty"`
}

// SimulatePolicyResponse is the response from POST /policies/simulate.
type SimulatePolicyResponse struct {
	Decision        string   `json:"decision"`
	MatchedPolicies []Policy `json:"matchedPolicies,omitempty"`
	Reason          *string  `json:"reason,omitempty"`
}

// EvaluatePolicyInput is the request body for POST /policies/evaluate.
type EvaluatePolicyInput struct {
	WorkspaceID string                 `json:"workspaceId"`
	AgentID     string                 `json:"agentId"`
	Resource    string                 `json:"resource"`
	Action      string                 `json:"action"`
	Context     map[string]interface{} `json:"context,omitempty"`
}

// EvaluatePolicyResponse is the response from POST /policies/evaluate.
type EvaluatePolicyResponse struct {
	Allowed         bool     `json:"allowed"`
	MatchedPolicies []Policy `json:"matchedPolicies,omitempty"`
	Reason          *string  `json:"reason,omitempty"`
}

// ---------------------------------------------------------------------------
// MCP (Model Context Protocol)
// ---------------------------------------------------------------------------

// McpServer represents a registered MCP server.
type McpServer struct {
	ID          string                 `json:"id"`
	WorkspaceID string                 `json:"workspaceId"`
	Name        string                 `json:"name"`
	Description *string                `json:"description,omitempty"`
	URL         string                 `json:"url"`
	Status      string                 `json:"status"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   string                 `json:"createdAt"`
	UpdatedAt   string                 `json:"updatedAt"`
}

// RegisterMcpServerInput is the request body for POST /mcp/servers.
type RegisterMcpServerInput struct {
	WorkspaceID string                 `json:"workspaceId"`
	Name        string                 `json:"name"`
	Description *string                `json:"description,omitempty"`
	URL         string                 `json:"url"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ListMcpServersInput contains query parameters for GET /mcp/servers.
type ListMcpServersInput struct {
	WorkspaceID string `json:"workspaceId"`
	Page        *int   `json:"page,omitempty"`
	Limit       *int   `json:"limit,omitempty"`
}

// UpdateMcpServerInput is the request body for PATCH /mcp/servers/:serverId.
type UpdateMcpServerInput struct {
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	URL         *string                `json:"url,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// McpTool represents a tool registered on an MCP server.
type McpTool struct {
	ID          string                 `json:"id"`
	ServerID    string                 `json:"serverId"`
	Name        string                 `json:"name"`
	Description *string                `json:"description,omitempty"`
	InputSchema map[string]interface{} `json:"inputSchema,omitempty"`
	CreatedAt   string                 `json:"createdAt"`
	UpdatedAt   string                 `json:"updatedAt"`
}

// RegisterMcpToolInput is the request body for POST /mcp/servers/:serverId/tools.
type RegisterMcpToolInput struct {
	Name        string                 `json:"name"`
	Description *string                `json:"description,omitempty"`
	InputSchema map[string]interface{} `json:"inputSchema,omitempty"`
}

// McpProxyInput is the user-facing input for proxy calls.
// The SDK builds the proper JSON-RPC 2.0 request internally.
type McpProxyInput struct {
	ServerID string                 `json:"-"` // Used to populate _authora.mcpServerId
	Method   string                 `json:"-"` // JSON-RPC method
	Params   map[string]interface{} `json:"-"` // JSON-RPC params (merged with _authora)
}

// mcpProxyJsonRpc is the actual JSON-RPC 2.0 body sent to POST /mcp/proxy.
type mcpProxyJsonRpc struct {
	Jsonrpc string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	ID      int                    `json:"id"`
	Params  map[string]interface{} `json:"params,omitempty"`
}

// McpProxyResponse is the response from POST /mcp/proxy.
type McpProxyResponse struct {
	Result interface{} `json:"result"`
	Error  interface{} `json:"error,omitempty"`
}

// ---------------------------------------------------------------------------
// Audit
// ---------------------------------------------------------------------------

// AuditEvent represents an audit log entry.
type AuditEvent struct {
	ID          string                 `json:"id"`
	OrgID       string                 `json:"orgId"`
	WorkspaceID *string                `json:"workspaceId,omitempty"`
	AgentID     *string                `json:"agentId,omitempty"`
	Type        string                 `json:"type"`
	Action      string                 `json:"action"`
	Resource    *string                `json:"resource,omitempty"`
	Result      string                 `json:"result"`
	Details     map[string]interface{} `json:"details,omitempty"`
	IPAddress   *string                `json:"ipAddress,omitempty"`
	UserAgent   *string                `json:"userAgent,omitempty"`
	Timestamp   string                 `json:"timestamp"`
}

// ListAuditEventsInput contains query parameters for GET /audit/events.
type ListAuditEventsInput struct {
	OrgID       *string `json:"orgId,omitempty"`
	WorkspaceID *string `json:"workspaceId,omitempty"`
	AgentID     *string `json:"agentId,omitempty"`
	Type        *string `json:"type,omitempty"`
	DateFrom    *string `json:"dateFrom,omitempty"`
	DateTo      *string `json:"dateTo,omitempty"`
	Resource    *string `json:"resource,omitempty"`
	Result      *string `json:"result,omitempty"`
	Page        *int    `json:"page,omitempty"`
	Limit       *int    `json:"limit,omitempty"`
}

// CreateAuditReportInput is the request body for POST /audit/reports.
type CreateAuditReportInput struct {
	OrgID    string `json:"orgId"`
	DateFrom string `json:"dateFrom"`
	DateTo   string `json:"dateTo"`
}

// AuditReport is the response from POST /audit/reports.
type AuditReport struct {
	ID        string `json:"id"`
	OrgID     string `json:"orgId"`
	DateFrom  string `json:"dateFrom"`
	DateTo    string `json:"dateTo"`
	Status    string `json:"status"`
	URL       string `json:"url,omitempty"`
	CreatedAt string `json:"createdAt"`
}

// AuditMetricsInput contains query parameters for GET /audit/metrics.
type AuditMetricsInput struct {
	OrgID       string  `json:"orgId"`
	WorkspaceID *string `json:"workspaceId,omitempty"`
	AgentID     *string `json:"agentId,omitempty"`
	DateFrom    *string `json:"dateFrom,omitempty"`
	DateTo      *string `json:"dateTo,omitempty"`
}

// AuditMetricRow represents a single row of audit metrics data.
type AuditMetricRow struct {
	Day             string `json:"day"`
	OrgID           string `json:"org_id"`
	WorkspaceID     string `json:"workspace_id"`
	AgentID         string `json:"agent_id"`
	TotalActions    int    `json:"total_actions"`
	AllowedActions  int    `json:"allowed_actions"`
	DeniedActions   int    `json:"denied_actions"`
	UniqueResources int    `json:"unique_resources"`
}

// ---------------------------------------------------------------------------
// Notifications
// ---------------------------------------------------------------------------

// Notification represents a notification.
type Notification struct {
	ID             string                 `json:"id"`
	OrganizationID string                 `json:"organizationId"`
	UserID         *string                `json:"userId,omitempty"`
	EventID        *string                `json:"eventId,omitempty"`
	Type           string                 `json:"type"`
	Title          string                 `json:"title"`
	Body           string                 `json:"body"`
	Read           bool                   `json:"read"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt      string                 `json:"createdAt"`
}

// ListNotificationsInput contains query parameters for GET /notifications.
type ListNotificationsInput struct {
	OrganizationID string  `json:"organizationId"`
	UserID         *string `json:"userId,omitempty"`
	UnreadOnly     *bool   `json:"unreadOnly,omitempty"`
	Limit          *int    `json:"limit,omitempty"`
	Offset         *int    `json:"offset,omitempty"`
}

// UnreadCountInput contains query parameters for GET /notifications/unread-count.
type UnreadCountInput struct {
	OrganizationID string  `json:"organizationId"`
	UserID         *string `json:"userId,omitempty"`
}

// UnreadCountResponse is the response from GET /notifications/unread-count.
type UnreadCountResponse struct {
	Count int `json:"count"`
}

// MarkAllReadInput is the request body for PATCH /notifications/read-all.
type MarkAllReadInput struct {
	OrganizationID string  `json:"organizationId"`
	UserID         *string `json:"userId,omitempty"`
}

// ---------------------------------------------------------------------------
// Webhooks
// ---------------------------------------------------------------------------

// Webhook represents a webhook configuration.
type Webhook struct {
	ID             string   `json:"id"`
	OrganizationID string   `json:"organizationId"`
	URL            string   `json:"url"`
	EventTypes     []string `json:"eventTypes"`
	SecretHash     string   `json:"secretHash,omitempty"`
	Enabled        bool     `json:"enabled"`
	CreatedAt      string   `json:"createdAt"`
	UpdatedAt      string   `json:"updatedAt"`
}

// CreateWebhookInput is the request body for POST /webhooks.
type CreateWebhookInput struct {
	OrganizationID string   `json:"organizationId"`
	URL            string   `json:"url"`
	EventTypes     []string `json:"eventTypes"`
	Secret         string   `json:"secret"`
}

// ListWebhooksInput contains query parameters for GET /webhooks.
type ListWebhooksInput struct {
	OrganizationID string `json:"organizationId"`
}

// UpdateWebhookInput is the request body for PATCH /webhooks/:webhookId.
type UpdateWebhookInput struct {
	URL        *string  `json:"url,omitempty"`
	EventTypes []string `json:"eventTypes,omitempty"`
	Secret     *string  `json:"secret,omitempty"`
	Enabled    *bool    `json:"enabled,omitempty"`
}

// ---------------------------------------------------------------------------
// Alerts
// ---------------------------------------------------------------------------

// Alert represents an alert configuration.
type Alert struct {
	ID             string                 `json:"id"`
	OrganizationID string                 `json:"organizationId"`
	Name           string                 `json:"name"`
	EventTypes     []string               `json:"eventTypes"`
	Conditions     map[string]interface{} `json:"conditions,omitempty"`
	Channels       []string               `json:"channels"`
	Enabled        bool                   `json:"enabled"`
	CreatedAt      string                 `json:"createdAt"`
	UpdatedAt      string                 `json:"updatedAt"`
}

// CreateAlertInput is the request body for POST /alerts.
type CreateAlertInput struct {
	OrganizationID string                 `json:"organizationId"`
	Name           string                 `json:"name"`
	EventTypes     []string               `json:"eventTypes"`
	Conditions     map[string]interface{} `json:"conditions"`
	Channels       []string               `json:"channels"`
}

// ListAlertsInput contains query parameters for GET /alerts.
type ListAlertsInput struct {
	OrganizationID string `json:"organizationId"`
}

// UpdateAlertInput is the request body for PATCH /alerts/:alertId.
type UpdateAlertInput struct {
	Name       *string                `json:"name,omitempty"`
	EventTypes []string               `json:"eventTypes,omitempty"`
	Conditions map[string]interface{} `json:"conditions,omitempty"`
	Channels   []string               `json:"channels,omitempty"`
	Enabled    *bool                  `json:"enabled,omitempty"`
}

// ---------------------------------------------------------------------------
// API Keys
// ---------------------------------------------------------------------------

// APIKey represents an API key.
type APIKey struct {
	ID             string   `json:"id"`
	OrganizationID string   `json:"organizationId"`
	Name           string   `json:"name"`
	KeyPrefix      string   `json:"keyPrefix"`
	Scopes         []string `json:"scopes,omitempty"`
	CreatedBy      string   `json:"createdBy"`
	ExpiresAt      *string  `json:"expiresAt,omitempty"`
	LastUsedAt     *string  `json:"lastUsedAt,omitempty"`
	CreatedAt      string   `json:"createdAt"`
}

// CreateAPIKeyInput is the request body for POST /api-keys.
type CreateAPIKeyInput struct {
	OrganizationID string   `json:"organizationId"`
	Name           string   `json:"name"`
	Scopes         []string `json:"scopes,omitempty"`
	CreatedBy      string   `json:"createdBy"`
	ExpiresInDays  *int     `json:"expiresInDays,omitempty"`
}

// CreateAPIKeyResponse is the response from POST /api-keys.
// The response is flat: the API key fields plus a rawKey field containing the plaintext key.
type CreateAPIKeyResponse struct {
	ID             string   `json:"id"`
	OrganizationID string   `json:"organizationId"`
	Name           string   `json:"name"`
	HashedKey      string   `json:"hashedKey,omitempty"`
	Scopes         []string `json:"scopes,omitempty"`
	CreatedBy      string   `json:"createdBy"`
	ExpiresAt      *string  `json:"expiresAt,omitempty"`
	LastUsedAt     *string  `json:"lastUsedAt,omitempty"`
	CreatedAt      string   `json:"createdAt"`
	RawKey         string   `json:"rawKey"`
}

// ListAPIKeysInput contains query parameters for GET /api-keys.
type ListAPIKeysInput struct {
	OrganizationID string `json:"organizationId"`
}

// ---------------------------------------------------------------------------
// Organizations
// ---------------------------------------------------------------------------

// Organization represents an organization.
type Organization struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// CreateOrganizationInput is the request body for POST /organizations.
type CreateOrganizationInput struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// ListOrganizationsInput contains query parameters for GET /organizations.
type ListOrganizationsInput struct {
	Page  *int `json:"page,omitempty"`
	Limit *int `json:"limit,omitempty"`
}

// ---------------------------------------------------------------------------
// Workspaces
// ---------------------------------------------------------------------------

// Workspace represents a workspace within an organization.
type Workspace struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organizationId"`
	Name           string `json:"name"`
	Slug           string `json:"slug"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
}

// CreateWorkspaceInput is the request body for POST /workspaces.
type CreateWorkspaceInput struct {
	OrganizationID string `json:"organizationId"`
	Name           string `json:"name"`
	Slug           string `json:"slug"`
}

// ListWorkspacesInput contains query parameters for GET /workspaces.
type ListWorkspacesInput struct {
	OrganizationID string `json:"organizationId"`
	Page           *int   `json:"page,omitempty"`
	Limit          *int   `json:"limit,omitempty"`
}
