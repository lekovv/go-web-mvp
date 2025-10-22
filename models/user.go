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
	RoleID       uuid.UUID `gorm:"type:uuid;not null" json:"-"`
	Role         Role      `gorm:"foreignKey:RoleID" json:"-"`
	FirstName    string    `gorm:"not null" json:"first_name"`
	LastName     string    `gorm:"not null" json:"last_name"`
	MiddleName   *string   `json:"middle_name"`
	IsActive     bool      `gorm:"not null" json:"is_active"`
	PhoneNumber  string    `gorm:"not null" json:"phone_number"`
	Created      time.Time `gorm:"not null" json:"-"`
	Updated      time.Time `gorm:"not null" json:"-"`
	Patient      *Patient  `gorm:"foreignKey:UserID" json:"-"`
	Doctor       *Doctor   `gorm:"foreignKey:UserID" json:"-"`
}

type UserResponse struct {
	ID              uuid.UUID `json:"id"`
	Email           string    `json:"email"`
	Gender          string    `json:"gender"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	MiddleName      *string   `json:"middle_name,omitempty"`
	IsActive        bool      `json:"is_active"`
	PhoneNumber     string    `json:"phone_number"`
	BirthDate       *string   `json:"birth_date,omitempty"`
	Specialization  *string   `json:"specialization,omitempty"`
	Bio             *string   `json:"bio,omitempty"`
	ExperienceYears *int      `json:"experience_years,omitempty"`
	Price           *int      `json:"price,omitempty"`
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
	Token        *string       `json:"token,omitempty"`
	UserResponse *UserResponse `json:"user,omitempty"`
}
