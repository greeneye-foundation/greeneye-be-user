package utils

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateStruct validates a struct using validator tags
func ValidateStruct(obj interface{}) error {
	return validate.Struct(obj)
}