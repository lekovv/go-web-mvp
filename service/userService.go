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
	user, err := s.userRepo.GetUserById(id)
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

	switch user.Role.Name {
	case "patient":
		patient, err := s.userRepo.GetPatientByUserId(user.ID)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, Error.WrapError(
					err,
					Error.ErrorTypeInternal,
					"Failed to get patient data",
				)
			}
		} else {
			birthDate := patient.BirthDate.Format("2006-01-02")
			resp.BirthDate = &birthDate
		}

	case "doctor":
		doctor, err := s.userRepo.GetDoctorByUserId(user.ID)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, Error.WrapError(
					err,
					Error.ErrorTypeInternal,
					"Failed to get doctor data",
				)
			}
		} else {
			resp.Specialization = &doctor.Specialization.Name
			resp.Bio = doctor.Bio
			resp.ExperienceYears = doctor.ExperienceYears
			resp.Price = &doctor.Price
		}

	default:
	}

	return resp, nil
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
