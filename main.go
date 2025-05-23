package main

import (
	"go-sqlserver-demo/routes"

	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	r := gin.Default()

	// Health check không cần DB
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Đăng ký các route cần DB
	routes.RegisterUserRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
