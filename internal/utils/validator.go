package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidationError represents a single validation error for a field.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidateStruct validates a struct and returns a slice of ValidationError if validation fails.
// Returns nil if validation passes.
func ValidateStruct(s interface{}) []ValidationError {
	var errors []ValidationError
	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			// Customize the error message based on the tag
			var msg string
			switch err.Tag() {
			case "required":
				msg = fmt.Sprintf("The '%s' field is required.", err.Field())
			case "email":
				msg = fmt.Sprintf("The '%s' field must be a valid email address.", err.Field())
			case "gt":
				msg = fmt.Sprintf("The '%s' field must be greater than %s.", err.Field(), err.Param())
			default:
				msg = fmt.Sprintf("The '%s' field failed on the '%s' validation.", err.Field(), err.Tag())
			}
			errors = append(errors, ValidationError{
				Field:   err.Field(),
				Message: msg,
			})
		}
	}
	if len(errors) > 0 {
		return errors
	}
	return nil
}