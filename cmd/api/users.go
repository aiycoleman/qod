// Filename: cmd/api/users.go
package main

import (
	"errors"
	"net/http"

	"github.com/aiycoleman/qod/internal/data"
	"github.com/aiycoleman/qod/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter,
	r *http.Request) {
	// Get the passed in data from the request body and store in a temporary struct
	var incomingData struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &incomingData)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// we will add the password later after we have hashed it
	user := &data.User{
		Username:  incomingData.Username,
		Email:     incomingData.Email,
		Activated: false,
	}

	// hash the password and store it along with the cleartext version
	err = user.Password.Set(incomingData.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Perform validation for the User
	v := validator.New()

	data.ValidateUser(v, *user)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.userModel.Insert(user) // we will add userModel to main() later
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	data := envelope{
		"user": user,
	}
	// Status code 201 resource created
	err = app.writeJSON(w, http.StatusCreated, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
