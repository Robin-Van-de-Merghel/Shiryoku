package utils

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// FieldError represents a validation error for a single field
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrorResponse is your 422 response structure
type ValidationErrorResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Errors  []FieldError      `json:"errors"`
	Schema  map[string]string `json:"schema,omitempty"`
}

// ValidateAndRespond validates a struct and returns 422 if invalid
// Call this at the START of your handler
func ValidateAndRespond(c *gin.Context, data any, schema map[string]string) bool {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("json")
		if name == "" || name == "-" {
			return fld.Name
		}
		return name
	})

	if err := v.Struct(data); err != nil {
		var fieldErrors []FieldError
		var validationErrors validator.ValidationErrors

		if errors.As(err, &validationErrors) {
			for _, vErr := range validationErrors {
				fieldErrors = append(fieldErrors, FieldError{
					Field:   vErr.Field(),
					Message: getValidationMessage(vErr),
				})
			}
		}

		response := ValidationErrorResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: "Request validation failed",
			Errors:  fieldErrors,
			Schema:  schema,
		}

		c.JSON(http.StatusUnprocessableEntity, response)
		return false // Stop handler execution
	}

	return true // Continue
}

func getValidationMessage(err validator.FieldError) string {
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

// ParseJSONError handles JSON unmarshaling errors (400)
func ParseJSONError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"code":    http.StatusBadRequest,
		"message": "Invalid JSON format",
		"error":   err.Error(),
	})
}
