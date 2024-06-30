package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/captainmango/greenlight/internal/data"
)

func (a *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var movie data.Movie

	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		a.errorResponse(w, r, http.StatusBadRequest, err.Error())

		return
	}

	err = a.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)

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

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Up",
		Year:      2011,
		Runtime:   190,
		Genres:    []string{"kids", "comedy"},
		Version:   1,
	}

	err = a.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)

	if err != nil {
		a.logError(r, err)

		a.serverErrorResponse(w, r, err)
	}
}
