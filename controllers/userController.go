package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	AppErrors "github.com/lekovv/go-web-mvp/errors"
	"github.com/lekovv/go-web-mvp/models"
	"github.com/lekovv/go-web-mvp/service"
	"github.com/lekovv/go-web-mvp/utils"
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
		return AppErrors.NewUnauthorizedError("Invalid token claims")
	}

	userID := claims.UserID

	user, err := ctrl.userService.GetUserById(userID)
	if err != nil {
		return err
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
		return AppErrors.NewBadRequestError("Invalid user ID format")
	}

	var payload *models.UpdateUserDTO
	if err := c.BodyParser(&payload); err != nil {
		return AppErrors.NewBadRequestError("Invalid request body: " + err.Error())
	}

	if validationErr := utils.ValidateStruct(payload); validationErr != nil {
		return validationErr
	}

	updatedUser, err := ctrl.userService.UpdateUser(id, payload)
	if err != nil {
		return err
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
		return AppErrors.NewBadRequestError("Invalid user ID format")
	}

	err = ctrl.userService.DeleteUserById(id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User deleted successfully",
	})
}
