package repository

import (
	"github.com/lekovv/go-crud-simple/db"
	"github.com/lekovv/go-crud-simple/user/model"
)

type UserRepoInterface interface {
	CreateUser(user *model.User) error
}

type UserRepository struct{}

func NewUserRepository() UserRepoInterface {
	return &UserRepository{}
}

func (r *UserRepository) CreateUser(user *model.User) error {
	result := db.DB.Create(user)
	return result.Error
}
