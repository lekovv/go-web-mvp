package utils

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/lekovv/go-web-mvp/errors"
)

var validate = validator.New()

type ValidationErrorDetail struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value,omitempty"`
	Message string `json:"message"`
}

func ValidateStruct[T any](payload T) *errors.AppError {
	err := validate.Struct(payload)
	if err != nil {
		var details []ValidationErrorDetail
		var messages []string

		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return errors.NewValidationError("Invalid validation data", []string{err.Error()})
		}

		for _, err := range validationErrors {
			fieldName := getFieldName(err.StructNamespace())
			message := getValidationMessage(err)

			detail := ValidationErrorDetail{
				Field:   fieldName,
				Tag:     err.Tag(),
				Value:   err.Param(),
				Message: message,
			}
			details = append(details, detail)
			messages = append(messages, fmt.Sprintf("%s: %s", fieldName, message))
		}

		return &errors.AppError{
			Type:       errors.ErrorTypeValidation,
			Message:    "Validation failed",
			StatusCode: http.StatusBadRequest,
			Details:    details,
		}
	}
	return nil
}

func getFieldName(namespace string) string {
	parts := strings.Split(namespace, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return namespace
}

func getValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "Field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Minimum length is %s characters", err.Param())
	case "max":
		return fmt.Sprintf("Maximum length is %s characters", err.Param())
	case "len":
		return fmt.Sprintf("Length must be exactly %s characters", err.Param())
	case "numeric":
		return "Must be a number"
	case "alpha":
		return "Must contain only letters"
	case "alphanum":
		return "Must contain only letters and numbers"
	case "uuid":
		return "Invalid UUID format"
	case "oneof":
		return fmt.Sprintf("Must be one of: %s", err.Param())
	default:
		return fmt.Sprintf("Validation failed for tag '%s'", err.Tag())
	}
}
