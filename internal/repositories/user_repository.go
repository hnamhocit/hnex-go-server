package repositories

import (
	"github.com/hnamhocit/go-learning/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func (r *UserRepository) GetUserByEmail(email string) *models.User {
	var user models.User
	err := r.DB.Where("email = ?", email).First(&user).Error

	if err != nil {
		return nil
	}

	return &user
}

func (r *UserRepository) UpdateUser(id uint, user *models.User) (*models.User, error) {
	err := r.DB.Model(user).Where("id = ?", id).Updates(user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *UserRepository) GetUserById(id uint) (*models.User, error) {
	var user models.User
	err := r.DB.Where("id = ?", id).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) CreateUser(user *models.User) (*models.User, error) {
	err := r.DB.Create(&user).Error

	if err != nil {
		return nil, err
	}

	return user, nil
}
