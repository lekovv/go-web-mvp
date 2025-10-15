package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Patient struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserId    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	BirthDate time.Time `gorm:"not null" json:"birth_date"`
	Created   time.Time `gorm:"not null" json:"created"`
	Updated   time.Time `gorm:"not null" json:"updated"`
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
