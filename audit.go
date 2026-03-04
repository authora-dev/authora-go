package authora

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type AuditService struct {
	client *httpClient
}

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

func (s *AuditService) GetEvent(ctx context.Context, eventID string) (*AuditEvent, error) {
	var resp AuditEvent
	if err := s.client.request(ctx, http.MethodGet, "/audit/events/"+eventID, nil, &resp); err != nil {
		return nil, fmt.Errorf("audit.GetEvent: %w", err)
	}
	return &resp, nil
}

func (s *AuditService) CreateReport(ctx context.Context, input *CreateAuditReportInput) (*AuditReport, error) {
	var resp AuditReport
	if err := s.client.request(ctx, http.MethodPost, "/audit/reports", input, &resp); err != nil {
		return nil, fmt.Errorf("audit.CreateReport: %w", err)
	}
	return &resp, nil
}

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

func (s *AuditService) StreamEvents(ctx context.Context, onEvent func(AuditEvent)) error {
	url := s.client.baseURL + "/audit/stream"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("audit.StreamEvents: %w", err)
	}
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Authorization", "Bearer "+s.client.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("audit.StreamEvents: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("audit.StreamEvents: unexpected status %d", resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)
	var eventType, data string

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "event: ") {
			eventType = strings.TrimSpace(line[7:])
		} else if strings.HasPrefix(line, "data: ") {
			data = line[6:]
		} else if line == "" && data != "" {
			if eventType == "audit" {
				var ev AuditEvent
				if json.Unmarshal([]byte(data), &ev) == nil {
					onEvent(ev)
				}
			}
			eventType = ""
			data = ""
		}
	}

	return scanner.Err()
}
