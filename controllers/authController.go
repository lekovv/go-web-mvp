package controllers

import (
	"github.com/gofiber/fiber/v2"
	AppErrors "github.com/lekovv/go-web-mvp/errors"
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
		return AppErrors.NewBadRequestError("Invalid request body: " + err.Error())
	}

	if validationErr := utils.ValidateStruct(payload); validationErr != nil {
		return validationErr
	}

	response, err := ctrl.authService.RegisterPatient(c.Context(), payload)
	if err != nil {
		return err
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

func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	var payload *models.LoginDTO

	if err := c.BodyParser(&payload); err != nil {
		return AppErrors.NewBadRequestError("Invalid request body: " + err.Error())
	}

	if validationErr := utils.ValidateStruct(payload); validationErr != nil {
		return validationErr
	}

	response, err := ctrl.authService.Login(c.Context(), payload)
	if err != nil {
		return err
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
		return AppErrors.NewUnauthorizedError("Invalid token claims")
	}

	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return AppErrors.NewUnauthorizedError("Authorization header is required")
	}

	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return AppErrors.NewUnauthorizedError("Invalid authorization header format")
	}

	token := authHeader[7:]

	err := ctrl.authService.Logout(c.Context(), token, user.UserID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Logout successful",
	})
}
