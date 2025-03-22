package models

type Profile struct {
	Model
	BackgroundURL *string `json:"background_url"`
	Bio           *string `json:"bio"`
	Address       *string `json:"address"`
	PhoneNumber   *string `json:"phone_number"`

	UserID uint `json:"user_id" gorm:"unique"`
}
