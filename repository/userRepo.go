package repository

import (
	"github.com/google/uuid"
	"github.com/lekovv/go-web-mvp/models"
	"gorm.io/gorm"
)

type UserRepoInterface interface {
	CreateUser(user *models.User) error
	GetUserById(id uuid.UUID, user *models.User) error
	UpdateUser(id uuid.UUID, updates *models.UpdateUserDTO) error
	DeleteUser(id uuid.UUID, user *models.User) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepoInterface {
	return &UserRepository{db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	result := r.db.Create(user)
	return result.Error
}

func (r *UserRepository) GetUserById(id uuid.UUID, user *models.User) error {
	result := r.db.First(user, "id = ?", id)
	return result.Error
}

func (r *UserRepository) UpdateUser(id uuid.UUID, updates *models.UpdateUserDTO) error {
	result := r.db.Model(&models.User{}).Where("id = ?", id).Updates(updates)
	return result.Error
}

func (r *UserRepository) DeleteUser(id uuid.UUID, user *models.User) error {
	result := r.db.First(user, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	result = r.db.Delete(user, "id =?", id)
	return result.Error
}
