package router

import (
	"film-rental/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	filmRoutes := r.Group("films")
	{
		filmRoutes.GET("", handler.GetFilms)
		filmRoutes.GET("/:id", handler.GetFilmDetail)
		filmRoutes.POST("", handler.AddFilm)
	}
}
