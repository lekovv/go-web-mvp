package layers

import (
	"github.com/lekovv/go-web-mvp/controllers"
	"github.com/lekovv/go-web-mvp/repository"
	"github.com/lekovv/go-web-mvp/service"
	"gorm.io/gorm"
)

type AppContainer struct {
	UserController *controllers.UserController
}

func NewAppContainer(db *gorm.DB) *AppContainer {
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	return &AppContainer{
		UserController: userController,
	}
}
