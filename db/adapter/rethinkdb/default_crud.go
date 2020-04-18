package rethinkdb

import (
	"context"
	"fmt"

	"go.zenithar.org/pkg/db"

	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

// Default contains the basic implementation of the EntityCRUD interface
type Default struct {
	table   string
	db      string
	session *r.Session
}

// NewCRUDTable sets up a new Default struct
func NewCRUDTable(session *r.Session, db, table string) *Default {
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

// GetTable returns no table
func (d *Default) GetTable() interface{} {
	return r.Table(d.table)
}

// GetSession returns the current session
func (d *Default) GetSession() interface{} {
	return d.session
}

// -----------------------------------------------------------------------------

// Insert inserts a document into the database
func (d *Default) Insert(ctx context.Context, data interface{}) error {
	_, err := r.Table(d.table).Insert(data).RunWrite(d.session, r.RunOpts{
		Context: ctx,
	})
	if err != nil {
		if err == r.ErrEmptyResult {
			return db.ErrNoResult
		}
		return fmt.Errorf("rethinkdb: unable to execute query: %w", err)
	}

	return nil
}

// InsertOrUpdate a document occording to ID presence in database
func (d *Default) InsertOrUpdate(ctx context.Context, id interface{}, data interface{}) error {
	_, err := r.Table(d.table).Insert(data, r.InsertOpts{Conflict: "update"}).RunWrite(d.session, r.RunOpts{
		Context: ctx,
	})
	if err != nil {
		if err == r.ErrEmptyResult {
			return db.ErrNoResult
		}
		return fmt.Errorf("rethinkdb: unable to execute query: %w", err)
	}

	return nil
}

// Find a document match given id
func (d *Default) Find(ctx context.Context, id interface{}, value interface{}) error {
	cursor, err := r.Table(d.table).Get(id).Run(d.session, r.RunOpts{
		Context: ctx,
	})
	if err != nil {
		return fmt.Errorf("rethinkdb: unable to execute query: %w", err)
	}

	if err := cursor.One(value); err != nil {
		if err == r.ErrEmptyResult {
			return db.ErrNoResult
		}
		return fmt.Errorf("rethinkdb: unable to retrieve query result: %w", err)
	}

	return nil
}

// FindOneBy a couple (k = v) in the database
func (d *Default) FindOneBy(ctx context.Context, key string, value interface{}, result interface{}) error {
	cursor, err := r.Table(d.table).GetAllByIndex(key, value).Run(d.session, r.RunOpts{
		Context: ctx,
	})
	if err != nil {
		return fmt.Errorf("rethinkdb: unable to execute query: %w", err)
	}

	if err := cursor.One(result); err != nil {
		if err == r.ErrEmptyResult {
			return db.ErrNoResult
		}
		return fmt.Errorf("rethinkdb: unable to retrieve query result: %w", err)
	}

	return nil
}

// FindBy all couples (k = v) in the database
func (d *Default) FindBy(ctx context.Context, key string, value interface{}, results interface{}) error {
	cursor, err := r.Table(d.table).Filter(func(row r.Term) r.Term {
		return row.Field(key).Eq(value)
	}).Run(d.session, r.RunOpts{
		Context: ctx,
	})
	if err != nil {
		return fmt.Errorf("rethinkdb: unable to execute query: %w", err)
	}

	if err := cursor.All(results); err != nil {
		if err == r.ErrEmptyResult {
			return db.ErrNoResult
		}
		return fmt.Errorf("rethinkdb: unable to retrieve query result: %w", err)
	}

	return nil
}

// FindByAndCount is used to count object that matchs the (key = value) predicate
func (d *Default) FindByAndCount(ctx context.Context, key string, value interface{}) (int, error) {
	cursor, err := r.Table(d.table).Filter(func(row r.Term) r.Term {
		return row.Field(key).Eq(value)
	}).Count().Run(d.session, r.RunOpts{
		Context: ctx,
	})
	if err != nil {
		return 0, fmt.Errorf("rethinkdb: unable to execute query: %w", err)
	}

	var count int
	if err := cursor.One(&count); err != nil {
		if err == r.ErrEmptyResult {
			return 0, db.ErrNoResult
		}
		return 0, fmt.Errorf("rethinkdb: unable to retrieve query result: %w", err)
	}

	return count, nil
}

// Where is used to fetch documents that match th filter from the database
func (d *Default) Where(ctx context.Context, filter interface{}, results interface{}) error {
	cursor, err := r.Table(d.table).Filter(filter).Run(d.session, r.RunOpts{
		Context: ctx,
	})
	if err != nil {
		return fmt.Errorf("rethinkdb: unable to execute query: %w", err)
	}

	if err := cursor.All(results); err != nil {
		if err == r.ErrEmptyResult {
			return db.ErrNoResult
		}
		return fmt.Errorf("rethinkdb: unable to retrieve query result: %w", err)
	}

	return nil
}

// WhereCount returns the document count that match the filter
func (d *Default) WhereCount(ctx context.Context, filter interface{}) (int, error) {
	cursor, err := r.Table(d.table).Filter(filter).Count().Run(d.session, r.RunOpts{
		Context: ctx,
	})
	if err != nil {
		return 0, fmt.Errorf("rethinkdb: unable to execute query: %w", err)
	}

	var count int
	if err := cursor.One(&count); err != nil {
		if err == r.ErrEmptyResult {
			return 0, db.ErrNoResult
		}
		return 0, fmt.Errorf("rethinkdb: unable to retrieve query result: %w", err)
	}

	return count, nil
}

// WhereAndFetchOne returns one document that match the filter
func (d *Default) WhereAndFetchOne(ctx context.Context, filter interface{}, result interface{}) error {
	cursor, err := r.Table(d.table).Filter(filter).Run(d.session, r.RunOpts{
		Context: ctx,
	})
	if err != nil {
		return fmt.Errorf("rethinkdb: unable to execute query: %w", err)
	}

	if err := cursor.One(result); err != nil {
		if err == r.ErrEmptyResult {
			return db.ErrNoResult
		}
		return fmt.Errorf("rethinkdb: unable to retrieve query result: %w", err)
	}

	return nil
}

// WhereAndFetchLimit returns paginated list of document
func (d *Default) WhereAndFetchLimit(ctx context.Context, filter interface{}, paginator *db.Pagination, results interface{}) error {
	cursor, err := r.Table(d.table).Filter(filter).Run(d.session, r.RunOpts{
		Context: ctx,
	})
	if err != nil {
		return fmt.Errorf("rethinkdb: unable to execute query: %w", err)
	}

	if err := cursor.All(results); err != nil {
		if err == r.ErrEmptyResult {
			return db.ErrNoResult
		}
		return fmt.Errorf("rethinkdb: unable to retrieve query result: %w", err)
	}

	return nil
}

// Update a document that match the selector
func (d *Default) Update(ctx context.Context, selector interface{}, data interface{}) error {
	_, err := r.Table(d.table).Filter(selector).Update(data).RunWrite(d.session, r.RunOpts{
		Context: ctx,
	})
	if err != nil {
		if err == r.ErrEmptyResult {
			return db.ErrNoResult
		}
		return fmt.Errorf("rethinkdb: unable to execute query: %w", err)
	}

	return nil
}

// UpdateID updates a document using his id
func (d *Default) UpdateID(ctx context.Context, id interface{}, data interface{}) error {
	_, err := r.Table(d.table).Get(id).Update(data).RunWrite(d.session, r.RunOpts{
		Context: ctx,
	})
	if err != nil {
		if err == r.ErrEmptyResult {
			return db.ErrNoResult
		}
		return fmt.Errorf("rethinkdb: unable to execute query: %w", err)
	}

	return nil
}

// DeleteAll documents from the database
func (d *Default) DeleteAll(ctx context.Context, pred interface{}) error {
	_, err := r.Table(d.table).Filter(pred).Delete().RunWrite(d.session, r.RunOpts{
		Context: ctx,
	})
	if err != nil {
		if err == r.ErrEmptyResult {
			return db.ErrNoResult
		}
		return fmt.Errorf("rethinkdb: unable to execute query: %w", err)
	}

	return nil
}

// Delete a document from the database
func (d *Default) Delete(ctx context.Context, id interface{}) error {
	_, err := r.Table(d.table).Get(id).Delete().RunWrite(d.session, r.RunOpts{
		Context: ctx,
	})
	if err != nil {
		if err == r.ErrEmptyResult {
			return db.ErrNoResult
		}
		return fmt.Errorf("rethinkdb: unable to execute query: %w", err)
	}

	return nil
}

// List all entities from the database
func (d *Default) List(ctx context.Context, results interface{}, sortParams *db.SortParameters, pagination *db.Pagination) error {
	return d.Search(ctx, results, nil, sortParams, pagination)
}

// Search all entities in the database
func (d *Default) Search(ctx context.Context, results interface{}, filter interface{}, sortParams *db.SortParameters, pagination *db.Pagination) error {
	term := r.Table(d.table)

	// Filter
	if filter != nil {
		term = term.Filter(filter)
	}

	// Get total
	if pagination != nil {
		total, err := d.WhereCount(ctx, filter)
		if err != nil {
			return fmt.Errorf("rethinkdb: unable to execute query: %w", err)
		}
		pagination.SetTotal(uint(total))
	}

	// Sort
	if sortParams != nil {
		term = term.OrderBy(ConvertSortParameters(*sortParams)...)
	}

	// Slice result
	if pagination != nil {
		term = term.Slice(pagination.Offset(), pagination.Offset()+pagination.PerPage)
	}

	// Run the query
	cursor, err := term.Run(d.session, r.RunOpts{
		Context: ctx,
	})
	if err != nil {
		return fmt.Errorf("rethinkdb: unable to execute query: %w", err)
	}

	// Fetch cursor
	err = cursor.All(results)
	if err != nil {
		if err == r.ErrEmptyResult {
			return db.ErrNoResult
		}
		return fmt.Errorf("rethinkdb: unable to retrieve query result: %w", err)
	}

	return nil
}
