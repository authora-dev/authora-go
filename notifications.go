package authora

import (
	"context"
	"fmt"
	"net/http"
)

type NotificationsService struct {
	client *httpClient
}

func (s *NotificationsService) List(ctx context.Context, input *ListNotificationsInput) (*PaginatedResponse[Notification], error) {
	q := queryString(map[string]interface{}{
		"organizationId": input.OrganizationID,
		"userId":         input.UserID,
		"unreadOnly":     input.UnreadOnly,
		"limit":          input.Limit,
		"offset":         input.Offset,
	})
	var resp PaginatedResponse[Notification]
	if err := s.client.request(ctx, http.MethodGet, "/notifications"+q, nil, &resp); err != nil {
		return nil, fmt.Errorf("notifications.List: %w", err)
	}
	return &resp, nil
}

func (s *NotificationsService) UnreadCount(ctx context.Context, input *UnreadCountInput) (*UnreadCountResponse, error) {
	q := queryString(map[string]interface{}{
		"organizationId": input.OrganizationID,
		"userId":         input.UserID,
	})
	var resp UnreadCountResponse
	if err := s.client.request(ctx, http.MethodGet, "/notifications/unread-count"+q, nil, &resp); err != nil {
		return nil, fmt.Errorf("notifications.UnreadCount: %w", err)
	}
	return &resp, nil
}

func (s *NotificationsService) MarkRead(ctx context.Context, notificationID string) (*Notification, error) {
	var resp Notification
	if err := s.client.request(ctx, http.MethodPatch, "/notifications/"+notificationID+"/read", nil, &resp); err != nil {
		return nil, fmt.Errorf("notifications.MarkRead: %w", err)
	}
	return &resp, nil
}

func (s *NotificationsService) MarkAllRead(ctx context.Context, input *MarkAllReadInput) error {
	if err := s.client.request(ctx, http.MethodPatch, "/notifications/read-all", input, nil); err != nil {
		return fmt.Errorf("notifications.MarkAllRead: %w", err)
	}
	return nil
}
