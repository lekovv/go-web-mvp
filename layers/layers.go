package layers

import (
	"github.com/lekovv/go-web-mvp/config"
	"github.com/lekovv/go-web-mvp/controllers"
	"github.com/lekovv/go-web-mvp/repository"
	"github.com/lekovv/go-web-mvp/service"
	"gorm.io/gorm"
)

type AppContainer struct {
	DB             *gorm.DB
	UserController *controllers.UserController
	AuthController *controllers.AuthController
	AuthService    service.AuthServiceInterface
}

func NewAppContainer(db *gorm.DB, env *config.Env) *AppContainer {
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	specializationRepo := repository.NewSpecializationRepository(db)
	authRepo := repository.NewAuthRepository(db)

	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(
		specializationRepo,
		userRepo,
		roleRepo,
		authRepo,
		env,
	)

	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(authService)

	return &AppContainer{
		DB:             db,
		UserController: userController,
		AuthController: authController,
		AuthService:    authService,
	}
}
