package database

import (
	"fmt"
	"go-sqlserver-demo/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func Connect() error {
	// Load .env (nếu đang local dev)
	_ = godotenv.Load()

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN environment variable is not set")
	}

	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.User{}, &models.Product{}); err != nil {
		return err
	}

	DB = db
	return nil
}

func LazyConnect() (*gorm.DB, error) {
	_ = godotenv.Load()

	if DB != nil {
		sqlDB, err := DB.DB()
		if err == nil && sqlDB.Ping() == nil {
			return DB, nil
		}
	}

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		return nil, fmt.Errorf("DB_DSN is not set")
	}

	fmt.Println("Connecting to DB with DSN:", dsn)
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	// Migrate models
	if err := db.AutoMigrate(&models.User{}, &models.Product{}, &models.CartItem{}); err != nil {
		return nil, err
	}

	DB = db.Debug()
	return DB, nil
}
