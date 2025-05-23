package main

import (
	"go-sqlserver-demo/database"
	"go-sqlserver-demo/routes"

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

	database.Connect()
	routes.RegisterUserRoutes(r)
	r.Run(":8080")
}
