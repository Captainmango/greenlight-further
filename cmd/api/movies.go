package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/captainmango/greenlight/internal/data"
)


func (a *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a movie")
}

func (a *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIdParam(r)

	if err != nil {
		a.notFoundResponse(w, r)

		return
	}

	movie := data.Movie{
		ID: id,
		CreatedAt: time.Now(),
		Title: "Up",
		Year: 2011,
		Runtime: 190,
		Genres: []string{"kids", "comedy"},
		Version: 1,
	}

	err = a.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)

	if err != nil {
		a.logError(r, err)

		a.serverErrorResponse(w, r, err)
	}
}
