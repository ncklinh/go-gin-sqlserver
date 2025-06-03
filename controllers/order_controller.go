package controllers

import (
	"go-sqlserver-demo/dtos/request"
	"go-sqlserver-demo/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateOrder(context *gin.Context, db *gorm.DB) {
	// temporary // TODO: replace with user id from jwt
	var createOrderReq request.CreateOrderRequest
	if err := context.ShouldBindJSON(&createOrderReq); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid userId"})
		return
	}

	var cartItems []models.CartItem
	db.Preload("Product").Where("user_id = ?", createOrderReq.UserId).Find(&cartItems)

	var items []models.OrderItem
	var totalAmount float64

	for _, cartItem := range cartItems {
		var product models.Product
		if err := db.First(&product, cartItem.ProductId).Error; err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid product ID"})
			return
		}

		items = append(items, models.OrderItem{
			ProductId: cartItem.ProductId,
			Product:   cartItem.Product,
			Quantity:  cartItem.Quantity,
			UnitPrice: product.Price,
		})
		totalAmount += (product.Price * float64(cartItem.Quantity))
	}

	order := models.Order{
		Items:         items,
		TotalAmount:   totalAmount,
		PaymentStatus: models.PaymentWaiting,
	}

	if err := db.Create(&order).Error; err != nil {
		log.Printf("❌ Failed to create order: %v", err)
	}

	for i := range cartItems {
		item := cartItems[i]
		// db.Where("id = ?", item.Id).Delete(&item)
		db.Delete(&models.CartItem{}, item.Id)

	}
	context.JSON(http.StatusOK, gin.H{"msg": "order created"})
}
