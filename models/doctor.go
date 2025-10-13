package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Doctor struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID           uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	SpecializationID uuid.UUID `gorm:"type:uuid;not null" json:"specialization_id"`
	Bio              *string   `json:"bio"`
	ExperienceYears  *int      `json:"experience_years"`
	Price            int       `gorm:"not null" json:"price"`
	Created          time.Time `gorm:"not null" json:"created_at"`
	Updated          time.Time `gorm:"not null" json:"updated_at"`
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
