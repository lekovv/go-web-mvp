package service

import (
	"github.com/lekovv/go-crud-simple/user/model"
	"github.com/lekovv/go-crud-simple/user/repository"
)

type UserServiceInterface interface {
	CreateUser(payload *model.CreateUserDTO) (*model.User, error)
}

type UserService struct {
	repo repository.UserRepoInterface
}

func NewUserService(repo repository.UserRepoInterface) UserServiceInterface {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(payload *model.CreateUserDTO) (*model.User, error) {
	newUser := &model.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		IsActive:  *payload.IsActive,
	}

	if err := s.repo.CreateUser(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}
