package repositories

import (
	"hnex_server/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func (r *UserRepository) GetUser(id uint) (*models.User, error) {
	var user models.User
	if err := r.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUsers() ([]*models.User, error) {
	var users []*models.User
	if err := r.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) UpdateUser(id uint, user *models.User) error {
	if err := r.DB.Model(&models.User{}).Where("id = ?", id).Updates(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) DeleteUser(id uint) error {
	if err := r.DB.Delete(&models.User{}, id).Error; err != nil {
		return err
	}
	return nil
}
