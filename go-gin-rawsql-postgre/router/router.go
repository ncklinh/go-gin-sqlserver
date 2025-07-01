package router

import (
	"film-rental/handler"
	"film-rental/middleware"
	"film-rental/model"
	"film-rental/token"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, jwtMaker *token.JWTMaker) {
	// Public routes (no authentication required)
	filmRoutes := r.Group("films")
	{
		filmRoutes.GET("", handler.GetFilms)          // Public read access
		filmRoutes.GET("/:id", handler.GetFilmDetail) // Public read access
	}

	// Protected routes (authentication required)
	authMiddleware := middleware.AuthMiddleware(jwtMaker)

	// Film management (requires specific permissions)
	filmProtectedRoutes := r.Group("films").Use(authMiddleware)
	{
		filmProtectedRoutes.POST("", middleware.RequirePermission(model.PermissionFilmCreate), handler.AddFilm)
		filmProtectedRoutes.PUT("/:id", middleware.RequirePermission(model.PermissionFilmUpdate), handler.UpdateFilm)
		filmProtectedRoutes.DELETE("/:id", middleware.RequirePermission(model.PermissionFilmDelete), handler.DeleteFilm)
	}

	// Staff management (admin only)
	staffRoutes := r.Group("/staff").Use(authMiddleware)
	{
		staffRoutes.GET("", middleware.RequirePermission(model.PermissionStaffRead), handler.GetStaffs)
		staffRoutes.POST("", middleware.RequirePermission(model.PermissionStaffCreate), handler.AddStaff)
	}

	// User authentication
	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/login", handler.LoginStaff(jwtMaker))
		userRoutes.POST("/refresh", handler.RefreshToken(jwtMaker))
	}
}
