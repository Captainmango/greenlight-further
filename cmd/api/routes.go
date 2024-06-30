package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.PanicHandler = a.panicHandler

	router.NotFound = http.HandlerFunc(a.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(a.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", a.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movies", a.createMovieHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", a.showMovieHandler)

	return router
}

func (a *application) panicHandler(w http.ResponseWriter, r *http.Request, rcv any) {
	w.Header().Set("Connection", "close")

	a.serverErrorResponse(w, r, fmt.Errorf("%s", rcv))
}
