package layers

import (
	"github.com/lekovv/go-web-mvp/config"
	"github.com/lekovv/go-web-mvp/controllers"
	"github.com/lekovv/go-web-mvp/repository"
	"github.com/lekovv/go-web-mvp/service"
	"gorm.io/gorm"
)

type AppContainer struct {
	UserController *controllers.UserController
	AuthController *controllers.AuthController
}

func NewAppContainer(db *gorm.DB, config *config.Env) *AppContainer {
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)

	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, roleRepo, config.JWTSecret, config.JWTExpire, db)

	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(authService)

	return &AppContainer{
		UserController: userController,
		AuthController: authController,
	}
}
