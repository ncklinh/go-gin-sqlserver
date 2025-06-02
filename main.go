package main

import (
	"go-sqlserver-demo/routes"
	"os"

	"github.com/gin-gonic/gin"

	_ "go-sqlserver-demo/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	r := gin.Default()

	// Health check without DB
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Register routes using DB
	routes.RegisterUserRoutes(r)
	routes.RegisterProductRoutes(r)
	routes.RegisterCartRoutes(r)
	routes.RegisterPaymentRoutes(r)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
