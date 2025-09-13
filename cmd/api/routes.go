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
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/quotes", app.createQuoteHandler)
	router.HandlerFunc(http.MethodGet, "/v1/quotes/:id", app.displayQuoteHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/quotes/:id", app.updateQuoteHandler)

	return app.recoverPanic(router)
}
