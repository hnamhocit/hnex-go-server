package handlers

import (
	"errors"
	"hnex_server/internal/models"
	"hnex_server/internal/repositories"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	Repo *repositories.UserRepository
}

type ActivateAccountData struct {
	Name           string
	ActivationCode string
	ExpiresAt      string
}

type PasswordResetData struct {
	Name string
	ID   uint
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	idUint, err := strconv.ParseUint(id, 10, 64) // Convert uint64 to uint for compatibility with GetUser method
	if err != nil {
		c.JSON(400, gin.H{"code": 0, "msg": "Invalid ID"})
		return
	}

	user, err := h.Repo.GetUser(uint(idUint))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"code": 0, "msg": "User not found"})
			return
		}

		c.JSON(500, gin.H{"code": 0, "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 1, "msg": "User found", "data": user})
}

func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	sub, ok := c.Get("sub")
	if !ok {
		c.JSON(400, gin.H{"code": 0, "msg": "Invalid sub"})
		return
	}

	user, err := h.Repo.GetUser(sub.(uint))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"code": 0, "msg": "User not found"})
			return
		}

		c.JSON(500, gin.H{"code": 0, "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 1, "msg": "User found", "data": user})
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.Repo.GetUsers()
	if err != nil {
		c.JSON(500, gin.H{"code": 0, "msg": "Internal server error"})
		return
	}

	c.JSON(200, gin.H{"code": 1, "msg": "Users found", "data": users})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"code": 0, "msg": "Invalid request"})
		return
	}

	uintId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"code": 0, "msg": "Invalid ID"})
		return
	}

	err = h.Repo.UpdateUser(uint(uintId), &user)
	if err != nil {
		c.JSON(500, gin.H{"code": 0, "msg": "Internal server error"})
		return
	}

	c.JSON(201, gin.H{"code": 1, "msg": "User updated", "data": id})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	parseUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"code": 0, "msg": "Invalid ID"})
		return
	}

	err = h.Repo.DeleteUser(uint(parseUint))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"code": 0, "msg": "User not found"})
			return
		}

		c.JSON(500, gin.H{"code": 0, "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 1, "msg": "User deleted", "data": gin.H{"id": parseUint}})
}
