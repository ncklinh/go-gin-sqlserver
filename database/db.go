package database

import (
	"log"
	"os"
	"go-sqlserver-demo/models"

    "github.com/joho/godotenv"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
    // Load .env file (chỉ trong dev, prod thì thường Cloud Run sẽ cung cấp biến môi trường)
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, reading configuration from environment variables")
    }

    dsn := os.Getenv("DB_DSN")
    if dsn == "" {
        log.Fatal("DB_DSN environment variable is not set")
    }

    db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Failed to connect to database: " + err.Error())
    }

    db.AutoMigrate(&models.User{})

    DB = db
}
