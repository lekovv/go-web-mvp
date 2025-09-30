package repository

import (
	"github.com/lekovv/go-crud-simple/db"
	"github.com/lekovv/go-crud-simple/user/model"
)

type UserRepository struct{}

func NewRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) CreateUser(user *model.User) error {
	result := db.DB.Create(user)
	return result.Error
}
