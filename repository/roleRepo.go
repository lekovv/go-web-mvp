package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/lekovv/go-web-mvp/models"
	"gorm.io/gorm"
)

type RoleRepoInterface interface {
	GetRoleByName(ctx context.Context, name string) (*models.Role, error)
	GetRoleByID(ctx context.Context, id uuid.UUID) (*models.Role, error)
}

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepoInterface {
	return &RoleRepository{db}
}

func (r *RoleRepository) GetRoleByName(ctx context.Context, name string) (*models.Role, error) {
	var role models.Role

	if err := r.db.WithContext(ctx).First(&role, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) GetRoleByID(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	var role models.Role
	if err := r.db.WithContext(ctx).First(&role, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}
