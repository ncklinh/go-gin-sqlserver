package handler

import (
	"film-rental/internal/repository"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetFilms(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		writeError(c, http.StatusInternalServerError, "invalid page value", err)
		return

	}
	limit, err1 := strconv.Atoi(c.Query("limit"))
	if err1 != nil {
		writeError(c, http.StatusInternalServerError, "invalid limit value", err)
		return

	}

	films, count, err2 := repository.GetAllFilms(page, limit)
	pageCount := math.Ceil(float64(count) / float64(limit))

	pagination := PaginationMeta{
		Limit:      limit,
		Page:       page,
		TotalCount: count,
		TotalPage:  int(pageCount),
	}
	if err2 != nil {
		writeError(c, http.StatusInternalServerError, "Failed to get films", err)
		return
	}
	writeSuccessWithMeta(c, http.StatusOK, "Success", pagination, films)
}
