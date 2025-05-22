package database

import (
	"go-sqlserver-demo/models"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := "sqlserver://sa:YourStrong@Passw0rd@localhost:1433?database=master"
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// AutoMigrate
	db.AutoMigrate(&models.User{})

	DB = db
}
