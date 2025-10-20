package service

import (
	"errors"

	"github.com/google/uuid"
	Error "github.com/lekovv/go-web-mvp/errors"
	"github.com/lekovv/go-web-mvp/models"
	"github.com/lekovv/go-web-mvp/repository"
	"gorm.io/gorm"
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, Error.NewNotFoundError("User not found")
		}
		return nil, Error.WrapError(
			err,
			Error.ErrorTypeInternal,
			"Failed to get user",
		)
	}

	return response, nil
}

func (s *UserService) UpdateUser(id uuid.UUID, payload *models.UpdateUserDTO) (*models.UserResponse, error) {
	_, err := s.userRepo.GetUserById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, Error.NewNotFoundError("User not found")
		}
		return nil, Error.WrapError(
			err,
			Error.ErrorTypeInternal,
			"Failed to check user existence",
		)
	}

	if err := s.userRepo.UpdateUser(id, payload); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, Error.NewNotFoundError("User not found")
		}
		return nil, Error.WrapError(
			err,
			Error.ErrorTypeInternal,
			"Failed to update user",
		)
	}

	return s.GetUserById(id)
}

func (s *UserService) DeleteUserById(id uuid.UUID) error {
	var user models.User

	if err := s.userRepo.DeleteUserById(id, &user); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Error.NewNotFoundError("User not found")
		}
		return Error.WrapError(
			err,
			Error.ErrorTypeInternal,
			"Failed to delete user",
		)
	}

	return nil
}
