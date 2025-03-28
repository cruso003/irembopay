package irembopay

import (
	"fmt"
)

// IremboPayError represents an error from the IremboPay API
type IremboPayError struct {
	StatusCode int
	Message    string
	Details    string
}

// Error implements the error interface
func (e *IremboPayError) Error() string {
	return fmt.Sprintf("IremboPay API error (HTTP %d): %s", e.StatusCode, e.Message)
}

// NewIremboPayError creates a new IremboPayError
func NewIremboPayError(statusCode int, message, details string) *IremboPayError {
	return &IremboPayError{
		StatusCode: statusCode,
		Message:    message,
		Details:    details,
	}
}

// IsNotFoundError checks if the error is a 404 Not Found error
func IsNotFoundError(err error) bool {
	if iremboErr, ok := err.(*IremboPayError); ok {
		return iremboErr.StatusCode == 404
	}
	return false
}

// IsBadRequestError checks if the error is a 400 Bad Request error
func IsBadRequestError(err error) bool {
	if iremboErr, ok := err.(*IremboPayError); ok {
		return iremboErr.StatusCode == 400
	}
	return false
}
