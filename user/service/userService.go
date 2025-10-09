package service

import (
	"github.com/google/uuid"
	"github.com/lekovv/go-web-mvp/user/model"
	"github.com/lekovv/go-web-mvp/user/repository"
)

type UserServiceInterface interface {
	CreateUser(payload *model.CreateUserDTO) (*model.User, error)
	GetUserById(id uuid.UUID) (*model.User, error)
	UpdateUser(id uuid.UUID, payload *model.UpdateUserDTO) (*model.User, error)
	DeleteUser(id uuid.UUID) error
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

func (s *UserService) GetUserById(id uuid.UUID) (*model.User, error) {
	var user model.User

	if err := s.repo.GetUserById(id, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UpdateUser(id uuid.UUID, payload *model.UpdateUserDTO) (*model.User, error) {
	if err := s.repo.UpdateUser(id, payload); err != nil {
		return nil, err
	}

	return s.GetUserById(id)
}

func (s *UserService) DeleteUser(id uuid.UUID) error {
	var user model.User

	if err := s.repo.DeleteUser(id, &user); err != nil {
		return err
	}

	return nil
}
