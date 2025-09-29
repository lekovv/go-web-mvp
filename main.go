package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/lekovv/go-crud-simple/http"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())
	http.RegisterRoutes(app)
	log.Fatal(app.Listen(":8080"))
}
