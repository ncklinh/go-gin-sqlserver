package main

import (
	"film-rental/internal/router"
	token "film-rental/internal/token"
	dbOrm "film-rental/pkg/db/gorm"
	dbRaw "film-rental/pkg/db/raw-sql"
	"film-rental/pkg/kafka"
	"film-rental/pkg/mqtt"
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

	dbRaw.InitDB(os.Getenv("DATABASE_URL"))
	dbOrm.Connect()

	kafka.InitKafkaProducer()

	// Start consumers
	go kafka.StartFilmConsumer("Consumer-1")
	go kafka.StartFilmConsumer("Consumer-2")
	go kafka.StartMetricsServer()
	go mqtt.StartMQTTSubscriber()

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
