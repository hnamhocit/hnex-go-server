package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hnamhocit/go-learning/internal/utils"
)

func AuthMiddleware(c *gin.Context) {
	authorization := c.Request.Header.Get("Authorization")
	tokens := strings.Split(authorization, " ")

	accessToken := tokens[1]

	if accessToken == "" {
		c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized!"})
		return
	}

	token, err := utils.VerifyToken(accessToken, "JWT_ACCESS_SECRET")
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token claims!"})
		return
	}

	sub := uint(claims["sub"].(float64))

	c.Set("sub", sub)
	c.Next()
}
