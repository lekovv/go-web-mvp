package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/lekovv/go-web-mvp/config"
	"github.com/lekovv/go-web-mvp/models"
	"github.com/lekovv/go-web-mvp/repository"
	"github.com/lekovv/go-web-mvp/utils"
	"gorm.io/gorm"
)

type AuthServiceInterface interface {
	RegisterPatient(payload *models.PatientRegistrationDTO) (*models.AuthResponse, error)
	CreateDoctor(payload *models.DoctorRegistrationDTO) (*models.AuthResponse, error)
	Login(payload *models.LoginDTO) (*models.AuthResponse, error)
	Logout(token string, userId uuid.UUID) error
	IsTokenBlacklisted(tokenHash string) (bool, error)
	DeleteExpiredTokens() error
}

type AuthService struct {
	specializationRepo repository.SpecializationRepoInterface
	userRepo           repository.UserRepoInterface
	roleRepo           repository.RoleRepoInterface
	authRepo           repository.AuthRepoInterface

	env *config.Env
	db  *gorm.DB
}

func NewAuthService(
	specializationRepo repository.SpecializationRepoInterface,
	userRepo repository.UserRepoInterface,
	roleRepo repository.RoleRepoInterface,
	authRepo repository.AuthRepoInterface,
	env *config.Env,
	db *gorm.DB,
) AuthServiceInterface {
	return &AuthService{
		specializationRepo: specializationRepo,
		userRepo:           userRepo,
		roleRepo:           roleRepo,
		authRepo:           authRepo,
		env:                env,
		db:                 db,
	}
}

func (s *AuthService) RegisterPatient(payload *models.PatientRegistrationDTO) (*models.AuthResponse, error) {
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
			UserID:    user.ID,
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

func (s *AuthService) CreateDoctor(payload *models.DoctorRegistrationDTO) (*models.AuthResponse, error) {
	existingUser, _ := s.userRepo.GetUserByEmail(payload.Email)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return nil, err
	}

	doctorRole, err := s.roleRepo.GetRoleByName("doctor")
	if err != nil {
		return nil, err
	}

	specialization, err := s.specializationRepo.GetSpecializationByName(payload.Specialization)
	if err != nil {
		return nil, err
	}

	var createdUser *models.User

	err = s.db.Transaction(func(tx *gorm.DB) error {
		user := &models.User{
			Email:        payload.Email,
			Gender:       payload.Gender,
			PasswordHash: hashedPassword,
			RoleID:       doctorRole.ID,
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

		doctor := &models.Doctor{
			UserID:           user.ID,
			SpecializationID: specialization.ID,
			Bio:              payload.Bio,
			ExperienceYears:  payload.ExperienceYears,
			Price:            payload.Price,
		}

		if err := tx.Create(doctor).Error; err != nil {
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
		s.env.JWTSecret,
		s.env.JWTExpire,
	)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{Token: &token}, nil
}

func (s *AuthService) Logout(token string, userId uuid.UUID) error {
	claims, err := utils.ValidateJWT(token, s.env.JWTSecret)
	if err != nil {
		return err
	}

	hashedToken := utils.HashToken(token, s.env.JWTSecret)

	blacklistToken := &models.BlacklistToken{
		TokenHash: hashedToken,
		UserID:    userId,
		Expires:   claims.ExpiresAt.Time,
	}

	if err := s.authRepo.AddToBlacklist(blacklistToken); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) IsTokenBlacklisted(tokenHash string) (bool, error) {
	return s.authRepo.IsTokenBlacklisted(tokenHash)
}

func (s *AuthService) DeleteExpiredTokens() error {
	return s.authRepo.DeleteExpiredTokens()
}
