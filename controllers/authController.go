package controllers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	Errors "github.com/lekovv/go-web-mvp/middleware"
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

func (ctrl *AuthController) RegisterPatient(c *fiber.Ctx) error {
	var payload *models.PatientRegistrationDTO

	if err := c.BodyParser(&payload); err != nil {
		Errors.ThrowBadRequestError("Invalid request body: " + err.Error())
	}

	if validationErr := utils.ValidateStruct(payload); validationErr != nil {
		Errors.ThrowError(validationErr)
	}

	response, err := ctrl.authService.RegisterPatient(payload)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			Errors.ThrowConflictError("User already exists")
		}
		Errors.ThrowInternalError("Failed to register patient: " + err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Patient registration successful",
		"data":    response,
	})
}

func (ctrl *AuthController) CreateDoctor(c *fiber.Ctx) error {
	var payload *models.DoctorRegistrationDTO

	if err := c.BodyParser(&payload); err != nil {
		Errors.ThrowBadRequestError("Invalid request body: " + err.Error())
	}

	if validationErr := utils.ValidateStruct(payload); validationErr != nil {
		Errors.ThrowError(validationErr)
	}

	response, err := ctrl.authService.CreateDoctor(payload)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			Errors.ThrowConflictError("User already exists")
		}
		if strings.Contains(err.Error(), "specialization") {
			Errors.ThrowNotFoundError("Specialization not found")
		}
		Errors.ThrowInternalError("Failed to create doctor: " + err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Doctor registration successful",
		"data":    response,
	})
}

func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	var payload *models.LoginDTO

	if err := c.BodyParser(&payload); err != nil {
		Errors.ThrowBadRequestError("Invalid request body: " + err.Error())
	}

	if validationErr := utils.ValidateStruct(payload); validationErr != nil {
		Errors.ThrowError(validationErr)
	}

	response, err := ctrl.authService.Login(payload)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			Errors.ThrowNotFoundError("User not found")
		}
		if strings.Contains(err.Error(), "invalid password") {
			Errors.ThrowUnauthorizedError("Invalid email or password")
		}
		if strings.Contains(err.Error(), "not active") {
			Errors.ThrowForbiddenError("User account is deactivated")
		}
		Errors.ThrowInternalError("Login failed: " + err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Login successful",
		"data":    response,
	})
}

func (ctrl *AuthController) Logout(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*utils.JWTClaims)
	if !ok {
		Errors.ThrowUnauthorizedError("Invalid token claims")
	}

	authHeader := c.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		Errors.ThrowUnauthorizedError("Authorization header is required")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	err := ctrl.authService.Logout(token, user.UserID)
	if err != nil {
		Errors.ThrowInternalError("Failed to logout: " + err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Logout successful",
	})
}
