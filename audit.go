package authora

import (
	"context"
	"fmt"
	"net/http"
)

// AuditService handles audit-related API endpoints.
type AuditService struct {
	client *httpClient
}

// ListEvents returns a paginated list of audit events. GET /audit/events
func (s *AuditService) ListEvents(ctx context.Context, input *ListAuditEventsInput) (*PaginatedResponse[AuditEvent], error) {
	q := ""
	if input != nil {
		q = queryString(map[string]interface{}{
			"orgId":       input.OrgID,
			"workspaceId": input.WorkspaceID,
			"agentId":     input.AgentID,
			"type":        input.Type,
			"dateFrom":    input.DateFrom,
			"dateTo":      input.DateTo,
			"resource":    input.Resource,
			"result":      input.Result,
			"page":        input.Page,
			"limit":       input.Limit,
		})
	}
	var resp PaginatedResponse[AuditEvent]
	if err := s.client.request(ctx, http.MethodGet, "/audit/events"+q, nil, &resp); err != nil {
		return nil, fmt.Errorf("audit.ListEvents: %w", err)
	}
	return &resp, nil
}

// GetEvent retrieves a single audit event by ID. GET /audit/events/:eventId
func (s *AuditService) GetEvent(ctx context.Context, eventID string) (*AuditEvent, error) {
	var resp AuditEvent
	if err := s.client.request(ctx, http.MethodGet, "/audit/events/"+eventID, nil, &resp); err != nil {
		return nil, fmt.Errorf("audit.GetEvent: %w", err)
	}
	return &resp, nil
}

// CreateReport generates an audit report. POST /audit/reports
func (s *AuditService) CreateReport(ctx context.Context, input *CreateAuditReportInput) (*AuditReport, error) {
	var resp AuditReport
	if err := s.client.request(ctx, http.MethodPost, "/audit/reports", input, &resp); err != nil {
		return nil, fmt.Errorf("audit.CreateReport: %w", err)
	}
	return &resp, nil
}

// GetMetrics returns audit metrics. GET /audit/metrics
func (s *AuditService) GetMetrics(ctx context.Context, input *AuditMetricsInput) ([]AuditMetricRow, error) {
	q := queryString(map[string]interface{}{
		"orgId":       input.OrgID,
		"workspaceId": input.WorkspaceID,
		"agentId":     input.AgentID,
		"dateFrom":    input.DateFrom,
		"dateTo":      input.DateTo,
	})
	var resp []AuditMetricRow
	if err := s.client.request(ctx, http.MethodGet, "/audit/metrics"+q, nil, &resp); err != nil {
		return nil, fmt.Errorf("audit.GetMetrics: %w", err)
	}
	return resp, nil
}
