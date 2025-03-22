package middlewares

import (
	"hnex_server/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func RefreshTokenMiddleware(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		c.JSON(401, gin.H{"code": 0, "msg": "Authorization header is missing!"})
		c.Abort()
		return
	}

	token := strings.TrimPrefix(header, "Bearer ")
	if token == "" {
		c.JSON(401, gin.H{"code": 0, "msg": "Bearer token is missing!"})
		c.Abort()
		return
	}

	tokenClaims, err := utils.ValidateToken(token, "JWT_REFRESH_SECRET")
	if err != nil {
		c.JSON(401, gin.H{"code": 0, "msg": err.Error()})
		c.Abort()
		return
	}

	c.Set("sub", tokenClaims.Sub)
	c.Next()
}

func AccessTokenMiddleware(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		c.JSON(401, gin.H{"code": 0, "msg": "Authorization header is missing!"})
		c.Abort()
		return
	}

	token := strings.TrimPrefix(header, "Bearer ")
	if token == "" {
		c.JSON(401, gin.H{"code": 0, "msg": "Bearer token is missing!"})
		c.Abort()
		return
	}

	tokenClaims, err := utils.ValidateToken(token, "JWT_ACCESS_SECRET")
	if err != nil {
		c.JSON(401, gin.H{"code": 0, "msg": err.Error()})
		c.Abort()
		return
	}

	c.Set("sub", tokenClaims.Sub)
	c.Next()
}
