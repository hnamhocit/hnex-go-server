package models

type Role string
type Theme string

type User struct {
	Model
	Password        string  `gorm:"column:password;not null" json:"password"`
	DisplayName     string  `gorm:"column:display_name;not null" json:"display_name"`
	Email           string  `gorm:"unique,column:email;not null" json:"email"`
	Username        string  `gorm:"unique,column:username;not null" json:"username"`
	Bio             *string `gorm:"column:bio" json:"bio"`
	PhoneNumber     *string `gorm:"column:phone_number" json:"phone_number"`
	PhotoURL        *string `gorm:"column:photo_url" json:"photo_url"`
	BackgroundURL   *string `gorm:"column:background_url" json:"background_url"`
	RefreshToken    *string `gorm:"column:refresh_token" json:"refresh_token"`
	Theme           Theme   `gorm:"column:theme,default:'LIGHT',type:ENUM('LIGHT', 'DARK')" json:"theme"`
	Role            Role    `gorm:"column:role,default:'USER',type:ENUM('USER', 'ADMIN')" json:"role"`
	IsEmailVerified bool    `gorm:"column:is_email_verified;default:false" json:"is_email_verified"`
	IsBanned        bool    `gorm:"column:is_banned;default:false" json:"is_banned"`
}
