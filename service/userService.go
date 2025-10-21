package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	AppErrors "github.com/lekovv/go-web-mvp/errors"
	"github.com/lekovv/go-web-mvp/models"
	"github.com/lekovv/go-web-mvp/repository"
	"gorm.io/gorm"
)

type UserServiceInterface interface {
	GetUserById(ctx context.Context, id uuid.UUID) (*models.UserResponse, error)
	UpdateUser(ctx context.Context, id uuid.UUID, payload *models.UpdateUserDTO) (*models.UserResponse, error)
	DeleteUserById(ctx context.Context, id uuid.UUID) error
}

type UserService struct {
	userRepo repository.UserRepoInterface
}

func NewUserService(userRepo repository.UserRepoInterface) UserServiceInterface {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUserById(ctx context.Context, id uuid.UUID) (*models.UserResponse, error) {
	user, err := s.userRepo.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, AppErrors.NewNotFoundError("User not found")
		}
		return nil, AppErrors.WrapError(err, AppErrors.ErrorTypeInternal, "Failed to get user")
	}

	resp := &models.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		Gender:      user.Gender,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		MiddleName:  user.MiddleName,
		IsActive:    user.IsActive,
		PhoneNumber: user.PhoneNumber,
	}

	if user.Patient != nil {
		birthDate := user.Patient.BirthDate.Format("2006-01-02")
		resp.BirthDate = &birthDate
	}

	if user.Doctor != nil {
		resp.Specialization = &user.Doctor.Specialization.Name
		resp.Bio = user.Doctor.Bio
		resp.ExperienceYears = user.Doctor.ExperienceYears
		resp.Price = &user.Doctor.Price
	}

	return resp, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id uuid.UUID, payload *models.UpdateUserDTO) (*models.UserResponse, error) {
	if err := s.userRepo.UpdateUser(ctx, id, payload); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, AppErrors.NewNotFoundError("User not found")
		}
		return nil, AppErrors.WrapError(
			err,
			AppErrors.ErrorTypeInternal,
			"Failed to update user",
		)
	}

	return s.GetUserById(ctx, id)
}

func (s *UserService) DeleteUserById(ctx context.Context, id uuid.UUID) error {
	var user models.User

	if err := s.userRepo.DeleteUserById(ctx, id, &user); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return AppErrors.NewNotFoundError("User not found")
		}
		return AppErrors.WrapError(
			err,
			AppErrors.ErrorTypeInternal,
			"Failed to delete user",
		)
	}

	return nil
}
