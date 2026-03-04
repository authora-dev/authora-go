package authora

import (
	"context"
	"fmt"
	"net/http"
)

type OrganizationsService struct {
	client *httpClient
}

func (s *OrganizationsService) Create(ctx context.Context, input *CreateOrganizationInput) (*Organization, error) {
	var resp Organization
	if err := s.client.request(ctx, http.MethodPost, "/organizations", input, &resp); err != nil {
		return nil, fmt.Errorf("organizations.Create: %w", err)
	}
	return &resp, nil
}

func (s *OrganizationsService) Get(ctx context.Context, orgID string) (*Organization, error) {
	var resp Organization
	if err := s.client.request(ctx, http.MethodGet, "/organizations/"+orgID, nil, &resp); err != nil {
		return nil, fmt.Errorf("organizations.Get: %w", err)
	}
	return &resp, nil
}

func (s *OrganizationsService) List(ctx context.Context, input *ListOrganizationsInput) (*PaginatedResponse[Organization], error) {
	q := ""
	if input != nil {
		q = queryString(map[string]interface{}{
			"page":  input.Page,
			"limit": input.Limit,
		})
	}
	var resp PaginatedResponse[Organization]
	if err := s.client.request(ctx, http.MethodGet, "/organizations"+q, nil, &resp); err != nil {
		return nil, fmt.Errorf("organizations.List: %w", err)
	}
	return &resp, nil
}
