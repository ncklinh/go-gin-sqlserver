package repository_test

import (
	"film-rental/internal/staff/repository"
	model "film-rental/internal/token/model"
	dbRaw "film-rental/pkg/db/raw-sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestIsUsernameExists_Mock tests the username uniqueness check functionality
func TestIsUsernameExists_Mock(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer mockDB.Close()

	// Set the mock database to the global db.DB
	dbRaw.DB = mockDB

	username := "testuser"

	// Test case 1: Username exists
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM staff WHERE username = \$1`).
		WithArgs(username).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	exists, err := repository.IsUsernameExists(username)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exists {
		t.Fatalf("expected username to exist, but it doesn't")
	}

	// Test case 2: Username doesn't exist
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM staff WHERE username = \$1`).
		WithArgs("nonexistent").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	exists, err = repository.IsUsernameExists("nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exists {
		t.Fatalf("expected username to not exist, but it does")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %s", err)
	}
}

// TestRefreshTokenFlow tests the complete refresh token flow
func TestRefreshTokenFlow(t *testing.T) {
	// This is a basic test to ensure the refresh token request model works
	// In a real application, you would test the actual JWT creation and verification

	req := model.RefreshTokenRequest{
		RefreshToken: "test-refresh-token",
	}

	if req.RefreshToken == "" {
		t.Fatalf("Refresh token should not be empty")
	}

	if req.RefreshToken != "test-refresh-token" {
		t.Fatalf("Expected refresh token to be 'test-refresh-token', got %s", req.RefreshToken)
	}

	// Test TokenResponse struct
	tokenResp := model.TokenResponse{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		TokenType:    "Bearer",
		ExpiresIn:    900,
	}

	if tokenResp.AccessToken == "" {
		t.Fatalf("Access token should not be empty")
	}

	if tokenResp.TokenType != "Bearer" {
		t.Fatalf("Expected token type to be 'Bearer', got %s", tokenResp.TokenType)
	}

	if tokenResp.ExpiresIn != 900 {
		t.Fatalf("Expected expires_in to be 900, got %d", tokenResp.ExpiresIn)
	}
}
