package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	AppErrors "github.com/lekovv/go-web-mvp/errors"
	"github.com/lekovv/go-web-mvp/models"
	"github.com/lekovv/go-web-mvp/service"
	"github.com/lekovv/go-web-mvp/utils"
)

type AdminController struct {
	authService service.AuthServiceInterface
	userService service.UserServiceInterface
}

func NewAdminController(authService service.AuthServiceInterface, userService service.UserServiceInterface) *AdminController {
	return &AdminController{
		authService: authService,
		userService: userService,
	}
}

func (ctrl *AdminController) CreateDoctor(c *fiber.Ctx) error {
	var payload *models.DoctorRegistrationDTO

	if err := c.BodyParser(&payload); err != nil {
		return AppErrors.NewBadRequestError("Invalid request body: " + err.Error())
	}

	if validationErr := utils.ValidateStruct(payload); validationErr != nil {
		return validationErr
	}

	response, err := ctrl.authService.CreateDoctor(c.Context(), payload)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Doctor registration successful",
		"data":    response,
	})
}

func (ctrl *AdminController) UpdateUser(c *fiber.Ctx) error {
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

	updatedUser, err := ctrl.userService.UpdateUser(c.Context(), id, payload)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   updatedUser,
	})
}

func (ctrl *AdminController) DeleteUserById(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return AppErrors.NewBadRequestError("Invalid user ID format")
	}

	err = ctrl.userService.DeleteUserById(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User deleted successfully",
	})
}
