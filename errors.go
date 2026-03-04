package authora

import (
	"errors"
	"fmt"
)

type AuthoraError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Code       string `json:"code"`
}

func (e *AuthoraError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("authora: %s (code=%s, status=%d)", e.Message, e.Code, e.StatusCode)
	}
	return fmt.Sprintf("authora: %s (status=%d)", e.Message, e.StatusCode)
}

func IsNotFoundError(err error) bool {
	var ae *AuthoraError
	if errors.As(err, &ae) {
		return ae.StatusCode == 404
	}
	return false
}

func IsAuthenticationError(err error) bool {
	var ae *AuthoraError
	if errors.As(err, &ae) {
		return ae.StatusCode == 401
	}
	return false
}

func IsRateLimitError(err error) bool {
	var ae *AuthoraError
	if errors.As(err, &ae) {
		return ae.StatusCode == 429
	}
	return false
}

func IsForbiddenError(err error) bool {
	var ae *AuthoraError
	if errors.As(err, &ae) {
		return ae.StatusCode == 403
	}
	return false
}

func IsValidationError(err error) bool {
	var ae *AuthoraError
	if errors.As(err, &ae) {
		return ae.StatusCode == 400
	}
	return false
}

func IsConflictError(err error) bool {
	var ae *AuthoraError
	if errors.As(err, &ae) {
		return ae.StatusCode == 409
	}
	return false
}
