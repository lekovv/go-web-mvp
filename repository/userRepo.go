package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/lekovv/go-web-mvp/models"
	"gorm.io/gorm"
)

type UserRepoInterface interface {
	CreateUserWithPatient(ctx context.Context, user *models.User, patient *models.Patient) error
	CreateUserWithDoctor(ctx context.Context, user *models.User, doctor *models.Doctor) error
	GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, updates *models.UpdateUserDTO) error
	DeleteUserById(ctx context.Context, id uuid.UUID, user *models.User) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepoInterface {
	return &UserRepository{db}
}

func (r *UserRepository) CreateUserWithPatient(ctx context.Context, user *models.User, patient *models.Patient) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		patient.UserID = user.ID

		if err := tx.Create(patient).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *UserRepository) CreateUserWithDoctor(ctx context.Context, user *models.User, doctor *models.Doctor) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		doctor.UserID = user.ID

		if err := tx.Create(doctor).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *UserRepository) GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Patient").
		Preload("Doctor.Specialization").
		First(&user, "id = ?", id).Error
	return &user, err
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	if err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, id uuid.UUID, updates *models.UpdateUserDTO) error {
	result := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *UserRepository) DeleteUserById(ctx context.Context, id uuid.UUID, user *models.User) error {
	result := r.db.WithContext(ctx).Delete(user, "id =?", id)
	return result.Error
}
