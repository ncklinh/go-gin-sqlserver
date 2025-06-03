package routes

import (
	"go-sqlserver-demo/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(r *gin.Engine) {
	orderRoutes := r.Group("orders")
	{
		orderRoutes.POST("/create", wrapDB(controllers.CreateOrder))
	}
}
