package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BlacklistToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TokenHash string    `gorm:"not null;unique" json:"token_hash"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Expires   time.Time `gorm:"not null" json:"expires"`
	Created   time.Time `gorm:"not null" json:"created"`
}

func (token *BlacklistToken) BeforeCreate(tx *gorm.DB) (err error) {
	token.ID = uuid.New()
	token.Created = time.Now()
	return nil
}
