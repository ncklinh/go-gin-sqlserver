package database

import (
	"go-sqlserver-demo/models"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var (
	DB   *gorm.DB
	once sync.Once
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

	if err := db.AutoMigrate(&models.User{}); err != nil {
		return err
	}

	DB = db
	return nil
}

func LazyConnect() (*gorm.DB, error) {
	var err error
	once.Do(func() {
		err = Connect()
	})
	return DB, err
}
