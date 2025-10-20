package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	Errors "github.com/lekovv/go-web-mvp/middleware"
	"github.com/lekovv/go-web-mvp/models"
	"github.com/lekovv/go-web-mvp/service"
	"github.com/lekovv/go-web-mvp/utils"
	_ "gorm.io/gorm"
)

type UserController struct {
	userService service.UserServiceInterface
}

func NewUserController(userService service.UserServiceInterface) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (ctrl *UserController) GetUserById(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(*utils.JWTClaims)
	if !ok || claims == nil {
		Errors.ThrowUnauthorizedError("Invalid token claims")
	}

	userID := claims.UserID

	user, err := ctrl.userService.GetUserById(userID)
	if err != nil {
		Errors.ThrowInternalError("Failed to get user: " + err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}

func (ctrl *UserController) UpdateUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		Errors.ThrowBadRequestError("Invalid user ID format")
	}

	var payload *models.UpdateUserDTO
	if err := c.BodyParser(&payload); err != nil {
		Errors.ThrowBadRequestError("Invalid request body: " + err.Error())
	}

	if validationErr := utils.ValidateStruct(payload); validationErr != nil {
		Errors.ThrowError(validationErr)
	}

	updatedUser, err := ctrl.userService.UpdateUser(id, payload)
	if err != nil {
		Errors.ThrowInternalError("Failed to update user: " + err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   updatedUser,
	})
}

func (ctrl *UserController) DeleteUserById(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		Errors.ThrowBadRequestError("Invalid user ID format")
	}

	err = ctrl.userService.DeleteUserById(id)
	if err != nil {
		Errors.ThrowInternalError("Failed to delete user: " + err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   fmt.Sprintf("User with id '%s' deleted successfully", id.String()),
	})
}
