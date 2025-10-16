package service

import (
	"github.com/google/uuid"
	"github.com/lekovv/go-web-mvp/models"
	"github.com/lekovv/go-web-mvp/repository"
)

type UserServiceInterface interface {
	GetUserById(id uuid.UUID) (*models.UserResponse, error)
	UpdateUser(id uuid.UUID, payload *models.UpdateUserDTO) (*models.UserResponse, error)
	DeleteUserById(id uuid.UUID) error
}

type UserService struct {
	userRepo repository.UserRepoInterface
}

func NewUserService(userRepo repository.UserRepoInterface) UserServiceInterface {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUserById(id uuid.UUID) (*models.UserResponse, error) {
	response, err := s.userRepo.GetUserById(id)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *UserService) UpdateUser(id uuid.UUID, payload *models.UpdateUserDTO) (*models.UserResponse, error) {
	if err := s.userRepo.UpdateUser(id, payload); err != nil {
		return nil, err
	}

	return s.GetUserById(id)
}

func (s *UserService) DeleteUserById(id uuid.UUID) error {
	var user models.User

	if err := s.userRepo.DeleteUserById(id, &user); err != nil {
		return err
	}

	return nil
}
