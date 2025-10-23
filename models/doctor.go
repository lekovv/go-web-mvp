package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Doctor struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID           uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	SpecializationID uuid.UUID      `gorm:"type:uuid;not null" json:"specialization_id"`
	Specialization   Specialization `gorm:"foreignKey:SpecializationID" json:"-"`
	Bio              *string        `json:"bio"`
	ExperienceYears  *int           `json:"experience_years"`
	Price            int            `gorm:"not null" json:"price"`
	Created          time.Time      `gorm:"not null" json:"created"`
	Updated          time.Time      `gorm:"not null" json:"updated"`
}

type DoctorRegistrationDTO struct {
	Email           string  `json:"email" validate:"required,email"`
	Password        string  `json:"password" validate:"required,min=8,max=72"`
	FirstName       string  `json:"first_name" validate:"required"`
	LastName        string  `json:"last_name" validate:"required"`
	MiddleName      *string `json:"middle_name,omitempty"`
	PhoneNumber     string  `json:"phone_number" validate:"required,min=10,max=15"`
	Gender          string  `json:"gender" validate:"required,oneof=male female other"`
	Specialization  string  `json:"specialization" validate:"required"`
	Bio             *string `json:"bio"`
	ExperienceYears *int    `json:"experience_years"`
	Price           int     `json:"price" validate:"required"`
}

func (doctor *Doctor) BeforeCreate(tx *gorm.DB) (err error) {
	doctor.ID = uuid.New()
	doctor.Created = time.Now()
	doctor.Updated = time.Now()
	return nil
}

func (doctor *Doctor) BeforeUpdate(tx *gorm.DB) (err error) {
	doctor.Updated = time.Now()
	return nil
}
