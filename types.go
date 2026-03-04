package authora

import (
	"encoding/json"
	"net/http"
	"time"
)

type PaginatedResponse[T any] struct {
	Items      []T  `json:"items"`
	Total      int  `json:"total"`
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	TotalPages int  `json:"totalPages"`
	HasMore    bool `json:"hasMore"`
}

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

type CreateAgentInput struct {
	WorkspaceID string                 `json:"workspaceId"`
	Name        string                 `json:"name"`
	Description *string                `json:"description,omitempty"`
	CreatedBy   string                 `json:"createdBy"`
	ExpiresIn   *string                `json:"expiresIn,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type CreateAgentResponse = Agent

type ListAgentsInput struct {
	WorkspaceID string  `json:"workspaceId"`
	Status      *string `json:"status,omitempty"`
	Page        *int    `json:"page,omitempty"`
	Limit       *int    `json:"limit,omitempty"`
}

type VerifyAgentResponse struct {
	Valid  bool   `json:"valid"`
	Status string `json:"status"`
	Agent  *Agent `json:"agent,omitempty"`
}

type ActivateAgentInput struct {
	PublicKey string `json:"publicKey"`
}

type RotateKeyInput struct {
	PublicKey string `json:"publicKey"`
}

type RotateKeyResponse struct {
	Agent     Agent  `json:"agent"`
	NewAPIKey string `json:"newApiKey"`
}

type Role struct {
	ID                 string   `json:"id"`
	WorkspaceID        string   `json:"workspaceId"`
	Name               string   `json:"name"`
	Description        *string  `json:"description,omitempty"`
	Permissions        []string `json:"permissions"`
	DenyPermissions    []string `json:"denyPermissions,omitempty"`
	Stage              *string  `json:"stage,omitempty"`
	MaxSessionDuration *int     `json:"maxSessionDuration,omitempty"`
	IsBuiltin          bool     `json:"isBuiltin"`
	CreatedAt          string   `json:"createdAt"`
	UpdatedAt          string   `json:"updatedAt"`
}

type CreateRoleInput struct {
	WorkspaceID        string   `json:"workspaceId"`
	Name               string   `json:"name"`
	Description        *string  `json:"description,omitempty"`
	Permissions        []string `json:"permissions"`
	DenyPermissions    []string `json:"denyPermissions,omitempty"`
	Stage              *string  `json:"stage,omitempty"`
	MaxSessionDuration *int     `json:"maxSessionDuration,omitempty"`
}

type ListRolesInput struct {
	WorkspaceID string `json:"workspaceId"`
	Page        *int   `json:"page,omitempty"`
	Limit       *int   `json:"limit,omitempty"`
}

type UpdateRoleInput struct {
	Name               *string  `json:"name,omitempty"`
	Description        *string  `json:"description,omitempty"`
	Permissions        []string `json:"permissions,omitempty"`
	DenyPermissions    []string `json:"denyPermissions,omitempty"`
	Stage              *string  `json:"stage,omitempty"`
	MaxSessionDuration *int     `json:"maxSessionDuration,omitempty"`
}

type AgentRole struct {
	ID        string     `json:"id"`
	AgentID   string     `json:"agentId"`
	RoleID    string     `json:"roleId"`
	GrantedBy *string    `json:"grantedBy,omitempty"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
	CreatedAt string     `json:"createdAt"`
	Role      *Role      `json:"role,omitempty"`
}

