package authora

import (
	"context"
	"fmt"
	"net/http"
)

type RolesService struct {
	client *httpClient
}

func (s *RolesService) Create(ctx context.Context, input *CreateRoleInput) (*Role, error) {
	var resp Role
	if err := s.client.request(ctx, http.MethodPost, "/roles", input, &resp); err != nil {
		return nil, fmt.Errorf("roles.Create: %w", err)
	}
	return &resp, nil
}

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

func (s *RolesService) Get(ctx context.Context, roleID string) (*Role, error) {
	var resp Role
	if err := s.client.request(ctx, http.MethodGet, "/roles/"+roleID, nil, &resp); err != nil {
		return nil, fmt.Errorf("roles.Get: %w", err)
	}
	return &resp, nil
}

func (s *RolesService) Update(ctx context.Context, roleID string, input *UpdateRoleInput) (*Role, error) {
	var resp Role
	if err := s.client.request(ctx, http.MethodPatch, "/roles/"+roleID, input, &resp); err != nil {
		return nil, fmt.Errorf("roles.Update: %w", err)
	}
	return &resp, nil
}

func (s *RolesService) Delete(ctx context.Context, roleID string) error {
	if err := s.client.request(ctx, http.MethodDelete, "/roles/"+roleID, nil, nil); err != nil {
		return fmt.Errorf("roles.Delete: %w", err)
	}
	return nil
}

func (s *RolesService) AssignToAgent(ctx context.Context, agentID string, input *AssignRoleInput) (*AgentRole, error) {
	var resp AgentRole
	if err := s.client.request(ctx, http.MethodPost, "/agents/"+agentID+"/roles", input, &resp); err != nil {
		return nil, fmt.Errorf("roles.AssignToAgent: %w", err)
	}
	return &resp, nil
}

func (s *RolesService) UnassignFromAgent(ctx context.Context, agentID, roleID string) error {
	if err := s.client.request(ctx, http.MethodDelete, "/agents/"+agentID+"/roles/"+roleID, nil, nil); err != nil {
		return fmt.Errorf("roles.UnassignFromAgent: %w", err)
	}
	return nil
}

type AgentRolesResponse struct {
	AgentID string      `json:"agentId"`
	Roles   []AgentRole `json:"roles"`
}

func (s *RolesService) ListAgentRoles(ctx context.Context, agentID string) (*AgentRolesResponse, error) {
	var resp AgentRolesResponse
	if err := s.client.request(ctx, http.MethodGet, "/agents/"+agentID+"/roles", nil, &resp); err != nil {
		return nil, fmt.Errorf("roles.ListAgentRoles: %w", err)
	}
	return &resp, nil
}
