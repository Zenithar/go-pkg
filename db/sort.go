package db

import "strings"

// SortDirection is the enumeration for sort
type SortDirection int

const (
	// Ascending sort from bottom to up
	Ascending SortDirection = iota + 1
	// Descending sort from up to bottom
	Descending
)

var sortDirections = [...]string{
	"asc",
	"desc",
}

func (m SortDirection) String() string {
	return sortDirections[m-1]
}

// -----------------------------------------------------------------------------

// SortParameters contains a hashmap of field name with sort direction
type SortParameters map[string]SortDirection

// SortConverter convert a list of string to a SortParameters instance
func SortConverter(sorts []string) *SortParameters {
	params := SortParameters{}

	if len(sorts) > 0 {
		for _, cond := range sorts {
			if len(strings.TrimSpace(cond)) > 0 {
				switch cond[0] {
				case '-':
					params[cond[1:]] = Descending
				case '+', ' ':
					params[cond[1:]] = Ascending
				default:
					params[cond] = Ascending
				}
			}
		}
	}

	return &params
}
