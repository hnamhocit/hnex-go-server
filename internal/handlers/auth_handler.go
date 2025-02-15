package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hnamhocit/go-learning/internal/dtos"
	"github.com/hnamhocit/go-learning/internal/models"
	"github.com/hnamhocit/go-learning/internal/repositories"
	"github.com/hnamhocit/go-learning/internal/utils"
)

type AuthHandler struct {
	Repo *repositories.UserRepository
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input dtos.RegisterDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	existingUser := h.Repo.GetUserByEmail(input.Email)

	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists!"})
		return
	}

	hashedPassword, err := utils.Hash(input.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user := &models.User{Email: input.Email, DisplayName: input.DisplayName, Password: hashedPassword, Username: strings.Split(input.Email, "@")[1]}
	data, err := h.Repo.CreateUser(user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tokens, err := utils.GenerateTokens(data.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	hashedRefreshToken, err := utils.Hash(tokens.RefreshToken)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.Repo.UpdateUser(data.ID, &models.User{RefreshToken: &hashedRefreshToken})

	c.JSON(http.StatusOK, gin.H{
		"data": tokens,
	})
}

func (r *AuthHandler) Login(c *gin.Context) {
	var input dtos.LoginDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user := r.Repo.GetUserByEmail(input.Email)

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found!"})
		return
	}

	ok, err := utils.Verify(input.Password, user.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password is incorrect!"})
		return
	}

	tokens, err := utils.GenerateTokens(user.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	hashedRefreshToken, err := utils.Hash(tokens.RefreshToken)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	r.Repo.UpdateUser(user.ID, &models.User{RefreshToken: &hashedRefreshToken})

	c.JSON(http.StatusOK, gin.H{
		"data": tokens,
	})
}

func (r *AuthHandler) Logout(c *gin.Context) {
	sub, ok := c.Get("sub")

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized!"})
		return
	}

	r.Repo.UpdateUser(sub.(uint), &models.User{RefreshToken: nil})

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"success": true}})
}

func (r *AuthHandler) Refresh(c *gin.Context) {
	authorization := c.Request.Header.Get("Authorization")
	tokens := strings.Split(authorization, " ")
	refreshToken := tokens[1]

	if refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized!"})
		return
	}

	sub, ok := c.Get("sub")

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized!"})
		return
	}

	user, err := r.Repo.GetUserById(sub.(uint))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	isMatch, err := utils.Verify(refreshToken, *user.RefreshToken)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if !isMatch {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token!"})
		return
	}

	newTokens, err := utils.GenerateTokens(sub.(uint))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	hashedRefreshToken, err := utils.Hash(newTokens.RefreshToken)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	r.Repo.UpdateUser(sub.(uint), &models.User{RefreshToken: &hashedRefreshToken})

	c.JSON(http.StatusOK, gin.H{
		"data": newTokens,
	})
}
