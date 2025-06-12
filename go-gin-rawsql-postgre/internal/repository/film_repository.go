package repository

import (
	"film-rental/internal/model"
	"film-rental/pkg/db"
)

func GetAllFilms() ([]model.Film, error) {
	rows, err := db.DB.Query("SELECT film_id, title, description, release_year, rental_duration, rental_rate, length, replacement_cost, rating, last_update FROM film")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var films []model.Film
	for rows.Next() {
		var f model.Film
		if err := rows.Scan(&f.ID, &f.Title, &f.Description, &f.ReleaseYear, &f.RentalDuration, &f.RentalRate, &f.Length, &f.ReplacementCost, &f.Rating, &f.LastUpdate); err != nil {
			continue
		}
		films = append(films, f)
	}
	return films, nil
}
