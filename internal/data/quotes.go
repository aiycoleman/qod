// Filename: internal/data/quotes.go
package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

// The QuoteModel expects a connection pool
type QuoteModel struct {
	DB *sql.DB
}

// Insert a new row in the quotes table
// A pointer to the quote
func (q QuoteModel) Insert(quote *Quote) error {
	// SQL statement to be executed
	query := `
		INSERT INTO quotes (content, author)
		VALUES ($1, $2)
		RETURNING id, created_at, version
		`
	// values to replace the $1 and $2
	args := []any{quote.Content, quote.Author}

	// Context with a 3-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// execute query against the database
	return q.DB.QueryRowContext(ctx, query, args...).Scan(&quote.ID, &quote.CreatedAt, &quote.Version)
}

// Get a specific quote from the quote table
func (q QuoteModel) Get(id int64) (*Quote, error) {
	// check if the id is valid
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	// the SQL query to be executed against the database table
	query := `
		SELECT id, content, author, created_at, version
		FROM quotes
		WHERE id = $1
		`
	// Declare a variable of type Quote to store the returned quote
	var quote Quote

	// Set a 3-second context/timer
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := q.DB.QueryRowContext(ctx, query, id).Scan(&quote.ID,
		&quote.Content,
		&quote.Author,
		&quote.CreatedAt,
		&quote.Version)

	// check for which type error
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &quote, nil
}

// Update a specific quote from the db
func (q QuoteModel) Update(quote *Quote) error {
	// The SQL query to be executed against the database table
	// Every time we make an update, we increment the version number
	query := `
        UPDATE quotes
        SET content = $1, author = $2, version = version + 1
        WHERE id = $3
        RETURNING version
		`
	// values to replace the $1 and $2
	args := []any{quote.Content, quote.Author, quote.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return q.DB.QueryRowContext(ctx, query, args...).Scan(&quote.Version)

}

// Delete a specific Quote
func (q QuoteModel) Delete(id int64) error {

	// check if the id is valid
	if id < 1 {
		return ErrRecordNotFound
	}

	// the SQL query to be executed against the database table
	query := `
        DELETE FROM quotes
        WHERE id = $1
      `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := q.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	// Check if any rows were delete
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

// Get all quotes
func (q QuoteModel) GetAll(content string, author string, filters Filters) ([]*Quote, Metadata, error) {
	query := fmt.Sprintf(`
        SELECT COUNT(*) OVER(), id, created_at, content, author, version
        FROM quotes
        WHERE (to_tsvector('simple', content) @@
              plainto_tsquery('simple', $1) OR $1 = '') 
        AND (to_tsvector('simple', author) @@ 
             plainto_tsquery('simple', $2) OR $2 = '') 
        ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := q.DB.QueryContext(ctx, query, content, author, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()
	totalRecords := 0

	quotes := []*Quote{}

	for rows.Next() {
		var quote Quote
		err := rows.Scan(
			&totalRecords,
			&quote.ID,
			&quote.CreatedAt,
			&quote.Content,
			&quote.Author,
			&quote.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		quotes = append(quotes, &quote)
	}

	// check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return quotes, metadata, nil
}
