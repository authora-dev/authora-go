//go:build integration

package authora

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

const (
	testAPIKey      = "authora_live_076270f52d3fc0fe9af9d08fe49b2803eb8b64ba5132fc76"
	testBaseURL     = "https://api.authora.dev/api/v1"
	testWorkspaceID = "ws_a7067ccce35d36b5"
	testOrgID       = "org_92582b4a512e52ff"
)

func newTestClient() *Client {
	return NewClient(testAPIKey, WithBaseURL(testBaseURL))
}

func ptr[T any](v T) *T {
	return &v
}

func uniqueName(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

// ---------------------------------------------------------------------------
// 1. Agent Lifecycle
// ---------------------------------------------------------------------------

func TestAgentLifecycle(t *testing.T) {
	client := newTestClient()
	ctx := context.Background()
	var agentID string

	t.Cleanup(func() {
		if agentID != "" {
			_, _ = client.Agents.Revoke(ctx, agentID)
		}
	})

	t.Run("Create", func(t *testing.T) {
		agent, err := client.Agents.Create(ctx, &CreateAgentInput{
			WorkspaceID: testWorkspaceID,
			Name:        uniqueName("go-integ-agent"),
			CreatedBy:   "integration-test",
		})
		if err != nil {
			t.Fatalf("Create agent: %v", err)
		}
		if agent.ID == "" {
			t.Fatal("expected non-empty agent ID")
		}
		if agent.WorkspaceID != testWorkspaceID {
			t.Fatalf("expected workspaceId %s, got %s", testWorkspaceID, agent.WorkspaceID)
		}
		agentID = agent.ID
		t.Logf("created agent: %s", agentID)
	})

	t.Run("Get", func(t *testing.T) {
		if agentID == "" {
			t.Skip("no agent created")
		}
		agent, err := client.Agents.Get(ctx, agentID)
		if err != nil {
			t.Fatalf("Get agent: %v", err)
		}
		if agent.ID != agentID {
			t.Fatalf("expected agent ID %s, got %s", agentID, agent.ID)
		}
	})

	t.Run("List", func(t *testing.T) {
		list, err := client.Agents.List(ctx, &ListAgentsInput{
			WorkspaceID: testWorkspaceID,
			Limit:       ptr(10),
		})
		if err != nil {
			t.Fatalf("List agents: %v", err)
		}
		if list.Total < 1 {
			t.Fatal("expected at least 1 agent in total count")
		}
		if len(list.Items) < 1 {
			t.Fatal("expected at least 1 agent in Items")
		}
	})

	t.Run("Revoke", func(t *testing.T) {
		if agentID == "" {
			t.Skip("no agent created")
		}
		agent, err := client.Agents.Revoke(ctx, agentID)
		if err != nil {
			t.Fatalf("Revoke agent: %v", err)
		}
		if agent.Status != "REVOKED" {
			t.Fatalf("expected status 'REVOKED', got %q", agent.Status)
		}
		// Already cleaned up; prevent double-revoke.
		agentID = ""
	})
}

// ---------------------------------------------------------------------------
// 1b. Agent Security Lifecycle (activate, suspend, rotate-key, verify)
// ---------------------------------------------------------------------------

func TestAgentSecurityLifecycle(t *testing.T) {
	client := newTestClient()
	ctx := context.Background()
	var agentID string

	t.Cleanup(func() {
		if agentID != "" {
			_, _ = client.Agents.Revoke(ctx, agentID)
		}
	})

	// Create
	t.Run("Create", func(t *testing.T) {
		agent, err := client.Agents.Create(ctx, &CreateAgentInput{
			WorkspaceID: testWorkspaceID,
			Name:        uniqueName("go-security-agent"),
			CreatedBy:   "integration-test",
		})
		if err != nil {
			t.Fatalf("Create: %v", err)
		}
		agentID = agent.ID
		if agent.Status != "PENDING" {
			t.Fatalf("expected PENDING, got %q", agent.Status)
		}
	})

	// Activate with public key
	t.Run("Activate", func(t *testing.T) {
		if agentID == "" {
			t.Skip("no agent")
		}
		agent, err := client.Agents.Activate(ctx, agentID, &ActivateAgentInput{
			PublicKey: "test-pubkey-go-" + fmt.Sprintf("%d", time.Now().UnixNano()),
		})
		if err != nil {
			t.Fatalf("Activate: %v", err)
		}
		if agent.Status != "ACTIVE" {
			t.Fatalf("expected ACTIVE, got %q", agent.Status)
		}
	})

	// Rotate key (must be ACTIVE)
	t.Run("RotateKey", func(t *testing.T) {
		if agentID == "" {
			t.Skip("no agent")
		}
		resp, err := client.Agents.RotateKey(ctx, agentID, &RotateKeyInput{
			PublicKey: "rotated-pubkey-go-" + fmt.Sprintf("%d", time.Now().UnixNano()),
		})
		if err != nil {
			t.Fatalf("RotateKey: %v", err)
		}
		_ = resp // response shape may vary
	})

	// Suspend
	t.Run("Suspend", func(t *testing.T) {
		if agentID == "" {
			t.Skip("no agent")
		}
		agent, err := client.Agents.Suspend(ctx, agentID)
		if err != nil {
			t.Fatalf("Suspend: %v", err)
		}
		if agent.Status != "SUSPENDED" {
			t.Fatalf("expected SUSPENDED, got %q", agent.Status)
		}
	})

	// Revoke
	t.Run("Revoke", func(t *testing.T) {
		if agentID == "" {
			t.Skip("no agent")
		}
		agent, err := client.Agents.Revoke(ctx, agentID)
		if err != nil {
			t.Fatalf("Revoke: %v", err)
		}
		if agent.Status != "REVOKED" {
			t.Fatalf("expected REVOKED, got %q", agent.Status)
		}
		agentID = ""
	})
}

// ---------------------------------------------------------------------------
// 2. RBAC Flow
// ---------------------------------------------------------------------------

func TestRBACFlow(t *testing.T) {
	client := newTestClient()
	ctx := context.Background()
	var (
		agentID string
		roleID  string
	)

	t.Cleanup(func() {
		if agentID != "" && roleID != "" {
			_ = client.Roles.UnassignFromAgent(ctx, agentID, roleID)
		}
		if roleID != "" {
			_ = client.Roles.Delete(ctx, roleID)
		}
		if agentID != "" {
			_, _ = client.Agents.Revoke(ctx, agentID)
		}
	})

	// Step 1: Create an agent for role assignment.
	t.Run("SetupAgent", func(t *testing.T) {
		agent, err := client.Agents.Create(ctx, &CreateAgentInput{
			WorkspaceID: testWorkspaceID,
			Name:        uniqueName("go-rbac-agent"),
			CreatedBy:   "integration-test",
		})
		if err != nil {
			t.Fatalf("Create agent: %v", err)
		}
		agentID = agent.ID
	})

	// Step 2: Create a role with permissions.
	t.Run("CreateRole", func(t *testing.T) {
		role, err := client.Roles.Create(ctx, &CreateRoleInput{
			WorkspaceID: testWorkspaceID,
			Name:        uniqueName("go-rbac-role"),
			Description: ptr("test role for RBAC flow"),
			Permissions: []string{"files:read", "files:write", "db:read"},
		})
		if err != nil {
			t.Fatalf("Create role: %v", err)
		}
		if role.ID == "" {
			t.Fatal("expected non-empty role ID")
		}
		roleID = role.ID
		t.Logf("created role: %s", roleID)
	})

	// Step 3: Assign the role to the agent.
	t.Run("AssignRole", func(t *testing.T) {
		if agentID == "" || roleID == "" {
			t.Skip("prerequisites missing")
		}
		ar, err := client.Roles.AssignToAgent(ctx, agentID, &AssignRoleInput{
			RoleID:    roleID,
			GrantedBy: ptr("integration-test"),
		})
		if err != nil {
			t.Fatalf("AssignToAgent: %v", err)
		}
		if ar.AgentID != agentID {
			t.Fatalf("expected agentId %s, got %s", agentID, ar.AgentID)
		}
	})

	// Step 4: List agent roles.
	t.Run("ListAgentRoles", func(t *testing.T) {
		if agentID == "" {
			t.Skip("no agent")
		}
		resp, err := client.Roles.ListAgentRoles(ctx, agentID)
		if err != nil {
			t.Fatalf("ListAgentRoles: %v", err)
		}
		if resp.AgentID != agentID {
			t.Fatalf("expected agentId %s, got %s", agentID, resp.AgentID)
		}
		if len(resp.Roles) < 1 {
			t.Fatal("expected at least 1 role assignment")
		}
	})

	// Step 5: Check single permission (should be allowed via role).
	// The permission check matches resource against the role's permission strings,
	// so resource must be "files:read" (the full permission string).
	t.Run("CheckPermission", func(t *testing.T) {
		if agentID == "" {
			t.Skip("no agent")
		}
		resp, err := client.Permissions.Check(ctx, &CheckPermissionInput{
			AgentID:  agentID,
			Resource: "files:read",
			Action:   "read",
		})
		if err != nil {
			t.Fatalf("Check permission: %v", err)
		}
		if !resp.Allowed {
			t.Fatal("expected permission to be allowed")
		}
	})

	// Step 6: Batch check.
	t.Run("BatchCheck", func(t *testing.T) {
		if agentID == "" {
			t.Skip("no agent")
		}
		resp, err := client.Permissions.CheckBatch(ctx, &BatchCheckInput{
			AgentID: agentID,
			Checks: []BatchCheckItem{
				{Resource: "files:read", Action: "read"},
				{Resource: "files:write", Action: "write"},
				{Resource: "files:delete", Action: "delete"},
			},
		})
		if err != nil {
			t.Fatalf("BatchCheck: %v", err)
		}
		if len(resp.Results) != 3 {
			t.Fatalf("expected 3 results, got %d", len(resp.Results))
		}
		// files:read and files:write should be allowed; files:delete should be denied.
		if !resp.Results[0].Allowed {
			t.Error("files:read should be allowed")
		}
		if !resp.Results[1].Allowed {
			t.Error("files:write should be allowed")
		}
		if resp.Results[2].Allowed {
			t.Error("files:delete should NOT be allowed")
		}
	})

	// Step 7: Unassign role.
	t.Run("UnassignRole", func(t *testing.T) {
		if agentID == "" || roleID == "" {
			t.Skip("prerequisites missing")
		}
		if err := client.Roles.UnassignFromAgent(ctx, agentID, roleID); err != nil {
			t.Fatalf("UnassignFromAgent: %v", err)
		}
	})

	// Step 8: Delete role.
	t.Run("DeleteRole", func(t *testing.T) {
		if roleID == "" {
			t.Skip("no role")
		}
		if err := client.Roles.Delete(ctx, roleID); err != nil {
			t.Fatalf("Delete role: %v", err)
		}
		roleID = "" // prevent double-delete in cleanup
	})
}

// ---------------------------------------------------------------------------
// 3. Policy Flow
// ---------------------------------------------------------------------------

func TestPolicyFlow(t *testing.T) {
	client := newTestClient()
	ctx := context.Background()
	var policyID string

	t.Cleanup(func() {
		if policyID != "" {
			_ = client.Policies.Delete(ctx, policyID)
		}
	})

	t.Run("Create", func(t *testing.T) {
		policy, err := client.Policies.Create(ctx, &CreatePolicyInput{
			WorkspaceID: testWorkspaceID,
			Name:        uniqueName("go-integ-policy"),
			Description: ptr("integration test policy"),
			Effect:      "ALLOW",
			Principals: map[string]interface{}{
				"roles": []string{"editor"},
			},
			Resources: []string{"files:*"},
			Actions:   []string{"read", "write"},
			Priority:  ptr(10),
			Enabled:   ptr(true),
		})
		if err != nil {
			t.Fatalf("Create policy: %v", err)
		}
		if policy.ID == "" {
			t.Fatal("expected non-empty policy ID")
		}
		if policy.Effect != "ALLOW" {
			t.Fatalf("expected effect ALLOW, got %q", policy.Effect)
		}
		policyID = policy.ID
		t.Logf("created policy: %s", policyID)
	})

	t.Run("List", func(t *testing.T) {
		policies, err := client.Policies.List(ctx, &ListPoliciesInput{
			WorkspaceID: testWorkspaceID,
		})
		if err != nil {
			t.Fatalf("List policies: %v", err)
		}
		if len(policies) < 1 {
			t.Fatal("expected at least 1 policy")
		}
	})

	t.Run("Update", func(t *testing.T) {
		if policyID == "" {
			t.Skip("no policy")
		}
		updated, err := client.Policies.Update(ctx, policyID, &UpdatePolicyInput{
			Description: ptr("updated description"),
			Priority:    ptr(20),
		})
		if err != nil {
			t.Fatalf("Update policy: %v", err)
		}
		if updated.Description == nil || *updated.Description != "updated description" {
			t.Fatal("description not updated")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		if policyID == "" {
			t.Skip("no policy")
		}
		if err := client.Policies.Delete(ctx, policyID); err != nil {
			t.Fatalf("Delete policy: %v", err)
		}
		policyID = "" // prevent double-delete
	})
}

// ---------------------------------------------------------------------------
// 4. Delegation Flow
// ---------------------------------------------------------------------------

func TestDelegationFlow(t *testing.T) {
	client := newTestClient()
	ctx := context.Background()
	var (
		issuerID     string
		delegateeID  string
		roleID       string
		delegationID string
	)

	t.Cleanup(func() {
		if delegationID != "" {
			_, _ = client.Delegations.Revoke(ctx, delegationID)
		}
		if issuerID != "" && roleID != "" {
			_ = client.Roles.UnassignFromAgent(ctx, issuerID, roleID)
		}
		if roleID != "" {
			_ = client.Roles.Delete(ctx, roleID)
		}
		if issuerID != "" {
			_, _ = client.Agents.Revoke(ctx, issuerID)
		}
		if delegateeID != "" {
			_, _ = client.Agents.Revoke(ctx, delegateeID)
		}
	})

	// Create two agents: issuer and delegatee.
	t.Run("SetupAgents", func(t *testing.T) {
		a1, err := client.Agents.Create(ctx, &CreateAgentInput{
			WorkspaceID: testWorkspaceID,
			Name:        uniqueName("go-deleg-issuer"),
			CreatedBy:   "integration-test",
		})
		if err != nil {
			t.Fatalf("Create issuer: %v", err)
		}
		issuerID = a1.ID

		a2, err := client.Agents.Create(ctx, &CreateAgentInput{
			WorkspaceID: testWorkspaceID,
			Name:        uniqueName("go-deleg-delegatee"),
			CreatedBy:   "integration-test",
		})
		if err != nil {
			t.Fatalf("Create delegatee: %v", err)
		}
		delegateeID = a2.ID
	})

	// Create a role and assign to issuer so they hold the permissions.
	t.Run("SetupIssuerRole", func(t *testing.T) {
		if issuerID == "" {
			t.Skip("no issuer")
		}
		role, err := client.Roles.Create(ctx, &CreateRoleInput{
			WorkspaceID: testWorkspaceID,
			Name:        uniqueName("go-deleg-role"),
			Permissions: []string{"files:read", "files:write"},
		})
		if err != nil {
			t.Fatalf("Create role: %v", err)
		}
		roleID = role.ID

		_, err = client.Roles.AssignToAgent(ctx, issuerID, &AssignRoleInput{
			RoleID:    roleID,
			GrantedBy: ptr("integration-test"),
		})
		if err != nil {
			t.Fatalf("Assign role to issuer: %v", err)
		}
	})

	// Create the delegation from issuer to delegatee.
	t.Run("CreateDelegation", func(t *testing.T) {
		if issuerID == "" || delegateeID == "" {
			t.Skip("prerequisites missing")
		}
		d, err := client.Delegations.Create(ctx, &CreateDelegationInput{
			IssuerAgentID: issuerID,
			TargetAgentID: delegateeID,
			Permissions:   []string{"files:read"},
			Constraints: map[string]interface{}{
				"maxDepth": 1,
			},
		})
		if err != nil {
			t.Fatalf("Create delegation: %v", err)
		}
		if d.ID == "" {
			t.Fatal("expected non-empty delegation ID")
		}
		delegationID = d.ID
		t.Logf("created delegation: %s", delegationID)
	})

	t.Run("Get", func(t *testing.T) {
		if delegationID == "" {
			t.Skip("no delegation")
		}
		d, err := client.Delegations.Get(ctx, delegationID)
		if err != nil {
			t.Fatalf("Get delegation: %v", err)
		}
		if d.ID != delegationID {
			t.Fatalf("expected delegation ID %s, got %s", delegationID, d.ID)
		}
		if d.IssuerAgentID != issuerID {
			t.Fatalf("expected issuerAgentId %s, got %s", issuerID, d.IssuerAgentID)
		}
	})

	t.Run("List", func(t *testing.T) {
		delegations, err := client.Delegations.List(ctx, &ListDelegationsInput{
			Status: ptr("active"),
		})
		if err != nil {
			t.Fatalf("List delegations: %v", err)
		}
		if len(delegations) < 1 {
			t.Fatal("expected at least 1 delegation")
		}
	})

	t.Run("Revoke", func(t *testing.T) {
		if delegationID == "" {
			t.Skip("no delegation")
		}
		d, err := client.Delegations.Revoke(ctx, delegationID)
		if err != nil {
			t.Fatalf("Revoke delegation: %v", err)
		}
		if d.Status != "REVOKED" {
			t.Fatalf("expected status 'REVOKED', got %q", d.Status)
		}
		delegationID = "" // prevent double-revoke in cleanup
	})
}

// ---------------------------------------------------------------------------
// 5. Audit Flow
// ---------------------------------------------------------------------------

func TestAuditFlow(t *testing.T) {
	client := newTestClient()
	ctx := context.Background()

	t.Run("ListEvents", func(t *testing.T) {
		limit := 5
		events, err := client.Audit.ListEvents(ctx, &ListAuditEventsInput{
			Limit: &limit,
		})
		if err != nil {
			t.Fatalf("ListEvents: %v", err)
		}
		// The test workspace may or may not have events, so just check no error.
		t.Logf("audit events total: %d, returned: %d", events.Total, len(events.Items))
	})

	t.Run("GetMetrics", func(t *testing.T) {
		rows, err := client.Audit.GetMetrics(ctx, &AuditMetricsInput{
			OrgID: testOrgID,
		})
		if err != nil {
			t.Fatalf("GetMetrics: %v", err)
		}
		t.Logf("audit metric rows: %d", len(rows))
	})
}

// ---------------------------------------------------------------------------
// 6. Webhook Flow
// ---------------------------------------------------------------------------

func TestWebhookFlow(t *testing.T) {
	client := newTestClient()
	ctx := context.Background()
	var webhookID string

	t.Cleanup(func() {
		if webhookID != "" {
			_ = client.Webhooks.Delete(ctx, webhookID)
		}
	})

	t.Run("Create", func(t *testing.T) {
		wh, err := client.Webhooks.Create(ctx, &CreateWebhookInput{
			OrganizationID: testOrgID,
			URL:            "https://example.com/webhook-go-test",
			EventTypes:     []string{"agent.created", "agent.revoked"},
			Secret:         "test-secret-go-sdk",
		})
		if err != nil {
			t.Fatalf("Create webhook: %v", err)
		}
		if wh.ID == "" {
			t.Fatal("expected non-empty webhook ID")
		}
		webhookID = wh.ID
		t.Logf("created webhook: %s", webhookID)
	})

	t.Run("List", func(t *testing.T) {
		webhooks, err := client.Webhooks.List(ctx, &ListWebhooksInput{
			OrganizationID: testOrgID,
		})
		if err != nil {
			t.Fatalf("List webhooks: %v", err)
		}
		if len(webhooks) < 1 {
			t.Fatal("expected at least 1 webhook")
		}
	})

	t.Run("Update", func(t *testing.T) {
		if webhookID == "" {
			t.Skip("no webhook")
		}
		updated, err := client.Webhooks.Update(ctx, webhookID, &UpdateWebhookInput{
			URL:     ptr("https://example.com/webhook-go-updated"),
			Enabled: ptr(false),
		})
		if err != nil {
			t.Fatalf("Update webhook: %v", err)
		}
		if updated.URL != "https://example.com/webhook-go-updated" {
			t.Fatalf("expected updated URL, got %q", updated.URL)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		if webhookID == "" {
			t.Skip("no webhook")
		}
		if err := client.Webhooks.Delete(ctx, webhookID); err != nil {
			t.Fatalf("Delete webhook: %v", err)
		}
		webhookID = "" // prevent double-delete
	})
}

// ---------------------------------------------------------------------------
// 7. Alert Flow
// ---------------------------------------------------------------------------

func TestAlertFlow(t *testing.T) {
	client := newTestClient()
	ctx := context.Background()
	var alertID string

	t.Cleanup(func() {
		if alertID != "" {
			_ = client.Alerts.Delete(ctx, alertID)
		}
	})

	t.Run("Create", func(t *testing.T) {
		alert, err := client.Alerts.Create(ctx, &CreateAlertInput{
			OrganizationID: testOrgID,
			Name:           uniqueName("go-integ-alert"),
			EventTypes:     []string{"agent.revoked", "permission.denied"},
			Conditions: map[string]interface{}{
				"severity": "high",
			},
			Channels: []string{"email"},
		})
		if err != nil {
			t.Fatalf("Create alert: %v", err)
		}
		if alert.ID == "" {
			t.Fatal("expected non-empty alert ID")
		}
		alertID = alert.ID
		t.Logf("created alert: %s", alertID)
	})

	t.Run("List", func(t *testing.T) {
		alerts, err := client.Alerts.List(ctx, &ListAlertsInput{
			OrganizationID: testOrgID,
		})
		if err != nil {
			t.Fatalf("List alerts: %v", err)
		}
		if len(alerts) < 1 {
			t.Fatal("expected at least 1 alert")
		}
	})

	t.Run("Update", func(t *testing.T) {
		if alertID == "" {
			t.Skip("no alert")
		}
		updated, err := client.Alerts.Update(ctx, alertID, &UpdateAlertInput{
			Name:    ptr("go-integ-alert-updated"),
			Enabled: ptr(false),
		})
		if err != nil {
			t.Fatalf("Update alert: %v", err)
		}
		if updated.Name != "go-integ-alert-updated" {
			t.Fatalf("expected updated name, got %q", updated.Name)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		if alertID == "" {
			t.Skip("no alert")
		}
		if err := client.Alerts.Delete(ctx, alertID); err != nil {
			t.Fatalf("Delete alert: %v", err)
		}
		alertID = "" // prevent double-delete
	})
}

// ---------------------------------------------------------------------------
// 8. API Key Flow
// ---------------------------------------------------------------------------

func TestAPIKeyFlow(t *testing.T) {
	client := newTestClient()
	ctx := context.Background()
	var createdKeyID string

	t.Cleanup(func() {
		if createdKeyID != "" {
			_ = client.APIKeys.Revoke(ctx, createdKeyID)
		}
	})

	t.Run("List", func(t *testing.T) {
		keys, err := client.APIKeys.List(ctx, &ListAPIKeysInput{
			OrganizationID: testOrgID,
		})
		if err != nil {
			t.Fatalf("List API keys: %v", err)
		}
		// At minimum the key we are using exists.
		t.Logf("found %d API keys", len(keys))
	})

	t.Run("Create", func(t *testing.T) {
		resp, err := client.APIKeys.Create(ctx, &CreateAPIKeyInput{
			OrganizationID: testOrgID,
			Name:           uniqueName("go-integ-apikey"),
			CreatedBy:      "integration-test",
			ExpiresInDays:  ptr(1),
		})
		if err != nil {
			t.Fatalf("Create API key: %v", err)
		}
		if resp.RawKey == "" {
			t.Fatal("expected non-empty rawKey value")
		}
		if resp.ID == "" {
			t.Fatal("expected non-empty API key ID")
		}
		createdKeyID = resp.ID
		t.Logf("created API key: %s", createdKeyID)
	})

	t.Run("Delete", func(t *testing.T) {
		if createdKeyID == "" {
			t.Skip("no key created")
		}
		if err := client.APIKeys.Revoke(ctx, createdKeyID); err != nil {
			t.Fatalf("Revoke API key: %v", err)
		}
		createdKeyID = "" // prevent double-delete
	})
}

// ---------------------------------------------------------------------------
// 9. Notification Flow
// ---------------------------------------------------------------------------

func TestNotificationFlow(t *testing.T) {
	client := newTestClient()
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		resp, err := client.Notifications.List(ctx, &ListNotificationsInput{
			OrganizationID: testOrgID,
		})
		if err != nil {
			t.Fatalf("List notifications: %v", err)
		}
		t.Logf("found %d notifications (total: %d)", len(resp.Items), resp.Total)
	})

	t.Run("UnreadCount", func(t *testing.T) {
		resp, err := client.Notifications.UnreadCount(ctx, &UnreadCountInput{
			OrganizationID: testOrgID,
		})
		if err != nil {
			t.Fatalf("UnreadCount: %v", err)
		}
		if resp.Count < 0 {
			t.Fatal("count should be >= 0")
		}
		t.Logf("unread count: %d", resp.Count)
	})

	t.Run("MarkAllRead", func(t *testing.T) {
		err := client.Notifications.MarkAllRead(ctx, &MarkAllReadInput{
			OrganizationID: testOrgID,
		})
		if err != nil {
			t.Fatalf("MarkAllRead: %v", err)
		}
	})
}

// ---------------------------------------------------------------------------
// 10. Organization & Workspace
// ---------------------------------------------------------------------------

func TestOrgAndWorkspace(t *testing.T) {
	client := newTestClient()
	ctx := context.Background()

	t.Run("GetOrganization", func(t *testing.T) {
		org, err := client.Organizations.Get(ctx, testOrgID)
		if err != nil {
			t.Fatalf("Get org: %v", err)
		}
		if org.ID != testOrgID {
			t.Fatalf("expected org ID %s, got %s", testOrgID, org.ID)
		}
		t.Logf("org: %s (%s)", org.Name, org.ID)
	})

	t.Run("ListWorkspaces", func(t *testing.T) {
		list, err := client.Workspaces.List(ctx, &ListWorkspacesInput{
			OrganizationID: testOrgID,
		})
		if err != nil {
			t.Fatalf("List workspaces: %v", err)
		}
		if list.Total < 1 {
			t.Fatal("expected at least 1 workspace")
		}
		if len(list.Items) < 1 {
			t.Fatal("expected at least 1 workspace in Items")
		}
		t.Logf("workspaces: total=%d, returned=%d", list.Total, len(list.Items))
	})

	t.Run("GetWorkspace", func(t *testing.T) {
		ws, err := client.Workspaces.Get(ctx, testWorkspaceID)
		if err != nil {
			t.Fatalf("Get workspace: %v", err)
		}
		if ws.ID != testWorkspaceID {
			t.Fatalf("expected workspace ID %s, got %s", testWorkspaceID, ws.ID)
		}
		t.Logf("workspace: %s (%s)", ws.Name, ws.ID)
	})
}

// ---------------------------------------------------------------------------
// 11. MCP Server & Tool Registration + Proxy
// ---------------------------------------------------------------------------

func TestMcpFlow(t *testing.T) {
	client := newTestClient()
	ctx := context.Background()
	var serverID string

	t.Cleanup(func() {
		// No delete endpoint for MCP servers; cleanup is a no-op.
	})

	// Register MCP server
	t.Run("RegisterServer", func(t *testing.T) {
		server, err := client.Mcp.RegisterServer(ctx, &RegisterMcpServerInput{
			WorkspaceID: testWorkspaceID,
			Name:        uniqueName("go-mcp-server"),
			URL:         "http://127.0.0.1:9100",
			Description: ptr("Go SDK integration test MCP server"),
		})
		if err != nil {
			t.Fatalf("RegisterServer: %v", err)
		}
		if server.ID == "" {
			t.Fatal("expected non-empty server ID")
		}
		serverID = server.ID
		t.Logf("registered MCP server: %s", serverID)
	})

	// List servers
	t.Run("ListServers", func(t *testing.T) {
		servers, err := client.Mcp.ListServers(ctx, &ListMcpServersInput{
			WorkspaceID: testWorkspaceID,
		})
		if err != nil {
			t.Fatalf("ListServers: %v", err)
		}
		if len(servers.Items) < 1 {
			t.Fatal("expected at least 1 MCP server")
		}
		found := false
		for _, s := range servers.Items {
			if s.ID == serverID {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("registered server not found in list")
		}
	})

	// Get server
	t.Run("GetServer", func(t *testing.T) {
		if serverID == "" {
			t.Skip("no server")
		}
		server, err := client.Mcp.GetServer(ctx, serverID)
		if err != nil {
			t.Fatalf("GetServer: %v", err)
		}
		if server.ID != serverID {
			t.Fatalf("expected server ID %s, got %s", serverID, server.ID)
		}
	})

	// Update server
	t.Run("UpdateServer", func(t *testing.T) {
		if serverID == "" {
			t.Skip("no server")
		}
		updated, err := client.Mcp.UpdateServer(ctx, serverID, &UpdateMcpServerInput{
			Description: ptr("Updated by Go SDK integration test"),
		})
		if err != nil {
			t.Fatalf("UpdateServer: %v", err)
		}
		if updated.ID != serverID {
			t.Fatalf("expected server ID %s, got %s", serverID, updated.ID)
		}
	})

	// Register tool
	t.Run("RegisterTool", func(t *testing.T) {
		if serverID == "" {
			t.Skip("no server")
		}
		tool, err := client.Mcp.RegisterTool(ctx, serverID, &RegisterMcpToolInput{
			Name:        "echo",
			Description: ptr("Echo tool for testing"),
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"message": map[string]interface{}{"type": "string"},
				},
			},
		})
		if err != nil {
			t.Fatalf("RegisterTool: %v", err)
		}
		if tool.ID == "" {
			t.Fatal("expected non-empty tool ID")
		}
		if tool.Name != "echo" {
			t.Fatalf("expected tool name 'echo', got %q", tool.Name)
		}
	})

	// List tools
	t.Run("ListTools", func(t *testing.T) {
		if serverID == "" {
			t.Skip("no server")
		}
		tools, err := client.Mcp.ListTools(ctx, serverID)
		if err != nil {
			t.Fatalf("ListTools: %v", err)
		}
		if len(tools) < 1 {
			t.Fatal("expected at least 1 tool")
		}
		found := false
		for _, tool := range tools {
			if tool.Name == "echo" {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("echo tool not found in list")
		}
	})

	// Proxy: call echo tool through authorization pipeline
	// Need an agent with MCP permissions
	t.Run("Proxy", func(t *testing.T) {
		if serverID == "" {
			t.Skip("no server")
		}

		// Create agent + role for proxy auth
		proxyAgent, err := client.Agents.Create(ctx, &CreateAgentInput{
			WorkspaceID: testWorkspaceID,
			Name:        uniqueName("go-mcp-proxy-agent"),
			CreatedBy:   "integration-test",
		})
		if err != nil {
			t.Fatalf("Create proxy agent: %v", err)
		}
		proxyRole, err := client.Roles.Create(ctx, &CreateRoleInput{
			WorkspaceID: testWorkspaceID,
			Name:        uniqueName("go-mcp-proxy-role"),
			Permissions: []string{fmt.Sprintf("mcp:%s:tool.*", serverID)},
		})
		if err != nil {
			t.Fatalf("Create proxy role: %v", err)
		}
		_, err = client.Roles.AssignToAgent(ctx, proxyAgent.ID, &AssignRoleInput{
			RoleID:    proxyRole.ID,
			GrantedBy: ptr("integration-test"),
		})
		if err != nil {
			t.Fatalf("Assign proxy role: %v", err)
		}

		resp, err := client.Mcp.Proxy(ctx, &McpProxyInput{
			ServerID: serverID,
			Method:   "tools/call",
			Params: map[string]interface{}{
				"name":      "echo",
				"arguments": map[string]interface{}{"message": "hello-from-go-sdk"},
				"_authora": map[string]interface{}{
					"mcpServerId": serverID,
					"agentId":     proxyAgent.ID,
					"timestamp":   time.Now().UTC().Format(time.RFC3339),
				},
			},
		})
		if err != nil {
			t.Fatalf("Proxy: %v", err)
		}
		resultStr := fmt.Sprintf("%v", resp.Result)
		if resultStr == "" {
			t.Fatal("proxy returned empty result")
		}
		t.Logf("proxy result: %s", resultStr)

		// Cleanup
		_ = client.Roles.UnassignFromAgent(ctx, proxyAgent.ID, proxyRole.ID)
		_ = client.Roles.Delete(ctx, proxyRole.ID)
		_, _ = client.Agents.Revoke(ctx, proxyAgent.ID)
	})
}

// ---------------------------------------------------------------------------
// 12. Policy Simulate & Evaluate
// ---------------------------------------------------------------------------

func TestPolicySimulateEvaluate(t *testing.T) {
	client := newTestClient()
	ctx := context.Background()
	var (
		agentID  string
		roleID   string
		policyID string
	)

	t.Cleanup(func() {
		if policyID != "" {
			_ = client.Policies.Delete(ctx, policyID)
		}
		if agentID != "" && roleID != "" {
			_ = client.Roles.UnassignFromAgent(ctx, agentID, roleID)
		}
		if roleID != "" {
			_ = client.Roles.Delete(ctx, roleID)
		}
		if agentID != "" {
			_, _ = client.Agents.Revoke(ctx, agentID)
		}
	})

	// Setup: agent + role + DENY policy
	t.Run("Setup", func(t *testing.T) {
		agent, err := client.Agents.Create(ctx, &CreateAgentInput{
			WorkspaceID: testWorkspaceID,
			Name:        uniqueName("go-poleval-agent"),
			CreatedBy:   "integration-test",
		})
		if err != nil {
			t.Fatalf("Create agent: %v", err)
		}
		agentID = agent.ID

		role, err := client.Roles.Create(ctx, &CreateRoleInput{
			WorkspaceID: testWorkspaceID,
			Name:        uniqueName("go-poleval-role"),
			Permissions: []string{"docs:*:read"},
		})
		if err != nil {
			t.Fatalf("Create role: %v", err)
		}
		roleID = role.ID

		_, err = client.Roles.AssignToAgent(ctx, agentID, &AssignRoleInput{
			RoleID:    roleID,
			GrantedBy: ptr("integration-test"),
		})
		if err != nil {
			t.Fatalf("Assign role: %v", err)
		}

		policy, err := client.Policies.Create(ctx, &CreatePolicyInput{
			WorkspaceID: testWorkspaceID,
			Name:        uniqueName("go-deny-policy"),
			Effect:      "DENY",
			Principals: map[string]interface{}{
				"roles": []string{role.Name},
			},
			Resources: []string{"docs:secret"},
			Actions:   []string{"read"},
			Priority:  ptr(100),
			Enabled:   ptr(true),
		})
		if err != nil {
			t.Fatalf("Create DENY policy: %v", err)
		}
		policyID = policy.ID
	})

	// Simulate
	t.Run("Simulate", func(t *testing.T) {
		if agentID == "" {
			t.Skip("no agent")
		}
		resp, err := client.Policies.Simulate(ctx, &SimulatePolicyInput{
			WorkspaceID: testWorkspaceID,
			AgentID:     agentID,
			Resource:    "docs:secret",
			Action:      "read",
		})
		if err != nil {
			t.Fatalf("Simulate: %v", err)
		}
		t.Logf("simulate decision: %s, reason: %v", resp.Decision, resp.Reason)
	})

	// Evaluate
	t.Run("Evaluate", func(t *testing.T) {
		if agentID == "" {
			t.Skip("no agent")
		}
		resp, err := client.Policies.Evaluate(ctx, &EvaluatePolicyInput{
			WorkspaceID: testWorkspaceID,
			AgentID:     agentID,
			Resource:    "docs:secret",
			Action:      "read",
		})
		if err != nil {
			t.Fatalf("Evaluate: %v", err)
		}
		t.Logf("evaluate allowed: %v, reason: %v", resp.Allowed, resp.Reason)
	})
}
