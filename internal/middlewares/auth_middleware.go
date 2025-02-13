package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hnamhocit/go-learning/internal/utils"
)

func AuthMiddleware(c *gin.Context) {
	tokenString, err := c.Cookie("access_token")

	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.VerifyToken(tokenString, "JWT_ACCESS_SECRET")

	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token claims!"})
		return
	}

	c.Set("sub", uint(claims["sub"].(float64)))

	c.Next()
}
