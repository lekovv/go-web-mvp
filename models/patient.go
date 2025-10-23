package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Patient struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	BirthDate time.Time `gorm:"not null" json:"birth_date"`
	Created   time.Time `gorm:"not null" json:"created"`
	Updated   time.Time `gorm:"not null" json:"updated"`
}

type PatientRegistrationDTO struct {
	Email       string  `json:"email" validate:"required,email"`
	Password    string  `json:"password" validate:"required,min=8,max=72"`
	FirstName   string  `json:"first_name" validate:"required"`
	LastName    string  `json:"last_name" validate:"required"`
	MiddleName  *string `json:"middle_name,omitempty"`
	PhoneNumber string  `json:"phone_number" validate:"required,min=10,max=15"`
	BirthDate   string  `json:"birth_date" validate:"required"`
	Gender      string  `json:"gender" validate:"required,oneof=male female other"`
}

func (patient *Patient) BeforeCreate(tx *gorm.DB) (err error) {
	patient.ID = uuid.New()
	patient.Created = time.Now()
	patient.Updated = time.Now()
	return nil
}

func (patient *Patient) BeforeUpdate(tx *gorm.DB) (err error) {
	patient.Updated = time.Now()
	return nil
}
