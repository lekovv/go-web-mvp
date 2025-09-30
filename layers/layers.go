package layers

import (
	"github.com/lekovv/go-crud-simple/controllers"
	"github.com/lekovv/go-crud-simple/user/repository"
	"github.com/lekovv/go-crud-simple/user/service"
)

type AppContainer struct {
	UserController *controllers.UserController
}

func NewAppContainer() *AppContainer {
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	return &AppContainer{
		UserController: userController,
	}
}
