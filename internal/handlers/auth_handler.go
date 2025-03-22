package handlers

import (
	"hnex_server/internal/models"
	"hnex_server/internal/repositories"
	"hnex_server/internal/utils"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	Repo        *repositories.AuthRepository
	UserRepo    *repositories.UserRepository
	ProfileRepo *repositories.ProfileRepository
}

type LoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=100"`
}

type RegisterDTO struct {
	DisplayName string `json:"display_name" binding:"required,min=2,max=32"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8,max=100"`
}

type ActivateAccountDTO struct {
	ActivationCode string `json:"activation_code" binding:"required"`
	UserID         uint   `json:"user_id" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input LoginDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(500, gin.H{"code": 0, "msg": err.Error()})
		return
	}

	var existingUser *models.User
	if result := h.Repo.DB.Where("email = ?", input.Email).First(&existingUser); result.Error != nil {
		c.JSON(500, gin.H{"code": 0, "msg": result.Error.Error()})
		return
	}

	if existingUser == nil {
		c.JSON(500, gin.H{"code": 0, "msg": "User not found!"})
		return
	}

	ok, err := utils.Verify(input.Password, existingUser.Password)
	if err != nil {
		c.JSON(500, gin.H{"code": 0, "msg": err.Error()})
		return
	}

	if !ok {
		c.JSON(500, gin.H{"code": 0, "msg": "Password is incorrect!"})
		return
	}

	tokens, err := utils.GenerateTokens(existingUser.ID)
	if err != nil {
		c.JSON(500, gin.H{"code": 0, "msg": err.Error()})
		return
	}

	updateErr := h.Repo.UpdateRefreshToken(existingUser.ID, &tokens.RefreshToken)
	if updateErr != nil {
		c.JSON(500, gin.H{"code": 0, "msg": updateErr.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 1, "msg": "Login successful!", "data": tokens})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input RegisterDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(500, gin.H{"code": 0, "msg": err.Error()})
		return
	}

	existingUser := &models.User{}
	result := h.Repo.DB.Where("email = ?", input.Email).First(existingUser)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		c.JSON(500, gin.H{"code": 0, "msg": result.Error.Error()})
		return
	}

	if existingUser.ID != 0 {
		c.JSON(500, gin.H{"code": 0, "msg": "Email already exists!"})
		return
	}

	hashedPassword, err := utils.Hash(input.Password)
	if err != nil {
		c.JSON(500, gin.H{"code": 0, "msg": err.Error()})
		return
	}

	newUser := &models.User{
		Email:       input.Email,
		Password:    hashedPassword,
		DisplayName: input.DisplayName,
	}

	if err := h.Repo.DB.Create(newUser).Error; err != nil {
		c.JSON(500, gin.H{"code": 0, "msg": err.Error()})
		return
	}

	tokens, err := utils.GenerateTokens(newUser.ID)
	if err != nil {
		c.JSON(500, gin.H{"code": 0, "msg": err.Error()})
		return
	}

	updateUserErr := h.Repo.UpdateRefreshToken(newUser.ID, &tokens.RefreshToken)
	if updateUserErr != nil {
		c.JSON(500, gin.H{"code": 0, "msg": updateUserErr.Error()})
		return
	}

	createProfileErr := h.ProfileRepo.Create(newUser.ID)
	if createProfileErr != nil {
		c.JSON(500, gin.H{"code": 0, "msg": createProfileErr.Error()})
		return
	}

	go func() {
		type ActivationCodeData struct {
			DisplayName    string
			ActivationCode string
			ExpiresAt      string
		}
		activationCode, expiresAt := utils.GenerateActivationCode()

		var profile models.Profile
		err := h.Repo.DB.Model(&models.Profile{}).Where("user_id = ?", newUser.ID).First(&profile).Error
		if err != nil {
			log.Printf("[PROFILE] Error fetching profile: %v", err.Error())
		}

		data := ActivationCodeData{
			DisplayName:    newUser.DisplayName,
			ActivationCode: activationCode,
			ExpiresAt:      expiresAt.Format(time.RFC3339),
		}

		err = h.Repo.UpdateActivationCode(newUser.ID, activationCode, expiresAt)
		if err != nil {
			log.Printf("[PROFILE] Error updating profile activation code: %v", err)
		}

		err = SendMail("Activate Account", newUser.Email, "activation_email", data)
		if err != nil {
			log.Printf("[MAIL] Error sending activation email: %v", err)
		}
	}()

	c.JSON(200, gin.H{"code": 1, "msg": "Registration successful!", "data": tokens})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	sub, ok := c.Get("sub")
	if !ok {
		c.JSON(500, gin.H{"code": 0, "msg": "Invalid token!"})
		return
	}

	user := &models.User{}
	result := h.Repo.DB.Where("id = ?", sub).First(user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(500, gin.H{"code": 0, "msg": "User not found!"})
			return
		}

		c.JSON(500, gin.H{"code": 0, "msg": result.Error.Error()})
		return
	}

	tokens, err := utils.GenerateTokens(user.ID)
	if err != nil {
		c.JSON(500, gin.H{"code": 0, "msg": err.Error()})
		return
	}

	updateErr := h.Repo.UpdateRefreshToken(user.ID, &tokens.RefreshToken)
	if updateErr != nil {
		c.JSON(500, gin.H{"code": 0, "msg": updateErr.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 1, "msg": "Refresh successful!", "data": tokens})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	sub, ok := c.Get("sub")
	if !ok {
		c.JSON(500, gin.H{"code": 0, "msg": "Invalid token!"})
		return
	}

	user := &models.User{}
	result := h.Repo.DB.Where("id = ?", sub).First(user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(500, gin.H{"code": 0, "msg": "User not found!"})
			return
		}

		c.JSON(500, gin.H{"code": 0, "msg": result.Error.Error()})
		return
	}

	updateErr := h.Repo.UpdateRefreshToken(user.ID, nil)
	if updateErr != nil {
		c.JSON(500, gin.H{"code": 0, "msg": updateErr.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 1, "msg": "Logout successful!"})
}

func (h *AuthHandler) ActivateAccount(c *gin.Context) {
	var input ActivateAccountDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"code": 0, "msg": err.Error()})
		return
	}

	user, err := h.UserRepo.GetUser(input.UserID)
	if err != nil {
		c.JSON(500, gin.H{"code": 0, "msg": err.Error()})
		return
	}

	if user.ActivationCodeExpiresAt.Before(time.Now()) {
		c.JSON(400, gin.H{"code": 0, "msg": "Activation code expired!"})
		return
	}

	if *user.ActivationCode != input.ActivationCode {
		c.JSON(400, gin.H{"code": 0, "msg": "Invalid activation code!"})
		return
	}

	emailVerifiedErr := h.Repo.EmailVerified(input.UserID)
	if emailVerifiedErr != nil {
		c.JSON(500, gin.H{"code": 0, "msg": emailVerifiedErr.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 1, "msg": "Account activated!", "data": gin.H{
		"is_email_verified": true,
	}})
}

type RefreshActivationCodeDTO struct {
	UserID uint `json:"user_id" binding:"required"`
}

func (h *AuthHandler) RefreshActivateCode(c *gin.Context) {
	var input RefreshActivationCodeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"code": 0, "msg": err.Error()})
		return
	}

	activationCode, expiresAt := utils.GenerateActivationCode()
	if err := h.Repo.UpdateActivationCode(input.UserID, activationCode, expiresAt); err != nil {
		c.JSON(500, gin.H{"code": 0, "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 1, "msg": "Activation code refreshed!", "data": gin.H{
		"activation_code": activationCode,
		"expires_at":      expiresAt,
	}})
}
