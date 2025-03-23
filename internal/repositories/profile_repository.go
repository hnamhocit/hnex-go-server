package repositories

import (
	"hnex_server/internal/models"

	"gorm.io/gorm"
)

type ProfileRepository struct {
	DB *gorm.DB
}

func (r *ProfileRepository) Create(userId uint) error {
	return r.DB.Create(&models.Profile{UserID: userId}).Error
}

func (r *ProfileRepository) Get(userId uint) (*models.Profile, error) {
	var profile models.Profile
	if err := r.DB.Where("user_id = ?", userId).First(&profile).Error; err != nil {
		return nil, err
	}

	return &profile, nil
}
