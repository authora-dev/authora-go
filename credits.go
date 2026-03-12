package authora

import (
	"context"
	"fmt"
	"net/http"
)

type CreditsService struct {
	client *httpClient
}

func (s *CreditsService) Balance(ctx context.Context) (*CreditBalance, error) {
	var resp CreditBalance
	if err := s.client.request(ctx, http.MethodGet, "/credits", nil, &resp); err != nil {
		return nil, fmt.Errorf("credits.Balance: %w", err)
	}
	return &resp, nil
}

func (s *CreditsService) Transactions(ctx context.Context, input *ListCreditTransactionsInput) (*PaginatedResponse[CreditTransaction], error) {
	q := queryString(map[string]interface{}{
		"type":   input.Type,
		"limit":  input.Limit,
		"offset": input.Offset,
	})
	var resp PaginatedResponse[CreditTransaction]
	if err := s.client.request(ctx, http.MethodGet, "/credits/transactions"+q, nil, &resp); err != nil {
		return nil, fmt.Errorf("credits.Transactions: %w", err)
	}
	return &resp, nil
}

func (s *CreditsService) Checkout(ctx context.Context, pack string) (*CreditCheckoutResult, error) {
	body := map[string]string{"pack": pack}
	var resp CreditCheckoutResult
	if err := s.client.request(ctx, http.MethodPost, "/credits/checkout", body, &resp); err != nil {
		return nil, fmt.Errorf("credits.Checkout: %w", err)
	}
	return &resp, nil
}
