// Filename: cmd/api/healthcheck.go

package main

import (
	"encoding/json"
	"net/http"
)

// healthcheckHandler gives us the health of the system
func (app *application) healthcheckHandler(w http.ResponseWriter,
	r *http.Request) {

	// // Create a small JSON string with dynamic values from config
	// js := `{"status": "available", "environment": %q, "version": %q}`
	// js = fmt.Sprintf(js, app.config.env, app.config.version)

	// // Content-Type is text/plain by default
	// // Tell the client the response is JSON (not plain text)
	// w.Header().Set("Content-Type", "application/json")

	// // Write the JSON as the HTTP response body.
	// w.Write([]byte(js))

	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     app.config.version,
	}

	jsResponse, err := json.Marshal(data)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		return
	}

	jsResponse = append(jsResponse, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsResponse)

}
