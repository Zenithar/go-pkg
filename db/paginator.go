package db

import (
	"math"
)

const (
	// DefaultPerPage defines the default value for pagination
	DefaultPerPage uint = 20
)

// Pagination is a pagination calcul handler for database request.
type Pagination struct {
	Page    uint
	PerPage uint
	total   uint
}

// SetTotal is used to defines the total count of paginated values.
func (p *Pagination) SetTotal(total uint) {
	p.total = total
}

// NumPages returns the total number of pages
func (p *Pagination) NumPages() uint {
	return maxuint(1, uint(math.Ceil(float64(p.total)/float64(p.PerPage))))
}

// Total returns the total number of items
func (p *Pagination) Total() uint {
	return p.total
}

// Offset returns the offset of first element
func (p *Pagination) Offset() uint {
	offset := (p.Page - 1) * p.PerPage
	// a couple reasonable boundaries
	offset = minuint(offset, p.total)
	offset = maxuint(offset, 0)
	return offset
}

// PrevPage returns the page number for the page before this
// bottoms out at the first page
func (p *Pagination) PrevPage() uint {
	return maxuint(p.Page-1, 1)
}

// HasPrev returns the status if current page has a previous one
func (p *Pagination) HasPrev() bool {
	return p.Page > 1
}

// NextPage returns the page number for the next page. won't go past the end
func (p *Pagination) NextPage() uint {
	return minuint(p.Page+1, p.NumPages())
}

// HasNext returns the status if current page has a next one
func (p *Pagination) HasNext() bool {
	return p.Page+1 <= p.NumPages()
}

// HasOtherPages returns the status of having previous or next pages
func (p *Pagination) HasOtherPages() bool {
	return p.HasPrev() || p.HasNext()
}

// CurrentPageCount returns the element count of the current page
func (p *Pagination) CurrentPageCount() uint {
	return minuint((p.total - p.Offset()), p.PerPage)
}

// NewPaginator returns a pagination holder
func NewPaginator(page, perPage uint) *Pagination {
	// Sanitize inputs
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = DefaultPerPage
	}

	// Return paginator instance
	return &Pagination{
		Page:    page,
		PerPage: perPage,
		total:   0,
	}
}

// -----------------------------------------------------------------------------

func minuint(a, b uint) uint {
	return uint(math.Min(float64(a), float64(b)))
}

func maxuint(a, b uint) uint {
	return uint(math.Max(float64(a), float64(b)))
}
