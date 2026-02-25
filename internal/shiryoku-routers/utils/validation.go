package utils

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	shiryoku_errors "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/errors"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
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
					Message: shiryoku_errors.GetValidationMessage(vErr),
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

// ParseJSONError handles JSON unmarshaling errors (400)
func ParseJSONError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"code":    http.StatusBadRequest,
		"message": "Invalid JSON format",
		"error":   err.Error(),
	})
}

// FieldTypeInfo stores precomputed type info for a field
type FieldTypeInfo struct {
	JSONName string
	GORMName string
	GoType   reflect.Type
	JSONKind reflect.Kind // expected JSON kind (string, float64, bool)
}

// buildFieldTypeMap returns JSON kind info for a struct
func buildFieldTypeMap(s any) map[string]FieldTypeInfo {
	fields := make(map[string]FieldTypeInfo)
	val := reflect.TypeOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return fields
	}

	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)
		if f.PkgPath != "" {
			continue
		}

		jsonTag := f.Tag.Get("json")
		if idx := strings.Index(jsonTag, ","); idx != -1 {
			jsonTag = jsonTag[:idx]
		}
		if jsonTag == "" || jsonTag == "-" {
			jsonTag = f.Name
		}

		gormTag := f.Tag.Get("gorm")
		column := ""
		for _, part := range strings.Split(gormTag, ";") {
			if strings.HasPrefix(part, "column:") {
				column = strings.TrimPrefix(part, "column:")
				break
			}
		}

		kind := f.Type.Kind()
		if kind == reflect.Int || kind == reflect.Int64 || kind == reflect.Float32 || kind == reflect.Float64 {
			kind = reflect.Float64 // JSON numbers are float64
		}

		fields[jsonTag] = FieldTypeInfo{
			JSONName: jsonTag,
			GORMName: column,
			GoType:   f.Type,
			JSONKind: kind,
		}
	}

	return fields
}

func ValidateSearchParamTypesPrecomputed(params *models.SearchParams, allowedMaps ...map[string]FieldTypeInfo) error {
	for _, spec := range params.Search {
		// Scalar
		if spec.Scalar != nil {
			var found FieldTypeInfo
			var ok bool
			for _, m := range allowedMaps {
				if info, exists := m[spec.Scalar.Parameter]; exists {
					found = info
					ok = true
					break
				}
			}
			if !ok {
				return fmt.Errorf("invalid search parameter: %s", spec.Scalar.Parameter)
			}
			if !isJSONValueCompatible(spec.Scalar.Value, found.JSONKind) {
				return fmt.Errorf("value for %s must be %s", spec.Scalar.Parameter, found.JSONKind)
			}
		}

		// Vector
		if spec.Vector != nil {
			var found FieldTypeInfo
			var ok bool
			for _, m := range allowedMaps {
				if info, exists := m[spec.Vector.Parameter]; exists {
					found = info
					ok = true
					break
				}
			}
			if !ok {
				return fmt.Errorf("invalid search parameter: %s", spec.Vector.Parameter)
			}
			values, ok := spec.Vector.Values.([]any)
			if !ok {
				return fmt.Errorf("values for %s must be array", spec.Vector.Parameter)
			}
			for _, v := range values {
				if !isJSONValueCompatible(v, found.JSONKind) {
					return fmt.Errorf("value in %s must be %s", spec.Vector.Parameter, found.JSONKind)
				}
			}
		}
	}
	return nil
}

// Check if JSON input is compatible with expected kind
func isJSONValueCompatible(v any, kind reflect.Kind) bool {
	if v == nil {
		return true
	}
	switch kind {
	case reflect.String:
		_, ok := v.(string)
		return ok
	case reflect.Float64:
		_, isFloat := v.(float64)
		_, isInt := v.(int)
		return isFloat || isInt
	case reflect.Bool:
		_, ok := v.(bool)
		return ok
	default:
		return true
	}
}
