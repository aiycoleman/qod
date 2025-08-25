// Filename: cmd/api/routes.go

package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// routes specifies our routes
func (app *application) routes() http.Handler {
	// create a new router
	router := httprouter.New()

	// Define a GET route for health check
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	return router
}
