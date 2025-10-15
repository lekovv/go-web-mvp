package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name    string    `gorm:"unique;not null" json:"name"`
	Created time.Time `gorm:"not null" json:"created"`
	Updated time.Time `gorm:"not null" json:"updated"`
}

func (role *Role) BeforeCreate(tx *gorm.DB) (err error) {
	role.ID = uuid.New()
	role.Created = time.Now()
	role.Updated = time.Now()
	return nil
}

func (role *Role) BeforeUpdate(tx *gorm.DB) (err error) {
	role.Updated = time.Now()
	return nil
}
