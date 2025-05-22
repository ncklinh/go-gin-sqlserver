package routes

import (
	"go-sqlserver-demo/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine) {
	userGroup := r.Group("/users")
	{
		userGroup.GET("/", controllers.GetUsers)
		userGroup.POST("/", controllers.CreateUser)
		userGroup.GET("/:username", controllers.GetUserByUsername)
		userGroup.PUT("/:username", controllers.UpdateUser)
		userGroup.DELETE("/:username", controllers.DeleteUser)
	}
}
