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
	"sync"
	"time"

	// Load postgresql drivers
	_ "github.com/jackc/pgx"
	_ "github.com/jackc/pgx/pgtype"
	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/lib/pq"

	"go.zenithar.org/pkg/log"

	"github.com/jmoiron/sqlx"
	"github.com/opencensus-integrations/ocsql"
	"github.com/pkg/errors"
	try "gopkg.in/matryer/try.v1"
)

var (
	once sync.Once
	conn *sqlx.DB
)

// Configuration represents database connection configuration
type Configuration struct {
	AutoMigrate      bool
	ConnectionString string
	Username         string
	Password         string
}

// Connection provides Wire provider for a PostgreSQL database connection
func Connection(ctx context.Context, cfg *Configuration) (*sqlx.DB, error) {

	err := try.Do(func(attempt int) (bool, error) {
		var err error

		connStr, err := ParseURL(cfg.ConnectionString)
		if err != nil {
			return false, errors.Wrap(err, "PosgreSQL error")
		}

		defaultDriver := "postgres"
		// Check driver option presence
		if drv, ok := connStr.Options["driver"]; ok {
			switch drv {
			case "postgres", "pgx":
				defaultDriver = drv
			default:
				return false, errors.New("invalid 'driver' option value, 'postgres' or 'pgx' supported")
			}
		}

		// Overrides settings
		connStr.User = cfg.Username
		connStr.Password = cfg.Password

		// Instrument with opentracing
		driverName, err := ocsql.Register(
			defaultDriver,
			ocsql.WithOptions(ocsql.TraceOptions{
				AllowRoot:    false,
				Ping:         true,
				RowsNext:     true,
				RowsClose:    true,
				RowsAffected: true,
				LastInsertID: true,
				Query:        true,
				QueryParams:  false,
			}),
		)
		if err != nil {
			return false, errors.Wrap(err, "failed to register ocsql driver")
		}

		// Connect to database
		conn, err = sqlx.Open(driverName, connStr.String())
		if err != nil {
			return attempt < 10, errors.Wrap(err, "PostgreSQL error : "+connStr.String())
		}

		// Check connection
		if err = conn.Ping(); err != nil {
			return attempt < 10, errors.Wrap(err, "PostgreSQL error : "+connStr.String())
		}

		// Update connection pool settings
		conn.SetConnMaxLifetime(5 * time.Minute)
		conn.SetMaxIdleConns(0)
		conn.SetMaxOpenConns(95)

		log.For(ctx).Info("PostGreSQL connected !")

		return false, nil
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	once.Do(func() {
		// Start statistic puller
		dbstatsCloser := ocsql.RecordStats(conn.DB, 5*time.Second)

		go func() {
			select {
			case <-ctx.Done():
				dbstatsCloser()
				log.SafeClose(conn, "Unable to close database connection")
			}
		}()
	})

	// Return connection
	return conn, nil
}
