package handlers

import (
	"hnex_server/internal/repositories"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Repo *repositories.UserRepository
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	idUint, err := strconv.ParseUint(id, 10, 64)// Convert uint64 to uint for compatibility with GetUser method
	if err != nil {
		c.JSON(400, gin.H{"code": 0, "msg": "Invalid ID"})
		return
	}

	user, err := h.Repo.GetUser(uint(idUint))
	if err != nil {
		c.JSON(404, gin.H{"code": 0, "msg": "User not found"})
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
		c.JSON(404, gin.H{"code": 0, "msg": "User not found"})
		return
	}

	c.JSON(200, gin.H{"code": 1, "msg": "User found", "data": user})
}
