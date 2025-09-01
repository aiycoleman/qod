// Filename: cmd/api/comments.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/aiycoleman/qod/internal/data"
)

func (app *application) createCommentHandler(w http.ResponseWriter,
	r *http.Request) {
	// create a struct to hold a comment
	// struct tags[“] to make the names display in lowercase
	var incomingData struct {
		Content string `json:"content"`
		Author  string `json:"author"`
	}
	// perform the decoding
	err := json.NewDecoder(r.Body).Decode(&incomingData)
	if err != nil {
		app.errorResponseJSON(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// for now display the result
	fmt.Fprintf(w, "%+v\n", incomingData)
}
