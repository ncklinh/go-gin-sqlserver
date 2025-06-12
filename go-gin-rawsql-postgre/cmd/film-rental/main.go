package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
	"film-rental/pkg/db"
	"film-rental/internal/handler"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.InitDB(os.Getenv("DATABASE_URL"))

	r := gin.Default()
	r.GET("/films", handler.GetFilms)

	r.Run(":8080")
}
