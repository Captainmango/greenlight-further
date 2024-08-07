package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	e := envelope{
		"status": "available",
		"system_info": envelope{
			"environment": app.config.env,
			"version": version,
		},
	}

	err := app.writeJSON(w, http.StatusOK, e, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
