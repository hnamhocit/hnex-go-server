package models

type User struct {
	Model
	Email           string  `json:"email" gorm:"unique"`
	Password        string  `json:"password" gorm:"not null"`
	DisplayName     string  `json:"display_name" gorm:"not null"`
	IsEmailVerified bool    `json:"is_email_verified" gorm:"default:false"`
	PhotoURL        *string `json:"photo_url"`
	BackgroundURL   *string `json:"background_url"`
	Bio             *string `json:"bio"`
	Address         *string `json:"address"`
	PhoneNumber     *string `json:"phone_number"`
	RefreshToken    *string `json:"refresh_token"`
}
