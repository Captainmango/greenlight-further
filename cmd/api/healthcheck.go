package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	res := map[string]string {
		"status": "available",
		"environment": app.config.env,
		"version": version,
	}

	err := app.writeJSON(w, http.StatusOK, res, nil)

	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "The server ecountered an error", http.StatusInternalServerError)
	}
}
