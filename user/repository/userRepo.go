package repository

import (
	"github.com/google/uuid"
	"github.com/lekovv/go-crud-simple/user/model"
	"gorm.io/gorm"
)

type UserRepoInterface interface {
	CreateUser(user *model.User) error
	GetUserById(id uuid.UUID, user *model.User) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepoInterface {
	return &UserRepository{db}
}

func (r *UserRepository) CreateUser(user *model.User) error {
	result := r.db.Create(user)
	return result.Error
}

func (r *UserRepository) GetUserById(id uuid.UUID, user *model.User) error {
	result := r.db.First(user, "id = ?", id)
	return result.Error
}
