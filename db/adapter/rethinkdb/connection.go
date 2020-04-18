package rethinkdb

import (
	"context"
	"fmt"

	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
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
		return nil, fmt.Errorf("rethinkdb: unable to connect to server: %w", err)
	}

	// Return connection
	return conn, nil
}
