package handlers

import (
	"errors"
	"fmt"
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

func (h *AuthHandler) Login(c *gin.Context) {
	var input LoginDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(500, gin.H{"code": 0, "msg": err.Error()})
		return
	}

	var existingUser *models.User
	if err := h.Repo.DB.Where("email = ?", input.Email).First(&existingUser).Error; err != nil {
		c.JSON(500, gin.H{"code": 0, "msg": fmt.Sprintf("[USER]: %v", err.Error())})
		return
	}

	ok, err := utils.Verify(input.Password, existingUser.Password)
	if err != nil {
		c.JSON(500, gin.H{"code": 0, "msg": fmt.Sprintf("[ARGON2]: %v", err.Error())})
		return
	}

	if !ok {
		c.JSON(500, gin.H{"code": 0, "msg": "Password is incorrect!"})
		return
	}

	tokens, err := utils.GenerateTokens(existingUser.ID)
	if err != nil {
		c.JSON(500, gin.H{"code": 0, "msg": fmt.Sprintf("[TOKEN]: %v", err.Error())})
		return
	}

	updateErr := h.Repo.UpdateRefreshToken(existingUser.ID, &tokens.RefreshToken)
	if updateErr != nil {
		c.JSON(500, gin.H{"code": 0, "msg": fmt.Sprintf("[USER]: %v", updateErr.Error())})
		return
	}

	c.JSON(200, gin.H{"code": 1, "msg": "Login successful!", "data": tokens})
}

type RegisterDTO struct {
	DisplayName string `json:"display_name" binding:"required,min=2,max=32"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8,max=100"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input RegisterDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(500, gin.H{"code": 0, "msg": err.Error()})
		return
	}

	existingUser := &models.User{}
	err := h.Repo.DB.Where("email = ?", input.Email).First(existingUser).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(500, gin.H{"code": 0, "msg": fmt.Sprintf("[USER]: %v", err)})
			return
		}
	}

	hashedPassword, err := utils.Hash(input.Password)
	if err != nil {
		c.JSON(500, gin.H{"code": 0, "msg": fmt.Sprintf("[PASSWORD]: %v", err)})
		return
	}

	newUser := &models.User{
		Email:       input.Email,
		Password:    hashedPassword,
		DisplayName: input.DisplayName,
	}

	if err := h.Repo.DB.Create(newUser).Error; err != nil {
		c.JSON(500, gin.H{"code": 0, "msg": fmt.Sprintf("[USER]: %v", err)})
		return
	}

	tokens, err := utils.GenerateTokens(newUser.ID)
	if err != nil {
		c.JSON(500, gin.H{"code": 0, "msg": fmt.Sprintf("[TOKENS]: %v", err)})
		return
	}

	updateUserErr := h.Repo.UpdateRefreshToken(newUser.ID, &tokens.RefreshToken)
	if updateUserErr != nil {
		c.JSON(500, gin.H{"code": 0, "msg": fmt.Sprintf("[USER]: %v", updateUserErr)})
		return
	}

	createProfileErr := h.ProfileRepo.Create(newUser.ID)
	if createProfileErr != nil {
		c.JSON(500, gin.H{"code": 0, "msg": fmt.Sprintf("[PROFILE]: %v", createProfileErr)})
		return
	}

	go func() {
		type ActivationCodeData struct {
			DisplayName    string
			ActivationCode string
			ExpiresAt      string
		}
		activationCode, expiresAt := utils.GenerateActivationCode()

		data := ActivationCodeData{
			DisplayName:    newUser.DisplayName,
			ActivationCode: activationCode,
			ExpiresAt:      expiresAt.Format(time.RFC3339),
		}

		err = h.Repo.UpdateActivationCode(newUser.ID, activationCode, expiresAt)
		if err != nil {
			log.Printf("[PROFILE]: %v", err.Error())
		}

		err = SendMail("Activate Account", newUser.Email, "activate_account", data)
		if err != nil {
			log.Printf("[MAIL]: %v", err.Error())
		}
	}()

	c.JSON(200, gin.H{"code": 1, "msg": "Registration successful!", "data": tokens})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	sub, ok := c.Get("sub")
	if !ok {
		c.JSON(400, gin.H{"code": 0, "msg": "Invalid id!"})
		return
	}

	userID, ok := sub.(uint)
	if !ok {
		c.JSON(400, gin.H{"code": 0, "msg": "Invalid id!"})
		return
	}

	user, err := h.UserRepo.GetUser(userID)
	if err != nil {
		c.JSON(500, gin.H{"code": 0, "msg": fmt.Sprintf("[USER]: %v", err.Error())})
		return
	}

	tokens, err := utils.GenerateTokens(user.ID)
	if err != nil {
		c.JSON(500, gin.H{"code": 0, "msg": fmt.Sprintf("[TOKEN]: %v", err.Error())})
		return
	}

	updateErr := h.Repo.UpdateRefreshToken(user.ID, &tokens.RefreshToken)
	if updateErr != nil {
		c.JSON(500, gin.H{"code": 0, "msg": fmt.Sprintf("[UPDATE]: %v", updateErr.Error())})
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
	err := h.Repo.DB.Where("id = ?", sub).First(user).Error
	if err != nil {
		c.JSON(500, gin.H{"code": 0, "msg": fmt.Sprintf("[USER]: %v", err.Error())})
		return
	}

	updateErr := h.Repo.UpdateRefreshToken(user.ID, nil)
	if updateErr != nil {
		c.JSON(500, gin.H{"code": 0, "msg": fmt.Sprintf("[USER]: %v", updateErr.Error())})
		return
	}

	c.JSON(200, gin.H{"code": 1, "msg": "Logout successful!"})
}

type ActivateAccountDTO struct {
	ActivationCode string `json:"activation_code" binding:"required"`
}

func (h *AuthHandler) ActivateAccount(c *gin.Context) {
	var input ActivateAccountDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"code": 0, "msg": err.Error()})
		return
	}

	sub, ok := c.Get("sub")
	if !ok {
		c.JSON(400, gin.H{"code": 0, "msg": "Invalid id!"})
		return
	}

	userID, ok := sub.(uint)
	if !ok {
		c.JSON(400, gin.H{"code": 0, "msg": "Invalid id!"})
		return
	}

	user, err := h.UserRepo.GetUser(userID)
	if err != nil {
		c.JSON(500, gin.H{"code": 0, "msg": fmt.Sprintf("[USER]: %v", err)})
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

	emailVerifiedErr := h.Repo.EmailVerified(userID)
	if emailVerifiedErr != nil {
		c.JSON(500, gin.H{"code": 0, "msg": fmt.Sprintf("[EMAIL]: %v", emailVerifiedErr.Error())})
		return
	}

	c.JSON(200, gin.H{"code": 1, "msg": "Account activated!", "data": gin.H{
		"is_email_verified": true,
	}})
}

func (h *AuthHandler) RefreshActivateCode(c *gin.Context) {
	sub, ok := c.Get("sub")
	if !ok {
		c.JSON(400, gin.H{"code": 0, "msg": "Invalid id!"})
		return
	}

	userID, ok := sub.(uint)
	if !ok {
		c.JSON(400, gin.H{"code": 0, "msg": "Invalid id!"})
		return
	}

	activationCode, expiresAt := utils.GenerateActivationCode()
	if err := h.Repo.UpdateActivationCode(userID, activationCode, expiresAt); err != nil {
		c.JSON(500, gin.H{"code": 0, "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 1, "msg": "Activation code refreshed!", "data": gin.H{
		"activation_code": activationCode,
		"expires_at":      expiresAt,
	}})
}
