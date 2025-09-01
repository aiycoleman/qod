// Filename: internal/data/quotes.go
package data

import (
	"time"
)

// Uppercase allows them to be exportable/public
type Quote struct {
	ID        int64     // unique value for each quote
	Content   string    // the quote data
	Author    string    // the person who wrote the quote
	CreatedAt time.Time // database timestamp
	Version   int32     // incremented on each update
}
