package controllers

import (
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

func AddCartItem(context *gin.Context, db *gorm.DB) {
	var requestItem models.CartItem
	if err := context.ShouldBindJSON(&requestItem); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	db.Create(&requestItem)
	context.JSON(http.StatusOK, gin.H{"msg": "Cart item added"})
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
