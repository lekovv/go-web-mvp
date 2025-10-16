package repository

import (
	"time"

	"github.com/lekovv/go-web-mvp/models"
	"gorm.io/gorm"
)

type AuthRepoInterface interface {
	AddToBlacklist(blt *models.BlacklistToken) error
	IsTokenBlacklisted(token string) (bool, error)
	DeleteExpiredTokens() error
}

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db}
}

func (r *AuthRepository) AddToBlacklist(blt *models.BlacklistToken) error {
	result := r.db.Create(blt)
	return result.Error
}

func (r *AuthRepository) IsTokenBlacklisted(token string) (bool, error) {
	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM blacklist_tokens WHERE token_hash = ? AND expires > ?)`
	err := r.db.Raw(query, token, time.Now()).Scan(&exists).Error
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *AuthRepository) DeleteExpiredTokens() error {
	result := r.db.Where("expires <= ?", time.Now()).Delete(&models.BlacklistToken{})
	return result.Error
}
