package handler

import (
	"database/sql"
	"film-rental/kafka"
	"film-rental/model"
	"film-rental/repository"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// validateFilmFields validates all required film fields
func validateFilmFields(film model.Film) (string, error) {
	if film.Title == "" {
		return "Title is required", nil
	}
	if film.Description == "" {
		return "Description is required", nil
	}
	if film.ReleaseYear < 1888 || film.ReleaseYear > 2030 {
		return "Release year must be between 1888 and 2030", nil
	}
	if film.RentalDuration <= 0 {
		return "Rental duration must be greater than 0", nil
	}
	if film.RentalRate < 0 {
		return "Rental rate cannot be negative", nil
	}
	if film.Length <= 0 {
		return "Length must be greater than 0", nil
	}
	if film.ReplacementCost < 0 {
		return "Replacement cost cannot be negative", nil
	}
	if film.LanguageId <= 0 {
		return "Language ID must be greater than 0", nil
	}
	return "", nil
}

func GetFilms(c *gin.Context) {
	log.Println("GET /films called")
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

	// ✅ Publish Kafka event (non-blocking)
	go kafka.PublishRentalEvent(fmt.Sprintf("Films list viewed at %v", time.Now()))

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
		writeError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	message, _ := validateFilmFields(film)
	if message != "" {
		writeError(c, http.StatusBadRequest, message, nil)
		return
	}

	id, err := repository.InsertFilm(film)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "Failed to insert film", err)
		return
	}
	writeSuccess(c, http.StatusCreated, "Success", map[string]any{"id": id})
}

func UpdateFilm(c *gin.Context) {
	filmId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "Invalid film ID", err)
		return
	}

	var film model.Film
	if err := c.ShouldBindJSON(&film); err != nil {
		writeError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	message, _ := validateFilmFields(film)
	if message != "" {
		writeError(c, http.StatusBadRequest, message, nil)
		return
	}

	err = repository.UpdateFilm(filmId, film)
	if err != nil {
		if err == sql.ErrNoRows {
			writeError(c, http.StatusNotFound, "Film not found", err)
			return
		}
		writeError(c, http.StatusInternalServerError, "Failed to update film", err)
		return
	}

	writeSuccess(c, http.StatusOK, "Film updated successfully", nil)
}

func DeleteFilm(c *gin.Context) {
	filmId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "Invalid film ID", err)
		return
	}

	err = repository.DeleteFilm(filmId)
	if err != nil {
		if err == sql.ErrNoRows {
			writeError(c, http.StatusNotFound, "Film not found", err)
			return
		}
		writeError(c, http.StatusInternalServerError, "Failed to delete film", err)
		return
	}

	writeSuccess(c, http.StatusOK, "Film deleted successfully", nil)
}
