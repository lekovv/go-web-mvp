package repository

import (
	"context"
	"time"

	"github.com/lekovv/go-web-mvp/models"
	"gorm.io/gorm"
)

type AuthRepoInterface interface {
	AddToBlacklist(ctx context.Context, blt *models.BlacklistToken) error
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
	DeleteExpiredTokens(ctx context.Context) error
}

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db}
}

func (r *AuthRepository) AddToBlacklist(ctx context.Context, blt *models.BlacklistToken) error {
	result := r.db.WithContext(ctx).Create(blt)
	return result.Error
}

func (r *AuthRepository) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM blacklist_tokens WHERE token_hash = ? AND expires > ?)`
	err := r.db.WithContext(ctx).Raw(query, token, time.Now()).Scan(&exists).Error
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *AuthRepository) DeleteExpiredTokens(ctx context.Context) error {
	result := r.db.WithContext(ctx).Where("expires <= ?", time.Now()).Delete(&models.BlacklistToken{})
	return result.Error
}
