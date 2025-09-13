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

	// Add the quote to the database table
	err = app.quoteModel.Insert(quote)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Set a location header (the path to the newly created quote)
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/quotes/%d", quote.ID))

	// Send a JSON response with 201 (new resource createed) status code
	data := envelope{
		"quote": quote,
	}

	err = app.writeJSON(w, http.StatusCreated, data, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
