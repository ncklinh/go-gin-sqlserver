package controllers

import (
	"go-sqlserver-demo/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPaymentMethods(context *gin.Context, db *gorm.DB) {
	var paymentMethods []models.PaymentMethod
	if err := db.Find(&paymentMethods).Error; err != nil {
		context.JSON(http.StatusBadGateway, err)
		return
	}
	context.JSON(http.StatusOK, gin.H{"data": paymentMethods})
}

func CreatePaymentMethod(context *gin.Context, db *gorm.DB) {
	var paymentMethod models.PaymentMethod
	if err := context.ShouldBindJSON(&paymentMethod); err != nil {
		context.JSON(http.StatusBadRequest, err)
		return
	}
	db.Where("code = ?", paymentMethod.Code).FirstOrCreate(&paymentMethod)
	context.JSON(http.StatusOK, gin.H{"msg": "Method added"})
}

func MakePayment(context *gin.Context, db *gorm.DB) {
	var paymentTransaction models.PaymentTransaction

	if err := context.ShouldBindJSON(&paymentTransaction); err != nil {
		context.JSON(http.StatusBadRequest, err)
		return
	}
	db.Create(&paymentTransaction)
	context.JSON(http.StatusOK, gin.H{"msg": "Paid"})

}
