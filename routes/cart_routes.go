package routes

import (
	"go-sqlserver-demo/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterCartRoutes(r *gin.Engine) {
	cartRoutes := r.Group("cart")
	{
		cartRoutes.GET("/items", wrapDB(controllers.GetCartItems))
		cartRoutes.POST("/items", wrapDB(controllers.AddCartItem))
		cartRoutes.DELETE("/items/:id", wrapDB(controllers.DeleteCartItem))
	}
}
