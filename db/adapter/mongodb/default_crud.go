// MIT License
//
// Copyright (c) 2019 Thibault NORMAND
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package mongodb

import (
	"context"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pkg/errors"

	"go.zenithar.org/pkg/db"
	"go.zenithar.org/pkg/log"
)

// Default contains the basic implementation of the MongoCRUD interface
type Default struct {
	table   string
	db      string
	session *mongo.Client
}

// NewCRUDTable sets up a new Default struct
func NewCRUDTable(session *mongo.Client, db, table string) *Default {
	return &Default{
		db:      db,
		table:   table,
		session: session,
	}
}

// -----------------------------------------------------------------------------

// GetTableName returns table's name
func (d *Default) GetTableName() string {
	return d.table
}

// GetDBName returns database's name
func (d *Default) GetDBName() string {
	return d.db
}

// GetTable returns the table as a mgo.Collection
func (d *Default) GetTable() interface{} {
	return d.session.Database(d.db).Collection(d.table)
}

// GetSession returns the current session
func (d *Default) GetSession() interface{} {
	return d.session
}

// -----------------------------------------------------------------------------

// Insert inserts a document into the database
func (d *Default) Insert(ctx context.Context, data interface{}) error {
	// Run in transaction
	return Transaction(ctx, d.session, func() error {
		_, err := d.session.Database(d.db).Collection(d.table).InsertOne(ctx, data)
		return err
	})
}

// InsertOrUpdate inserts or update document if exists
func (d *Default) InsertOrUpdate(ctx context.Context, id interface{}, data interface{}) error {
	// Run in transaction
	return Transaction(ctx, d.session, func() error {
		_, err := d.session.Database(d.db).Collection(d.table).UpdateOne(ctx, id, data)
		return err
	})
}

// Update performs an update on an existing resource according to passed data
func (d *Default) Update(ctx context.Context, selector interface{}, data interface{}) error {
	// Run in transaction
	return Transaction(ctx, d.session, func() error {
		_, err := d.session.Database(d.db).Collection(d.table).UpdateMany(ctx, selector, data)
		return err
	})
}

// UpdateID performs an update on an existing resource with ID that equals the id argument
func (d *Default) UpdateID(ctx context.Context, id interface{}, data interface{}) error {
	// Run in transaction
	return Transaction(ctx, d.session, func() error {
		_, err := d.session.Database(d.db).Collection(d.table).UpdateOne(ctx, bson.M{
			"_id": id,
		}, data)
		return err
	})
}

// DeleteAll deletes resources that match the passed filter
func (d *Default) DeleteAll(ctx context.Context, pred interface{}) error {
	// Run in transaction
	return Transaction(ctx, d.session, func() error {
		_, err := d.session.Database(d.db).Collection(d.table).DeleteMany(ctx, pred)
		return err
	})
}

// Delete deletes a resource with specified ID
func (d *Default) Delete(ctx context.Context, id interface{}) error {
	// Instrument opentracing
	if span := opentracing.SpanFromContext(ctx); span != nil {
		ext.DBType.Set(span, "mongo")
		ext.DBInstance.Set(span, d.db)
		ext.PeerHostname.Set(span, d.session.ConnectionString())
	}

	// Run in transaction
	return Transaction(ctx, d.session, func() error {
		_, err := d.session.Database(d.db).Collection(d.table).DeleteOne(ctx, id)
		return err
	})
}

// Find searches for a resource in the database and then returns a cursor
func (d *Default) Find(ctx context.Context, id interface{}, value interface{}) error {
	// Run in transaction
	return Transaction(ctx, d.session, func() error {
		res, err := d.session.Database(d.db).Collection(d.table).Find(ctx, bson.M{
			"_id": id,
		})
		if err != nil {
			return err
		}
		return res.Decode(value)
	})
}

// FindFetchOne searches for a resource and then unmarshals the first row into value
func (d *Default) FindFetchOne(ctx context.Context, id string, value interface{}) error {
	// Run in transaction
	return Transaction(ctx, d.session, func() error {
		res := d.session.Database(d.db).Collection(d.table).FindOne(ctx, id)
		return res.Decode(value)
	})
}

// FindOneBy is an utility for fetching values if they are stored in a key-value manenr.
func (d *Default) FindOneBy(ctx context.Context, key string, value interface{}, result interface{}) error {
	// Run in transaction
	return Transaction(ctx, d.session, func() error {
		res := d.session.Database(d.db).Collection(d.table).FindOne(ctx, bson.M{
			key: value,
		})
		return res.Decode(result)
	})
}

// FindBy is an utility for fetching values if they are stored in a key-value manenr.
func (d *Default) FindBy(ctx context.Context, key string, value interface{}, results interface{}) error {
	// Run in transaction
	return Transaction(ctx, d.session, func() error {
		res, err := d.session.Database(d.db).Collection(d.table).Find(ctx, bson.M{
			key: value,
		})
		if err != nil {
			return err
		}
		return res.Decode(results)
	})
}

