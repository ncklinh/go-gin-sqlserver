package handler

import (
	"bytes"
	"encoding/json"
	"film-rental/internal/token"
	"film-rental/internal/token/model"
	dbRaw "film-rental/pkg/db/raw-sql"
	"film-rental/pkg/middleware"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRouter() (*gin.Engine, *token.JWTMaker) {
	gin.SetMode(gin.TestMode)

	// Setup JWT maker
	secretKey := "12345678901234567890123456789012"
	jwtMaker, _ := token.NewJWTMaker(secretKey)

	// Setup router
	router := gin.New()

	// Public routes
	filmRoutes := router.Group("films")
	{
		filmRoutes.GET("", GetFilms)
		filmRoutes.GET("/:id", GetFilmDetail)
	}

	return router, jwtMaker
}

func setupProtectedTestRouter() (*gin.Engine, *token.JWTMaker) {
	gin.SetMode(gin.TestMode)

	// Setup JWT maker
	secretKey := "12345678901234567890123456789012"
	jwtMaker, _ := token.NewJWTMaker(secretKey)

	// Setup router with auth middleware
	router := gin.New()
	authMiddleware := middleware.AuthMiddleware(jwtMaker)

	// Protected routes
	filmProtectedRoutes := router.Group("films").Use(authMiddleware)
	{
		filmProtectedRoutes.POST("", middleware.RequirePermission(model.PermissionFilmCreate), AddFilm)
		filmProtectedRoutes.PUT("/:id", middleware.RequirePermission(model.PermissionFilmUpdate), UpdateFilm)
		filmProtectedRoutes.DELETE("/:id", middleware.RequirePermission(model.PermissionFilmDelete), DeleteFilm)
	}

	return router, jwtMaker
}

// setupTestDB initializes a mock database for testing
func setupTestDB(t *testing.T) (sqlmock.Sqlmock, func()) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}

	// Set the mock database to the global db.DB
	dbRaw.DB = mockDB

	// Return cleanup function
	cleanup := func() {
		mockDB.Close()
	}

	return mock, cleanup
}

func TestGetFilms(t *testing.T) {
	mock, cleanup := setupTestDB(t)
	defer cleanup()

	// Set up mock expectations for GetAllFilms
	// First query: SELECT film_id, title, description, release_year, rental_duration, rental_rate, length, replacement_cost, rating, last_update, language_id FROM film ORDER BY film_id DESC LIMIT $1 OFFSET $2
	rows := sqlmock.NewRows([]string{"film_id", "title", "description", "release_year", "rental_duration", "rental_rate", "length", "replacement_cost", "rating", "last_update", "language_id"}).
		AddRow(1, "Test Film 1", "Test Description 1", 2020, 3, 2.99, 120, 19.99, "PG", sqlmock.AnyArg(), 1).
		AddRow(2, "Test Film 2", "Test Description 2", 2021, 3, 3.99, 130, 24.99, "PG-13", sqlmock.AnyArg(), 1)

	mock.ExpectQuery("SELECT (.+) FROM film ORDER BY film_id DESC LIMIT").WithArgs(25, 0).WillReturnRows(rows)

	// Second query: SELECT COUNT (*) FROM film
	mock.ExpectQuery("SELECT COUNT").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	router, _ := setupTestRouter()

	// Test successful request
	req, err := http.NewRequest("GET", "/films", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response, "data")
}

func TestGetFilmDetail(t *testing.T) {
	mock, cleanup := setupTestDB(t)
	defer cleanup()

	router, _ := setupTestRouter()

	tests := []struct {
		name           string
		filmID         string
		expectedStatus int
		setupMock      func()
	}{
		{
			name:           "Valid film ID",
			filmID:         "1",
			expectedStatus: http.StatusOK,
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"film_id", "title", "description", "release_year", "rental_duration", "rental_rate", "length", "replacement_cost", "rating", "last_update", "language_id"}).
					AddRow(1, "Test Film 1", "Test Description 1", 2020, 3, 2.99, 120, 19.99, "PG", sqlmock.AnyArg(), 1)
				mock.ExpectQuery("SELECT (.+) FROM film WHERE film_id").WithArgs(1).WillReturnRows(rows)
			},
		},
		{
			name:           "Non-numeric film ID",
			filmID:         "abc",
			expectedStatus: http.StatusBadRequest,
			setupMock:      func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			req, err := http.NewRequest("GET", "/films/"+tt.filmID, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Contains(t, response, "data")
			}
		})
	}
}

