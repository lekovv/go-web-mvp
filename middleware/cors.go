package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/lekovv/go-crud-simple/config"
)

func CORS(origin *config.Env) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     origin.FrontendUrl,
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST, PATCH, DELETE",
		AllowCredentials: true,
	})
}
