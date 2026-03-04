package authora

import (
	"context"
	"fmt"
	"net/http"
)

type AlertsService struct {
	client *httpClient
}

func (s *AlertsService) Create(ctx context.Context, input *CreateAlertInput) (*Alert, error) {
	var resp Alert
	if err := s.client.request(ctx, http.MethodPost, "/alerts", input, &resp); err != nil {
		return nil, fmt.Errorf("alerts.Create: %w", err)
	}
	return &resp, nil
}

func (s *AlertsService) List(ctx context.Context, input *ListAlertsInput) ([]Alert, error) {
	q := queryString(map[string]interface{}{
		"organizationId": input.OrganizationID,
	})
	var resp []Alert
	if err := s.client.request(ctx, http.MethodGet, "/alerts"+q, nil, &resp); err != nil {
		return nil, fmt.Errorf("alerts.List: %w", err)
	}
	return resp, nil
}

func (s *AlertsService) Update(ctx context.Context, alertID string, input *UpdateAlertInput) (*Alert, error) {
	var resp Alert
	if err := s.client.request(ctx, http.MethodPatch, "/alerts/"+alertID, input, &resp); err != nil {
		return nil, fmt.Errorf("alerts.Update: %w", err)
	}
	return &resp, nil
}

func (s *AlertsService) Delete(ctx context.Context, alertID string) error {
	if err := s.client.request(ctx, http.MethodDelete, "/alerts/"+alertID, nil, nil); err != nil {
		return fmt.Errorf("alerts.Delete: %w", err)
	}
	return nil
}
