package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hnamhocit/go-learning/internal/repositories"
)

type UserHandler struct {
	Repo *repositories.UserRepository
}

func (r *UserHandler) GetProfile(c *gin.Context) {
	id, ok := c.Get("sub")

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID!"})
		return
	}

	user, err := r.Repo.GetUserById(id.(uint))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}
