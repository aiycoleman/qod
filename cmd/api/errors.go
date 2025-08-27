// Filename: cmd/api/errors.go
package main

import (
	"fmt"
	"net/http"
)

// log an error message
func (a *application) logError(r *http.Request, err error) {
	method := r.Method
	uri := r.URL.RequestURI()
	a.logger.Error(err.Error(), "method", method, "uri", uri)
}

// send an error response in JSON
func (a *application) errorResponseJSON(w http.ResponseWriter, r *http.Request, status int, message any) {
	errorData := envelope{"error": message}
	err := a.writeJSON(w, status, errorData, nil)
	if err != nil {
		a.logError(r, err)
		w.WriteHeader(500)
	}
}

// send an error response if our server messes up
func (a *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	// 1st log the error message
	a.logError(r, err)
	// prepare message to response to send to the client
	message := "the server encountered a problem and could not process your request"
	a.errorResponseJSON(w, r, http.StatusInternalServerError, message)
}

// send an error response if our client messes up with a 404
func (a *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	// Only log server errors, not client errors
	// Prepare a response to send to the client
	message := "the requested resource could not be found"
	a.errorResponseJSON(w, r, http.StatusNotFound, message)
}

// send an error response if our client messes up a 405
func (a *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	// we don't log, since its a client error
	// Prepare a formatted response to send to the client
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)

	a.errorResponseJSON(w, r, http.StatusMethodNotAllowed, message)
}