func TestAddFilmWithAuth(t *testing.T) {
	mock, cleanup := setupTestDB(t)
	defer cleanup()

	router, jwtMaker := setupProtectedTestRouter()

	tests := []struct {
		name           string
		username       string
		role           string
		filmData       map[string]interface{}
		expectedStatus int
		setupMock      func()
	}{
		{
			name:     "Admin can create film",
			username: "admin",
			role:     model.RoleAdmin,
			filmData: map[string]interface{}{
				"title":            "New Film",
				"description":      "New Description",
				"release_year":     2023,
				"rental_rate":      4.99,
				"rental_duration":  3,
				"length":           120,
				"replacement_cost": 19.99,
				"language_id":      1,
			},
			expectedStatus: http.StatusCreated,
			setupMock: func() {
				mock.ExpectQuery("INSERT INTO film").WillReturnRows(sqlmock.NewRows([]string{"film_id"}).AddRow(1))
			},
		},
		{
			name:     "User cannot create film",
			username: "user",
			role:     model.RoleUser,
			filmData: map[string]interface{}{
				"title":            "New Film",
				"description":      "New Description",
				"release_year":     2023,
				"rental_rate":      4.99,
				"rental_duration":  3,
				"length":           120,
				"replacement_cost": 19.99,
				"language_id":      1,
			},
			expectedStatus: http.StatusForbidden,
			setupMock:      func() {},
		},
		{
			name:     "Invalid film data",
			username: "admin",
			role:     model.RoleAdmin,
			filmData: map[string]interface{}{
				"title": "", // Empty title should fail validation
			},
			expectedStatus: http.StatusBadRequest,
			setupMock:      func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			// Create token
			accessToken, err := jwtMaker.CreateToken(tt.username, tt.role, time.Hour, token.TokenTypeAccessToken)
			require.NoError(t, err)

			// Prepare request body
			jsonData, err := json.Marshal(tt.filmData)
			require.NoError(t, err)

			// Create request
			req, err := http.NewRequest("POST", "/films", bytes.NewBuffer(jsonData))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+accessToken)

			// Make request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Contains(t, response, "data")
			}
		})
	}
}

func TestUpdateFilmWithAuth(t *testing.T) {
	mock, cleanup := setupTestDB(t)
	defer cleanup()

	router, jwtMaker := setupProtectedTestRouter()

	tests := []struct {
		name           string
		username       string
		role           string
		filmID         string
		filmData       map[string]interface{}
		expectedStatus int
		setupMock      func()
	}{
		{
			name:     "Admin can update film",
			username: "admin",
			role:     model.RoleAdmin,
			filmID:   "1",
			filmData: map[string]interface{}{
				"title":            "Updated Film",
				"description":      "Updated Description",
				"release_year":     2024,
				"rental_rate":      5.99,
				"rental_duration":  3,
				"length":           120,
				"replacement_cost": 19.99,
				"language_id":      1,
			},
			expectedStatus: http.StatusOK,
			setupMock: func() {
				mock.ExpectExec("UPDATE film SET").
					WithArgs("Updated Film", "Updated Description", 2024, 3, sqlmock.AnyArg(), 120, 19.99, "", "", sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:     "User cannot update film",
			username: "user",
			role:     model.RoleUser,
			filmID:   "1",
			filmData: map[string]interface{}{
				"title":            "Updated Film",
				"description":      "Updated Description",
				"release_year":     2024,
				"rental_rate":      5.99,
				"rental_duration":  3,
				"length":           120,
				"replacement_cost": 19.99,
				"language_id":      1,
			},
			expectedStatus: http.StatusForbidden,
			setupMock:      func() {},
		},
		{
			name:     "Invalid film data",
			username: "admin",
			role:     model.RoleAdmin,
			filmID:   "1",
			filmData: map[string]interface{}{
				"title": "", // Empty title should fail validation
			},
			expectedStatus: http.StatusBadRequest,
			setupMock:      func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			// Create token
			accessToken, err := jwtMaker.CreateToken(tt.username, tt.role, time.Hour, token.TokenTypeAccessToken)
			require.NoError(t, err)

			// Prepare request body
			jsonData, err := json.Marshal(tt.filmData)
			require.NoError(t, err)

			// Create request
			req, err := http.NewRequest("PUT", "/films/"+tt.filmID, bytes.NewBuffer(jsonData))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+accessToken)

			// Make request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Contains(t, response, "data")
			}
		})
	}
}

