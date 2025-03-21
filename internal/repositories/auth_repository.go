package repositories

import (
	"hnex_server/internal/models"
	"hnex_server/internal/utils"

	"gorm.io/gorm"
)

type AuthRepository struct {
	DB *gorm.DB
}

func (r *AuthRepository) UpdateRefreshToken(userID uint, newToken *string) error {
	if newToken == nil {
		return r.DB.Model(&models.User{}).Where("id = ?", userID).Update("refresh_token", nil).Error
	}

	hashedToken, err := utils.Hash(*newToken)

	if err != nil {
		return err
	}

	return r.DB.Model(&models.User{}).Where("id = ?", userID).Update("refresh_token", hashedToken).Error
}
