package middleware

import (
	"film-rental/internal/token"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	secretKey := "12345678901234567890123456789012" // 32 characters
	jwtMaker, err := token.NewJWTMaker(secretKey)
	require.NoError(t, err)

	tests := []struct {
		name           string
		setupAuth      func() string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Valid token",
			setupAuth: func() string {
				token, _ := jwtMaker.CreateToken("testuser", "admin", time.Hour, token.TokenTypeAccessToken)
				return "Bearer " + token
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"success"}`,
		},
		{
			name: "Missing authorization header",
			setupAuth: func() string {
				return ""
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"authorization header is not provided"}`,
		},
		{
			name: "Invalid authorization header format",
			setupAuth: func() string {
				return "InvalidFormat"
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"invalid authorization header format"}`,
		},
		{
			name: "Unsupported authorization type",
			setupAuth: func() string {
				return "Basic dGVzdDp0ZXN0"
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"unsupported authorization type basic"}`,
		},
		{
			name: "Invalid token",
			setupAuth: func() string {
				return "Bearer invalid.token.here"
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"token is invalid"}`,
		},
		{
			name: "Expired token",
			setupAuth: func() string {
				token, _ := jwtMaker.CreateToken("testuser", "admin", -time.Hour, token.TokenTypeAccessToken)
				return "Bearer " + token
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"token has expired"}`,
		},
		{
			name: "Wrong token type",
			setupAuth: func() string {
				token, _ := jwtMaker.CreateToken("testuser", "admin", time.Hour, token.TokenTypeRefreshToken)
				return "Bearer " + token
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"token is invalid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new gin router for each test
			router := gin.New()
			authMiddleware := AuthMiddleware(jwtMaker)

			// Add a test endpoint that uses the auth middleware
			router.GET("/test", authMiddleware, func(c *gin.Context) {
				// Check if payload is set in context
				payload, exists := c.Get(authorizationPayloadKey)
				if exists {
					assert.NotNil(t, payload)
					tokenPayload, ok := payload.(*token.Payload)
					assert.True(t, ok)
					assert.Equal(t, "testuser", tokenPayload.Username)
					assert.Equal(t, "admin", tokenPayload.Role)
				}
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			// Create request
			req, err := http.NewRequest("GET", "/test", nil)
			require.NoError(t, err)

			// Set authorization header if provided
			if authHeader := tt.setupAuth(); authHeader != "" {
				req.Header.Set(authorizationHeaderKey, authHeader)
			}

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

func TestAuthMiddlewareWithPayload(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	secretKey := "12345678901234567890123456789012"
	jwtMaker, err := token.NewJWTMaker(secretKey)
	require.NoError(t, err)

	// Create a valid token
	accessToken, err := jwtMaker.CreateToken("testuser", "admin", time.Hour, token.TokenTypeAccessToken)
	require.NoError(t, err)

	// Create router with auth middleware
	router := gin.New()
	authMiddleware := AuthMiddleware(jwtMaker)

	// Add test endpoint
	router.GET("/test", authMiddleware, func(c *gin.Context) {
		payload, exists := c.Get(authorizationPayloadKey)
		assert.True(t, exists)

		tokenPayload, ok := payload.(*token.Payload)
		assert.True(t, ok)
		assert.Equal(t, "testuser", tokenPayload.Username)
		assert.Equal(t, "admin", tokenPayload.Role)
		assert.Equal(t, uint8(token.TokenTypeAccessToken), uint8(tokenPayload.Type))

		c.JSON(http.StatusOK, gin.H{
			"username": tokenPayload.Username,
			"role":     tokenPayload.Role,
		})
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
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "testuser")
	assert.Contains(t, w.Body.String(), "admin")
}
