// Filename: cmd/api/healthcheck.go

package main

import (
	"net/http"
)

// healthcheckHandler gives us the health of the system
func (app *application) healthcheckHandler(w http.ResponseWriter,
	r *http.Request) {

	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     app.config.version,
	}

	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		return
	}
}
