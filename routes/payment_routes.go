package routes

import (
	"go-sqlserver-demo/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterPaymentRoutes(r *gin.Engine) {
	paymentRoutes := r.Group("payments")
	{
		paymentRoutes.GET("/", wrapDB(controllers.GetPaymentMethods))
		paymentRoutes.POST("/add-method", wrapDB(controllers.CreatePaymentMethod))
		paymentRoutes.POST("/make", wrapDB(controllers.MakePayment))
	}
}
