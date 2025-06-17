package model

import "time"

type Film struct {
	ID              int       `json:"film_id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	ReleaseYear     int       `json:"release_year"`
	RentalDuration  int       `json:"rental_duration"`
	RentalRate      float32   `json:"rental_rate"`
	Length          int       `json:"length"`
	ReplacementCost float32   `json:"replacement_cost"`
	Rating          string    `json:"rating"`
	LastUpdate      time.Time `json:"last_update"`
	LanguageId      int       `json:"language_id"`
}
