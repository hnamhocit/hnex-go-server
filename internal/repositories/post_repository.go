package repositories

import (
	"github.com/hnamhocit/go-learning/internal/models"
	"gorm.io/gorm"
)

type PostRepository struct {
	DB *gorm.DB
}

func (h *PostRepository) CreatePost(post *models.Post) (*models.Post, error) {
	if err := h.DB.Create(post).Error; err != nil {
		return nil, err
	}

	return post, nil
}

func (h *PostRepository) GetPostById(id int) (*models.Post, error) {
	var post models.Post
	if err := h.DB.Where("id = ?", id).First(&post).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

func (h *PostRepository) GetPosts() ([]models.Post, error) {
	var posts []models.Post

	if err := h.DB.Preload("Author").Preload("Media").Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}
