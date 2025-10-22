package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lekovv/go-web-mvp/config"
	"github.com/lekovv/go-web-mvp/controllers"
	"github.com/lekovv/go-web-mvp/layers"
	"github.com/lekovv/go-web-mvp/middleware"
)

func RegisterRoutes(
	app *fiber.App,
	appContainer *layers.AppContainer,
	env *config.Env,
) {
	api := app.Group("/api")

	api.Use(middleware.InjectAuthService(appContainer.AuthService))

	setupAuthRoutes(api, appContainer.AuthController, env)
	setupUserRoutes(api, appContainer.UserController, env)
	setupAdminRoutes(api, appContainer.AdminController, env)

	app.Get("/health", healthHandler)
}

func setupAdminRoutes(
	api fiber.Router,
	controller *controllers.AdminController,
	env *config.Env,
) {
	adminRoutes := api.Group("/admin")
	adminRoutes.Use(middleware.JWTAuth(env))
	adminRoutes.Use(middleware.RequireRole("admin"))

	adminRoutes.Post("/create-doctor", controller.CreateDoctor)
	adminRoutes.Patch("/update-user/:id", controller.UpdateUser)
	adminRoutes.Delete("/delete-user/:id", controller.DeleteUserById)
}

func setupAuthRoutes(
	api fiber.Router,
	controller *controllers.AuthController,
	env *config.Env,
) {
	authRoutes := api.Group("/auth")
	authRoutes.Post("/registration", controller.RegisterPatient)
	authRoutes.Post("/login", controller.Login)
	authRoutes.Post("/logout", middleware.JWTAuth(env), controller.Logout)
}

func setupUserRoutes(
	api fiber.Router,
	controller *controllers.UserController,
	env *config.Env,
) {
	userRoutes := api.Group("/user")
	userRoutes.Use(middleware.JWTAuth(env))

	userRoutes.Get("/get-user", controller.GetUserById)
}

func healthHandler(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"status":  "ok",
		"message": "GO web-app on PostgreSQL using Golang REST API",
	})
}
