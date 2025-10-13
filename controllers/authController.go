package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lekovv/go-web-mvp/models"
	"github.com/lekovv/go-web-mvp/service"
	"github.com/lekovv/go-web-mvp/utils"
)

type AuthController struct {
	authService service.AuthServiceInterface
}

func NewAuthController(authService service.AuthServiceInterface) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (ctrl *AuthController) RegisterUser(c *fiber.Ctx) error {
	var payload *models.UserRegistrationDTO

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	response, err := ctrl.authService.RegisterUser(payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "registration successful",
		"data":    response,
	})
}
