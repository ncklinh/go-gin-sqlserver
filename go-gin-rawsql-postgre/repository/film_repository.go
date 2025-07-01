package repository

import (
	"database/sql"
	"film-rental/db"
	"film-rental/model"
	"fmt"
)

const columnQuery = "film_id, title, description, release_year, rental_duration, rental_rate, length, replacement_cost, rating, last_update, language_id"

func scanFilmRow(scanner interface {
	Scan(dest ...any) error
}) (*model.Film, error) {
	var f model.Film
	err := scanner.Scan(
		&f.ID, &f.Title, &f.Description, &f.ReleaseYear,
		&f.RentalDuration, &f.RentalRate, &f.Length,
		&f.ReplacementCost, &f.Rating, &f.LastUpdate, &f.LanguageId,
	)
	return &f, err
}

func GetAllFilms(page int, limit int) ([]*model.Film, int, error) {
	queryStr := `SELECT ` + columnQuery + ` FROM film ORDER BY film_id DESC LIMIT $1 OFFSET $2`

	rows, err := db.DB.Query(queryStr, limit, (page-1)*limit)

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	rowCount := db.DB.QueryRow("SELECT COUNT (*) FROM film")

	var totalCount int
	if err := rowCount.Scan(&totalCount); err != nil {
		return nil, 0, err
	}

	var films []*model.Film
	for rows.Next() {
		if f, err := scanFilmRow(rows); err != nil {
			continue
		} else {
			films = append(films, f)

		}
	}

	return films, totalCount, nil
}

func GetFilmDetail(filmId int) (*model.Film, error) {
	queryStr := fmt.Sprintf(`SELECT %s FROM film WHERE film_id = $1`, columnQuery)

	f, err := scanFilmRow(db.DB.QueryRow(queryStr, filmId))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return f, nil
}

func InsertFilm(film model.Film) (int64, error) {
	query := `
		INSERT INTO film (
			title, description, release_year,
			rental_duration, rental_rate, length,
			replacement_cost, rating, last_update, language_id
		) VALUES  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	    RETURNING film_id
	`

	var lastID int64
	err := db.DB.QueryRow(query,
		film.Title,
		film.Description,
		film.ReleaseYear,
		film.RentalDuration,
		film.RentalRate,
		film.Length,
		film.ReplacementCost,
		film.Rating,
		film.LastUpdate,
		film.LanguageId,
	).Scan(&lastID)

	if err != nil {
		return 0, err
	}

	return lastID, nil
}

func UpdateFilm(filmId int, film model.Film) error {
	query := `
		UPDATE film SET 
			title = $1, description = $2, release_year = $3,
			rental_duration = $4, rental_rate = $5, length = $6,
			replacement_cost = $7, rating = $8, last_update = $9, language_id = $10
		WHERE film_id = $11
	`

	result, err := db.DB.Exec(query,
		film.Title,
		film.Description,
		film.ReleaseYear,
		film.RentalDuration,
		film.RentalRate,
		film.Length,
		film.ReplacementCost,
		film.Rating,
		film.LastUpdate,
		film.LanguageId,
		filmId,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func DeleteFilm(filmId int) error {
	query := `DELETE FROM film WHERE film_id = $1`

	result, err := db.DB.Exec(query, filmId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
