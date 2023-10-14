package main

import (
	"net/http"
)

func (app application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := envelope{
		"status": "available",
		"system_info": map[string]string{
			"env":     app.config.env,
			"version": version,
		},
	}

	if err := app.writeJSON(w, data); err != nil {
		app.writeError(w, http.StatusInternalServerError)
	}
}
