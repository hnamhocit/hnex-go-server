package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hnamhocit/go-learning/internal/repositories"
	"github.com/hnamhocit/go-learning/internal/utils"
)

type MediaHandler struct {
	Repo *repositories.MediaRepository
}

func (h *MediaHandler) UploadFile(c *gin.Context) {
	fileModel, file, err := utils.GetFile(c, nil)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := h.Repo.UploadFile(fileModel)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dst := utils.GenerateUniqueName(file.Filename, data.ID)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *MediaHandler) UploadFiles(c *gin.Context) {
	fileModels, files, err := utils.GetFiles(c, nil)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := h.Repo.UploadFiles(fileModels)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i, file := range files {
		dst := utils.GenerateUniqueName(file.Filename, data[i].ID)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}
