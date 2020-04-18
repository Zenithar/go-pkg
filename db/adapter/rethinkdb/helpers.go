package rethinkdb

import (
	"strings"

	"go.zenithar.org/pkg/db"

	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

// ConvertSortParameters to rethinkdb query string
func ConvertSortParameters(params db.SortParameters) []interface{} {

	var sorts []interface{}
	for k, v := range params {
		switch v {
		case db.Ascending:
			sorts = append(sorts, r.Asc(strings.ToLower(k)))
			break
		case db.Descending:
			sorts = append(sorts, r.Desc(strings.ToLower(k)))
			break
		default:
			sorts = append(sorts, r.Desc(strings.ToLower(k)))
		}
	}

	// Apply default sort
	if len(sorts) == 0 {
		sorts = append(sorts, r.Asc("id"))
	}

	return sorts
}
