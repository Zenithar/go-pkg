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

package postgresql

import (
	"context"
	"database/sql"
	"reflect"

	"go.zenithar.org/pkg/db"
	"go.zenithar.org/pkg/log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pkg/errors"
)

// Default contains the basic implementation of the SQL interface
type Default struct {
	table   string
	db      string
	session *sqlx.DB

	mapper          *reflectx.Mapper
	columns         []string
	sortableColumns map[string]bool
}

// NewCRUDTable sets up a new Default struct
func NewCRUDTable(session *sqlx.DB, db, table string, columns, sortable []string) *Default {
	sortableColumns := map[string]bool{}
	for _, column := range sortable {
		sortableColumns[column] = true
	}

	return &Default{
		db:              db,
		table:           table,
		session:         session,
		mapper:          reflectx.NewMapper("db"),
		columns:         columns,
		sortableColumns: sortableColumns,
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

// GetSession returns the current session
func (d *Default) GetSession() interface{} {
	return d.session
}

// -----------------------------------------------------------------------------

// Create a record
func (d *Default) Create(ctx context.Context, data interface{}) error {

	// Extract columns and values
	columns, values := d.extractColumnPairs(data)

	// Prepare query
	query := sq.Insert(d.table).
		Columns(columns...).
		Values(values...).
		PlaceholderFormat(sq.Dollar)

	// Build sql query
	q, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "Database error")
	}

	// Instrument opentracing
	if span := opentracing.SpanFromContext(ctx); span != nil {
		ext.DBType.Set(span, "sql")
		ext.DBStatement.Set(span, q)
	}

	// Prepare the statement
	stmt, err := d.session.PreparexContext(ctx, q)
	if err != nil {
		return errors.Wrap(err, "Database error")
	}
	defer func(stmt *sqlx.Stmt) {
		log.SafeClose(stmt, "Unable to close statement")
	}(stmt)

	// Do the insert query
	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		return errors.Wrap(err, "Database error")
	}

	return nil
}

// WhereCount is used to cound resultset elements from the given filter
func (d *Default) WhereCount(ctx context.Context, filter interface{}) (int, error) {
	// Prepare query
	qb := sq.Select("COUNT(*) as count").
		From(d.table).
		PlaceholderFormat(sq.Dollar)

	if filter != nil {
		qb = qb.Where(filter)
	}

	// Build sql query
	q, args, err := qb.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "SQL error")
	}

	// Prepare the statement
	stmt, err := d.session.PreparexContext(ctx, q)
	if err != nil {
		return 0, errors.Wrap(err, "Database error")
	}
	defer func(stmt *sqlx.Stmt) {
		log.SafeClose(stmt, "Unable to close statement")
	}(stmt)

	var count int
	if err := stmt.QueryRowContext(ctx, args...).Scan(&count); err == sql.ErrNoRows {
		return 0, db.ErrNoResult
	} else if err != nil {
		return 0, errors.Wrap(err, "Database error")
	}

	// Return no error
	return count, nil
}

// WhereAndFetchOne returns only one element from the given filter
func (d *Default) WhereAndFetchOne(ctx context.Context, filter interface{}, result interface{}) error {
	// Prepare query
	qb := sq.Select(d.columns...).
		From(d.table).
		Where(filter).
		Limit(1).
		PlaceholderFormat(sq.Dollar)

	// Build sql query
	q, args, err := qb.ToSql()
	if err != nil {
		return errors.Wrap(err, "SQL error")
	}

	// Instrument opentracing
	if span := opentracing.SpanFromContext(ctx); span != nil {
		ext.DBType.Set(span, "sql")
		ext.DBStatement.Set(span, q)
	}

	// Prepare the statement
	stmt, err := d.session.PreparexContext(ctx, q)
	if err != nil {
		return errors.Wrap(err, "Database error")
	}
	defer func(stmt *sqlx.Stmt) {
		log.SafeClose(stmt, "Unable to close statement")
	}(stmt)

	// Do the insert query
	err = d.session.QueryRowxContext(ctx, q, args...).StructScan(result)
	if err == sql.ErrNoRows {
		return db.ErrNoResult
	} else if err != nil {
		return errors.Wrap(err, "Database error")
	}

	// Return no error
	return nil
}

