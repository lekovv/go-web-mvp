package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Specialization struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name    string    `gorm:"unique;not null" json:"name"`
	Created time.Time `gorm:"not null" json:"created_at"`
	Updated time.Time `gorm:"not null" json:"updated_at"`
}

func (specialization *Specialization) BeforeCreate(tx *gorm.DB) (err error) {
	specialization.ID = uuid.New()
	specialization.Created = time.Now()
	specialization.Updated = time.Now()
	return nil
}

func (specialization *Specialization) BeforeUpdate(tx *gorm.DB) (err error) {
	specialization.Updated = time.Now()
	return nil
}
