package authora

import (
	"context"
	"fmt"
	"net/http"
)

// UserDelegationsService manages user-to-agent delegation grants.
type UserDelegationsService struct {
	client *httpClient
}

// Create creates a user delegation grant.
func (s *UserDelegationsService) Create(ctx context.Context, input *CreateUserDelegationInput) (*UserDelegationGrant, error) {
	var resp UserDelegationGrant
	if err := s.client.request(ctx, http.MethodPost, "/user-delegations", input, &resp); err != nil {
		return nil, fmt.Errorf("userDelegations.Create: %w", err)
	}
	return &resp, nil
}

// Get retrieves a specific delegation grant by ID.
func (s *UserDelegationsService) Get(ctx context.Context, grantID string) (*UserDelegationGrant, error) {
	var resp UserDelegationGrant
	if err := s.client.request(ctx, http.MethodGet, "/user-delegations/"+grantID, nil, &resp); err != nil {
		return nil, fmt.Errorf("userDelegations.Get: %w", err)
	}
	return &resp, nil
}

// ListByUser lists delegation grants by user.
func (s *UserDelegationsService) ListByUser(ctx context.Context, userID string, input *ListUserDelegationInput) ([]UserDelegationGrant, error) {
	q := ""
	if input != nil {
		q = queryString(map[string]interface{}{
			"status": input.Status,
		})
	}
	var resp []UserDelegationGrant
	if err := s.client.request(ctx, http.MethodGet, "/user-delegations/by-user/"+userID+q, nil, &resp); err != nil {
		return nil, fmt.Errorf("userDelegations.ListByUser: %w", err)
	}
	return resp, nil
}

// ListByAgent lists delegation grants by agent.
func (s *UserDelegationsService) ListByAgent(ctx context.Context, agentID string, input *ListUserDelegationInput) ([]UserDelegationGrant, error) {
	q := ""
	if input != nil {
		q = queryString(map[string]interface{}{
			"status": input.Status,
		})
	}
	var resp []UserDelegationGrant
	if err := s.client.request(ctx, http.MethodGet, "/user-delegations/by-agent/"+agentID+q, nil, &resp); err != nil {
		return nil, fmt.Errorf("userDelegations.ListByAgent: %w", err)
	}
	return resp, nil
}

// ListByOrg lists delegation grants by organization.
func (s *UserDelegationsService) ListByOrg(ctx context.Context, orgID string, input *ListUserDelegationOrgInput) (*UserDelegationOrgResponse, error) {
	q := ""
	if input != nil {
		q = queryString(map[string]interface{}{
			"status": input.Status,
			"page":   input.Page,
			"limit":  input.Limit,
		})
	}
	var resp UserDelegationOrgResponse
	if err := s.client.request(ctx, http.MethodGet, "/user-delegations/by-org/"+orgID+q, nil, &resp); err != nil {
		return nil, fmt.Errorf("userDelegations.ListByOrg: %w", err)
	}
	return &resp, nil
}

// Revoke revokes a delegation grant.
func (s *UserDelegationsService) Revoke(ctx context.Context, grantID string, input *RevokeUserDelegationInput) (*UserDelegationGrant, error) {
	var resp UserDelegationGrant
	if err := s.client.request(ctx, http.MethodPost, "/user-delegations/"+grantID+"/revoke", input, &resp); err != nil {
		return nil, fmt.Errorf("userDelegations.Revoke: %w", err)
	}
	return &resp, nil
}

// IssueToken issues a fresh delegation JWT from a grant.
func (s *UserDelegationsService) IssueToken(ctx context.Context, grantID string, input *IssueUserDelegationTokenInput) (*UserDelegationToken, error) {
	var resp UserDelegationToken
	if err := s.client.request(ctx, http.MethodPost, "/user-delegations/"+grantID+"/token", input, &resp); err != nil {
		return nil, fmt.Errorf("userDelegations.IssueToken: %w", err)
	}
	return &resp, nil
}

// RefreshToken refreshes a delegation JWT.
func (s *UserDelegationsService) RefreshToken(ctx context.Context, grantID string, input *RefreshUserDelegationTokenInput) (*UserDelegationToken, error) {
	var resp UserDelegationToken
	if err := s.client.request(ctx, http.MethodPost, "/user-delegations/"+grantID+"/refresh", input, &resp); err != nil {
		return nil, fmt.Errorf("userDelegations.RefreshToken: %w", err)
	}
	return &resp, nil
}

// VerifyToken verifies a delegation JWT.
func (s *UserDelegationsService) VerifyToken(ctx context.Context, input *VerifyUserDelegationTokenInput) (*VerifyUserDelegationTokenResult, error) {
	var resp VerifyUserDelegationTokenResult
	if err := s.client.request(ctx, http.MethodPost, "/user-delegations/tokens/verify", input, &resp); err != nil {
		return nil, fmt.Errorf("userDelegations.VerifyToken: %w", err)
	}
	return &resp, nil
}
