package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lekovv/go-crud-simple/user/model"
	"github.com/lekovv/go-crud-simple/user/service"
	"github.com/lekovv/go-crud-simple/utils"
)

type UserController struct {
	service service.UserServiceInterface
}

func NewUserController(service service.UserServiceInterface) *UserController {
	return &UserController{
		service: service,
	}
}

func (ctrl *UserController) CreateUser(c *fiber.Ctx) error {
	var payload *model.CreateUserDTO

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail", "message": err.Error(),
		})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	newUser, err := ctrl.service.CreateUser(payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error", "message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success", "data": fiber.Map{"value": newUser},
	})
}
