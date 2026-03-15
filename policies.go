package authora

import (
	"context"
	"fmt"
	"net/http"
)

type PoliciesService struct {
	client *httpClient
}

func (s *PoliciesService) Create(ctx context.Context, input *CreatePolicyInput) (*Policy, error) {
	var resp Policy
	if err := s.client.request(ctx, http.MethodPost, "/policies", input, &resp); err != nil {
		return nil, fmt.Errorf("policies.Create: %w", err)
	}
	return &resp, nil
}

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

func (s *PoliciesService) Update(ctx context.Context, policyID string, input *UpdatePolicyInput) (*Policy, error) {
	var resp Policy
	if err := s.client.request(ctx, http.MethodPatch, "/policies/"+policyID, input, &resp); err != nil {
		return nil, fmt.Errorf("policies.Update: %w", err)
	}
	return &resp, nil
}

func (s *PoliciesService) Delete(ctx context.Context, policyID string) error {
	if err := s.client.request(ctx, http.MethodDelete, "/policies/"+policyID, nil, nil); err != nil {
		return fmt.Errorf("policies.Delete: %w", err)
	}
	return nil
}

func (s *PoliciesService) Simulate(ctx context.Context, input *SimulatePolicyInput) (*SimulatePolicyResponse, error) {
	var resp SimulatePolicyResponse
	if err := s.client.request(ctx, http.MethodPost, "/policies/simulate", input, &resp); err != nil {
		return nil, fmt.Errorf("policies.Simulate: %w", err)
	}
	return &resp, nil
}

func (s *PoliciesService) Evaluate(ctx context.Context, input *EvaluatePolicyInput) (*EvaluatePolicyResponse, error) {
	var resp EvaluatePolicyResponse
	if err := s.client.request(ctx, http.MethodPost, "/policies/evaluate", input, &resp); err != nil {
		return nil, fmt.Errorf("policies.Evaluate: %w", err)
	}
	return &resp, nil
}

func (s *PoliciesService) AttachToTarget(ctx context.Context, input *AttachPolicyInput) (*PolicyAttachment, error) {
	var resp PolicyAttachment
	if err := s.client.request(ctx, http.MethodPost, "/policies/attachments", input, &resp); err != nil {
		return nil, fmt.Errorf("policies.AttachToTarget: %w", err)
	}
	return &resp, nil
}

func (s *PoliciesService) DetachFromTarget(ctx context.Context, input *AttachPolicyInput) error {
	if err := s.client.request(ctx, http.MethodPost, "/policies/detach", input, nil); err != nil {
		return fmt.Errorf("policies.DetachFromTarget: %w", err)
	}
	return nil
}

func (s *PoliciesService) DetachByID(ctx context.Context, attachmentID string) error {
	if err := s.client.request(ctx, http.MethodDelete, "/policies/attachments/"+attachmentID, nil, nil); err != nil {
		return fmt.Errorf("policies.DetachByID: %w", err)
	}
	return nil
}

func (s *PoliciesService) ListAttachments(ctx context.Context, input *ListAttachmentsInput) ([]PolicyAttachment, error) {
	q := queryString(map[string]interface{}{
		"targetType": input.TargetType,
		"targetId":   input.TargetID,
	})
	var resp []PolicyAttachment
	if err := s.client.request(ctx, http.MethodGet, "/policies/attachments"+q, nil, &resp); err != nil {
		return nil, fmt.Errorf("policies.ListAttachments: %w", err)
	}
	return resp, nil
}

func (s *PoliciesService) ListPolicyTargets(ctx context.Context, policyID string) ([]PolicyAttachment, error) {
	var resp []PolicyAttachment
	if err := s.client.request(ctx, http.MethodGet, "/policies/"+policyID+"/attachments", nil, &resp); err != nil {
		return nil, fmt.Errorf("policies.ListPolicyTargets: %w", err)
	}
	return resp, nil
}

func (s *PoliciesService) AddPermission(ctx context.Context, input *AddPermissionInput) (*Policy, error) {
	var resp Policy
	if err := s.client.request(ctx, http.MethodPost, "/policies/add-permission", input, &resp); err != nil {
		return nil, fmt.Errorf("policies.AddPermission: %w", err)
	}
	return &resp, nil
}

func (s *PoliciesService) RemovePermission(ctx context.Context, input *RemovePermissionInput) (*Policy, error) {
	var resp Policy
	if err := s.client.request(ctx, http.MethodPost, "/policies/remove-permission", input, &resp); err != nil {
		return nil, fmt.Errorf("policies.RemovePermission: %w", err)
	}
	return &resp, nil
}
