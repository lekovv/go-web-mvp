package repository

import (
	"github.com/google/uuid"
	"github.com/lekovv/go-web-mvp/models"
	"gorm.io/gorm"
)

type UserRepoInterface interface {
	GetUserById(id uuid.UUID) (*models.UserResponse, error)
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

func (r *UserRepository) GetUserById(id uuid.UUID) (*models.UserResponse, error) {
	var resp models.UserResponse

	selectFields := `
        users.id,
        users.email,
        users.gender,
        users.first_name,
        users.last_name,
        users.middle_name,
        users.is_active,
        users.phone_number,
        to_char(patients.birth_date, 'YYYY-MM-DD') AS birth_date,
        specializations.name AS specialization,
        doctors.bio AS bio,
        doctors.experience_years AS experience_years,
        doctors.price AS price
    `
	err := r.db.
		Model(&models.User{}).
		Select(selectFields).
		Joins("LEFT JOIN patients ON patients.user_id = users.id").
		Joins("LEFT JOIN doctors ON doctors.user_id = users.id").
		Joins("LEFT JOIN specializations ON specializations.id = doctors.specialization_id").
		Where("users.id = ?", id).
		Scan(&resp).
		Error

	if err != nil {
		return nil, err
	}
	return &resp, nil
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
