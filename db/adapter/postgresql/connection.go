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
	"time"

	"go.zenithar.org/pkg/log"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	try "gopkg.in/matryer/try.v1"
)

// Configuration repesents database connection configuration
type Configuration struct {
	AutoMigrate      bool
	ConnectionString string
	Username         string
	Password         string
}

// Connection provides Wire provider for a PostgreSQL database connection
func Connection(ctx context.Context, cfg *Configuration) (*sqlx.DB, error) {
	var conn *sqlx.DB
	err := try.Do(func(attempt int) (bool, error) {
		var err error

		connStr, err := ParseURL(cfg.ConnectionString)
		if err != nil {
			return false, errors.Wrap(err, "PosgreSQL error")
		}

		// Overrides settings
		connStr.User = cfg.Username
		connStr.Password = cfg.Password

		// Connect to database
		conn, err = sqlx.Open("postgres", connStr.String())
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
		return nil, err
	}

	// Return connection
	return conn, nil
}
