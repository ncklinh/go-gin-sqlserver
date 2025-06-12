package controllers

import (
	"errors"
	"fmt"
	"go-sqlserver-demo/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetCartItems(context *gin.Context, db *gorm.DB) {
	userIdStr := context.Query("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid userId"})
		return
	}

	var items []models.CartItem
	db.Preload("Product").Where("user_id = ?", userId).Find(&items)
	context.JSON(http.StatusOK, gin.H{"data": items})
}

func CartItemExisted(db *gorm.DB, item *models.CartItem) (*models.CartItem, error) {
	var existingItem models.CartItem
	err := db.Where(&models.CartItem{UserId: item.UserId, ProductId: item.ProductId}).Preload("Product").First(&existingItem).Error
	return &existingItem, err
}

func AddCartItem(context *gin.Context, db *gorm.DB) {
	var requestItem models.CartItem
	if err := context.ShouldBindJSON(&requestItem); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	existed, err := CartItemExisted(db, &requestItem)

	if err == nil {
		totalCount := requestItem.Quantity + existed.Quantity
		if totalCount > existed.Product.Stock {
			context.JSON(http.StatusNotImplemented, gin.H{"Error": fmt.Sprintf("Only %d item(s) left in stock. You have added %d", existed.Product.Stock, totalCount)})
			return
		}

		db.Model(&models.CartItem{}).Where("product_id = ?", requestItem.ProductId).Update("quantity", (totalCount))
		context.JSON(http.StatusOK, gin.H{"msg": "Cart item added"})
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		var product models.Product
		db.Where("id = ?", requestItem.ProductId).First(&product)
		if requestItem.Quantity > product.Stock {
			context.JSON(http.StatusNotImplemented, gin.H{"Error": fmt.Sprintf("Only %d item(s) left in stock. You have added %d", existed.Product.Stock, requestItem.Quantity)})
			return
		}
		db.Create(&requestItem)
		context.JSON(http.StatusOK, gin.H{"msg": "Cart item added"})
		return
	}
	context.JSON(http.StatusNotImplemented, gin.H{"error": err.Error()})
} 

// TODO: update case amount > stock? merge function?
func UpdateCartItem(context *gin.Context, db *gorm.DB) {
	cartIdStr := context.Param("id")
	cartId, err := strconv.Atoi(cartIdStr)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid cartId"})
		return

	}
	var reqBody models.CartItem
	if err := context.ShouldBindJSON(&reqBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	if err := db.Model(&models.CartItem{}).Where("id = ?", cartId).Update("quantity", reqBody.Quantity).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"msg": "Cart item updated"})
}

func DeleteCartItem(context *gin.Context, db *gorm.DB) {
	cartIdStr := context.Param("id")
	cartId, err := strconv.Atoi(cartIdStr)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid cartId"})
		return

	}
	db.Where("id = ?", cartId).Delete(&models.CartItem{})
	context.JSON(http.StatusOK, gin.H{"msg": "Cart item deleted"})
}

func GetCartSummary(context *gin.Context, db *gorm.DB) {
	userIdStr := context.Query("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid userId"})
		return
	}

	var items []models.CartItem
	db.Preload("Product").Where("user_id = ?", userId).Find(&items)
	totalPrice := 0.0
	totalItems := 0
	for _, item := range items {
		totalItems += item.Quantity
		totalPrice += (item.Product.Price * float64(item.Quantity))
	}
	context.JSON(http.StatusOK, gin.H{"data": models.CartSummary{
		UserId:     userId,
		TotalPrice: totalPrice,
		TotalItems: totalItems,
		Items:      items,
	}})
}
