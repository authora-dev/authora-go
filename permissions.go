package authora

import (
	"context"
	"fmt"
	"net/http"
)

type PermissionsService struct {
	client *httpClient
}

func (s *PermissionsService) Check(ctx context.Context, input *CheckPermissionInput) (*CheckPermissionResponse, error) {
	var resp CheckPermissionResponse
	if err := s.client.request(ctx, http.MethodPost, "/permissions/check", input, &resp); err != nil {
		return nil, fmt.Errorf("permissions.Check: %w", err)
	}
	return &resp, nil
}

func (s *PermissionsService) CheckBatch(ctx context.Context, input *BatchCheckInput) (*BatchCheckResponse, error) {
	var resp BatchCheckResponse
	if err := s.client.request(ctx, http.MethodPost, "/permissions/check-batch", input, &resp); err != nil {
		return nil, fmt.Errorf("permissions.CheckBatch: %w", err)
	}
	return &resp, nil
}

func (s *PermissionsService) Effective(ctx context.Context, agentID string) (*EffectivePermissionsResponse, error) {
	var resp EffectivePermissionsResponse
	if err := s.client.request(ctx, http.MethodGet, "/agents/"+agentID+"/permissions", nil, &resp); err != nil {
		return nil, fmt.Errorf("permissions.Effective: %w", err)
	}
	return &resp, nil
}
