package data

import (
	"time"

	"github.com/captainmango/greenlight/internal/validator"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

type MovieJSON struct {
	Title   string   `json:"title"`
	Year    int32    `json:"year,omitempty"`
	Runtime Runtime  `json:"runtime,omitempty"`
	Genres  []string `json:"genres,omitempty"`
}

func ValidateMovieJSON(v *validator.Validator, movie *MovieJSON) {
	// TITLE VALIDATIONS
	v.Check(movie.Title != "", "title", "title cannot be empty")
	v.Check(len(movie.Title) <= 500, "title", "title must be less than 500 bytes long")

	// YEAR VALIDATIONS
	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	// RUNTIME VALIDATIONS
	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")

	// GENRE VALIDATIONS
	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}
