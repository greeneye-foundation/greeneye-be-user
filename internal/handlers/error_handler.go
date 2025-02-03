package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error   string      `json:"error"`
	Details interface{} `json:"details,omitempty"`
}

type ValidationError struct {
	Details interface{}
}

func (v *ValidationError) Error() string {
	return "Validation Error"
}

func ErrorHandler(c *gin.Context, err error) {
	code := http.StatusInternalServerError
	response := ErrorResponse{
		Error: "Internal Server Error",
	}

	// Handle standard errors
	var httpError *gin.Error
	if errors.As(err, &httpError) {
		code = http.StatusInternalServerError  // Gin errors don't have status codes
		response.Error = httpError.Err.Error() // Use Err field instead of Error()
	}

	// Handle validation errors
	var validationError *ValidationError
	if errors.As(err, &validationError) {
		code = http.StatusBadRequest
		response.Error = "Validation Error"
		response.Details = validationError.Details
	}

	// Handle specific error types
	switch {
	case errors.Is(err, ErrUserNotFound):
		code = http.StatusNotFound
		response.Error = "User Not Found"
	case errors.Is(err, ErrInvalidCredentials):
		code = http.StatusUnauthorized
		response.Error = "Invalid Credentials"
	}

	// Log error if it's a server error
	if code == http.StatusInternalServerError {
		// Add your logging logic here
		// For example:
		// logger.Error("Internal Server Error", zap.Error(err))
	}

	// Send JSON response
	c.JSON(code, response)
}

// Global error middleware
func GlobalErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var e error
				switch x := err.(type) {
				case string:
					e = errors.New(x)
				case error:
					e = x
				default:
					e = errors.New("unknown panic")
				}

				ErrorHandler(c, e)
				c.Abort()
			}
		}()

		c.Next()
	}
}

// Example predefined errors
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
