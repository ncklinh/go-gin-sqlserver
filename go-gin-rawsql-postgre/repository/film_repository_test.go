package repository_test

import (
	"film-rental/db"
	"film-rental/model"
	"film-rental/repository"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestInsertFilm_Mock tests the film insertion functionality
func TestInsertFilm_Mock(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer mockDB.Close()

	// Set the mock database to the global db.DB
	db.DB = mockDB

	expectedID := int64(42)

	film := model.Film{
		Title:           "Test",
		Description:     "Desc",
		ReleaseYear:     2025,
		RentalDuration:  5,
		RentalRate:      2.99,
		Length:          90,
		ReplacementCost: 14.99,
		Rating:          "PG",
		LastUpdate:      time.Now(),
		LanguageId:      1,
	}
	mock.ExpectQuery(`INSERT INTO film .* RETURNING film_id`).
		WithArgs(
			film.Title,
			film.Description,
			film.ReleaseYear,
			film.RentalDuration,
			sqlmock.AnyArg(), // float
			film.Length,
			sqlmock.AnyArg(), // float
			film.Rating,
			sqlmock.AnyArg(), // time.Time
			film.LanguageId,
		).
		WillReturnRows(sqlmock.NewRows([]string{"film_id"}).AddRow(expectedID))

	id, err := repository.InsertFilm(film)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != expectedID {
		t.Fatalf("expected ID %d, got %d", expectedID, id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %s", err)
	}
}

// TestIsUsernameExists_Mock tests the username uniqueness check functionality
func TestIsUsernameExists_Mock(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer mockDB.Close()

	// Set the mock database to the global db.DB
	db.DB = mockDB

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
