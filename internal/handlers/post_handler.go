package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hnamhocit/go-learning/internal/models"
	"github.com/hnamhocit/go-learning/internal/repositories"
	"github.com/hnamhocit/go-learning/internal/utils"
)

type PostHandler struct {
	Repo      *repositories.PostRepository
	MediaRepo *repositories.MediaRepository
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	sub, ok := c.Get("sub")

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	content := c.PostForm("content")

	data, err := h.Repo.CreatePost(&models.Post{Content: content, AuthorId: sub.(uint)})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fileModels, files, err := utils.GetFiles(c, &data.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	media, err := h.MediaRepo.UploadFiles(fileModels)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i, file := range files {

		dst := utils.GenerateUniqueName(file.Filename, media[i].ID)
		log.Println(dst)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *PostHandler) GetPosts(c *gin.Context) {
	posts, err := h.Repo.GetPosts()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": posts})
}
