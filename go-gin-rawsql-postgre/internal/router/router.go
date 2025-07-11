package router

import (
	filmHandler "film-rental/internal/film/handler"
	staffHandler "film-rental/internal/staff/handler"
	"film-rental/internal/token"
	tokenModel "film-rental/internal/token/model"
	"film-rental/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, jwtMaker *token.JWTMaker) {
	// Public routes (no authentication required)
	filmRoutes := r.Group("films")
	{
		filmRoutes.GET("", filmHandler.GetFilms)
		filmRoutes.GET("/:id", filmHandler.GetFilmDetail)
	}

	// Protected routes (authentication required)
	authMiddleware := middleware.AuthMiddleware(jwtMaker)

	filmProtectedRoutes := r.Group("films").Use(authMiddleware)
	{
		filmProtectedRoutes.POST("", middleware.RequirePermission(tokenModel.PermissionFilmCreate), filmHandler.AddFilm)
		filmProtectedRoutes.PUT("/:id", middleware.RequirePermission(tokenModel.PermissionFilmUpdate), filmHandler.UpdateFilm)
		filmProtectedRoutes.DELETE("/:id", middleware.RequirePermission(tokenModel.PermissionFilmDelete), filmHandler.DeleteFilm)
	}

	staffRoutes := r.Group("/staff").Use(authMiddleware)
	{
		staffRoutes.GET("", middleware.RequirePermission(tokenModel.PermissionStaffRead), staffHandler.GetStaffs)
		staffRoutes.POST("", middleware.RequirePermission(tokenModel.PermissionStaffCreate), staffHandler.AddStaff)
	}

	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/login", staffHandler.LoginStaff(jwtMaker))
		userRoutes.POST("/refresh", staffHandler.RefreshToken(jwtMaker))
	}
}
