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

package rethinkdb

import (
	"context"

	"github.com/pkg/errors"

	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

// Configuration repesents database connection configuration
type Configuration struct {
	AutoMigrate bool
	Addresses   []string
	Database    string
	Username    string
	Password    string
	AuthKey     string
}

// Connection provides Wire provider for a PostgreSQL database connection
func Connection(ctx context.Context, cfg *Configuration) (*r.Session, error) {
	// Prepare options
	opts := r.ConnectOpts{
		Addresses:  cfg.Addresses,
		Database:   cfg.Database,
		InitialCap: 10,
		MaxOpen:    10,
		NumRetries: 5,
	}

	// Optional parameters
	if cfg.AuthKey != "" {
		opts.AuthKey = cfg.AuthKey
	}
	if cfg.Username != "" {
		opts.Username = cfg.Username
	}
	if cfg.Password != "" {
		opts.Password = cfg.Password
	}

	// Activate Opentracing
	opts.UseOpentracing = true

	// Initialize a new setup connection
	conn, err := r.Connect(opts)
	if err != nil {
		return nil, errors.Wrap(err, "RethinkDB error")
	}

	// Return connection
	return conn, nil
}
