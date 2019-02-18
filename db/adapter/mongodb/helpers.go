package mongodb

import (
	"fmt"
	"strings"

	"go.zenithar.org/pkg/db"
)

// ConvertSortParameters to sql query string
func ConvertSortParameters(params db.SortParameters) []string {

	var sorts []string
	for k, v := range params {
		switch v {
		case db.Ascending:
			sorts = append(sorts, strings.ToLower(k))
		case db.Descending:
			sorts = append(sorts, fmt.Sprintf("-%s", strings.ToLower(k)))
		default:
			sorts = append(sorts, strings.ToLower(k))
		}
	}

	return sorts

}
