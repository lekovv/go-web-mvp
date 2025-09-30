package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/lekovv/go-crud-simple/config"
	"github.com/lekovv/go-crud-simple/db"
	"github.com/lekovv/go-crud-simple/http"
)

func main() {
	env, err := config.LoadEnv(".")
	if err != nil {
		log.Fatal("Failed to load config", err)
	}
	db.ConnectDB(&env)

	app := fiber.New()
	app.Use(logger.New())
	http.RegisterRoutes(app)
	log.Fatal(app.Listen(":8080"))
}
