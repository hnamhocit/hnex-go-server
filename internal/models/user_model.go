package models

import "time"

const (
	USER  = iota
	ADMIN = iota
)

type User struct {
	Model
	Email                   string     `json:"email" gorm:"unique"`
	Password                string     `json:"password" gorm:"not null"`
	DisplayName             string     `json:"display_name" gorm:"not null"`
	IsEmailVerified         bool       `json:"is_email_verified" gorm:"default:false"`
	PhotoURL                *string    `json:"photo_url"`
	RefreshToken            *string    `json:"refresh_token"`
	ActivationCode          *string    `json:"activation_code"`
	ActivationCodeExpiresAt *time.Time `json:"activation_code_expires_at"`
	Role                    int        `json:"role" gorm:"default:0"`

	Profile Profile `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}
