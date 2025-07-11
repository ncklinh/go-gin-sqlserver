package response

import (
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type SuccessResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type MetaResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	PaginationMeta
	Data interface{} `json:"data"`
}

type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalCount int `json:"total_count"`
	TotalPage  int `json:"total_page"`
}

func WriteSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, SuccessResponse{
		Code:    statusCode,
		Message: message,
		Data:    data,
	})
}

func WriteError(c *gin.Context, statusCode int, message string, err error) {
	res := ErrorResponse{
		Code:    statusCode,
		Message: message,
	}
	if err != nil {
		res.Error = err.Error()
	}
	c.JSON(statusCode, res)
}

func WriteSuccessWithMeta(c *gin.Context, statusCode int, message string, paginationMeta PaginationMeta, data interface{}) {
	c.JSON(statusCode, MetaResponse{
		Code:           statusCode,
		Message:        message,
		PaginationMeta: paginationMeta,
		Data:           data,
	})
}
