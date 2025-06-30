package router

import (
	"film-rental/handler"
	"film-rental/middleware"
	"film-rental/token"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, jwtMaker *token.JWTMaker) {
	filmRoutes := r.Group("films")
	{
		filmRoutes.GET("", handler.GetFilms)
		filmRoutes.GET("/:id", handler.GetFilmDetail)
		filmRoutes.POST("", handler.AddFilm)
	}

	accountRoutes := r.Group("/accounts").Use(middleware.AuthMiddleware(jwtMaker))
	{
		accountRoutes.GET("", handler.GetStaffs)
		accountRoutes.POST("", handler.AddStaff)
	}

	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/login", handler.LoginStaff(jwtMaker))
	}
}
