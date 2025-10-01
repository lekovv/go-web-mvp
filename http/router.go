package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lekovv/go-crud-simple/controllers"
	"github.com/lekovv/go-crud-simple/layers"
)

func RegisterRoutes(app *fiber.App, appContainer *layers.AppContainer) {
	api := app.Group("/api")

	setupUserRoutes(api, appContainer.UserController)

	app.Get("/health", healthHandler)
}

func setupUserRoutes(api fiber.Router, controller *controllers.UserController) {
	userRoutes := api.Group("/user")
	userRoutes.Post("/create-user", controller.CreateUser)
	userRoutes.Get("/get-user-by-id", controller.GetUserById)
}

func healthHandler(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"status":  "ok",
		"message": "A simple CRUD project on PostgreSQL using Golang REST API",
	})
}
