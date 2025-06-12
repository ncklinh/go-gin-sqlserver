package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"film-rental/internal/repository"
)

func GetFilms(c *gin.Context) {
	films, err := repository.GetAllFilms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get films"})
		return
	}
	c.JSON(http.StatusOK, films)
}
