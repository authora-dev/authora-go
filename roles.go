package authora

import (
	"context"
	"fmt"
	"net/http"
)

// RolesService handles role-related API endpoints.
type RolesService struct {
	client *httpClient
}

// Create creates a new role. POST /roles
func (s *RolesService) Create(ctx context.Context, input *CreateRoleInput) (*Role, error) {
	var resp Role
	if err := s.client.request(ctx, http.MethodPost, "/roles", input, &resp); err != nil {
		return nil, fmt.Errorf("roles.Create: %w", err)
	}
	return &resp, nil
}

// List returns a paginated list of roles. GET /roles
func (s *RolesService) List(ctx context.Context, input *ListRolesInput) (*PaginatedResponse[Role], error) {
	q := queryString(map[string]interface{}{
		"workspaceId": input.WorkspaceID,
		"page":        input.Page,
		"limit":       input.Limit,
	})
	var resp PaginatedResponse[Role]
	if err := s.client.request(ctx, http.MethodGet, "/roles"+q, nil, &resp); err != nil {
		return nil, fmt.Errorf("roles.List: %w", err)
	}
	return &resp, nil
}

// Get retrieves a single role by ID. GET /roles/:roleId
func (s *RolesService) Get(ctx context.Context, roleID string) (*Role, error) {
	var resp Role
	if err := s.client.request(ctx, http.MethodGet, "/roles/"+roleID, nil, &resp); err != nil {
		return nil, fmt.Errorf("roles.Get: %w", err)
	}
	return &resp, nil
}

// Update modifies an existing role. PATCH /roles/:roleId
func (s *RolesService) Update(ctx context.Context, roleID string, input *UpdateRoleInput) (*Role, error) {
	var resp Role
	if err := s.client.request(ctx, http.MethodPatch, "/roles/"+roleID, input, &resp); err != nil {
		return nil, fmt.Errorf("roles.Update: %w", err)
	}
	return &resp, nil
}

// Delete removes a role. DELETE /roles/:roleId
func (s *RolesService) Delete(ctx context.Context, roleID string) error {
	if err := s.client.request(ctx, http.MethodDelete, "/roles/"+roleID, nil, nil); err != nil {
		return fmt.Errorf("roles.Delete: %w", err)
	}
	return nil
}

// AssignToAgent assigns a role to an agent. POST /agents/:agentId/roles
func (s *RolesService) AssignToAgent(ctx context.Context, agentID string, input *AssignRoleInput) (*AgentRole, error) {
	var resp AgentRole
	if err := s.client.request(ctx, http.MethodPost, "/agents/"+agentID+"/roles", input, &resp); err != nil {
		return nil, fmt.Errorf("roles.AssignToAgent: %w", err)
	}
	return &resp, nil
}

// UnassignFromAgent removes a role from an agent.
// DELETE /agents/:agentId/roles/:roleId
func (s *RolesService) UnassignFromAgent(ctx context.Context, agentID, roleID string) error {
	if err := s.client.request(ctx, http.MethodDelete, "/agents/"+agentID+"/roles/"+roleID, nil, nil); err != nil {
		return fmt.Errorf("roles.UnassignFromAgent: %w", err)
	}
	return nil
}

// AgentRolesResponse is the response from GET /agents/:agentId/roles.
type AgentRolesResponse struct {
	AgentID string      `json:"agentId"`
	Roles   []AgentRole `json:"roles"`
}

// ListAgentRoles returns all roles assigned to an agent.
// GET /agents/:agentId/roles
func (s *RolesService) ListAgentRoles(ctx context.Context, agentID string) (*AgentRolesResponse, error) {
	var resp AgentRolesResponse
	if err := s.client.request(ctx, http.MethodGet, "/agents/"+agentID+"/roles", nil, &resp); err != nil {
		return nil, fmt.Errorf("roles.ListAgentRoles: %w", err)
	}
	return &resp, nil
}
