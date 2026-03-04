package authora

import (
	"context"
	"fmt"
	"net/http"
)

type APIKeysService struct {
	client *httpClient
}

func (s *APIKeysService) Create(ctx context.Context, input *CreateAPIKeyInput) (*CreateAPIKeyResponse, error) {
	var resp CreateAPIKeyResponse
	if err := s.client.request(ctx, http.MethodPost, "/api-keys", input, &resp); err != nil {
		return nil, fmt.Errorf("apikeys.Create: %w", err)
	}
	return &resp, nil
}

func (s *APIKeysService) List(ctx context.Context, input *ListAPIKeysInput) ([]APIKey, error) {
	q := queryString(map[string]interface{}{
		"organizationId": input.OrganizationID,
	})
	var resp []APIKey
	if err := s.client.request(ctx, http.MethodGet, "/api-keys"+q, nil, &resp); err != nil {
		return nil, fmt.Errorf("apikeys.List: %w", err)
	}
	return resp, nil
}

func (s *APIKeysService) Revoke(ctx context.Context, keyID string) error {
	if err := s.client.request(ctx, http.MethodDelete, "/api-keys/"+keyID, nil, nil); err != nil {
		return fmt.Errorf("apikeys.Revoke: %w", err)
	}
	return nil
}
