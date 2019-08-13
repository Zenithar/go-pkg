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

	mongowrapper "github.com/opencensus-integrations/gomongowrapper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/xerrors"

	"go.zenithar.org/pkg/db"
	"go.zenithar.org/pkg/log"
)

// Default contains the basic implementation of the MongoCRUD interface
type Default struct {
	table   string
	db      string
	session *mongowrapper.WrappedClient
}

// NewCRUDTable sets up a new Default struct
func NewCRUDTable(session *mongowrapper.WrappedClient, db, table string) *Default {
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
func (d *Default) InsertOrUpdate(ctx context.Context, filter interface{}, data interface{}) error {
	// Run in transaction
	return Transaction(ctx, d.session, func() error {
		_, err := d.session.Database(d.db).Collection(d.table).UpdateOne(ctx, filter, bson.M{"$set": data}, options.Update().SetUpsert(true))
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
	// Run in transaction
	return Transaction(ctx, d.session, func() error {
		_, err := d.session.Database(d.db).Collection(d.table).DeleteOne(ctx, id)
		return err
	})
}

// Find searches for a resource in the database and then returns a cursor
func (d *Default) Find(ctx context.Context, id interface{}, value interface{}) error {
	res, err := d.session.Database(d.db).Collection(d.table).Find(ctx, bson.M{
		"_id": id,
	})
	if err != nil {
		return err
	}
	return res.Decode(value)
}

// FindFetchOne searches for a resource and then unmarshals the first row into value
func (d *Default) FindFetchOne(ctx context.Context, id string, value interface{}) error {
	res := d.session.Database(d.db).Collection(d.table).FindOne(ctx, id)
	return res.Decode(value)
}

// FindOneBy is an utility for fetching values if they are stored in a key-value manenr.
func (d *Default) FindOneBy(ctx context.Context, key string, value interface{}, result interface{}) error {
	res := d.session.Database(d.db).Collection(d.table).FindOne(ctx, bson.M{
		key: value,
	})
	return res.Decode(result)
}

// FindBy is an utility for fetching values if they are stored in a key-value manenr.
func (d *Default) FindBy(ctx context.Context, key string, value interface{}, results interface{}) error {
	res, err := d.session.Database(d.db).Collection(d.table).Find(ctx, bson.M{
		key: value,
	})
	if err != nil {
		return err
	}
	return res.Decode(results)
}

// FindByAndCount returns the number of elements that match the filter
func (d *Default) FindByAndCount(ctx context.Context, key string, value interface{}) (int64, error) {
	return d.session.Database(d.db).Collection(d.table).CountDocuments(ctx, bson.M{
		key: value,
	})
}

// FindByAndFetch retrieves a value by key and then fills results with the result.
func (d *Default) FindByAndFetch(ctx context.Context, key string, value interface{}, results interface{}) error {
	res, err := d.session.Database(d.db).Collection(d.table).Find(ctx, bson.M{
		key: value,
	})
	if err != nil {
		return err
	}
	return res.Decode(results)
}

// WhereCount allows counting with multiple fields
func (d *Default) WhereCount(ctx context.Context, filter interface{}) (int64, error) {
	return d.session.Database(d.db).Collection(d.table).CountDocuments(ctx, filter)
}

// Where allows filtering with multiple fields
func (d *Default) Where(ctx context.Context, filter interface{}, results interface{}) error {
	res, err := d.session.Database(d.db).Collection(d.table).Find(ctx, filter)
	if err != nil {
		return err
	}
	return res.Decode(results)
}

// WhereAndFetchLimit filters with multiple fields and then fills results with all found resources
func (d *Default) WhereAndFetchLimit(ctx context.Context, filter interface{}, paginator *db.Pagination, results interface{}) error {
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
}

// WhereAndFetchOne filters with multiple fields and then fills result with the first found resource
func (d *Default) WhereAndFetchOne(ctx context.Context, filter interface{}, result interface{}) error {
	res := d.session.Database(d.db).Collection(d.table).FindOne(ctx, filter)
	return res.Decode(result)
}

// List all entities from the database
func (d *Default) List(ctx context.Context, sortParams *db.SortParameters, pagination *db.Pagination, results interface{}) (uint, error) {
	// Run in transaction
	return d.Search(ctx, bson.M{}, sortParams, pagination, results)
}

// Search all entities from the database
func (d *Default) Search(ctx context.Context, filter interface{}, sortParams *db.SortParameters, pagination *db.Pagination, results interface{}) (uint, error) {
	// Apply Filter
	if filter == nil {
		filter = bson.M{}
	}

	// Get total
	count, err := d.WhereCount(ctx, filter)
	if err != nil {
		return 0, xerrors.Errorf("mongodb: unable to count element before search: %w", err)
	}
	// If no result skip data request
	if count == 0 {
		return 0, db.ErrNoResult
	}

	if pagination != nil {
		pagination.SetTotal(uint(count))
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

	// Do the query
	cur, err := d.session.Database(d.db).Collection(d.table).Find(ctx, filter, opts)
	if err != nil {
		return 0, xerrors.Errorf("mongodb: unable to query collection: %w", err)
	}

	// Extract all entities
	if err := cur.All(ctx, results); err != nil {
		return 0, xerrors.Errorf("mongodb: unable to extract entities: %w", err)
	}

	// Check if cursor has errors
	if err := cur.Err(); err != nil {
		return 0, xerrors.Errorf("mongodb: cursor has error: %w", err)
	}

	// Close the cursor
	if err := cur.Close(ctx); err != nil {
		return 0, xerrors.Errorf("mongodb: unable to close cursor: %w", err)
	}

	// Return no error
	return uint(count), nil
}

// -----------------------------------------------------------------------------

// TransactionFunc is the transaction handler closure contract
type TransactionFunc func() error

// Transaction runs the transactionfunc in a transaction
func Transaction(ctx context.Context, client *mongowrapper.WrappedClient, fn TransactionFunc) error {
	return client.UseSession(ctx, func(sctx mongo.SessionContext) error {
		// Start transaction
		if err := sctx.StartTransaction(); err != nil {
			return xerrors.Errorf("mongodb: %w", err)
		}

		// Run the closure
		if err := fn(); err != nil {
			log.CheckErrCtx(sctx, "Unable to abort transaction", sctx.AbortTransaction(sctx))
			return xerrors.Errorf("mongodb: %w", err)
		}

		for {
			// Commit the transaction
			err := sctx.CommitTransaction(sctx)
			if err == nil {
				return nil
			}

			if cmdErr, ok := err.(mongo.CommandError); ok && cmdErr.HasErrorLabel("TransientTransactionError") {
				log.For(ctx).Warn("Transient error occurred, retrying transaction ...")
				continue
			}
			return err
		}
	})
}