// Update the collection element with updates set matching the given filter
func (d *Default) Update(ctx context.Context, updates map[string]interface{}, filter interface{}) error {
	// Prepare query
	qb := sq.Update(d.table).
		SetMap(updates).
		Where(filter).
		PlaceholderFormat(sq.Dollar)

	// Build sql query
	q, args, err := qb.ToSql()
	if err != nil {
		return errors.Wrap(err, "SQL error")
	}

	// Instrument opentracing
	if span := opentracing.SpanFromContext(ctx); span != nil {
		ext.DBType.Set(span, "sql")
		ext.DBStatement.Set(span, q)
	}

	// Prepare the statement
	stmt, err := d.session.PreparexContext(ctx, q)
	if err != nil {
		return errors.Wrap(err, "Database error")
	}
	defer func(stmt *sqlx.Stmt) {
		log.SafeClose(stmt, "Unable to close statement")
	}(stmt)

	// Do the insert query
	res, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return errors.Wrap(err, "Database error")
	}

	// Check updates
	count, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "Database error")
	}

	// If no rows where affected return an handled error
	if count == 0 {
		return db.ErrNoModification
	}

	// Return no error
	return nil
}

// RemoveOne is used to remove one element from the collection that match the filter
func (d *Default) RemoveOne(ctx context.Context, filter interface{}) error {

	// Prepare query
	qb := sq.Delete(d.table).
		Where(filter).
		PlaceholderFormat(sq.Dollar)

	// Build sql query
	q, args, err := qb.ToSql()
	if err != nil {
		return errors.Wrap(err, "SQL error")
	}

	// Instrument opentracing
	if span := opentracing.SpanFromContext(ctx); span != nil {
		ext.DBType.Set(span, "sql")
		ext.DBStatement.Set(span, q)
	}

	// Prepare the statement
	stmt, err := d.session.PreparexContext(ctx, q)
	if err != nil {
		return errors.Wrap(err, "Database error")
	}
	defer func(stmt *sqlx.Stmt) {
		log.SafeClose(stmt, "Unable to close statement")
	}(stmt)

	// Do the insert query
	res, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return errors.Wrap(err, "Database error")
	}

	// Check updates
	count, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "Database error")
	}

	// If no rows where affected return an handled error
	if count == 0 {
		return db.ErrNoModification
	}

	// Return no error
	return nil
}

// Search for element in collection
func (d *Default) Search(ctx context.Context, filter interface{}, pagination *db.Pagination, sortParams *db.SortParameters, results interface{}) (int, error) {
	// Initialize statement
	q := sq.Select(d.columns...).
		From(d.table).
		PlaceholderFormat(sq.Dollar)

	// Count result set first
	count, err := d.WhereCount(ctx, filter)
	if err != nil {
		return 0, err
	}

	// If no result skip data request
	if count == 0 {
		return 0, db.ErrNoResult
	}

	if pagination != nil {
		pagination.SetTotal(uint(count))
	}

	if filter != nil {
		// Prepare the query
		q = q.Where(filter)
	}

	// Apply pagination on data query only
	if pagination != nil {
		q = q.Offset(uint64(pagination.Offset())).Limit(uint64(pagination.PerPage))
	}

	// Apply sort parameters
	if sortParams != nil {
		q = q.OrderBy(ConvertSortParameters(*sortParams, d.sortableColumns)...)
	}

	// Do the query
	sqlData, args, err := q.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "Unable to prepare search query")
	}

	// Prepare the statement
	stmt, err := d.session.PreparexContext(ctx, sqlData)
	if err != nil {
		return 0, errors.Wrap(err, "Database error")
	}
	defer func(stmt *sqlx.Stmt) {
		log.SafeClose(stmt, "Unable to close statement")
	}(stmt)

	if err := stmt.SelectContext(ctx, results, args...); err == sql.ErrNoRows {
		return 0, db.ErrNoResult
	} else if err != nil {
		return 0, errors.Wrap(err, "Database error")
	}

	// Return no error
	return count, nil
}

// -----------------------------------------------------------------------------

func (d *Default) extractColumnPairs(data interface{}) ([]string, []interface{}) {
	// Create type mapper
	valueMap := d.mapper.FieldMap(reflect.ValueOf(data))

	// Extract columns
	var columns []string
	var values []interface{}
	for column, value := range valueMap {
		columns = append(columns, column)
		values = append(values, value.Interface())
	}

	// Return all elements
	return columns, values
}
