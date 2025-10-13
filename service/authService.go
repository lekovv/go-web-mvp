package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/lekovv/go-web-mvp/models"
	"github.com/lekovv/go-web-mvp/repository"
	"github.com/lekovv/go-web-mvp/utils"
	"gorm.io/gorm"
)

type AuthServiceInterface interface {
	RegisterUser(payload *models.UserRegistrationDTO) (*models.AuthResponse, error)
	Login(payload *models.LoginDTO) (*models.AuthResponse, error)
	Logout(userId uuid.UUID) error
}

type AuthService struct {
	userRepo  repository.UserRepoInterface
	roleRepo  repository.RoleRepoInterface
	jwtSecret string
	jwtExpire int
	db        *gorm.DB
}

func NewAuthService(
	userRepo repository.UserRepoInterface,
	roleRepo repository.RoleRepoInterface,
	jwtSecret string,
	jwtExpire int,
	db *gorm.DB,
) AuthServiceInterface {
	return &AuthService{
		userRepo:  userRepo,
		roleRepo:  roleRepo,
		jwtSecret: jwtSecret,
		jwtExpire: jwtExpire,
		db:        db,
	}
}

func (s *AuthService) RegisterUser(payload *models.UserRegistrationDTO) (*models.AuthResponse, error) {
	existingUser, _ := s.userRepo.GetUserByEmail(payload.Email)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return nil, err
	}

	patientRole, err := s.roleRepo.GetRoleByName("patient")
	if err != nil {
		return nil, err
	}

	birthDate, err := utils.ParseDate(payload.BirthDate)
	if err != nil {
		return nil, err
	}

	var createdUser *models.User

	err = s.db.Transaction(func(tx *gorm.DB) error {
		user := &models.User{
			Email:        payload.Email,
			Gender:       payload.Gender,
			PasswordHash: hashedPassword,
			RoleID:       patientRole.ID,
			FirstName:    payload.FirstName,
			LastName:     payload.LastName,
			MiddleName:   payload.MiddleName,
			IsActive:     true,
			PhoneNumber:  payload.PhoneNumber,
		}

		if err := tx.Create(user).Error; err != nil {
			return err
		}
		createdUser = user

		patient := &models.Patient{
			UserId:    user.ID,
			BirthDate: birthDate,
		}

		if err := tx.Create(patient).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{User: createdUser}, nil
}

func (s *AuthService) Login(payload *models.LoginDTO) (*models.AuthResponse, error) {
	user, err := s.userRepo.GetUserByEmail(payload.Email)
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("user is not active")
	}

	if !utils.CheckPasswordHash(payload.Password, user.PasswordHash) {
		return nil, errors.New("invalid password")
	}

	token, err := utils.GenerateJWT(
		user.ID,
		user.Email,
		user.RoleID,
		s.jwtSecret,
		s.jwtExpire,
	)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{Token: &token}, nil
}

func (s *AuthService) Logout(userId uuid.UUID) error {
	return nil
}
