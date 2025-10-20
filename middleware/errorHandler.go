package middleware

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	Error "github.com/lekovv/go-web-mvp/errors"
	"gorm.io/gorm"
)

func ErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				var appErr *Error.AppError

				switch err := r.(type) {
				case *Error.AppError:
					appErr = err
				case error:
					appErr = Error.WrapError(
						err,
						Error.ErrorTypeInternal,
						"Internal server error")
				default:
					appErr = Error.NewInternalError("Unknown error occurred")
				}

				handleAppError(c, appErr)
			}
		}()

		err := c.Next()
		if err != nil {
			handleError(c, err)
		}

		return nil
	}
}

func handleError(c *fiber.Ctx, err error) {
	if appErr, ok := err.(*Error.AppError); ok {
		handleAppError(c, appErr)
		return
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		appErr := Error.NewNotFoundError("Record not found")
		handleAppError(c, appErr)
		return
	}

	if fiberErr, ok := err.(*fiber.Error); ok {
		var appErr *Error.AppError
		switch fiberErr.Code {
		case fiber.StatusNotFound:
			appErr = Error.NewNotFoundError(fiberErr.Message)
		case fiber.StatusUnauthorized:
			appErr = Error.NewUnauthorizedError(fiberErr.Message)
		case fiber.StatusForbidden:
			appErr = Error.NewForbiddenError(fiberErr.Message)
		case fiber.StatusBadRequest:
			appErr = Error.NewBadRequestError(fiberErr.Message)
		default:
			appErr = Error.WrapError(fiberErr, Error.ErrorTypeInternal, fiberErr.Message)
		}
		handleAppError(c, appErr)
		return
	}

	appErr := Error.WrapError(
		err,
		Error.ErrorTypeInternal,
		"Internal server error",
	)
	handleAppError(c, appErr)
}

func handleAppError(c *fiber.Ctx, appErr *Error.AppError) {
	if appErr.Type == Error.ErrorTypeInternal {
		log.Printf("Internal error: %s", appErr.Message)
	}

	response := Error.ErrorResponse{
		Status:  "error",
		Error:   appErr.Type,
		Message: appErr.Message,
	}

	if appErr.Type == Error.ErrorTypeValidation && appErr.Details != nil {
		response.Details = appErr.Details
	}

	_ = c.Status(appErr.StatusCode).JSON(response)
}

func ThrowError(appErr *Error.AppError) {
	panic(appErr)
}

func ThrowNotFoundError(message string) {
	panic(Error.NewNotFoundError(message))
}

func ThrowUnauthorizedError(message string) {
	panic(Error.NewUnauthorizedError(message))
}

func ThrowForbiddenError(message string) {
	panic(Error.NewForbiddenError(message))
}

func ThrowConflictError(message string) {
	panic(Error.NewConflictError(message))
}

func ThrowInternalError(message string) {
	panic(Error.NewInternalError(message))
}

func ThrowBadRequestError(message string) {
	panic(Error.NewBadRequestError(message))
}

func ThrowTooManyRequestsError(message string) {
	if message == "" {
		message = "Too many requests"
	}
	panic(&Error.AppError{
		Type:       "too_many_requests",
		Message:    message,
		StatusCode: 429,
	})
}
