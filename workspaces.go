package authora

import (
	"context"
	"fmt"
	"net/http"
)

type WorkspacesService struct {
	client *httpClient
}

func (s *WorkspacesService) Create(ctx context.Context, input *CreateWorkspaceInput) (*Workspace, error) {
	var resp Workspace
	if err := s.client.request(ctx, http.MethodPost, "/workspaces", input, &resp); err != nil {
		return nil, fmt.Errorf("workspaces.Create: %w", err)
	}
	return &resp, nil
}

func (s *WorkspacesService) Get(ctx context.Context, workspaceID string) (*Workspace, error) {
	var resp Workspace
	if err := s.client.request(ctx, http.MethodGet, "/workspaces/"+workspaceID, nil, &resp); err != nil {
		return nil, fmt.Errorf("workspaces.Get: %w", err)
	}
	return &resp, nil
}

func (s *WorkspacesService) List(ctx context.Context, input *ListWorkspacesInput) (*PaginatedResponse[Workspace], error) {
	q := queryString(map[string]interface{}{
		"organizationId": input.OrganizationID,
		"page":           input.Page,
		"limit":          input.Limit,
	})
	var resp PaginatedResponse[Workspace]
	if err := s.client.request(ctx, http.MethodGet, "/workspaces"+q, nil, &resp); err != nil {
		return nil, fmt.Errorf("workspaces.List: %w", err)
	}
	return &resp, nil
}

func (s *WorkspacesService) Update(ctx context.Context, workspaceID string, input *UpdateWorkspaceInput) (*Workspace, error) {
	var resp Workspace
	if err := s.client.request(ctx, http.MethodPatch, "/workspaces/"+workspaceID, input, &resp); err != nil {
		return nil, fmt.Errorf("workspaces.Update: %w", err)
	}
	return &resp, nil
}

func (s *WorkspacesService) Delete(ctx context.Context, workspaceID string) (*Workspace, error) {
	var resp Workspace
	if err := s.client.request(ctx, http.MethodDelete, "/workspaces/"+workspaceID, nil, &resp); err != nil {
		return nil, fmt.Errorf("workspaces.Delete: %w", err)
	}
	return &resp, nil
}

func (s *WorkspacesService) Restore(ctx context.Context, workspaceID string) (*Workspace, error) {
	var resp Workspace
	if err := s.client.request(ctx, http.MethodPost, "/workspaces/"+workspaceID+"/restore", nil, &resp); err != nil {
		return nil, fmt.Errorf("workspaces.Restore: %w", err)
	}
	return &resp, nil
}
