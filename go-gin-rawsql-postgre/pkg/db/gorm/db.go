package db

import (
	"log"

	monitoringModel "film-rental/pkg/monitoring/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func Connect(dsn string) error {
	// _ = godotenv.Load()

	// dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN environment variable is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	if err := db.AutoMigrate(&monitoringModel.EventLog{}); err != nil {
		return err
	}

	DB = db

	return nil
}
