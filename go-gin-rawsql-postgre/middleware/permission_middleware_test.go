package middleware

import (
	"film-rental/model"
	"film-rental/token"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequirePermission(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	secretKey := "12345678901234567890123456789012"
	jwtMaker, err := token.NewJWTMaker(secretKey)
	require.NoError(t, err)

	tests := []struct {
		name           string
		username       string
		role           string
		permission     string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Admin with film create permission",
			username:       "admin",
			role:           model.RoleAdmin,
			permission:     model.PermissionFilmCreate,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"success"}`,
		},
		{
			name:           "Admin with staff create permission",
			username:       "admin",
			role:           model.RoleAdmin,
			permission:     model.PermissionStaffCreate,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"success"}`,
		},
		{
			name:           "User with film read permission",
			username:       "user",
			role:           model.RoleUser,
			permission:     model.PermissionFilmRead,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"success"}`,
		},
		{
			name:           "User without film create permission",
			username:       "user",
			role:           model.RoleUser,
			permission:     model.PermissionFilmCreate,
			expectedStatus: http.StatusForbidden,
			expectedBody:   `{"error":"Insufficient permissions","message":"You don't have permission to perform this action","required_permission":"film:create","user_role":"user"}`,
		},
		{
			name:           "User without staff create permission",
			username:       "user",
			role:           model.RoleUser,
			permission:     model.PermissionStaffCreate,
			expectedStatus: http.StatusForbidden,
			expectedBody:   `{"error":"Insufficient permissions","message":"You don't have permission to perform this action","required_permission":"staff:create","user_role":"user"}`,
		},
		{
			name:           "User without film delete permission",
			username:       "user",
			role:           model.RoleUser,
			permission:     model.PermissionFilmDelete,
			expectedStatus: http.StatusForbidden,
			expectedBody:   `{"error":"Insufficient permissions","message":"You don't have permission to perform this action","required_permission":"film:delete","user_role":"user"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a valid token
			accessToken, err := jwtMaker.CreateToken(tt.username, tt.role, time.Hour, token.TokenTypeAccessToken)
			require.NoError(t, err)

			// Create router with auth and permission middleware
			router := gin.New()
			authMiddleware := AuthMiddleware(jwtMaker)
			permissionMiddleware := RequirePermission(tt.permission)

			// Add test endpoint
			router.GET("/test", authMiddleware, permissionMiddleware, func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			// Create request
			req, err := http.NewRequest("GET", "/test", nil)
			require.NoError(t, err)
			req.Header.Set(authorizationHeaderKey, "Bearer "+accessToken)

			// Create response recorder
			w := httptest.NewRecorder()

			// Serve request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedBody, strings.TrimSpace(w.Body.String()))
		})
	}
}

func TestRequirePermissionWithoutAuth(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Create router with only permission middleware (no auth)
	router := gin.New()
	permissionMiddleware := RequirePermission(model.PermissionFilmCreate)

	// Add test endpoint
	router.GET("/test", permissionMiddleware, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Create request
	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)

	// Create response recorder
	w := httptest.NewRecorder()

	// Serve request
	router.ServeHTTP(w, req)

	// Assertions - should fail because no auth payload is set
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authorization payload not found")
}

func TestRequirePermissionWithInvalidPayload(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Create router with permission middleware
	router := gin.New()
	permissionMiddleware := RequirePermission(model.PermissionFilmCreate)

	// Add test endpoint with invalid payload
	router.GET("/test", func(c *gin.Context) {
		// Set invalid payload type
		c.Set(authorizationPayloadKey, "invalid_payload")
	}, permissionMiddleware, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Create request
	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)

	// Create response recorder
	w := httptest.NewRecorder()

	// Serve request
	router.ServeHTTP(w, req)

	// Assertions - should fail because payload type is invalid
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid authorization payload")
}

func TestRequirePermissionWithAllPermissions(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	secretKey := "12345678901234567890123456789012"
	jwtMaker, err := token.NewJWTMaker(secretKey)
	require.NoError(t, err)

	// Test all permissions for admin role
	adminPermissions := []string{
		model.PermissionFilmRead,
		model.PermissionFilmCreate,
		model.PermissionFilmUpdate,
		model.PermissionFilmDelete,
		model.PermissionStaffRead,
		model.PermissionStaffCreate,
		model.PermissionStaffUpdate,
		model.PermissionStaffDelete,
		model.PermissionUserRead,
		model.PermissionUserCreate,
		model.PermissionUserUpdate,
		model.PermissionUserDelete,
	}

	for _, permission := range adminPermissions {
		t.Run("Admin_"+permission, func(t *testing.T) {
			// Create admin token
			accessToken, err := jwtMaker.CreateToken("admin", model.RoleAdmin, time.Hour, token.TokenTypeAccessToken)
			require.NoError(t, err)

			// Create router
			router := gin.New()
			authMiddleware := AuthMiddleware(jwtMaker)
			permissionMiddleware := RequirePermission(permission)

			// Add test endpoint
			router.GET("/test", authMiddleware, permissionMiddleware, func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			// Create request
			req, err := http.NewRequest("GET", "/test", nil)
			require.NoError(t, err)
			req.Header.Set(authorizationHeaderKey, "Bearer "+accessToken)

			// Create response recorder
			w := httptest.NewRecorder()

			// Serve request
			router.ServeHTTP(w, req)

			// Assertions - admin should have all permissions
			assert.Equal(t, http.StatusOK, w.Code)
		})
	}

	// Test limited permissions for user role
	userPermissions := []string{
		model.PermissionFilmRead,
		model.PermissionStaffRead,
		model.PermissionUserRead,
	}

	for _, permission := range userPermissions {
		t.Run("User_"+permission, func(t *testing.T) {
			// Create user token
			accessToken, err := jwtMaker.CreateToken("user", model.RoleUser, time.Hour, token.TokenTypeAccessToken)
			require.NoError(t, err)

			// Create router
			router := gin.New()
			authMiddleware := AuthMiddleware(jwtMaker)
			permissionMiddleware := RequirePermission(permission)

			// Add test endpoint
			router.GET("/test", authMiddleware, permissionMiddleware, func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			// Create request
			req, err := http.NewRequest("GET", "/test", nil)
			require.NoError(t, err)
			req.Header.Set(authorizationHeaderKey, "Bearer "+accessToken)

			// Create response recorder
			w := httptest.NewRecorder()

			// Serve request
			router.ServeHTTP(w, req)

			// Assertions - user should have these permissions
			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}
