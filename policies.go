package authora

import (
	"context"
	"fmt"
	"net/http"
)

// PoliciesService handles policy-related API endpoints.
type PoliciesService struct {
	client *httpClient
}

// Create creates a new policy. POST /policies
func (s *PoliciesService) Create(ctx context.Context, input *CreatePolicyInput) (*Policy, error) {
	var resp Policy
	if err := s.client.request(ctx, http.MethodPost, "/policies", input, &resp); err != nil {
		return nil, fmt.Errorf("policies.Create: %w", err)
	}
	return &resp, nil
}

// List returns policies for a workspace. GET /policies
func (s *PoliciesService) List(ctx context.Context, input *ListPoliciesInput) ([]Policy, error) {
	q := queryString(map[string]interface{}{
		"workspaceId": input.WorkspaceID,
		"page":        input.Page,
		"limit":       input.Limit,
	})
	var resp []Policy
	if err := s.client.request(ctx, http.MethodGet, "/policies"+q, nil, &resp); err != nil {
		return nil, fmt.Errorf("policies.List: %w", err)
	}
	return resp, nil
}

// Update modifies an existing policy. PATCH /policies/:policyId
func (s *PoliciesService) Update(ctx context.Context, policyID string, input *UpdatePolicyInput) (*Policy, error) {
	var resp Policy
	if err := s.client.request(ctx, http.MethodPatch, "/policies/"+policyID, input, &resp); err != nil {
		return nil, fmt.Errorf("policies.Update: %w", err)
	}
	return &resp, nil
}

// Delete removes a policy. DELETE /policies/:policyId
func (s *PoliciesService) Delete(ctx context.Context, policyID string) error {
	if err := s.client.request(ctx, http.MethodDelete, "/policies/"+policyID, nil, nil); err != nil {
		return fmt.Errorf("policies.Delete: %w", err)
	}
	return nil
}

// Simulate runs a policy simulation without enforcing. POST /policies/simulate
func (s *PoliciesService) Simulate(ctx context.Context, input *SimulatePolicyInput) (*SimulatePolicyResponse, error) {
	var resp SimulatePolicyResponse
	if err := s.client.request(ctx, http.MethodPost, "/policies/simulate", input, &resp); err != nil {
		return nil, fmt.Errorf("policies.Simulate: %w", err)
	}
	return &resp, nil
}

// Evaluate evaluates policies against a request. POST /policies/evaluate
func (s *PoliciesService) Evaluate(ctx context.Context, input *EvaluatePolicyInput) (*EvaluatePolicyResponse, error) {
	var resp EvaluatePolicyResponse
	if err := s.client.request(ctx, http.MethodPost, "/policies/evaluate", input, &resp); err != nil {
		return nil, fmt.Errorf("policies.Evaluate: %w", err)
	}
	return &resp, nil
}
