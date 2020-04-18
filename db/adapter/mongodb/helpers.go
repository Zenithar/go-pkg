package mongodb

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.zenithar.org/pkg/db"
)

// ConvertSortParameters to sql query string
func ConvertSortParameters(params db.SortParameters) bson.M {

	sorts := bson.M{}
	for k, v := range params {
		switch v {
		case db.Ascending:
			sorts[k] = 1
		case db.Descending:
			sorts[k] = -1
		default:
			sorts[k] = 1
		}
	}

	return sorts

}
