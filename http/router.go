package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lekovv/go-crud-simple/controllers"
)

func RegisterRoutes(app *fiber.App) {
	userController := controllers.NewUserController()

	api := app.Group("/api")
	userRoutes := api.Group("/user")

	userRoutes.Post("/create-user", userController.CreateUser)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "ok",
			"message": "A simple CRUD project on PostgreSQL using Golang REST API",
		})
	})
}
