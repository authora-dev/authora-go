package authora

import (
	"context"
	"fmt"
	"net/http"
)

type ApprovalsService struct {
	client *httpClient
}

func (s *ApprovalsService) Create(ctx context.Context, input *CreateApprovalInput) (*ApprovalChallenge, error) {
	var resp ApprovalChallenge
	if err := s.client.request(ctx, http.MethodPost, "/approvals", input, &resp); err != nil {
		return nil, fmt.Errorf("approvals.Create: %w", err)
	}
	return &resp, nil
}

func (s *ApprovalsService) List(ctx context.Context, input *ListApprovalsInput) (*PaginatedResponse[ApprovalChallenge], error) {
	q := queryString(map[string]interface{}{
		"status": input.Status,
		"limit":  input.Limit,
		"offset": input.Offset,
	})
	var resp PaginatedResponse[ApprovalChallenge]
	if err := s.client.request(ctx, http.MethodGet, "/approvals"+q, nil, &resp); err != nil {
		return nil, fmt.Errorf("approvals.List: %w", err)
	}
	return &resp, nil
}

func (s *ApprovalsService) Get(ctx context.Context, id string) (*ApprovalChallenge, error) {
	var resp ApprovalChallenge
	if err := s.client.request(ctx, http.MethodGet, "/approvals/"+id, nil, &resp); err != nil {
		return nil, fmt.Errorf("approvals.Get: %w", err)
	}
	return &resp, nil
}

func (s *ApprovalsService) GetStatus(ctx context.Context, id string) (*ApprovalStatusResponse, error) {
	var resp ApprovalStatusResponse
	if err := s.client.request(ctx, http.MethodGet, "/approvals/"+id+"/status", nil, &resp); err != nil {
		return nil, fmt.Errorf("approvals.GetStatus: %w", err)
	}
	return &resp, nil
}

func (s *ApprovalsService) Stats(ctx context.Context) (*ApprovalStats, error) {
	var resp ApprovalStats
	if err := s.client.request(ctx, http.MethodGet, "/approvals/stats", nil, &resp); err != nil {
		return nil, fmt.Errorf("approvals.Stats: %w", err)
	}
	return &resp, nil
}

func (s *ApprovalsService) Decide(ctx context.Context, id string, input *DecideApprovalInput) (*ApprovalChallenge, error) {
	var resp ApprovalChallenge
	if err := s.client.request(ctx, http.MethodPost, "/approvals/"+id+"/decide", input, &resp); err != nil {
		return nil, fmt.Errorf("approvals.Decide: %w", err)
	}
	return &resp, nil
}

func (s *ApprovalsService) BulkDecide(ctx context.Context, input *BulkDecideInput) (*BulkDecideResult, error) {
	var resp BulkDecideResult
	if err := s.client.request(ctx, http.MethodPost, "/approvals/bulk-decide", input, &resp); err != nil {
		return nil, fmt.Errorf("approvals.BulkDecide: %w", err)
	}
	return &resp, nil
}

func (s *ApprovalsService) Suggestions(ctx context.Context, id string) ([]PermissionSuggestion, error) {
	var resp []PermissionSuggestion
	if err := s.client.request(ctx, http.MethodPost, "/approvals/"+id+"/suggestions", nil, &resp); err != nil {
		return nil, fmt.Errorf("approvals.Suggestions: %w", err)
	}
	return resp, nil
}

func (s *ApprovalsService) GetSettings(ctx context.Context) (*ApprovalSettings, error) {
	var resp ApprovalSettings
	if err := s.client.request(ctx, http.MethodGet, "/approvals/settings", nil, &resp); err != nil {
		return nil, fmt.Errorf("approvals.GetSettings: %w", err)
	}
	return &resp, nil
}

func (s *ApprovalsService) UpdateSettings(ctx context.Context, input *UpdateApprovalSettingsInput) (*ApprovalSettings, error) {
	var resp ApprovalSettings
	if err := s.client.request(ctx, http.MethodPatch, "/approvals/settings", input, &resp); err != nil {
		return nil, fmt.Errorf("approvals.UpdateSettings: %w", err)
	}
	return &resp, nil
}

func (s *ApprovalsService) TestAi(ctx context.Context, input *TestAiInput) (*TestAiResult, error) {
	var resp TestAiResult
	if err := s.client.request(ctx, http.MethodPost, "/approvals/settings/test-ai", input, &resp); err != nil {
		return nil, fmt.Errorf("approvals.TestAi: %w", err)
	}
	return &resp, nil
}

func (s *ApprovalsService) ListPatterns(ctx context.Context, input *ListPatternsInput) ([]ApprovalPattern, error) {
	q := queryString(map[string]interface{}{
		"status":    input.Status,
		"readyOnly": input.ReadyOnly,
	})
	var resp []ApprovalPattern
	if err := s.client.request(ctx, http.MethodGet, "/approvals/patterns"+q, nil, &resp); err != nil {
		return nil, fmt.Errorf("approvals.ListPatterns: %w", err)
	}
	return resp, nil
}

