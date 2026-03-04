package authora

import (
	"context"
	"fmt"
	"net/http"
)

// AgentsService handles agent-related API endpoints.
type AgentsService struct {
	client *httpClient
}

// Create registers a new agent. POST /agents
func (s *AgentsService) Create(ctx context.Context, input *CreateAgentInput) (*CreateAgentResponse, error) {
	var resp CreateAgentResponse
	if err := s.client.request(ctx, http.MethodPost, "/agents", input, &resp); err != nil {
		return nil, fmt.Errorf("agents.Create: %w", err)
	}
	return &resp, nil
}

// List returns a paginated list of agents. GET /agents
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

// Get retrieves a single agent by ID. GET /agents/:agentId
func (s *AgentsService) Get(ctx context.Context, agentID string) (*Agent, error) {
	var resp Agent
	if err := s.client.request(ctx, http.MethodGet, "/agents/"+agentID, nil, &resp); err != nil {
		return nil, fmt.Errorf("agents.Get: %w", err)
	}
	return &resp, nil
}

// Verify checks whether an agent is valid. This is a public endpoint
// that does not require authentication. GET /agents/:agentId/verify
func (s *AgentsService) Verify(ctx context.Context, agentID string) (*VerifyAgentResponse, error) {
	var resp VerifyAgentResponse
	if err := s.client.requestNoAuth(ctx, http.MethodGet, "/agents/"+agentID+"/verify", nil, &resp); err != nil {
		return nil, fmt.Errorf("agents.Verify: %w", err)
	}
	return &resp, nil
}

// Activate activates a pending agent with a public key.
// POST /agents/:agentId/activate
func (s *AgentsService) Activate(ctx context.Context, agentID string, input *ActivateAgentInput) (*Agent, error) {
	var resp Agent
	if err := s.client.request(ctx, http.MethodPost, "/agents/"+agentID+"/activate", input, &resp); err != nil {
		return nil, fmt.Errorf("agents.Activate: %w", err)
	}
	return &resp, nil
}

// Suspend suspends an active agent. POST /agents/:agentId/suspend
func (s *AgentsService) Suspend(ctx context.Context, agentID string) (*Agent, error) {
	var resp Agent
	if err := s.client.request(ctx, http.MethodPost, "/agents/"+agentID+"/suspend", nil, &resp); err != nil {
		return nil, fmt.Errorf("agents.Suspend: %w", err)
	}
	return &resp, nil
}

// Revoke permanently revokes an agent. POST /agents/:agentId/revoke
func (s *AgentsService) Revoke(ctx context.Context, agentID string) (*Agent, error) {
	var resp Agent
	if err := s.client.request(ctx, http.MethodPost, "/agents/"+agentID+"/revoke", nil, &resp); err != nil {
		return nil, fmt.Errorf("agents.Revoke: %w", err)
	}
	return &resp, nil
}

// RotateKey rotates an agent's API key. POST /agents/:agentId/rotate-key
func (s *AgentsService) RotateKey(ctx context.Context, agentID string, input *RotateKeyInput) (*RotateKeyResponse, error) {
	var resp RotateKeyResponse
	if err := s.client.request(ctx, http.MethodPost, "/agents/"+agentID+"/rotate-key", input, &resp); err != nil {
		return nil, fmt.Errorf("agents.RotateKey: %w", err)
	}
	return &resp, nil
}
