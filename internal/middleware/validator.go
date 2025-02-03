package middleware

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateRequest(payload interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a new instance of the payload for each request
		payloadValue := reflect.New(reflect.TypeOf(payload).Elem()).Interface()

		// Bind JSON body to the payload
		if err := c.ShouldBindJSON(payloadValue); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request body",
			})
			c.Abort()
			return
		}

		// Validate the payload
		if err := validate.Struct(payloadValue); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			errors := make(map[string]string)

			for _, e := range validationErrors {
				errors[e.Field()] = e.Tag()
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": errors,
			})
			c.Abort()
			return
		}

		// Set the validated payload in the context
		c.Set("validated", payloadValue)
		c.Next()
	}
}

// Helper function to get validated payload from context
func GetValidatedPayload[T any](c *gin.Context) (T, bool) {
	var result T
	payload, exists := c.Get("validated")
	if !exists {
		return result, false
	}

	result, ok := payload.(T)
	return result, ok
}
