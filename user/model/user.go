package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	FirstName string    `gorm:"not null" json:"first_name"`
	LastName  string    `gorm:"not null" json:"last_name"`
	IsActive  bool      `gorm:"not null" json:"is_active"`
	Created   time.Time `gorm:"not null" json:"created"`
	Updated   time.Time `gorm:"not null" json:"updated"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.New()
	user.Created = time.Now()
	user.Updated = time.Now()
	return nil
}

func (user *User) BeforeUpdate(tx *gorm.DB) (err error) {
	user.Updated = time.Now()
	return nil
}
