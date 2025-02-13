package dtos

type RegisterDTO struct {
	DisplayName string `json:"displayName" binding:"required,min=6,max=50"`
	Password    string `json:"password" binding:"required,min=8"`
	Email       string `json:"email" binding:"required,email"`
}
