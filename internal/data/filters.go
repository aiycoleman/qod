// Filename: internal/data/filters.go
package data

import (
	"github.com/aiycoleman/qod/internal/validator"
)

type Filters struct {
	Page     int // page number the client wants
	PageSize int // number of records per page
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 500, "page", "must be a maximum of 500")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")
}

// Calculate how many records to send back
func (f Filters) limit() int {
	return f.PageSize
}

// Calculate the offset so that we remember how many records have been sent
// and how many remain to be sent
func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}
