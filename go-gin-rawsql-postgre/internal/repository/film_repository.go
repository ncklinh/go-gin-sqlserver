package repository

import (
	"film-rental/internal/model"
	"film-rental/pkg/db"
	"fmt"
)

func GetAllFilms(page int, limit int) ([]model.Film, int, error) {
	queryStr := fmt.Sprintf(`SELECT film_id, title, description, release_year, rental_duration, rental_rate, length, replacement_cost, rating, last_update FROM film LIMIT %d OFFSET %d`, limit, (page-1)*limit)

	rows, err := db.DB.Query(queryStr)

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	rowCount := db.DB.QueryRow("SELECT COUNT (*) FROM film")

	var totalCount int
	if err := rowCount.Scan(&totalCount); err != nil {

	}

	var films []model.Film
	for rows.Next() {
		var f model.Film
		if err := rows.Scan(&f.ID, &f.Title, &f.Description, &f.ReleaseYear, &f.RentalDuration, &f.RentalRate, &f.Length, &f.ReplacementCost, &f.Rating, &f.LastUpdate); err != nil {
			continue
		}
		films = append(films, f)
	}

	return films, totalCount, nil
}
