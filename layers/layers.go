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

func NewAppContainer(db *gorm.DB, env *config.Env) *AppContainer {
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	authRepo := repository.NewAuthRepository(db)

	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(
		userRepo,
		roleRepo,
		authRepo,
		env,
		db,
	)

	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(authService)

	return &AppContainer{
		UserController: userController,
		AuthController: authController,
	}
}
