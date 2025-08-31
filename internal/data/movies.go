package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/captainmango/greenlight/internal/validator"
	"github.com/lib/pq"
)

type MovieUpdateJSON struct {
	Title   *string  `json:"title"`
	Year    *int32   `json:"year"`
	Runtime *Runtime `json:"runtime"`
	Genres  []string `json:"genres"`
}

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

func ValidateMovieJSON(v *validator.Validator, movie *Movie) {
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (m MovieDAO) Get(id int64) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	// specify the columns so if we change the schema the query doesn't break
	query := `
SELECT id, created_at, title, year, runtime, genres, version 
FROM movies 
WHERE id = $1;
`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resultMovie := Movie{}
	// Scan puts the values in the box. It does not translate things to the struct fields (annoying!!)
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&resultMovie.ID,
		&resultMovie.CreatedAt,
		&resultMovie.Title,
		&resultMovie.Year,
		&resultMovie.Runtime,
		pq.Array(&resultMovie.Genres),
		&resultMovie.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &resultMovie, nil
}

func (m MovieDAO) Update(movie *Movie) error {
	// Declare the SQL query for updating the record and returning the new version
	// number.
	query := `
UPDATE movies
SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
WHERE id = $5 AND version = $6
RETURNING version`
	// Create an args slice containing the values for the placeholder parameters.
	args := []any{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
		movie.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&movie.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m MovieDAO) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
DELETE FROM movies
WHERE id = $1;
`
	args := []any{id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffecred, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffecred == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m MovieDAO) GetAll() ([]Movie, error) {
	query := `
SELECT id, created_at, title, year, runtime, genres, version 
FROM movies;
`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	movies := []Movie{}
	for rows.Next() {
		resultMovie := Movie{}
		rows.Scan(
			&resultMovie.ID,
			&resultMovie.CreatedAt,
			&resultMovie.Title,
			&resultMovie.Year,
			&resultMovie.Runtime,
			pq.Array(&resultMovie.Genres),
			&resultMovie.Version,
		)
		movies = append(movies, resultMovie)
	}

	return movies, nil
}
