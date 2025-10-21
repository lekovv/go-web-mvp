package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/lekovv/go-web-mvp/config"
	AppErrors "github.com/lekovv/go-web-mvp/errors"
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
	env                *config.Env
}

func NewAuthService(
	specializationRepo repository.SpecializationRepoInterface,
	userRepo repository.UserRepoInterface,
	roleRepo repository.RoleRepoInterface,
	authRepo repository.AuthRepoInterface,
	env *config.Env,
) AuthServiceInterface {
	return &AuthService{
		specializationRepo: specializationRepo,
		userRepo:           userRepo,
		roleRepo:           roleRepo,
		authRepo:           authRepo,
		env:                env,
	}
}

func (s *AuthService) RegisterPatient(payload *models.PatientRegistrationDTO) (*models.AuthResponse, error) {
	existingUser, err := s.userRepo.GetUserByEmail(payload.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, AppErrors.WrapError(
				err,
				AppErrors.ErrorTypeInternal,
				"Failed to check existing user",
			)
		}
	}

	if existingUser != nil {
		return nil, AppErrors.NewConflictError("User already exists")
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return nil, AppErrors.WrapError(
			err,
			AppErrors.ErrorTypeInternal,
			"Failed to hash password",
		)
	}

	patientRole, err := s.roleRepo.GetRoleByName("patient")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, AppErrors.NewNotFoundError("Patient role not found")
		}
		return nil, AppErrors.WrapError(
			err,
			AppErrors.ErrorTypeInternal,
			"Failed to get patient role",
		)
	}

	birthDate, err := utils.ParseDate(payload.BirthDate)
	if err != nil {
		return nil, AppErrors.NewValidationError(
			"Invalid birth date format",
			[]string{err.Error()},
		)
	}

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

	patient := &models.Patient{
		BirthDate: birthDate,
	}

	if err := s.userRepo.CreateUserWithPatient(user, patient); err != nil {
		return nil, AppErrors.WrapError(
			err,
			AppErrors.ErrorTypeInternal,
			"Failed to create user and patient",
		)
	}

	return &models.AuthResponse{User: user}, nil
}

func (s *AuthService) CreateDoctor(payload *models.DoctorRegistrationDTO) (*models.AuthResponse, error) {
	existingUser, err := s.userRepo.GetUserByEmail(payload.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, AppErrors.WrapError(
				err,
				AppErrors.ErrorTypeInternal,
				"Failed to check existing user",
			)
		}
	}

	if existingUser != nil {
		return nil, AppErrors.NewConflictError("User already exists")
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return nil, AppErrors.WrapError(
			err,
			AppErrors.ErrorTypeInternal,
			"Failed to hash password",
		)
	}

	doctorRole, err := s.roleRepo.GetRoleByName("doctor")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, AppErrors.NewNotFoundError("Doctor role not found")
		}
		return nil, AppErrors.WrapError(
			err,
			AppErrors.ErrorTypeInternal,
			"Failed to get doctor role",
		)
	}

	specialization, err := s.specializationRepo.GetSpecializationByName(payload.Specialization)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, AppErrors.NewNotFoundError("Specialization not found")
		}
		return nil, AppErrors.WrapError(
			err,
			AppErrors.ErrorTypeInternal,
			"Failed to get specialization",
		)
	}

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

	doctor := &models.Doctor{
		SpecializationID: specialization.ID,
		Bio:              payload.Bio,
		ExperienceYears:  payload.ExperienceYears,
		Price:            payload.Price,
	}

	if err := s.userRepo.CreateUserWithDoctor(user, doctor); err != nil {
		return nil, AppErrors.WrapError(
			err,
			AppErrors.ErrorTypeInternal,
			"Failed to create user and doctor",
		)
	}

	return &models.AuthResponse{User: user}, nil
}

func (s *AuthService) Login(payload *models.LoginDTO) (*models.AuthResponse, error) {
	user, err := s.userRepo.GetUserByEmail(payload.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, AppErrors.NewUnauthorizedError("Invalid email or password")
		}
		return nil, AppErrors.WrapError(
			err,
			AppErrors.ErrorTypeInternal,
			"Failed to get user",
		)
	}

	if !user.IsActive {
		return nil, AppErrors.NewForbiddenError("User account is deactivated")
	}

	if !utils.CheckPasswordHash(payload.Password, user.PasswordHash) {
		return nil, AppErrors.NewUnauthorizedError("Invalid email or password")
	}

	token, err := utils.GenerateJWT(
		user.ID,
		user.Email,
		user.RoleID,
		s.env.JWTSecret,
		s.env.JWTExpire,
	)
	if err != nil {
		return nil, AppErrors.WrapError(
			err,
			AppErrors.ErrorTypeInternal,
			"Failed to generate token",
		)
	}

	return &models.AuthResponse{Token: &token}, nil
}

func (s *AuthService) Logout(token string, userId uuid.UUID) error {
	claims, err := utils.ValidateJWT(token, s.env.JWTSecret)
	if err != nil {
		return AppErrors.WrapError(
			err,
			AppErrors.ErrorTypeUnauthorized,
			"Invalid token",
		)
	}

	hashedToken := utils.HashToken(token, s.env.JWTSecret)
	blacklistToken := &models.BlacklistToken{
		TokenHash: hashedToken,
		UserID:    userId,
		Expires:   claims.ExpiresAt.Time,
	}

	if err := s.authRepo.AddToBlacklist(blacklistToken); err != nil {
		return AppErrors.WrapError(
			err,
			AppErrors.ErrorTypeInternal,
			"Failed to add token to blacklist",
		)
	}

	return nil
}

func (s *AuthService) IsTokenBlacklisted(tokenHash string) (bool, error) {
	isBlacklisted, err := s.authRepo.IsTokenBlacklisted(tokenHash)
	if err != nil {
		return false, AppErrors.WrapError(
			err,
			AppErrors.ErrorTypeInternal,
			"Failed to check token blacklist",
		)
	}

	return isBlacklisted, nil
}

func (s *AuthService) DeleteExpiredTokens() error {
	if err := s.authRepo.DeleteExpiredTokens(); err != nil {
		return AppErrors.WrapError(
			err,
			AppErrors.ErrorTypeInternal,
			"Failed to delete expired tokens",
		)
	}

	return nil
}
