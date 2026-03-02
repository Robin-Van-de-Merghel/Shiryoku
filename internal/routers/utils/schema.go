package utils

import (
	"reflect"
	"strings"
)

// GenerateSchema creates a schema map from a struct using reflection
// It reads json tags and builds a human-readable schema
func GenerateSchema(data any) map[string]string {
	schema := make(map[string]string)
	t := reflect.TypeOf(data)

	// Handle pointers
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Iterate over struct fields
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		validateTag := field.Tag.Get("validate")

		// Extract field name from json tag
		fieldName := strings.Split(jsonTag, ",")[0]
		if fieldName == "" || fieldName == "-" {
			fieldName = field.Name
		}

		// Build type description
		typeStr := field.Type.String()
		if strings.Contains(typeStr, "[]") {
			typeStr = "array"
		}

		// Check if required
		required := "optional"
		if strings.Contains(validateTag, "required") {
			required = "required"
		}

		schema[fieldName] = typeStr + " - " + required
	}

	return schema
}
