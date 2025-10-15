package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lekovv/go-web-mvp/config"
	"github.com/lekovv/go-web-mvp/controllers"
	"github.com/lekovv/go-web-mvp/layers"
	"github.com/lekovv/go-web-mvp/middleware"
	"gorm.io/gorm"
)

func RegisterRoutes(
	app *fiber.App,
	appContainer *layers.AppContainer,
	env *config.Env,
) {
	api := app.Group("/api")

	setupAuthRoutes(api, appContainer.AuthController, env, appContainer.DB)
	setupUserRoutes(api, appContainer.UserController, env, appContainer.DB)

	app.Get("/health", healthHandler)
}

func setupAuthRoutes(
	api fiber.Router,
	controller *controllers.AuthController,
	env *config.Env,
	db *gorm.DB,
) {
	authRoutes := api.Group("/auth")
	authRoutes.Post("/registration", controller.RegisterUser)
	authRoutes.Post("/login", controller.Login)
	authRoutes.Post("/logout", middleware.JWTAuth(env, db), controller.Logout)
}

func setupUserRoutes(
	api fiber.Router,
	controller *controllers.UserController,
	env *config.Env,
	db *gorm.DB,
) {
	userRoutes := api.Group("/user")
	userRoutes.Use(middleware.JWTAuth(env, db))
	userRoutes.Get("/get-user-by-id", controller.GetUserById)
	userRoutes.Patch("/update-user/:id", controller.UpdateUser)
	userRoutes.Delete("/delete-user/:id", controller.DeleteUserById)
}

func healthHandler(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"status":  "ok",
		"message": "GO web-app on PostgreSQL using Golang REST API",
	})
}
