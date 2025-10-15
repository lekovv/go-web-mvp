package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/lekovv/go-web-mvp/config"
	"github.com/lekovv/go-web-mvp/utils"
)

func CORS(env *config.Env) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     env.FrontendUrl,
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

func JWTAuth(env *config.Env) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var token string

		authHeader := c.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}

		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "fail",
				"message": "authorization header is required",
			})
		}

		//hashedToken := utils.HashToken(token, env.JWTSecret)
		//
		//isBlacklisted, err := authRepo.IsTokenBlacklisted(hashedToken)
		//if err != nil {
		//	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		//		"status":  "fail",
		//		"message": "error checking token blacklist",
		//	})
		//}
		//
		//if isBlacklisted {
		//	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		//		"status":  "fail",
		//		"message": "token is blacklisted",
		//	})
		//}

		claims, err := utils.ValidateJWT(token, env.JWTSecret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "fail",
				"message": "invalid or expired token",
			})
		}

		c.Locals("user", claims)

		return c.Next()
	}
}
