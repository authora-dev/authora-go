package authora

import (
	"errors"
	"fmt"
)

// AuthoraError represents an API error returned by the Authora service.
type AuthoraError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Code       string `json:"code"`
}

// Error implements the error interface.
func (e *AuthoraError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("authora: %s (code=%s, status=%d)", e.Message, e.Code, e.StatusCode)
	}
	return fmt.Sprintf("authora: %s (status=%d)", e.Message, e.StatusCode)
}

// IsNotFoundError returns true if the error is a 404 Not Found response.
func IsNotFoundError(err error) bool {
	var ae *AuthoraError
	if errors.As(err, &ae) {
		return ae.StatusCode == 404
	}
	return false
}

// IsAuthenticationError returns true if the error is a 401 Unauthorized response.
func IsAuthenticationError(err error) bool {
	var ae *AuthoraError
	if errors.As(err, &ae) {
		return ae.StatusCode == 401
	}
	return false
}

// IsRateLimitError returns true if the error is a 429 Too Many Requests response.
func IsRateLimitError(err error) bool {
	var ae *AuthoraError
	if errors.As(err, &ae) {
		return ae.StatusCode == 429
	}
	return false
}

// IsForbiddenError returns true if the error is a 403 Forbidden response.
func IsForbiddenError(err error) bool {
	var ae *AuthoraError
	if errors.As(err, &ae) {
		return ae.StatusCode == 403
	}
	return false
}

// IsValidationError returns true if the error is a 400 Bad Request response.
func IsValidationError(err error) bool {
	var ae *AuthoraError
	if errors.As(err, &ae) {
		return ae.StatusCode == 400
	}
	return false
}

// IsConflictError returns true if the error is a 409 Conflict response.
func IsConflictError(err error) bool {
	var ae *AuthoraError
	if errors.As(err, &ae) {
		return ae.StatusCode == 409
	}
	return false
}
