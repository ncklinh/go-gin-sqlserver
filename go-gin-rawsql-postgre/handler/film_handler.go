package handler

import (
	"film-rental/model"
	"film-rental/repository"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetFilms(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1

	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 25
	}

	films, count, err := repository.GetAllFilms(page, limit)
	pageCount := math.Ceil(float64(count) / float64(limit))

	pagination := PaginationMeta{
		Limit:      limit,
		Page:       page,
		TotalCount: count,
		TotalPage:  int(pageCount),
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "Failed to get films", err)
		return
	}
	writeSuccessWithMeta(c, http.StatusOK, "Success", pagination, films)
}

func GetFilmDetail(c *gin.Context) {
	filmId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid id value", err)
		return
	}

	filmDetail, err := repository.GetFilmDetail(filmId)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "Failed to get film detail", err)
		return
	}
	writeSuccess(c, http.StatusOK, "Success", filmDetail)
}

func AddFilm(c *gin.Context) {
	var film model.Film
	if err := c.ShouldBindJSON(&film); err != nil {
		writeError(c, http.StatusBadRequest, "", err)
		return
	}
	id, err := repository.InsertFilm(film)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "Failed to insert film", err)
		return
	}
	writeSuccess(c, http.StatusCreated, "Success", map[string]any{"id": id})
}
