package main

import (
	"film-rental/db"
	"film-rental/kafka"
	"film-rental/router"
	token "film-rental/token"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Try to load .env file (for local development), but don't fail if it doesn't exist
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	db.InitDB(os.Getenv("DATABASE_URL"))
	go kafka.StartRentalConsumer()

	r := gin.Default()
	jwtMaker, err := token.NewJWTMaker(os.Getenv("TOKEN_SYMMETRIC_KEY"))
	if err != nil {
		log.Fatalf("Failed to create JWT maker: %v", err)
	}
	router.RegisterRoutes(r, jwtMaker)
	r.Use(CORSMiddleware())

	r.Run(":8080")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
