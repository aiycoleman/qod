// Filename: cmd/api/quotes.go
package main

import (
	"errors"
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

// Displays quotes
func (app *application) displayQuoteHandler(w http.ResponseWriter, r *http.Request) {
	// get the id from the url
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Call Get(to retrieve data based on id)
	quote, err := app.quoteModel.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// display quote
	data := envelope{
		"quote": quote,
	}
	err = app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// DEdit quotes
func (app *application) updateQuoteHandler(w http.ResponseWriter, r *http.Request) {
	// Get ID from the URL
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Call Get() to retrirve the comment with the specified ID
	quote, err := app.quoteModel.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// temporary incoming data struct
	var incomingData struct {
		Content *string `json:"content"`
		Author  *string `json:"author"`
	}

	// perform decoding
	err = app.readJSON(w, r, &incomingData)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Check to see which fields need to be updated
	// if  nill;, no update needed for any feild
	if incomingData.Content != nil {
		quote.Content = *incomingData.Content
	}

	if incomingData.Author != nil {
		quote.Author = *incomingData.Author
	}

	// Validate
	v := validator.New()
	data.ValidateQuote(v, quote)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Add the quote to the database table
	err = app.quoteModel.Update(quote)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	data := envelope{
		"quote": quote,
	}
	err = app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) deleteQuoteHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.quoteModel.Delete(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// display the quote
	data := envelope{"message": "quote successfully deleted"}
	err = app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) listQuotesHandler(w http.ResponseWriter, r *http.Request) {
	quotes, err := app.quoteModel.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	data := envelope{
		"quotes": quotes,
	}
	err = app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
