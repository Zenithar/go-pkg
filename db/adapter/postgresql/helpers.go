package postgresql

import (
	"fmt"
	"unicode"

	"go.zenithar.org/pkg/db"
)

// ToSnakeCase convert the given string to snake case following the Golang format:
// acronyms are converted to lower-case and preceded by an underscore.
func ToSnakeCase(in string) string {
	runes := []rune(in)

	var out []rune
	for i := 0; i < len(runes); i++ {
		if i > 0 && (unicode.IsUpper(runes[i]) || unicode.IsNumber(runes[i])) && ((i+1 < len(runes) && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}

// ConvertSortParameters to sql query string
func ConvertSortParameters(params db.SortParameters, sortableColumns map[string]bool) []string {

	var sorts []string
	for column, direction := range params {
		realColumn := ToSnakeCase(column)
		if _, ok := sortableColumns[realColumn]; ok {
			switch direction {
			case db.Ascending:
				sorts = append(sorts, fmt.Sprintf("%s asc", realColumn))
			case db.Descending:
				sorts = append(sorts, fmt.Sprintf("%s desc", realColumn))
			default:
			}
		}
	}

	return sorts
}
