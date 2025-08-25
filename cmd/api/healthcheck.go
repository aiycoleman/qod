// Filename: cmd/api/healthcheck.go

package main

import (
	"fmt"
	"net/http"
)

// healthcheckHandler gives us the health of the system
func (app *application) healthcheckHandler(w http.ResponseWriter,
	r *http.Request) {

	var version string // Version number of the API

	// Create a small JSON string with dynamic values from config
	js := `{"status": "available", "environment": %q, "version": %q}`
	js = fmt.Sprintf(js, app.config.env, version)

	// Content-Type is text/plain by default
	// Tell the client the response is JSON (not plain text)
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON as the HTTP response body.
	w.Write([]byte(js))

}
