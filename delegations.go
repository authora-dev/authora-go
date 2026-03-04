package authora

import (
	"context"
	"fmt"
	"net/http"
)

type DelegationsService struct {
	client *httpClient
}

func (s *DelegationsService) Create(ctx context.Context, input *CreateDelegationInput) (*Delegation, error) {
	var resp Delegation
	if err := s.client.request(ctx, http.MethodPost, "/delegations", input, &resp); err != nil {
		return nil, fmt.Errorf("delegations.Create: %w", err)
	}
	return &resp, nil
}

func (s *DelegationsService) Get(ctx context.Context, delegationID string) (*Delegation, error) {
	var resp Delegation
	if err := s.client.request(ctx, http.MethodGet, "/delegations/"+delegationID, nil, &resp); err != nil {
		return nil, fmt.Errorf("delegations.Get: %w", err)
	}
	return &resp, nil
}

func (s *DelegationsService) Revoke(ctx context.Context, delegationID string) (*Delegation, error) {
	var resp Delegation
	if err := s.client.request(ctx, http.MethodPost, "/delegations/"+delegationID+"/revoke", nil, &resp); err != nil {
		return nil, fmt.Errorf("delegations.Revoke: %w", err)
	}
	return &resp, nil
}

func (s *DelegationsService) Verify(ctx context.Context, input *VerifyDelegationInput) (*VerifyDelegationResponse, error) {
	var resp VerifyDelegationResponse
	if err := s.client.request(ctx, http.MethodPost, "/delegations/verify", input, &resp); err != nil {
		return nil, fmt.Errorf("delegations.Verify: %w", err)
	}
	return &resp, nil
}

func (s *DelegationsService) List(ctx context.Context, input *ListDelegationsInput) ([]Delegation, error) {
	q := ""
	if input != nil {
		q = queryString(map[string]interface{}{
			"status": input.Status,
			"page":   input.Page,
			"limit":  input.Limit,
		})
	}
	var resp []Delegation
	if err := s.client.request(ctx, http.MethodGet, "/delegations"+q, nil, &resp); err != nil {
		return nil, fmt.Errorf("delegations.List: %w", err)
	}
	return resp, nil
}

func (s *DelegationsService) ListByAgent(ctx context.Context, agentID string, input *ListAgentDelegationsInput) ([]Delegation, error) {
	q := ""
	if input != nil {
		q = queryString(map[string]interface{}{
			"direction": input.Direction,
			"status":    input.Status,
			"page":      input.Page,
			"limit":     input.Limit,
		})
	}
	var resp []Delegation
	if err := s.client.request(ctx, http.MethodGet, "/agents/"+agentID+"/delegations"+q, nil, &resp); err != nil {
		return nil, fmt.Errorf("delegations.ListByAgent: %w", err)
	}
	return resp, nil
}
