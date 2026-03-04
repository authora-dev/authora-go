package authora

import (
	"context"
	"fmt"
	"net/http"
)

// WebhooksService handles webhook-related API endpoints.
type WebhooksService struct {
	client *httpClient
}

// Create registers a new webhook. POST /webhooks
func (s *WebhooksService) Create(ctx context.Context, input *CreateWebhookInput) (*Webhook, error) {
	var resp Webhook
	if err := s.client.request(ctx, http.MethodPost, "/webhooks", input, &resp); err != nil {
		return nil, fmt.Errorf("webhooks.Create: %w", err)
	}
	return &resp, nil
}

// List returns webhooks for an organization. GET /webhooks
func (s *WebhooksService) List(ctx context.Context, input *ListWebhooksInput) ([]Webhook, error) {
	q := queryString(map[string]interface{}{
		"organizationId": input.OrganizationID,
	})
	var resp []Webhook
	if err := s.client.request(ctx, http.MethodGet, "/webhooks"+q, nil, &resp); err != nil {
		return nil, fmt.Errorf("webhooks.List: %w", err)
	}
	return resp, nil
}

// Update modifies an existing webhook. PATCH /webhooks/:webhookId
func (s *WebhooksService) Update(ctx context.Context, webhookID string, input *UpdateWebhookInput) (*Webhook, error) {
	var resp Webhook
	if err := s.client.request(ctx, http.MethodPatch, "/webhooks/"+webhookID, input, &resp); err != nil {
		return nil, fmt.Errorf("webhooks.Update: %w", err)
	}
	return &resp, nil
}

// Delete removes a webhook. DELETE /webhooks/:webhookId
func (s *WebhooksService) Delete(ctx context.Context, webhookID string) error {
	if err := s.client.request(ctx, http.MethodDelete, "/webhooks/"+webhookID, nil, nil); err != nil {
		return fmt.Errorf("webhooks.Delete: %w", err)
	}
	return nil
}
