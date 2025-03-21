package handlers

import (
	"hnex_server/internal/models"
	"hnex_server/internal/repositories"
	"hnex_server/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	Repo *repositories.AuthRepository
}

type LoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=100"`
}

type RegisterDTO struct {
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
		Email:    input.Email,
		Password: hashedPassword,
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

	updateErr := h.Repo.UpdateRefreshToken(newUser.ID, &tokens.RefreshToken)
	if updateErr != nil {
		c.JSON(500, gin.H{"code": 0, "msg": updateErr.Error()})
		return
	}

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
