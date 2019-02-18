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
