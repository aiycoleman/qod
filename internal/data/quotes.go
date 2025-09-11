// Filename: internal/data/quotes.go
package data

import (
	"time"

	"github.com/aiycoleman/qod/internal/validator"
)

// Uppercase allows them to be exportable/public
type Quote struct {
	ID        int64     `json:"id"`      // unique value for each quote
	Content   string    `json:"content"` // the quote data
	Author    string    `json:"author"`  // the person who wrote the quote
	CreatedAt time.Time `json:"-"`       // database timestamp
	Version   int32     `json:"version"` // incremented on each update
}

// Performs the validation checks
func ValidateQuote(v *validator.Validator, quote *Quote) {
	// check if the Content field is empty
	v.Check(quote.Content != "", "content", "must be provided")
	// check if the Author field us empty
	v.Check(quote.Author != "", "author", "must be provided")
	// chekc if the content in the field is empty
	v.Check(len(quote.Content) <= 100, "content", "must not be more than 100 bytes long")
	// check if the Author field is empty
	v.Check(len(quote.Author) <= 25, "author", "must not be more than 25 bytes long")
}
