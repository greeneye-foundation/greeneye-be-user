package errors

import (
	"net/http"
)

const (
	InvalidToken = 4001
)

// CustomError represents a structured error with more context
type CustomError struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

// Error implements the error interface
func (e *CustomError) Error() string {
	return e.Message
}

// Predefined errors
var (
	ErrBadRequest     = &CustomError{Code: http.StatusBadRequest, Message: "Bad Request"}
	ErrUnauthorized   = &CustomError{Code: http.StatusUnauthorized, Message: "Unauthorized"}
	ErrForbidden      = &CustomError{Code: http.StatusForbidden, Message: "Forbidden"}
	ErrNotFound       = &CustomError{Code: http.StatusNotFound, Message: "Not Found"}
	ErrInternalServer = &CustomError{Code: http.StatusInternalServerError, Message: "Internal Server Error"}
)

// New creates a new CustomError
func New(code int, message string, details ...string) *CustomError {
	return &CustomError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// GetHTTPStatusCode extracts the HTTP status code from the error if it's a CustomError
func GetHTTPStatusCode(err error) int {
	if customErr, ok := err.(*CustomError); ok {
		return customErr.Code
	}
	return http.StatusInternalServerError
}
