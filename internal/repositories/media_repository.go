package repositories

import (
	"github.com/hnamhocit/go-learning/internal/models"
	"gorm.io/gorm"
)

type MediaRepository struct {
	DB *gorm.DB
}

func (r *MediaRepository) UploadFile(file *models.Media) (*models.Media, error) {
	if err := r.DB.Create(file).Error; err != nil {
		return nil, err
	}

	return file, nil
}

func (r *MediaRepository) UploadFiles(files []*models.Media) ([]*models.Media, error) {
	if err := r.DB.Create(files).Error; err != nil {
		return nil, err
	}

	return files, nil
}

func (r *MediaRepository) GetUploadFile(id uint) (*models.Media, error) {
	var upload models.Media
	result := r.DB.First(&upload, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &upload, nil
}
