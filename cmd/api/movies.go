package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/captainmango/greenlight/internal/data"
	"github.com/captainmango/greenlight/internal/validator"
)

func (a *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var movieJson struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year,omitzero"`
		Runtime data.Runtime `json:"runtime,omitzero"`
		Genres  []string     `json:"genres,omitempty"`
	}

	err := a.readJSON(w, r, &movieJson)
	if err != nil {
		a.badRequestResponse(w, r, err)

		return
	}

	movie := &data.Movie{
		Title:   movieJson.Title,
		Year:    movieJson.Year,
		Runtime: movieJson.Runtime,
		Genres:  movieJson.Genres,
	}

	v := validator.New()

	if data.ValidateMovieJSON(v, movie); !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)

		return
	}

	if err = a.dao.Movies.Insert(movie); err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = a.writeJSON(w, http.StatusOK, envelope{"movie": movie}, headers)

	if err != nil {
		a.logError(r, err)

		a.serverErrorResponse(w, r, err)
	}
}

func (a *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIdParam(r)

	if err != nil {
		a.notFoundResponse(w, r)

		return
	}

	movie, err := a.dao.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}

		return
	}

	err = a.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)

	if err != nil {
		a.logError(r, err)

		a.serverErrorResponse(w, r, err)
	}
}

func (a *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	var updateMovieJson struct {
		Title   *string       `json:"title"`
		Year    *int32        `json:"year"`
		Runtime *data.Runtime `json:"runtime"`
		Genres  []string      `json:"genres"`
	}

	movieId, err := a.readIdParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	movie, err := a.dao.Movies.Get(movieId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
	}

	err = a.readJSON(w, r, &updateMovieJson)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}


	if updateMovieJson.Title != nil {
		movie.Title = *updateMovieJson.Title
	}

	if updateMovieJson.Genres != nil {
		movie.Genres = updateMovieJson.Genres
	}

	if updateMovieJson.Year != nil {
		movie.Year = *updateMovieJson.Year
	}

	if updateMovieJson.Runtime != nil {
		movie.Runtime = *updateMovieJson.Runtime
	}

	v := validator.New()
	if data.ValidateMovieJSON(v, movie); !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.dao.Movies.Update(movie)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			a.editConfilctResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}

		return
	}

	err = a.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

func (a *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	movieId, err := a.readIdParam(r)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	err = a.dao.Movies.Delete(movieId)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
			return
		default:
			a.serverErrorResponse(w, r, err)
			return
		}
	}

	err = a.writeJSON(w, http.StatusOK, envelope{"message": "successfully deleted movie"}, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *application) getMoviesHandler(w http.ResponseWriter, r *http.Request) {
	movies, err := a.dao.Movies.GetAll()
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}

	a.writeJSON(w, http.StatusOK, envelope{"movies": movies}, nil)
}
