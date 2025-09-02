package database

import (
	"fmt"
	"log"
	"os"

	"blog/pkg/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err any

func InitDB() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort)
	fmt.Println(dsn)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}

	err = DB.AutoMigrate(&models.AuthKey{}, &models.Tree{}, &models.Content{}, &models.Meta{})
	if err != nil {
		fmt.Println("Error migrating database:", err)
	}
}
