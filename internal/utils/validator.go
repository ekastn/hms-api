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
			case "lte":
				msg = fmt.Sprintf("The '%s' field must be less than or equal to %s.", err.Field(), err.Param())
			case "min":
				msg = fmt.Sprintf("The '%s' field must be at least %s characters long.", err.Field(), err.Param())
			case "max":
				msg = fmt.Sprintf("The '%s' field must be at most %s characters long.", err.Field(), err.Param())
			case "oneof":
				msg = fmt.Sprintf("The '%s' field must be one of: %s.", err.Field(), err.Param())
			case "datetime":
				msg = fmt.Sprintf("The '%s' field must be a valid date and time.", err.Field())
			case "e164":
				msg = fmt.Sprintf("The '%s' field must be a valid E.164 formatted phone number.", err.Field())
			case "mongodb":
				msg = fmt.Sprintf("The '%s' field must be a valid MongoDB ObjectID.", err.Field())
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
