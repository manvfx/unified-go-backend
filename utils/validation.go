package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct validates a struct based on the tags provided
func ValidateStruct(obj interface{}) error {
	return validate.Struct(obj)
}

// FormatValidationError formats the validation errors
func FormatValidationError(err error) map[string]string {
	errors := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, validationError := range validationErrors {
			fieldName := strings.ToLower(validationError.Field())
			errors[fieldName] = validationError.Tag()
		}
	}
	return errors
}