type AssignRoleInput struct {
	RoleID    string     `json:"roleId"`
	GrantedBy *string    `json:"grantedBy,omitempty"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
}

type CheckPermissionInput struct {
	AgentID  string                 `json:"agentId"`
	Resource string                 `json:"resource"`
	Action   string                 `json:"action"`
	Context  map[string]interface{} `json:"context,omitempty"`
}

type CheckPermissionResponse struct {
	Allowed bool     `json:"allowed"`
	Reason  *string  `json:"reason,omitempty"`
	Matched []string `json:"matched,omitempty"`
	Denied  []string `json:"denied,omitempty"`
}

type BatchCheckItem struct {
	Resource string                 `json:"resource"`
	Action   string                 `json:"action"`
	Context  map[string]interface{} `json:"context,omitempty"`
}

type BatchCheckInput struct {
	AgentID string           `json:"agentId"`
	Checks  []BatchCheckItem `json:"checks"`
}

type BatchCheckResult struct {
	Resource string   `json:"resource"`
	Action   string   `json:"action"`
	Allowed  bool     `json:"allowed"`
	Reason   *string  `json:"reason,omitempty"`
	Matched  []string `json:"matched,omitempty"`
	Denied   []string `json:"denied,omitempty"`
}

type BatchCheckResponse struct {
	Results []BatchCheckResult `json:"results"`
}

type EffectivePermission struct {
	Permission string `json:"permission"`
	Source     string `json:"source"`
	RoleID     string `json:"roleId,omitempty"`
	RoleName   string `json:"roleName,omitempty"`
}

type EffectivePermissionsResponse struct {
	AgentID     string                `json:"agentId"`
	Permissions []EffectivePermission `json:"permissions"`
}

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

type CreateDelegationInput struct {
	IssuerAgentID string                 `json:"issuerAgentId"`
	TargetAgentID string                 `json:"targetAgentId"`
	Permissions   []string               `json:"permissions"`
	Constraints   map[string]interface{} `json:"constraints,omitempty"`
	ExpiresIn     *string                `json:"expiresIn,omitempty"`
}

type VerifyDelegationInput struct {
	DelegationID string `json:"delegationId"`
}

type VerifyDelegationResponse struct {
	Valid      bool        `json:"valid"`
	Delegation *Delegation `json:"delegation,omitempty"`
	Reason     *string     `json:"reason,omitempty"`
}

type ListDelegationsInput struct {
	Status *string `json:"status,omitempty"`
	Page   *int    `json:"page,omitempty"`
	Limit  *int    `json:"limit,omitempty"`
}

type ListAgentDelegationsInput struct {
	Direction *string `json:"direction,omitempty"`
	Status    *string `json:"status,omitempty"`
	Page      *int    `json:"page,omitempty"`
	Limit     *int    `json:"limit,omitempty"`
}

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

type ListPoliciesInput struct {
	WorkspaceID string `json:"workspaceId"`
	Page        *int   `json:"page,omitempty"`
	Limit       *int   `json:"limit,omitempty"`
}

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

type SimulatePolicyInput struct {
	WorkspaceID string                 `json:"workspaceId"`
	AgentID     string                 `json:"agentId"`
	Resource    string                 `json:"resource"`
	Action      string                 `json:"action"`
	Context     map[string]interface{} `json:"context,omitempty"`
}

type SimulatePolicyResponse struct {
	Decision        string   `json:"decision"`
	MatchedPolicies []Policy `json:"matchedPolicies,omitempty"`
	Reason          *string  `json:"reason,omitempty"`
}

type EvaluatePolicyInput struct {
	WorkspaceID string                 `json:"workspaceId"`
	AgentID     string                 `json:"agentId"`
	Resource    string                 `json:"resource"`
	Action      string                 `json:"action"`
	Context     map[string]interface{} `json:"context,omitempty"`
}

type EvaluatePolicyResponse struct {
	Allowed         bool     `json:"allowed"`
	MatchedPolicies []Policy `json:"matchedPolicies,omitempty"`
	Reason          *string  `json:"reason,omitempty"`
}

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

type RegisterMcpServerInput struct {
	WorkspaceID string                 `json:"workspaceId"`
	Name        string                 `json:"name"`
	Description *string                `json:"description,omitempty"`
	URL         string                 `json:"url"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type ListMcpServersInput struct {
	WorkspaceID string `json:"workspaceId"`
	Page        *int   `json:"page,omitempty"`
	Limit       *int   `json:"limit,omitempty"`
}

type UpdateMcpServerInput struct {
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	URL         *string                `json:"url,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type McpTool struct {
	ID          string                 `json:"id"`
	ServerID    string                 `json:"serverId"`
	Name        string                 `json:"name"`
	Description *string                `json:"description,omitempty"`
	InputSchema map[string]interface{} `json:"inputSchema,omitempty"`
	CreatedAt   string                 `json:"createdAt"`
	UpdatedAt   string                 `json:"updatedAt"`
}

type RegisterMcpToolInput struct {
	Name        string                 `json:"name"`
	Description *string                `json:"description,omitempty"`
	InputSchema map[string]interface{} `json:"inputSchema,omitempty"`
}

type McpProxyInput struct {
	ServerID string                 `json:"-"`
	Method   string                 `json:"-"`
	Params   map[string]interface{} `json:"-"`
}

type mcpProxyJsonRpc struct {
	Jsonrpc string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	ID      int                    `json:"id"`
	Params  map[string]interface{} `json:"params,omitempty"`
}

type McpProxyResponse struct {
	Result interface{} `json:"result"`
	Error  interface{} `json:"error,omitempty"`
}

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

type CreateAuditReportInput struct {
	OrgID    string `json:"orgId"`
	DateFrom string `json:"dateFrom"`
	DateTo   string `json:"dateTo"`
}

type AuditReport struct {
	ID        string `json:"id"`
	OrgID     string `json:"orgId"`
	DateFrom  string `json:"dateFrom"`
	DateTo    string `json:"dateTo"`
	Status    string `json:"status"`
	URL       string `json:"url,omitempty"`
	CreatedAt string `json:"createdAt"`
}

type AuditMetricsInput struct {
	OrgID       string  `json:"orgId"`
	WorkspaceID *string `json:"workspaceId,omitempty"`
	AgentID     *string `json:"agentId,omitempty"`
	DateFrom    *string `json:"dateFrom,omitempty"`
	DateTo      *string `json:"dateTo,omitempty"`
}

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

type ListNotificationsInput struct {
	OrganizationID string  `json:"organizationId"`
	UserID         *string `json:"userId,omitempty"`
	UnreadOnly     *bool   `json:"unreadOnly,omitempty"`
	Limit          *int    `json:"limit,omitempty"`
	Offset         *int    `json:"offset,omitempty"`
}

type UnreadCountInput struct {
	OrganizationID string  `json:"organizationId"`
	UserID         *string `json:"userId,omitempty"`
}

type UnreadCountResponse struct {
	Count int `json:"count"`
}

type MarkAllReadInput struct {
	OrganizationID string  `json:"organizationId"`
	UserID         *string `json:"userId,omitempty"`
}

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

type CreateWebhookInput struct {
	OrganizationID string   `json:"organizationId"`
	URL            string   `json:"url"`
	EventTypes     []string `json:"eventTypes"`
	Secret         string   `json:"secret"`
}

type ListWebhooksInput struct {
	OrganizationID string `json:"organizationId"`
}

type UpdateWebhookInput struct {
	URL        *string  `json:"url,omitempty"`
	EventTypes []string `json:"eventTypes,omitempty"`
	Secret     *string  `json:"secret,omitempty"`
	Enabled    *bool    `json:"enabled,omitempty"`
}

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

type CreateAlertInput struct {
	OrganizationID string                 `json:"organizationId"`
	Name           string                 `json:"name"`
	EventTypes     []string               `json:"eventTypes"`
	Conditions     map[string]interface{} `json:"conditions"`
	Channels       []string               `json:"channels"`
}

type ListAlertsInput struct {
	OrganizationID string `json:"organizationId"`
}

type UpdateAlertInput struct {
	Name       *string                `json:"name,omitempty"`
	EventTypes []string               `json:"eventTypes,omitempty"`
	Conditions map[string]interface{} `json:"conditions,omitempty"`
	Channels   []string               `json:"channels,omitempty"`
	Enabled    *bool                  `json:"enabled,omitempty"`
}

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

type CreateAPIKeyInput struct {
	OrganizationID string   `json:"organizationId"`
	Name           string   `json:"name"`
	Scopes         []string `json:"scopes,omitempty"`
	CreatedBy      string   `json:"createdBy"`
	ExpiresInDays  *int     `json:"expiresInDays,omitempty"`
}

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

type ListAPIKeysInput struct {
	OrganizationID string `json:"organizationId"`
}

type Organization struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type CreateOrganizationInput struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type ListOrganizationsInput struct {
	Page  *int `json:"page,omitempty"`
	Limit *int `json:"limit,omitempty"`
}

type Workspace struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organizationId"`
	Name           string `json:"name"`
	Slug           string `json:"slug"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
}

type CreateWorkspaceInput struct {
	OrganizationID string `json:"organizationId"`
	Name           string `json:"name"`
	Slug           string `json:"slug"`
}

type ListWorkspacesInput struct {
	OrganizationID string `json:"organizationId"`
	Page           *int   `json:"page,omitempty"`
	Limit          *int   `json:"limit,omitempty"`
}

type SignedResponse struct {
	Data       json.RawMessage
	StatusCode int
	Headers    http.Header
}

type ToolCallParams struct {
	ToolName        string
	Arguments       map[string]interface{}
	Method          string
	ID              interface{}
	DelegationToken string
}

type DelegationConstraints struct {
	MaxDepth       *int     `json:"maxDepth,omitempty"`
	ExpiresAt      *string  `json:"expiresAt,omitempty"`
	SingleUse      *bool    `json:"singleUse,omitempty"`
	AllowedTargets []string `json:"allowedTargets,omitempty"`
}

type EffectivePermissionsData struct {
	AgentID         string   `json:"agentId"`
	Permissions     []string `json:"permissions"`
	DenyPermissions []string `json:"denyPermissions"`
}
