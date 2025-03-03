package binders

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ValidateBindJSONRequest validates a request and groups common validation errors.
func ValidateBindJSONRequest(c *gin.Context, req interface{}) bool {
	// Bind the JSON request body to the struct
	if err := c.ShouldBindJSON(&req); err != nil {
		// Check if the error is a validation error
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			// Map to hold field names and corresponding error messages
			validationMap := make(map[string][]string)
			// Map validation tags to error messages
			validationMessages := map[string]string{
				"required": "%s' is required",
				"min":      "%s' must have at least %s characters",
				"max":      "%s' must have no more than %s characters",
				"email":    "%s' must be a valid email address",
				"url":      "%s' must be a valid URL",
				"len":      "%s' must be exactly %s characters",
			}

			// Helper function to get the JSON tag name for the struct field
			getJSONFieldName := func(err validator.FieldError, req interface{}) string {
				t := reflect.TypeOf(req)
				field, _ := t.Elem().FieldByName(err.StructField())
				jsonTag := field.Tag.Get("json")
				return jsonTag
			}

			// Loop through the validation errors
			for _, fieldError := range validationErrors {
				// Get the JSON tag name for the field
				fieldName := getJSONFieldName(fieldError, req)

				// Format the validation error message
				message := formatValidationError(fieldError, validationMessages)

				// Append the error message to the map for the respective field
				validationMap[fieldName] = append(validationMap[fieldName], message)
			}

			// Return grouped validation errors as JSON
			c.JSON(http.StatusBadRequest, gin.H{"validation_errors": validationMap})
			return false
		}

		// Non-validation error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return false
	}

	// Validation passed
	return true
}

// Helper function to format validation error messages
func formatValidationError(err validator.FieldError, messages map[string]string) string {
	// Get the custom message based on the validation tag
	message := messages[err.Tag()]
	if message == "" {
		message = "Field '%s' is invalid" // Default message
	}

	// Format the message based on the parameter (e.g., min length)
	if err.Param() != "" {
		return fmt.Sprintf(message, err.Field(), err.Param())
	}
	return fmt.Sprintf(message, err.Field())
}
