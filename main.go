package main

import (
    "go-sqlserver-demo/routes"
	"go-sqlserver-demo/database"

	"github.com/gin-gonic/gin"
)


func main() {
    r := gin.Default()
    database.Connect()
    routes.RegisterUserRoutes(r)
    r.Run(":8080")
}