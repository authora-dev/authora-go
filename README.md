# Authora Go SDK

Official Go client for the [Authora](https://authora.dev) API -- identity and authorization for AI agents.

- Go 1.21+
- Zero external dependencies (`net/http` + `encoding/json` only)
- `context.Context` on every method
- Functional options for client configuration
- Typed errors with helper predicates

## Installation

```bash
go get github.com/authora-dev/authora-go@v0.1.0
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    authora "github.com/authora-dev/authora-go"
)

func main() {
    client := authora.NewClient("authora_live_...",
        authora.WithBaseURL("https://api.authora.dev/api/v1"),
        authora.WithTimeout(30*time.Second),
    )

    ctx := context.Background()

    // Create an agent
    resp, err := client.Agents.Create(ctx, &authora.CreateAgentInput{
        WorkspaceID: "ws_abc123",
        Name:        "my-agent",
        CreatedBy:   "user_xyz",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Agent %s created. API Key: %s\n", resp.Agent.ID, resp.APIKey)

    // Check a permission
    check, err := client.Permissions.Check(ctx, &authora.CheckPermissionInput{
        AgentID:  resp.Agent.ID,
        Resource: "files:*",
        Action:   "read",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Allowed: %t\n", check.Allowed)
}
```

## Client Configuration

```go
// Default configuration (base URL: https://api.authora.dev/api/v1, timeout: 30s)
client := authora.NewClient("authora_live_...")

// Custom base URL
client := authora.NewClient("authora_live_...",
    authora.WithBaseURL("https://custom.api.dev/api/v1"),
)

// Custom timeout
client := authora.NewClient("authora_live_...",
    authora.WithTimeout(10 * time.Second),
)

// Custom HTTP client (timeout option ignored when using this)
client := authora.NewClient("authora_live_...",
    authora.WithHTTPClient(&http.Client{
        Timeout:   60 * time.Second,
        Transport: customTransport,
    }),
)
```

## API Reference

All methods accept `context.Context` as the first parameter.

### Agents

```go
client.Agents.Create(ctx, &CreateAgentInput{...})          // POST   /agents
client.Agents.List(ctx, &ListAgentsInput{...})              // GET    /agents
client.Agents.Get(ctx, agentID)                             // GET    /agents/:agentId
client.Agents.Verify(ctx, agentID)                          // GET    /agents/:agentId/verify (no auth)
client.Agents.Activate(ctx, agentID, &ActivateAgentInput{}) // POST   /agents/:agentId/activate
client.Agents.Suspend(ctx, agentID)                         // POST   /agents/:agentId/suspend
client.Agents.Revoke(ctx, agentID)                          // POST   /agents/:agentId/revoke
client.Agents.RotateKey(ctx, agentID, &RotateKeyInput{})    // POST   /agents/:agentId/rotate-key
```

### Roles

```go
client.Roles.Create(ctx, &CreateRoleInput{...})             // POST   /roles
client.Roles.List(ctx, &ListRolesInput{...})                // GET    /roles
client.Roles.Get(ctx, roleID)                               // GET    /roles/:roleId
client.Roles.Update(ctx, roleID, &UpdateRoleInput{...})     // PATCH  /roles/:roleId
client.Roles.Delete(ctx, roleID)                            // DELETE /roles/:roleId
```

### Agent Role Assignments

```go
client.Roles.AssignToAgent(ctx, agentID, &AssignRoleInput{})        // POST   /agents/:agentId/roles
client.Roles.UnassignFromAgent(ctx, agentID, roleID)                // DELETE /agents/:agentId/roles/:roleId
client.Roles.ListAgentRoles(ctx, agentID)                           // GET    /agents/:agentId/roles
```

### Permissions

```go
client.Permissions.Check(ctx, &CheckPermissionInput{...})       // POST /permissions/check
client.Permissions.CheckBatch(ctx, &BatchCheckInput{...})       // POST /permissions/check-batch
client.Permissions.Effective(ctx, agentID)                      // GET  /agents/:agentId/permissions
```

### Delegations

```go
client.Delegations.Create(ctx, &CreateDelegationInput{...})         // POST /delegations
client.Delegations.Get(ctx, delegationID)                           // GET  /delegations/:delegationId
client.Delegations.Revoke(ctx, delegationID)                        // POST /delegations/:delegationId/revoke
client.Delegations.Verify(ctx, &VerifyDelegationInput{...})         // POST /delegations/verify
client.Delegations.List(ctx, &ListDelegationsInput{...})            // GET  /delegations
client.Delegations.ListByAgent(ctx, agentID, &ListAgentDel...{})    // GET  /agents/:agentId/delegations
```

### Policies

```go
client.Policies.Create(ctx, &CreatePolicyInput{...})            // POST   /policies
client.Policies.List(ctx, &ListPoliciesInput{...})              // GET    /policies
client.Policies.Update(ctx, policyID, &UpdatePolicyInput{...})  // PATCH  /policies/:policyId
client.Policies.Delete(ctx, policyID)                           // DELETE /policies/:policyId
client.Policies.Simulate(ctx, &SimulatePolicyInput{...})        // POST   /policies/simulate
client.Policies.Evaluate(ctx, &EvaluatePolicyInput{...})        // POST   /policies/evaluate
```

### MCP (Model Context Protocol)

```go
client.Mcp.RegisterServer(ctx, &RegisterMcpServerInput{...})         // POST  /mcp/servers
client.Mcp.ListServers(ctx, &ListMcpServersInput{...})               // GET   /mcp/servers
client.Mcp.GetServer(ctx, serverID)                                  // GET   /mcp/servers/:serverId
client.Mcp.UpdateServer(ctx, serverID, &UpdateMcpServerInput{...})   // PATCH /mcp/servers/:serverId
client.Mcp.ListTools(ctx, serverID)                                  // GET   /mcp/servers/:serverId/tools
client.Mcp.RegisterTool(ctx, serverID, &RegisterMcpToolInput{...})   // POST  /mcp/servers/:serverId/tools
client.Mcp.Proxy(ctx, &McpProxyInput{...})                           // POST  /mcp/proxy
```

### Audit

```go
client.Audit.ListEvents(ctx, &ListAuditEventsInput{...})    // GET  /audit/events
client.Audit.GetEvent(ctx, eventID)                         // GET  /audit/events/:eventId
client.Audit.CreateReport(ctx, &CreateAuditReportInput{})   // POST /audit/reports
client.Audit.GetMetrics(ctx, &AuditMetricsInput{...})       // GET  /audit/metrics
```

### Notifications

```go
client.Notifications.List(ctx, &ListNotificationsInput{...})    // GET   /notifications
client.Notifications.UnreadCount(ctx, &UnreadCountInput{...})   // GET   /notifications/unread-count
client.Notifications.MarkRead(ctx, notificationID)              // PATCH /notifications/:notificationId/read
client.Notifications.MarkAllRead(ctx, &MarkAllReadInput{...})   // PATCH /notifications/read-all
```

### Webhooks

```go
client.Webhooks.Create(ctx, &CreateWebhookInput{...})                // POST   /webhooks
client.Webhooks.List(ctx, &ListWebhooksInput{...})                   // GET    /webhooks
client.Webhooks.Update(ctx, webhookID, &UpdateWebhookInput{...})     // PATCH  /webhooks/:webhookId
client.Webhooks.Delete(ctx, webhookID)                               // DELETE /webhooks/:webhookId
```

### Alerts

```go
client.Alerts.Create(ctx, &CreateAlertInput{...})              // POST   /alerts
client.Alerts.List(ctx, &ListAlertsInput{...})                 // GET    /alerts
client.Alerts.Update(ctx, alertID, &UpdateAlertInput{...})     // PATCH  /alerts/:alertId
client.Alerts.Delete(ctx, alertID)                             // DELETE /alerts/:alertId
```

### API Keys

```go
client.APIKeys.Create(ctx, &CreateAPIKeyInput{...})   // POST   /api-keys
client.APIKeys.List(ctx, &ListAPIKeysInput{...})      // GET    /api-keys
client.APIKeys.Revoke(ctx, keyID)                     // DELETE /api-keys/:keyId
```

### Organizations

```go
client.Organizations.Create(ctx, &CreateOrganizationInput{...})   // POST /organizations
client.Organizations.Get(ctx, orgID)                              // GET  /organizations/:orgId
client.Organizations.List(ctx, &ListOrganizationsInput{...})      // GET  /organizations
```

### Workspaces

```go
client.Workspaces.Create(ctx, &CreateWorkspaceInput{...})   // POST /workspaces
client.Workspaces.Get(ctx, workspaceID)                     // GET  /workspaces/:wsId
client.Workspaces.List(ctx, &ListWorkspacesInput{...})      // GET  /workspaces
```

## Error Handling

All API errors are returned as `*authora.AuthoraError` with status code, message, and optional error code.

```go
agent, err := client.Agents.Get(ctx, "agt_nonexistent")
if err != nil {
    if authora.IsNotFoundError(err) {
        fmt.Println("Agent not found")
    } else if authora.IsAuthenticationError(err) {
        fmt.Println("Invalid API key")
    } else if authora.IsRateLimitError(err) {
        fmt.Println("Rate limited, retry later")
    } else if authora.IsForbiddenError(err) {
        fmt.Println("Insufficient permissions")
    } else if authora.IsValidationError(err) {
        fmt.Println("Bad request")
    } else {
        fmt.Printf("Unexpected error: %v\n", err)
    }
}
```

You can also inspect the error directly:

```go
var apiErr *authora.AuthoraError
if errors.As(err, &apiErr) {
    fmt.Printf("Status: %d, Code: %s, Message: %s\n",
        apiErr.StatusCode, apiErr.Code, apiErr.Message)
}
```

## Optional Fields

Optional fields use pointer types. Use the address-of operator to set them:

```go
desc := "My agent description"
agent, err := client.Agents.Create(ctx, &authora.CreateAgentInput{
    WorkspaceID: "ws_abc123",
    Name:        "my-agent",
    CreatedBy:   "user_xyz",
    Description: &desc,
    Tags:        []string{"production", "billing"},
})
```

## Pagination

List endpoints that support pagination return `*authora.PaginatedResponse[T]`:

```go
page := 1
limit := 25
agents, err := client.Agents.List(ctx, &authora.ListAgentsInput{
    WorkspaceID: "ws_abc123",
    Page:        &page,
    Limit:       &limit,
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Total: %d, Page %d/%d\n", agents.Total, agents.Page, agents.TotalPages)
for _, a := range agents.Data {
    fmt.Printf("  - %s (%s)\n", a.Name, a.Status)
}
```

## License

MIT
