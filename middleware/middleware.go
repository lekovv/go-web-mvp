package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/lekovv/go-web-mvp/config"
	AppErrors "github.com/lekovv/go-web-mvp/errors"
	"github.com/lekovv/go-web-mvp/service"
	"github.com/lekovv/go-web-mvp/utils"
)

func CORS(env *config.Env) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     env.FrontendUrl,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PATCH, DELETE, PUT",
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
		LimitReached: func(c *fiber.Ctx) error {
			return AppErrors.NewBadRequestError("Rate limit exceeded")
		},
	})
}

func InjectAuthService(authService service.AuthServiceInterface) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("authService", authService)
		return c.Next()
	}
}

func JWTAuth(env *config.Env) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var token string
		authHeader := c.Get("Authorization")

		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}

		if token == "" {
			return AppErrors.NewUnauthorizedError("Authorization header is required")
		}

		authService, ok := c.Locals("authService").(service.AuthServiceInterface)
		if !ok {
			return AppErrors.NewInternalError("Auth service not available")
		}

		hashedToken := utils.HashToken(token, env.JWTSecret)
		isBlacklisted, err := authService.IsTokenBlacklisted(hashedToken)
		if err != nil {
			return AppErrors.WrapError(err, AppErrors.ErrorTypeInternal, "Error checking token blacklist")
		}

		if isBlacklisted {
			return AppErrors.NewUnauthorizedError("Token is blacklisted")
		}

		claims, err := utils.ValidateJWT(token, env.JWTSecret)
		if err != nil {
			return AppErrors.NewUnauthorizedError("Invalid or expired token")
		}

		c.Locals("user", claims)
		return c.Next()
	}
}
