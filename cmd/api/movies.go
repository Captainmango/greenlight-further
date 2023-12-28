package main

import (
	"fmt"
	"net/http"
)


func (a *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a movie")
}

func (a *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIdParam(r)

	if err != nil {
		http.NotFound(w, r)

		return
	}

	fmt.Fprintf(w, "found movie with id of %d\n", id)
}
