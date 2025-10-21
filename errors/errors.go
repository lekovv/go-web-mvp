package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type AppError struct {
	Type       ErrorType   `json:"type"`
	Message    string      `json:"message"`
	StatusCode int         `json:"status_code"`
	Details    interface{} `json:"details,omitempty"`
	Err        error       `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

type ErrorType string

const (
	ErrorTypeValidation   ErrorType = "validation_error"
	ErrorTypeNotFound     ErrorType = "not_found"
	ErrorTypeUnauthorized ErrorType = "unauthorized"
	ErrorTypeForbidden    ErrorType = "forbidden"
	ErrorTypeConflict     ErrorType = "conflict"
	ErrorTypeInternal     ErrorType = "internal_error"
	ErrorTypeBadRequest   ErrorType = "bad_request"
)

type ErrorResponse struct {
	Status  string      `json:"status"`
	Error   ErrorType   `json:"error"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func NewValidationError(message string, details interface{}) *AppError {
	return &AppError{
		Type:       ErrorTypeValidation,
		Message:    message,
		StatusCode: http.StatusBadRequest,
		Details:    details,
	}
}

func NewNotFoundError(message string) *AppError {
	if message == "" {
		message = "Resource not found"
	}
	return &AppError{
		Type:       ErrorTypeNotFound,
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

func NewUnauthorizedError(message string) *AppError {
	if message == "" {
		message = "Unauthorized access"
	}
	return &AppError{
		Type:       ErrorTypeUnauthorized,
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

func NewForbiddenError(message string) *AppError {
	if message == "" {
		message = "Access forbidden"
	}
	return &AppError{
		Type:       ErrorTypeForbidden,
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
}

func NewConflictError(message string) *AppError {
	if message == "" {
		message = "Resource conflict"
	}
	return &AppError{
		Type:       ErrorTypeConflict,
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

func NewInternalError(message string) *AppError {
	if message == "" {
		message = "Internal server error"
	}
	return &AppError{
		Type:       ErrorTypeInternal,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}

func NewBadRequestError(message string) *AppError {
	if message == "" {
		message = "Bad request"
	}
	return &AppError{
		Type:       ErrorTypeBadRequest,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

func WrapError(err error, errorType ErrorType, message string) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	var statusCode int
	switch errorType {
	case ErrorTypeValidation, ErrorTypeBadRequest:
		statusCode = http.StatusBadRequest
	case ErrorTypeNotFound:
		statusCode = http.StatusNotFound
	case ErrorTypeUnauthorized:
		statusCode = http.StatusUnauthorized
	case ErrorTypeForbidden:
		statusCode = http.StatusForbidden
	case ErrorTypeConflict:
		statusCode = http.StatusConflict
	default:
		statusCode = http.StatusInternalServerError
		errorType = ErrorTypeInternal
	}

	if message == "" {
		message = err.Error()
	}

	return &AppError{
		Type:       errorType,
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

func (e *AppError) UnwrapError() error {
	return e.Err
}
