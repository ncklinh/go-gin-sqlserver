package handler

import (
	"database/sql"
	"encoding/json"
	"film-rental/internal/film/model"
	"film-rental/internal/film/repository"
	"film-rental/pkg/kafka"
	"film-rental/pkg/redis"
	"film-rental/pkg/response"
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

	pagination := response.PaginationMeta{
		Limit:      limit,
		Page:       page,
		TotalCount: count,
		TotalPage:  int(pageCount),
	}
	if err != nil {
		response.WriteError(c, http.StatusInternalServerError, "Failed to get films", err)
		return
	}

	response.WriteSuccessWithMeta(c, http.StatusOK, "Success", pagination, films)
}

func GetFilmDetail(c *gin.Context) {
	filmId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid id value", err)
		return
	}
	var filmDetail *model.Film

	cacheKey := fmt.Sprintf("film:%d", filmId)
	cached, err := redis.Rdb.Get(redis.Ctx, cacheKey).Result()
	if err == nil {
		err = json.Unmarshal([]byte(cached), &filmDetail)
		if err == nil {
			response.WriteSuccess(c, http.StatusOK, "Success", filmDetail)
			return
		}
	}

	filmDetail, err = repository.GetFilmDetail(filmId)
	if err != nil {
		response.WriteError(c, http.StatusInternalServerError, "Failed to get film detail", err)
		return
	}
	jsonData, _ := json.Marshal(filmDetail)
	redis.Rdb.Set(redis.Ctx, cacheKey, jsonData, 5*time.Minute)

	response.WriteSuccess(c, http.StatusOK, "Success", filmDetail)
}

func AddFilm(c *gin.Context) {
	var film model.Film
	if err := c.ShouldBindJSON(&film); err != nil {
		response.WriteError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	message, _ := validateFilmFields(film)
	if message != "" {
		response.WriteError(c, http.StatusBadRequest, message, nil)
		return
	}

	id, err := repository.InsertFilm(film)
	if err != nil {
		response.WriteError(c, http.StatusInternalServerError, "Failed to insert film", err)
		return
	}
	// Publish Kafka event (non-blocking)
	msg := fmt.Sprintf("Film added: %s at %s", film.Title, time.Now())
	kafka.PublishFilmEvent(msg)

	response.WriteSuccess(c, http.StatusCreated, "Success", map[string]any{"id": id})
}

func UpdateFilm(c *gin.Context) {
	filmId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "Invalid film ID", err)
		return
	}

	var film model.Film
	if err := c.ShouldBindJSON(&film); err != nil {
		response.WriteError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	message, _ := validateFilmFields(film)
	if message != "" {
		response.WriteError(c, http.StatusBadRequest, message, nil)
		return
	}

	err = repository.UpdateFilm(filmId, film)
	if err != nil {
		if err == sql.ErrNoRows {
			response.WriteError(c, http.StatusNotFound, "Film not found", err)
			return
		}
		response.WriteError(c, http.StatusInternalServerError, "Failed to update film", err)
		return
	}
	// Publish Kafka event (non-blocking)
	msg := fmt.Sprintf("Film updated: ID=%d at %s", film.ID, time.Now())
	kafka.PublishFilmEvent(msg)
	response.WriteSuccess(c, http.StatusOK, "Film updated successfully", nil)
}

func DeleteFilm(c *gin.Context) {
	filmId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "Invalid film ID", err)
		return
	}

	err = repository.DeleteFilm(filmId)
	if err != nil {
		if err == sql.ErrNoRows {
			response.WriteError(c, http.StatusNotFound, "Film not found", err)
			return
		}
		response.WriteError(c, http.StatusInternalServerError, "Failed to delete film", err)
		return
	}

	response.WriteSuccess(c, http.StatusOK, "Film deleted successfully", nil)
}
