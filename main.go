package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/lekovv/go-web-mvp/config"
	"github.com/lekovv/go-web-mvp/db"
	"github.com/lekovv/go-web-mvp/http"
	"github.com/lekovv/go-web-mvp/layers"
	"github.com/lekovv/go-web-mvp/middleware"
)

func main() {
	env, err := config.LoadEnv(".")
	if err != nil {
		log.Fatal("Failed to load config", err)
	}
	database := db.ConnectDB(&env)

	app := fiber.New()
	app.Use(logger.New())
	app.Use(middleware.Recover())
	app.Use(middleware.CORS(&env))
	app.Use(middleware.RateLimiter())

	appContainer := layers.NewAppContainer(database.DB, &env)
	http.RegisterRoutes(app, appContainer)
	log.Fatal(app.Listen(":" + env.ServerPort))
}
