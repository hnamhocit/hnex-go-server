package config

import (
	"fmt"
	"hnex_server/internal/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func LoadEnv() {
	err := godotenv.Load()

	if err != nil {
		panic(err.Error())
	}

	log.Println("[ENV] load environments variables successfully.")
}

func LoadDBConfig() *gorm.DB {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.User{})

	log.Println("[DB] migrate database successfully.")
	log.Println("[DB] database connection established successfully.")

	return db
}

func LoadConfig() *gorm.DB {
	LoadEnv()
	db := LoadDBConfig()

	return db
}
