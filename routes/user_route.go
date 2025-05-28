package routes

import (
	"go-sqlserver-demo/controllers"
	"go-sqlserver-demo/database"
	"net/http"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine) {
	userGroup := r.Group("/users")
	{
		userGroup.GET("/", wrapDB(controllers.GetUsers))
		userGroup.POST("/", wrapDB(controllers.CreateUser))
		userGroup.GET("/:username", wrapDB(controllers.GetUserByUsername))
		userGroup.PUT("/:username", wrapDB(controllers.UpdateUser))
		userGroup.DELETE("/:username", wrapDB(controllers.DeleteUser))
	}
}

func wrapDB(handler func(*gin.Context, *gorm.DB)) gin.HandlerFunc {
	return func(c *gin.Context) {
		db, err := database.LazyConnect()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		handler(c, db)
	}
}
