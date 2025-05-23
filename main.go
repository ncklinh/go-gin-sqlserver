package main

import (
	"go-sqlserver-demo/database"
	"go-sqlserver-demo/routes"

	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		err := database.Connect()
		if err == nil {
			fmt.Println("Connected to DB")
			break
		}
		fmt.Printf("Failed to connect to DB. Retrying in 5 seconds... (%d/%d)\n", i+1, maxRetries)
		time.Sleep(5 * time.Second)
		if i == maxRetries-1 {
			panic("Failed to connect to DB after retries")
		}
	}

	routes.RegisterUserRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default nếu biến PORT không có
	}

	r.Run(":" + port)
}
