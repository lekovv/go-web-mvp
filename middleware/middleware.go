package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
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

func Recover() fiber.Handler {
	return recover.New(recover.Config{
		EnableStackTrace: true,
	})
}

func RateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
	})
}
