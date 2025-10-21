package repository

import (
	"github.com/google/uuid"
	"github.com/lekovv/go-web-mvp/models"
	"gorm.io/gorm"
)

type UserRepoInterface interface {
	CreateUserWithPatient(user *models.User, patient *models.Patient) error
	CreateUserWithDoctor(user *models.User, doctor *models.Doctor) error
	GetUserById(id uuid.UUID) (*models.User, error)
	GetPatientByUserId(userId uuid.UUID) (*models.Patient, error)
	GetDoctorByUserId(userId uuid.UUID) (*models.Doctor, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(id uuid.UUID, updates *models.UpdateUserDTO) error
	DeleteUserById(id uuid.UUID, user *models.User) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepoInterface {
	return &UserRepository{db}
}

func (r *UserRepository) CreateUserWithPatient(user *models.User, patient *models.Patient) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		patient.UserID = user.ID

		if err := tx.Create(patient).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *UserRepository) CreateUserWithDoctor(user *models.User, doctor *models.Doctor) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		doctor.UserID = user.ID

		if err := tx.Create(doctor).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *UserRepository) GetUserById(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.Preload("Role").First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetPatientByUserId(userId uuid.UUID) (*models.Patient, error) {
	var patient models.Patient
	if err := r.db.First(&patient, "user_id = ?", userId).Error; err != nil {
		return nil, err
	}
	return &patient, nil
}

func (r *UserRepository) GetDoctorByUserId(userId uuid.UUID) (*models.Doctor, error) {
	var doctor models.Doctor
	if err := r.db.Preload("Specialization").First(&doctor, "user_id = ?", userId).Error; err != nil {
		return nil, err
	}
	return &doctor, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUser(id uuid.UUID, updates *models.UpdateUserDTO) error {
	result := r.db.Model(&models.User{}).Where("id = ?", id).Updates(updates)
	return result.Error
}

func (r *UserRepository) DeleteUserById(id uuid.UUID, user *models.User) error {
	result := r.db.First(user, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	result = r.db.Delete(user, "id =?", id)
	return result.Error
}
