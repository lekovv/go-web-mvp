package service

import (
	"github.com/google/uuid"
	"github.com/lekovv/go-web-mvp/models"
	"github.com/lekovv/go-web-mvp/repository"
)

type UserServiceInterface interface {
	CreateUser(payload *models.CreateUserDTO) (*models.User, error)
	GetUserById(id uuid.UUID) (*models.User, error)
	UpdateUser(id uuid.UUID, payload *models.UpdateUserDTO) (*models.User, error)
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

func (s *UserService) CreateUser(payload *models.CreateUserDTO) (*models.User, error) {
	newUser := &models.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		IsActive:  *payload.IsActive,
	}

	if err := s.repo.CreateUser(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *UserService) GetUserById(id uuid.UUID) (*models.User, error) {
	var user models.User

	if err := s.repo.GetUserById(id, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UpdateUser(id uuid.UUID, payload *models.UpdateUserDTO) (*models.User, error) {
	if err := s.repo.UpdateUser(id, payload); err != nil {
		return nil, err
	}

	return s.GetUserById(id)
}

func (s *UserService) DeleteUser(id uuid.UUID) error {
	var user models.User

	if err := s.repo.DeleteUser(id, &user); err != nil {
		return err
	}

	return nil
}
