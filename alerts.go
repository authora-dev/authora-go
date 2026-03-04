package authora

import (
	"context"
	"fmt"
	"net/http"
)

// AlertsService handles alert-related API endpoints.
type AlertsService struct {
	client *httpClient
}

// Create creates a new alert. POST /alerts
func (s *AlertsService) Create(ctx context.Context, input *CreateAlertInput) (*Alert, error) {
	var resp Alert
	if err := s.client.request(ctx, http.MethodPost, "/alerts", input, &resp); err != nil {
		return nil, fmt.Errorf("alerts.Create: %w", err)
	}
	return &resp, nil
}

// List returns alerts for an organization. GET /alerts
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

// Update modifies an existing alert. PATCH /alerts/:alertId
func (s *AlertsService) Update(ctx context.Context, alertID string, input *UpdateAlertInput) (*Alert, error) {
	var resp Alert
	if err := s.client.request(ctx, http.MethodPatch, "/alerts/"+alertID, input, &resp); err != nil {
		return nil, fmt.Errorf("alerts.Update: %w", err)
	}
	return &resp, nil
}

// Delete removes an alert. DELETE /alerts/:alertId
func (s *AlertsService) Delete(ctx context.Context, alertID string) error {
	if err := s.client.request(ctx, http.MethodDelete, "/alerts/"+alertID, nil, nil); err != nil {
		return fmt.Errorf("alerts.Delete: %w", err)
	}
	return nil
}
