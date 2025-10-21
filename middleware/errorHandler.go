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
		err := c.Next()

		if err == nil {
			return nil
		}

		return handleError(c, err)
	}
}

func handleError(c *fiber.Ctx, err error) error {
	var appErr *Error.AppError
	if errors.As(err, &appErr) {
		return sendErrorResponse(c, appErr)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		appErr = Error.NewNotFoundError("Record not found")
		return sendErrorResponse(c, appErr)
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		appErr = convertFiberError(fiberErr)
		return sendErrorResponse(c, appErr)
	}

	log.Printf("Unhandled error: %v", err)
	appErr = Error.NewInternalError("Internal server error")
	return sendErrorResponse(c, appErr)
}

func convertFiberError(fiberErr *fiber.Error) *Error.AppError {
	switch fiberErr.Code {
	case fiber.StatusNotFound:
		return Error.NewNotFoundError(fiberErr.Message)
	case fiber.StatusUnauthorized:
		return Error.NewUnauthorizedError(fiberErr.Message)
	case fiber.StatusForbidden:
		return Error.NewForbiddenError(fiberErr.Message)
	case fiber.StatusBadRequest:
		return Error.NewBadRequestError(fiberErr.Message)
	case fiber.StatusConflict:
		return Error.NewConflictError(fiberErr.Message)
	default:
		return Error.WrapError(fiberErr, Error.ErrorTypeInternal, fiberErr.Message)
	}
}

func sendErrorResponse(c *fiber.Ctx, appErr *Error.AppError) error {
	if appErr.Type == Error.ErrorTypeInternal {
		log.Printf("Internal error: %v", appErr)
	}

	response := Error.ErrorResponse{
		Status:  "error",
		Error:   appErr.Type,
		Message: appErr.Message,
	}

	if appErr.Type == Error.ErrorTypeValidation && appErr.Details != nil {
		response.Details = appErr.Details
	}

	return c.Status(appErr.StatusCode).JSON(response)
}