func TestDeleteFilmWithAuth(t *testing.T) {
	mock, cleanup := setupTestDB(t)
	defer cleanup()

	router, jwtMaker := setupProtectedTestRouter()

	tests := []struct {
		name           string
		username       string
		role           string
		filmID         string
		expectedStatus int
		setupMock      func()
	}{
		{
			name:           "Admin can delete film",
			username:       "admin",
			role:           model.RoleAdmin,
			filmID:         "1",
			expectedStatus: http.StatusOK,
			setupMock: func() {
				mock.ExpectExec("DELETE FROM film WHERE film_id").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:           "User cannot delete film",
			username:       "user",
			role:           model.RoleUser,
			filmID:         "1",
			expectedStatus: http.StatusForbidden,
			setupMock:      func() {},
		},
		{
			name:           "Invalid film ID",
			username:       "admin",
			role:           model.RoleAdmin,
			filmID:         "abc",
			expectedStatus: http.StatusBadRequest,
			setupMock:      func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			// Create token
			accessToken, err := jwtMaker.CreateToken(tt.username, tt.role, time.Hour, token.TokenTypeAccessToken)
			require.NoError(t, err)

			// Create request
			req, err := http.NewRequest("DELETE", "/films/"+tt.filmID, nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+accessToken)

			// Make request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Contains(t, response, "message")
			}
		})
	}
}

func TestFilmEndpointsWithoutAuth(t *testing.T) {
	router, _ := setupProtectedTestRouter()

	tests := []struct {
		name           string
		method         string
		url            string
		body           map[string]interface{}
		expectedStatus int
	}{
		{
			name:           "POST without auth",
			method:         "POST",
			url:            "/films",
			body:           map[string]interface{}{"title": "Test"},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "PUT without auth",
			method:         "PUT",
			url:            "/films/1",
			body:           map[string]interface{}{"title": "Test"},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "DELETE without auth",
			method:         "DELETE",
			url:            "/films/1",
			body:           nil,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			var err error

			if tt.body != nil {
				jsonData, _ := json.Marshal(tt.body)
				req, err = http.NewRequest(tt.method, tt.url, bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req, err = http.NewRequest(tt.method, tt.url, nil)
			}
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), "authorization header is not provided")
		})
	}
}

func TestFilmEndpointsWithInvalidToken(t *testing.T) {
	router, _ := setupProtectedTestRouter()

	tests := []struct {
		name           string
		method         string
		url            string
		token          string
		expectedStatus int
	}{
		{
			name:           "POST with invalid token",
			method:         "POST",
			url:            "/films",
			token:          "Bearer invalid.token.here",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "PUT with invalid token",
			method:         "PUT",
			url:            "/films/1",
			token:          "Bearer invalid.token.here",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "DELETE with invalid token",
			method:         "DELETE",
			url:            "/films/1",
			token:          "Bearer invalid.token.here",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.url, nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", tt.token)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), "token is invalid")
		})
	}
}
