package repository

import (
	"github.com/google/uuid"
	"github.com/lekovv/go-web-mvp/models"
	"gorm.io/gorm"
)

type UserRepoInterface interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUserById(id uuid.UUID) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(id uuid.UUID, updates *models.UpdateUserDTO) error
	DeleteUserById(id uuid.UUID, user *models.User) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepoInterface {
	return &UserRepository{db}
}

func (r *UserRepository) CreateUser(user *models.User) (*models.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserById(id uuid.UUID) (*models.User, error) {
	var user models.User

	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUser(id uuid.UUID, updates *models.UpdateUserDTO) error {
	result := r.db.Model(&models.User{}).Where("id = ?", id).Updates(updates)
	return result.Error
}

func (r *UserRepository) DeleteUserById(id uuid.UUID, user *models.User) error {
	result := r.db.First(user, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	result = r.db.Delete(user, "id =?", id)
	return result.Error
}
