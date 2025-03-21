package routes

import (
	"hnex_server/internal/config"
	"hnex_server/internal/handlers"
	"hnex_server/internal/middlewares"
	"hnex_server/internal/repositories"

	"github.com/gin-gonic/gin"
)

func InitRoutes(c *gin.Engine) {
	db := config.LoadConfig()

	api := c.Group("api")
	{
		auth := api.Group("auth")
		{
			authRepo := repositories.AuthRepository{DB: db}
			authHandler := handlers.AuthHandler{Repo: &authRepo}

			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.GET("/logout", middlewares.AccessTokenMiddleware, authHandler.Logout)
			auth.GET("/refresh", middlewares.RefreshTokenMiddleware, authHandler.Refresh)
		}

		users := api.Group("users")
		{
			userRepo := repositories.UserRepository{DB: db}
			userHandler := handlers.UserHandler{Repo: &userRepo}

			users.POST("/me", middlewares.AccessTokenMiddleware, userHandler.GetCurrentUser)
			users.GET("/:id", userHandler.GetUser)
		}
	}
}
