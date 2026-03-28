# Authora Go SDK

Authorization for AI agents -- identity, permissions, and delegation management.

## Quick Start

```go
// go get github.com/authora-dev/authora-go@v0.4.2
package main

import (
    "context"
    "fmt"
    "log"

    authora "github.com/authora-dev/authora-go"
)

func main() {
    client := authora.NewClient("authora_live_...")
    ctx := context.Background()

    // Check a permission
    check, err := client.Permissions.Check(ctx, &authora.CheckPermissionInput{
        AgentID: "agt_abc", Resource: "files:*", Action: "read",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Allowed: %t\n", check.Allowed)
}
```

## Installation

```bash
go get github.com/authora-dev/authora-go@v0.4.2
```

## Getting Credentials

**Automatic (IDE agents):** If you use Claude Code, Cursor, or OpenCode, credentials are created automatically on first run via browser sign-in. See [self-onboarding instructions](https://authora.dev/llms-onboard.txt).

**Manual:** Sign up at [authora.dev/get-started](https://authora.dev/get-started), then find your credentials:

| Value | Format | Where to find it |
|---|---|---|
| **API Key** | `authora_live_...` | [Dashboard](https://client.authora.dev) > API Keys |
| **Workspace ID** | `ws_...` | [Dashboard](https://client.authora.dev) > Settings |
| **User ID** | `usr_...` | [Dashboard](https://client.authora.dev) > Settings |
| **Organization ID** | `org_...` | [Dashboard](https://client.authora.dev) > Settings |

**Environment variables (Docker/CI):** Set `AUTHORA_API_KEY`, `AUTHORA_AGENT_ID`, `AUTHORA_ORG_ID`, `AUTHORA_WORKSPACE_ID`.

The `CreatedBy` parameter used when creating agents or API keys is your **User ID** (`usr_...`).

## Features

- Go 1.21+
- Zero external dependencies (`net/http` + `encoding/json` only)
- `context.Context` on every method
- Functional options for client configuration
- Typed errors with helper predicates

## Extended Quick Start

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
    client := authora.NewClient("authora_live_...", // from Account > API Keys
        authora.WithBaseURL("https://api.authora.dev/api/v1"),
        authora.WithTimeout(30*time.Second),
    )

    ctx := context.Background()

    // Create an agent
    resp, err := client.Agents.Create(ctx, &authora.CreateAgentInput{
        WorkspaceID: "ws_...",    // from Account > Profile
        Name:        "my-agent",
        CreatedBy:   "usr_...",   // your User ID
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

## Edge Endpoints

For high-availability scenarios, Authora provides an edge proxy at `https://edge.authora.dev` powered by Cloudflare Workers. Agent identity verification, JWT validation, and public key lookups are served from globally distributed edge caches with 24-hour survivability if the origin is unreachable. The edge proxy runs in parallel with the primary API -- no client changes required.

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
client.Policies.AttachToTarget(ctx, &AttachPolicyInput{...})         // POST   /policies/attachments
client.Policies.DetachFromTarget(ctx, &DetachPolicyInput{...})       // POST   /policies/detach
client.Policies.ListAttachments(ctx, &ListAttachmentsInput{...})     // GET    /policies/attachments
client.Policies.ListPolicyTargets(ctx, policyID)                     // GET    /policies/:policyId/attachments
client.Policies.AddPermission(ctx, &AddPermissionInput{...})         // POST   /policies/add-permission
client.Policies.RemovePermission(ctx, &RemovePermissionInput{...})   // POST   /policies/remove-permission
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

### Approvals

```go
client.Approvals.Create(ctx, &CreateApprovalInput{...})                   // POST   /approvals
client.Approvals.List(ctx, &ListApprovalsInput{...})                      // GET    /approvals
client.Approvals.Get(ctx, challengeID)                                    // GET    /approvals/:id
client.Approvals.Decide(ctx, challengeID, &DecideApprovalInput{...})      // POST   /approvals/:id/decide
client.Approvals.BulkDecide(ctx, &BulkDecideInput{...})                   // POST   /approvals/bulk-decide
client.Approvals.Stats(ctx)                                               // GET    /approvals/stats
client.Approvals.ListEscalationRules(ctx)                                 // GET    /approvals/escalation-rules
client.Approvals.CreateEscalationRule(ctx, &CreateEscalationRuleInput{})  // POST   /approvals/escalation-rules
client.Approvals.UpdateEscalationRule(ctx, ruleID, &UpdateEscalationRuleInput{})  // PATCH  /approvals/escalation-rules/:id
client.Approvals.DeleteEscalationRule(ctx, ruleID)                        // DELETE /approvals/escalation-rules/:id
```

### Credits

```go
client.Credits.Balance(ctx)                                          // GET  /credits
client.Credits.Transactions(ctx, &ListTransactionsInput{...})        // GET  /credits/transactions
client.Credits.Checkout(ctx, &CreditCheckoutInput{...})              // POST /credits/checkout
```

### User Delegations

```go
client.UserDelegations.Create(ctx, &CreateGrantInput{...})                  // POST /user-delegations
client.UserDelegations.Get(ctx, grantID)                                    // GET  /user-delegations/:grantId
client.UserDelegations.ListByUser(ctx, userID, &ListGrantsInput{...})       // GET  /user-delegations/by-user/:userId
client.UserDelegations.ListByAgent(ctx, agentID, &ListGrantsInput{...})     // GET  /user-delegations/by-agent/:agentId
client.UserDelegations.Revoke(ctx, grantID, &RevokeGrantInput{...})         // POST /user-delegations/:grantId/revoke
client.UserDelegations.IssueToken(ctx, grantID, &IssueTokenInput{...})      // POST /user-delegations/:grantId/token
client.UserDelegations.VerifyToken(ctx, &VerifyTokenInput{...})             // POST /user-delegations/tokens/verify
client.UserDelegations.CreateTrust(ctx, &CreateTrustInput{...})             // POST /user-delegations/trust
client.UserDelegations.ListTrust(ctx, orgID)                                // GET  /user-delegations/trust/by-org/:orgId
client.UserDelegations.ApproveTrust(ctx, trustID)                           // POST /user-delegations/trust/:trustId/approve
client.UserDelegations.RevokeTrust(ctx, trustID)                            // POST /user-delegations/trust/:trustId/revoke
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
    WorkspaceID: "ws_...",    // from Account > Profile
    Name:        "my-agent",
    CreatedBy:   "usr_...",   // your User ID
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

## Agent Runtime

The `AgentRuntime` provides a full agent runtime with Ed25519 signed requests, thread-safe permission caching, delegation, and MCP tool calls.

```go
package main

import (
    "context"
    "fmt"
    "log"

    authora "github.com/authora-dev/authora-go"
)

func main() {
    client := authora.NewClient("authora_live_...")
    ctx := context.Background()

    // Create + activate an agent (generates Ed25519 keypair locally)
    runtime, keyPair, err := client.CreateAgent(ctx, &authora.CreateAgentInput{
        WorkspaceID: "ws_...",        // from Account > Profile
        Name:        "data-processor",
        CreatedBy:   "usr_...",       // your User ID
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Agent: %s, Public Key: %s\n", runtime.AgentID(), keyPair.PublicKey)

    // All requests are Ed25519-signed automatically
    profile, err := runtime.GetProfile(ctx)
    doc, err := runtime.GetIdentityDocument(ctx)

    // Server-side permission check
    check, err := runtime.CheckPermission(ctx, "files:read", "read", nil)

    // Client-side cached check (deny-first, 5-minute TTL, sync.RWMutex)
    allowed, err := runtime.HasPermission(ctx, "mcp:server1:tool.query")
    if allowed {
        result, err := runtime.CallTool(ctx, &authora.ToolCallParams{
            ToolName:  "query",
            Arguments: map[string]interface{}{"sql": "SELECT 1"},
        })
    }

    // Delegate permissions
    delegation, err := runtime.Delegate(ctx, "agent_...",
        []string{"files:read"},
        &authora.DelegationConstraints{ExpiresIn: "1h"},
    )

    // Key rotation
    updatedAgent, newKeyPair, err := runtime.RotateKey(ctx)

    // Lifecycle
    _, err = runtime.Suspend(ctx)
    _, _, err = runtime.Reactivate(ctx)
    _, err = runtime.Revoke(ctx)
}
```

## Cryptography

Ed25519 key generation, signing, and verification via Go stdlib `crypto/ed25519`.

```go
import authora "github.com/authora-dev/authora-go"

// Generate Ed25519 keypair (base64url encoded)
keyPair, err := authora.GenerateKeyPair()

// Sign and verify
sig, err := authora.Sign("hello world", keyPair.PrivateKey)
valid := authora.Verify("hello world", sig, keyPair.PublicKey)

// Build canonical signature payload
payload := authora.BuildSignaturePayload("POST", "/api/v1/agents", timestamp, &body)
```

## License

MIT
