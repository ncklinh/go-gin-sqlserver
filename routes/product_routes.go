package routes

import (
	"go-sqlserver-demo/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(r *gin.Engine) {
	productRoutes := r.Group("products")
	{
		productRoutes.GET("", wrapDB(controllers.GetProducts))
		productRoutes.GET("/detail/:id", wrapDB(controllers.GetProductDetail))
		productRoutes.POST("/create", wrapDB(controllers.CreateProduct))
		productRoutes.POST("/bulk", wrapDB(controllers.CreateMultipleProduct))
		productRoutes.PUT("/update/:id", wrapDB(controllers.UpdateProduct))
		productRoutes.DELETE("/:id", wrapDB(controllers.DeleteProduct))
	}
}