func (s *ApprovalsService) DismissPattern(ctx context.Context, id string) error {
	if err := s.client.request(ctx, http.MethodPost, "/approvals/patterns/"+id+"/dismiss", nil, nil); err != nil {
		return fmt.Errorf("approvals.DismissPattern: %w", err)
	}
	return nil
}

func (s *ApprovalsService) CreatePolicyFromPattern(ctx context.Context, id string) (map[string]interface{}, error) {
	var resp map[string]interface{}
	if err := s.client.request(ctx, http.MethodPost, "/approvals/patterns/"+id+"/create-policy", nil, &resp); err != nil {
		return nil, fmt.Errorf("approvals.CreatePolicyFromPattern: %w", err)
	}
	return resp, nil
}

func (s *ApprovalsService) ListEscalationRules(ctx context.Context) ([]EscalationRule, error) {
	var resp []EscalationRule
	if err := s.client.request(ctx, http.MethodGet, "/approvals/escalation-rules", nil, &resp); err != nil {
		return nil, fmt.Errorf("approvals.ListEscalationRules: %w", err)
	}
	return resp, nil
}

func (s *ApprovalsService) GetEscalationRule(ctx context.Context, id string) (*EscalationRule, error) {
	var resp EscalationRule
	if err := s.client.request(ctx, http.MethodGet, "/approvals/escalation-rules/"+id, nil, &resp); err != nil {
		return nil, fmt.Errorf("approvals.GetEscalationRule: %w", err)
	}
	return &resp, nil
}

func (s *ApprovalsService) CreateEscalationRule(ctx context.Context, input *CreateEscalationRuleInput) (*EscalationRule, error) {
	var resp EscalationRule
	if err := s.client.request(ctx, http.MethodPost, "/approvals/escalation-rules", input, &resp); err != nil {
		return nil, fmt.Errorf("approvals.CreateEscalationRule: %w", err)
	}
	return &resp, nil
}

func (s *ApprovalsService) UpdateEscalationRule(ctx context.Context, id string, input *UpdateEscalationRuleInput) (*EscalationRule, error) {
	var resp EscalationRule
	if err := s.client.request(ctx, http.MethodPatch, "/approvals/escalation-rules/"+id, input, &resp); err != nil {
		return nil, fmt.Errorf("approvals.UpdateEscalationRule: %w", err)
	}
	return &resp, nil
}

func (s *ApprovalsService) DeleteEscalationRule(ctx context.Context, id string) error {
	if err := s.client.request(ctx, http.MethodDelete, "/approvals/escalation-rules/"+id, nil, nil); err != nil {
		return fmt.Errorf("approvals.DeleteEscalationRule: %w", err)
	}
	return nil
}

func (s *ApprovalsService) GetVapidKey(ctx context.Context) (*VapidKeyResponse, error) {
	var resp VapidKeyResponse
	if err := s.client.request(ctx, http.MethodGet, "/approvals/push/vapid-key", nil, &resp); err != nil {
		return nil, fmt.Errorf("approvals.GetVapidKey: %w", err)
	}
	return &resp, nil
}

func (s *ApprovalsService) SubscribePush(ctx context.Context, input *PushSubscribeInput) error {
	if err := s.client.request(ctx, http.MethodPost, "/approvals/push/subscribe", input, nil); err != nil {
		return fmt.Errorf("approvals.SubscribePush: %w", err)
	}
	return nil
}

func (s *ApprovalsService) UnsubscribePush(ctx context.Context, endpoint string) error {
	body := map[string]string{"endpoint": endpoint}
	if err := s.client.request(ctx, http.MethodPost, "/approvals/push/unsubscribe", body, nil); err != nil {
		return fmt.Errorf("approvals.UnsubscribePush: %w", err)
	}
	return nil
}

func (s *ApprovalsService) ListWebhooks(ctx context.Context) ([]ApprovalWebhook, error) {
	var resp []ApprovalWebhook
	if err := s.client.request(ctx, http.MethodGet, "/approvals/webhooks", nil, &resp); err != nil {
		return nil, fmt.Errorf("approvals.ListWebhooks: %w", err)
	}
	return resp, nil
}

func (s *ApprovalsService) CreateWebhook(ctx context.Context, input *CreateApprovalWebhookInput) (*ApprovalWebhook, error) {
	var resp ApprovalWebhook
	if err := s.client.request(ctx, http.MethodPost, "/approvals/webhooks", input, &resp); err != nil {
		return nil, fmt.Errorf("approvals.CreateWebhook: %w", err)
	}
	return &resp, nil
}

func (s *ApprovalsService) UpdateWebhook(ctx context.Context, id string, input *UpdateApprovalWebhookInput) (*ApprovalWebhook, error) {
	var resp ApprovalWebhook
	if err := s.client.request(ctx, http.MethodPatch, "/approvals/webhooks/"+id, input, &resp); err != nil {
		return nil, fmt.Errorf("approvals.UpdateWebhook: %w", err)
	}
	return &resp, nil
}

func (s *ApprovalsService) DeleteWebhook(ctx context.Context, id string) error {
	if err := s.client.request(ctx, http.MethodDelete, "/approvals/webhooks/"+id, nil, nil); err != nil {
		return fmt.Errorf("approvals.DeleteWebhook: %w", err)
	}
	return nil
}
