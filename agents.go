package authora

import (
	"context"
	"fmt"
	"net/http"
)

type AgentsService struct {
	client *httpClient
}

func (s *AgentsService) Create(ctx context.Context, input *CreateAgentInput) (*CreateAgentResponse, error) {
	var resp CreateAgentResponse
	if err := s.client.request(ctx, http.MethodPost, "/agents", input, &resp); err != nil {
		return nil, fmt.Errorf("agents.Create: %w", err)
	}
	return &resp, nil
}

func (s *AgentsService) List(ctx context.Context, input *ListAgentsInput) (*PaginatedResponse[Agent], error) {
	q := queryString(map[string]interface{}{
		"workspaceId": input.WorkspaceID,
		"status":      input.Status,
		"page":        input.Page,
		"limit":       input.Limit,
	})
	var resp PaginatedResponse[Agent]
	if err := s.client.request(ctx, http.MethodGet, "/agents"+q, nil, &resp); err != nil {
		return nil, fmt.Errorf("agents.List: %w", err)
	}
	return &resp, nil
}

func (s *AgentsService) Get(ctx context.Context, agentID string) (*Agent, error) {
	var resp Agent
	if err := s.client.request(ctx, http.MethodGet, "/agents/"+agentID, nil, &resp); err != nil {
		return nil, fmt.Errorf("agents.Get: %w", err)
	}
	return &resp, nil
}

func (s *AgentsService) Verify(ctx context.Context, agentID string) (*VerifyAgentResponse, error) {
	var resp VerifyAgentResponse
	if err := s.client.requestNoAuth(ctx, http.MethodGet, "/agents/"+agentID+"/verify", nil, &resp); err != nil {
		return nil, fmt.Errorf("agents.Verify: %w", err)
	}
	return &resp, nil
}

func (s *AgentsService) Activate(ctx context.Context, agentID string, input *ActivateAgentInput) (*Agent, error) {
	var resp Agent
	if err := s.client.request(ctx, http.MethodPost, "/agents/"+agentID+"/activate", input, &resp); err != nil {
		return nil, fmt.Errorf("agents.Activate: %w", err)
	}
	return &resp, nil
}

func (s *AgentsService) Suspend(ctx context.Context, agentID string) (*Agent, error) {
	var resp Agent
	if err := s.client.request(ctx, http.MethodPost, "/agents/"+agentID+"/suspend", nil, &resp); err != nil {
		return nil, fmt.Errorf("agents.Suspend: %w", err)
	}
	return &resp, nil
}

func (s *AgentsService) Revoke(ctx context.Context, agentID string) (*Agent, error) {
	var resp Agent
	if err := s.client.request(ctx, http.MethodPost, "/agents/"+agentID+"/revoke", nil, &resp); err != nil {
		return nil, fmt.Errorf("agents.Revoke: %w", err)
	}
	return &resp, nil
}

func (s *AgentsService) RotateKey(ctx context.Context, agentID string, input *RotateKeyInput) (*RotateKeyResponse, error) {
	var resp RotateKeyResponse
	if err := s.client.request(ctx, http.MethodPost, "/agents/"+agentID+"/rotate-key", input, &resp); err != nil {
		return nil, fmt.Errorf("agents.RotateKey: %w", err)
	}
	return &resp, nil
}
