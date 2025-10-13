package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Appointment struct {
	ID                  uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	DoctorID            uuid.UUID `gorm:"type:uuid;not null" json:"doctor_id"`
	PatientID           uuid.UUID `gorm:"type:uuid;not null" json:"patient_id"`
	AppointmentDate     time.Time `gorm:"not null" json:"appointment_date"`
	AppointmentStatusId uuid.UUID `gorm:"type:uuid;not null" json:"appointment_status_id"`
	Created             time.Time `gorm:"not null" json:"created_at"`
	Updated             time.Time `gorm:"not null" json:"updated_at"`
}

func (appointment *Appointment) BeforeCreate(tx *gorm.DB) (err error) {
	appointment.ID = uuid.New()
	appointment.Created = time.Now()
	appointment.Updated = time.Now()
	return nil
}

func (appointment *Appointment) BeforeUpdate(tx *gorm.DB) (err error) {
	appointment.Updated = time.Now()
	return nil
}

type AppointmentStatus struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name    string    `gorm:"not null" json:"name"`
	Created time.Time `gorm:"not null" json:"created_at"`
	Updated time.Time `gorm:"not null" json:"updated_at"`
}

func (appointmentStatus *AppointmentStatus) BeforeCreate(tx *gorm.DB) (err error) {
	appointmentStatus.ID = uuid.New()
	appointmentStatus.Created = time.Now()
	appointmentStatus.Updated = time.Now()
	return nil
}

func (appointmentStatus *AppointmentStatus) BeforeUpdate(tx *gorm.DB) (err error) {
	appointmentStatus.Updated = time.Now()
	return nil
}
