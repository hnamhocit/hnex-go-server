package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hnamhocit/go-learning/internal/handlers"
	"github.com/hnamhocit/go-learning/internal/middlewares"
	"github.com/hnamhocit/go-learning/internal/repositories"
	"gorm.io/gorm"
)

func LoadRoutes(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/api")

	userRepo := repositories.UserRepository{DB: db}
	authHandler := handlers.AuthHandler{Repo: &userRepo}
	userHandler := handlers.UserHandler{Repo: &userRepo}

	{
		auth := api.Group("auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/logout", middlewares.AuthMiddleware, authHandler.Logout)
			auth.GET("/refresh", middlewares.AuthMiddleware, authHandler.Refresh)
		}

		users := api.Group("users")
		{
			users.GET("/profile", middlewares.AuthMiddleware, userHandler.GetProfile)
		}
	}
}
