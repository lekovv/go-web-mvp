package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/lekovv/go-crud-simple/config"
	"github.com/lekovv/go-crud-simple/db"
	"github.com/lekovv/go-crud-simple/http"
	"github.com/lekovv/go-crud-simple/layers"
	"github.com/lekovv/go-crud-simple/middleware"
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

	appContainer := layers.NewAppContainer(database.DB)
	http.RegisterRoutes(app, appContainer)
	log.Fatal(app.Listen(":" + env.ServerPort))
}
