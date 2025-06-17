package main

import (
	"film-rental/db"
	"film-rental/router"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.InitDB(os.Getenv("DATABASE_URL"))

	r := gin.Default()
	router.RegisterRoutes(r)

	r.Run(":8080")
}
