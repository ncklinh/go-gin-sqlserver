package routes

import (
	"go-sqlserver-demo/database"
	"go-sqlserver-demo/models"

	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine) {
	r.GET("/users", func(c *gin.Context) {
		db, err := database.LazyConnect()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var users []models.User
		if err := db.Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, users)
	})
}

func GetUsers(c *gin.Context) {
	db, err := database.LazyConnect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB connect failed: " + err.Error()})
		return
	}

	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Query failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}
