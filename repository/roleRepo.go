package repository

import (
	"github.com/lekovv/go-web-mvp/models"
	"gorm.io/gorm"
)

type RoleRepoInterface interface {
	GetRoleByName(name string) (*models.Role, error)
}

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepoInterface {
	return &RoleRepository{db}
}

func (r *RoleRepository) GetRoleByName(name string) (*models.Role, error) {
	var role models.Role

	if err := r.db.First(&role, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return &role, nil
}
