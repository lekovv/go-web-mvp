package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Email        string    `gorm:"unique;not null" json:"email"`
	Gender       string    `gorm:"not null" json:"gender"`
	PasswordHash string    `gorm:"not null" json:"-"`
	RoleID       uuid.UUID `gorm:"type:uuid;not null" json:"role_id"`
	FirstName    string    `gorm:"not null" json:"first_name"`
	LastName     string    `gorm:"not null" json:"last_name"`
	MiddleName   *string   `json:"middle_name"`
	IsActive     bool      `gorm:"not null" json:"is_active"`
	PhoneNumber  string    `gorm:"not null" json:"phone_number"`
	Created      time.Time `gorm:"not null" json:"created"`
	Updated      time.Time `gorm:"not null" json:"updated"`
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

type UpdateUserDTO struct {
	FirstName  *string `json:"first_name,omitempty"`
	LastName   *string `json:"last_name,omitempty"`
	MiddleName *string `json:"middle_name,omitempty"`
	IsActive   *bool   `json:"is_active,omitempty"`
}

type LoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token *string `json:"token,omitempty"`
	User  *User   `json:"user,omitempty"`
}
