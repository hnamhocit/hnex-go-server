package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hnamhocit/go-learning/config"
	"github.com/hnamhocit/go-learning/internal/routes"
)

func main() {
	r := gin.Default()
	db := config.Load(r)
	PORT := os.Getenv("PORT")

	routes.LoadRoutes(r, db)

	r.Run(":" + PORT)
}
