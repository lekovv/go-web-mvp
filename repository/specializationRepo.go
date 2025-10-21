package repository

import (
	"context"

	"github.com/lekovv/go-web-mvp/models"
	"gorm.io/gorm"
)

type SpecializationRepoInterface interface {
	GetSpecializationByName(ctx context.Context, name string) (*models.Specialization, error)
}

type SpecializationRepository struct {
	db *gorm.DB
}

func NewSpecializationRepository(db *gorm.DB) SpecializationRepoInterface {
	return &SpecializationRepository{db}
}

func (r *SpecializationRepository) GetSpecializationByName(ctx context.Context, name string) (*models.Specialization, error) {
	var specialization models.Specialization

	if err := r.db.WithContext(ctx).First(&specialization, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return &specialization, nil
}
