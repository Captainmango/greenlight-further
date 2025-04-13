package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/captainmango/greenlight/internal/data"
	"github.com/captainmango/greenlight/internal/validator"
)

func (a *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var movieJson data.MovieJSON

	err := a.readJSON(w, r, &movieJson)
	if err != nil {
		a.badRequestResponse(w, r, err)

		return
	}

	v := validator.New()

	if data.ValidateMovieJSON(v, &movieJson); !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)

		return
	}

	movie := &data.Movie{
		Title:   movieJson.Title,
		Year:    movieJson.Year,
		Runtime: movieJson.Runtime,
		Genres:  movieJson.Genres,
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
