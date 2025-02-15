package utils

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hnamhocit/go-learning/internal/models"
)

func GetFile(c *gin.Context, postId *uint) (*models.Media, *multipart.FileHeader, error) {
	sub, ok := c.Get("sub")

	if !ok {
		return nil, nil, errors.New("unauthorized")
	}

	file, err := c.FormFile("file")

	if err != nil {
		return nil, nil, err
	}

	dst := "./uploads/" + file.Filename

	upload := &models.Media{
		Name:        file.Filename,
		Size:        file.Size,
		Path:        dst,
		ContentType: file.Header.Get("Content-Type"),
		UserID:      sub.(uint),
		PostId:      postId,
	}

	return upload, file, nil
}

func GetFiles(c *gin.Context, postId *uint) ([]*models.Media, []*multipart.FileHeader, error) {
	sub, ok := c.Get("sub")

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unauthorized!"})
		return nil, nil, errors.New("unauthorized")
	}

	form, err := c.MultipartForm()

	if err != nil {
		return nil, nil, err
	}

	files := form.File["files"]

	var uploads []*models.Media
	{
	}

	for _, file := range files {
		dst := "./uploads/" + file.Filename

		upload := &models.Media{
			Name:        file.Filename,
			Size:        file.Size,
			Path:        dst,
			ContentType: file.Header.Get("Content-Type"),
			UserID:      sub.(uint),
			PostId:      postId,
		}

		uploads = append(uploads, upload)

	}

	return uploads, files, nil
}

func GenerateUniqueName(fileName string, id uint) string {
	fileSplit := strings.Split(fileName, ".")
	name, ext := fileSplit[0], fileSplit[len(fileSplit)-1]
	dst := "uploads/" + name + "-" + fmt.Sprint(id) + "." + ext
	return dst
}
