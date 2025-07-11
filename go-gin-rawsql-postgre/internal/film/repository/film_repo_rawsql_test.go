package repository_test

import (
	"film-rental/internal/film/model"
	"film-rental/internal/film/repository"
	dbRaw "film-rental/pkg/db/raw-sql"
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
	dbRaw.DB = mockDB

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
