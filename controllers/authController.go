package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lekovv/go-web-mvp/models"
	"github.com/lekovv/go-web-mvp/service"
	"github.com/lekovv/go-web-mvp/utils"
	"gorm.io/gorm"
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
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "fail",
				"message": "User Not Found",
			})
		}
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

func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	var payload *models.LoginDTO

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

	response, err := ctrl.authService.Login(payload)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "fail",
				"message": "User Not Found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "login successful",
		"data":    response,
	})
}

func (ctrl *AuthController) Logout(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*utils.JWTClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "unauthorized",
		})
	}

	err := ctrl.authService.Logout(user.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "logout successful",
	})
}
