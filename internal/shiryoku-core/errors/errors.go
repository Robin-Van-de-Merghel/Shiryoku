package errors

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field   string
	Message string
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", v.Field, v.Message)
}

func GetValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "dive":
		return fmt.Sprintf("Invalid item in %s", err.Field())
	case "email":
		return "Must be a valid email"
	case "min":
		return fmt.Sprintf("Must be at least %s characters", err.Param())
	case "max":
		return fmt.Sprintf("Must be at most %s characters", err.Param())
	default:
		return fmt.Sprintf("Validation failed: %s", err.Tag())
	}
}

type NotFoundError struct {
	Resource string
	ID       string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %s not found", e.Resource, e.ID)
}
