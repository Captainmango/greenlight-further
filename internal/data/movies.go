package data

import (
	"database/sql"
	"time"

	"github.com/captainmango/greenlight/internal/validator"
	"github.com/lib/pq"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitzero"`
	Runtime   Runtime   `json:"runtime,omitzero"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

type MovieJSON struct {
	Title   string   `json:"title"`
	Year    int32    `json:"year,omitzero"`
	Runtime Runtime  `json:"runtime,omitzero"`
	Genres  []string `json:"genres,omitempty"`
}

type MovieDAO struct {
	DB *sql.DB
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

func (m MovieDAO) Insert(movie *Movie) error {
	// Define the SQL query for inserting a new record in the movies table and returning
	// the system-generated data.
	query := `
INSERT INTO movies (title, year, runtime, genres)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, version`
	// Create an args slice containing the values for the placeholder parameters from
	// the movie struct. Declaring this slice immediately next to our SQL query helps to
	// make it nice and clear *what values are being used where* in the query.
	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}
	// Use the QueryRow() method to execute the SQL query on our connection pool,
	// passing in the args slice as a variadic parameter and scanning the system-
	// generated id, created_at and version values into the movie struct.
	return m.DB.QueryRow(query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (m MovieDAO) Get(id int64) (*Movie, error) {
	return nil, nil
}

func (m MovieDAO) Update(movie *Movie) error {
	return nil
}

func (m MovieDAO) Delete(id int64) error {
	return nil
}
