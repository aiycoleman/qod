// Filename: cmd/api/routes.go

package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// routes specifies our routes
func (app *application) routes() http.Handler {
	// setup a new routes
	router := httprouter.New()

	// handle 404
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	// handle 405
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	// setup routes

	// Define a GET route for health check
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	return app.recoverPanic(router)
}