// FindByAndCount returns the number of elements that match the filter
func (d *Default) FindByAndCount(ctx context.Context, key string, value interface{}) (int64, error) {
	var count int64

	// Run in transaction
	if err := Transaction(ctx, d.session, func() error {
		n, err := d.session.Database(d.db).Collection(d.table).Count(ctx, bson.M{
			key: value,
		})
		count = n
		return err
	}); err != nil {
		return count, err
	}

	return count, nil
}

// FindByAndFetch retrieves a value by key and then fills results with the result.
func (d *Default) FindByAndFetch(ctx context.Context, key string, value interface{}, results interface{}) error {
	// Run in transaction
	return Transaction(ctx, d.session, func() error {
		res, err := d.session.Database(d.db).Collection(d.table).Find(ctx, bson.M{
			key: value,
		})
		if err != nil {
			return err
		}
		return res.Decode(results)
	})
}

// WhereCount allows counting with multiple fields
func (d *Default) WhereCount(ctx context.Context, filter interface{}) (int64, error) {
	var count int64

	// Run in transaction
	if err := Transaction(ctx, d.session, func() error {
		n, err := d.session.Database(d.db).Collection(d.table).Count(ctx, filter)
		count = n
		return err
	}); err != nil {
		return count, err
	}

	return count, nil
}

// Where allows filtering with multiple fields
func (d *Default) Where(ctx context.Context, filter interface{}, results interface{}) error {
	// Run in transaction
	return Transaction(ctx, d.session, func() error {
		res, err := d.session.Database(d.db).Collection(d.table).Find(ctx, filter)
		if err != nil {
			return err
		}
		return res.Decode(results)
	})
}

// WhereAndFetchLimit filters with multiple fields and then fills results with all found resources
func (d *Default) WhereAndFetchLimit(ctx context.Context, filter interface{}, paginator *db.Pagination, results interface{}) error {
	// Run in transaction
	return Transaction(ctx, d.session, func() error {
		limit := int64(paginator.PerPage)
		skip := int64(paginator.Offset())
		res, err := d.session.Database(d.db).Collection(d.table).Find(ctx, filter, &options.FindOptions{
			Limit: &limit,
			Skip:  &skip,
		})
		if err != nil {
			return err
		}
		return res.Decode(results)
	})
}

// WhereAndFetchOne filters with multiple fields and then fills result with the first found resource
func (d *Default) WhereAndFetchOne(ctx context.Context, filter interface{}, result interface{}) error {
	// Run in transaction
	return Transaction(ctx, d.session, func() error {
		res := d.session.Database(d.db).Collection(d.table).FindOne(ctx, filter)
		return res.Decode(result)
	})
}

// List all entities from the database
func (d *Default) List(ctx context.Context, results interface{}, sortParams *db.SortParameters, pagination *db.Pagination) error {
	// Run in transaction
	return d.Search(ctx, results, bson.M{}, sortParams, pagination)
}

// Search all entities from the database
func (d *Default) Search(ctx context.Context, results interface{}, filter interface{}, sortParams *db.SortParameters, pagination *db.Pagination) error {
	// Apply Filter
	if filter == nil {
		filter = bson.M{}
	}

	// Get total
	if pagination != nil {
		total, err := d.WhereCount(ctx, filter)
		if err != nil {
			return errors.Wrap(err, "Database error")
		}
		pagination.SetTotal(uint(total))
	}

	// Prepare the query
	opts := &options.FindOptions{}

	// Apply sorts
	if sortParams != nil {
		sort := ConvertSortParameters(*sortParams)
		if len(sort) > 0 {
			opts.SetSort(sort)
		}
	}

	// Paginate
	if pagination != nil {
		opts.SetLimit(int64(pagination.PerPage))
		opts.SetSkip(int64(pagination.Offset()))
	}

	return Transaction(ctx, d.session, func() error {
		res, err := d.session.Database(d.db).Collection(d.table).Find(ctx, filter, opts)
		if err != nil {
			return err
		}
		return res.Decode(results)
	})
}

// -----------------------------------------------------------------------------

// TransactionFunc is the transaction handler closure contract
type TransactionFunc func() error

// Transaction runs the transactionfunc in a transaction
func Transaction(ctx context.Context, client *mongo.Client, fn TransactionFunc) error {
	// Initialize a session
	session, err := client.StartSession()
	if err != nil {
		return errors.Wrap(err, "Database error")
	}
	defer session.EndSession(ctx)

	// Start transaction
	if err := session.StartTransaction(); err != nil {
		return errors.Wrap(err, "Database error")
	}

	// Run the closure
	if err := fn(); err != nil {
		log.CheckErrCtx(ctx, "Unable to abort transaction", session.AbortTransaction(ctx))
		return err
	}

	// Commit the transaction
	err = session.CommitTransaction(ctx)
	if err != nil {
		log.CheckErrCtx(ctx, "Unable to abort transaction", session.AbortTransaction(ctx))
		return err
	}

	// No error
	return nil
}
