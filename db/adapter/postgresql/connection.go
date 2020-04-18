package postgresql

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.zenithar.org/pkg/log"

	"github.com/jmoiron/sqlx"
	"github.com/opencensus-integrations/ocsql"
	try "gopkg.in/matryer/try.v1"

	// Load postgresql drivers
	_ "github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
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
			return false, fmt.Errorf("postgresql: %w", err)
		}

		defaultDriver := "postgres"
		// Check driver option presence
		if drv, ok := connStr.Options["driver"]; ok {

			// Remove from connection string
			delete(connStr.Options, "driver")

			// Check usages
			switch drv {
			case "postgres", "pgx":
				defaultDriver = drv
			default:
				return false, errors.New("postgresql: invalid 'driver' option value, 'postgres' or 'pgx' supported")
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
			return false, fmt.Errorf("postgresql: failed to register ocsql driver: %w", err)
		}

		// Connect to database
		conn, err = sqlx.Open(driverName, connStr.String())
		if err != nil {
			return attempt < 10, fmt.Errorf("postgresql: unable to open driver: %w", err)
		}

		// Check connection
		if err = conn.Ping(); err != nil {
			return attempt < 10, fmt.Errorf("postgresql: unable to ping database: %w", err)
		}

		// Update connection pool settings
		conn.SetConnMaxLifetime(5 * time.Minute)
		conn.SetMaxIdleConns(0)
		conn.SetMaxOpenConns(95)

		log.For(ctx).Info("PostGreSQL connected !")

		return false, nil
	})
	if err != nil {
		return nil, fmt.Errorf("postgresql: unable to connect to database: %w", err)
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
