// Filename: cmd/api/quotes.go
package main

import (
	"fmt"
	"net/http"

	"github.com/aiycoleman/qod/internal/data"
	"github.com/aiycoleman/qod/internal/validator"
)

func (app *application) createQuoteHandler(w http.ResponseWriter,
	r *http.Request) {
	// create a struct to hold a quote
	// struct tags[“] to make the names display in lowercase
	var incomingData struct {
		Content string `json:"content"`
		Author  string `json:"author"`
	}
	// perform the decoding
	err := app.readJSON(w, r, &incomingData)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	quote := &data.Quote{
		Content: incomingData.Content,
		Author:  incomingData.Author,
	}

	// Initialize a Validator instance
	v := validator.New()

	// Do the validation
	data.ValidateQuote(v, quote)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors) // implemented later
		return
	}

	// for now display the result
	fmt.Fprintf(w, "%+v\n", incomingData)
}
