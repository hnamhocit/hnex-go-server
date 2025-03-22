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

	authRepo := repositories.AuthRepository{DB: db}
	userRepo := repositories.UserRepository{DB: db}
	profileRepo := repositories.ProfileRepository{DB: db}

	api := c.Group("api")
	{
		auth := api.Group("auth")
		{
			authHandler := handlers.AuthHandler{Repo: &authRepo, UserRepo: &userRepo, ProfileRepo: &profileRepo}

			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.GET("/logout", middlewares.AccessTokenMiddleware, authHandler.Logout)
			auth.GET("/refresh", middlewares.RefreshTokenMiddleware, authHandler.Refresh)
			auth.PATCH("/activate", middlewares.AccessTokenMiddleware, authHandler.ActivateAccount)
			auth.GET("/activate/refresh", middlewares.AccessTokenMiddleware, authHandler.RefreshActivateCode)
		}

		users := api.Group("users")
		{
			userHandler := handlers.UserHandler{Repo: &userRepo}

			users.GET("/me", middlewares.AccessTokenMiddleware, userHandler.GetCurrentUser)
			users.GET("/:id", userHandler.GetUser)
			users.PATCH("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}
	}
}
