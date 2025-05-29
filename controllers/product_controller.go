package controllers

import (
	"fmt"
	"go-sqlserver-demo/dtos/request"
	"go-sqlserver-demo/dtos/response"
	"go-sqlserver-demo/models"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetProducts godoc
// @Param page query int true "Page"
// @Param limit query int true "Limit"
// @Summary Get list of products
// @Description Get all products from the database
// @Tags product-apis
// @Produce json
// @Success 200 {array} models.Product
// @Router /products [get]
func GetProducts(context *gin.Context, db *gorm.DB) {
	var products []models.Product

	//paging
	var pagination request.Pagination
	if err := context.ShouldBindQuery(&pagination); err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
		return
	}

	page := pagination.Page
	limit := pagination.Limit
	var totalCount int64
	db.Model(&models.Product{}).Count(&totalCount)

	offset := (page - 1) * limit
	db.Limit(pagination.Limit).Offset(offset).Find(&products)

	resp := response.PaginatedListResponse[[]models.Product]{
		Data:       products,
		Page:       page,
		Limit:      limit,
		TotalCount: totalCount,
		TotalPage:  int(math.Ceil(float64(totalCount) / float64(limit))),
	}
	context.JSON(http.StatusOK, resp)
}

// @Param id path int true "Product ID"
// @Summary Get product detail
// @Tags product-apis
// @Produce json
// @Success 200 {array} models.Product
// @Router /products/detail/{id} [get]
func GetProductDetail(context *gin.Context, db *gorm.DB) {
	var product models.Product
	idStr := context.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product id"})
		return
	}
	if err1 := db.Where("id = ?", id).First(&product).Error; err1 != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	context.JSON(http.StatusOK, product)
}

func CreateProduct(context *gin.Context, db *gorm.DB) {
	var product models.Product
	if err := context.ShouldBindJSON(&product); err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
		return
	}
	db.Create(&product)
	context.JSON(http.StatusOK, gin.H{"msg": "Product added", "prodInfo": product})
}

func CreateMultipleProduct(context *gin.Context, db *gorm.DB) {
	var products []models.Product
	if err := context.ShouldBindJSON(&products); err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var insertFailed []models.Product
	insertedCount := 0
	var reasonInsertFailed string

	for _, p := range products {
		if p.Name == "" || p.Price <= 0 {
			insertFailed = append(insertFailed, p)
			reasonInsertFailed = "Validation failed: name and price are required"
			continue
		}

		if err := db.Create(&p).Error; err != nil {
			insertFailed = append(insertFailed, p)
			fmt.Println(err.Error())
		} else {
			insertedCount++
		}
	}

	context.JSON(http.StatusOK, gin.H{"msg": "Task done", "inserted": insertedCount, "failed": gin.H{"items": insertFailed, "reason": reasonInsertFailed}})
}

func UpdateProduct(context *gin.Context, db *gorm.DB) {
	var product models.Product
	id := context.Param("id")

	if err := db.Where("id = ?", id).First(&product).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if err := context.ShouldBindJSON(&product); err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
		return
	}
	db.Save(&product)
	context.JSON(http.StatusOK, gin.H{"msg": "Product updated", "prodInfo": product})
}

func DeleteProduct(context *gin.Context, db *gorm.DB) {
	var product models.Product
	if err := context.ShouldBindJSON(&product); err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
		return
	}
	db.Delete(&product)
	context.JSON(http.StatusOK, gin.H{"msg": "Product deleted", "prodInfo": product})
}
