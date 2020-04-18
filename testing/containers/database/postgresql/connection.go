package postgresql

import (
	"context"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	dockertest "github.com/ory/dockertest/v3"
	"golang.org/x/xerrors"
)

// Connect returns a PostgreSQL connection form a container or a running instance
func Connect(_ context.Context) (*sqlx.DB, *Configuration, error) {

	// Check environment variable first
	if url := os.Getenv("TEST_DATABASE_POSTGRESQL"); url != "" {

		// Check URL syntax
		u, err := postgresql.ParseURL(url)
		if err != nil {
			return nil, nil, xerrors.Errorf("testing: unable to parse PostgreSQL DSN: %w", err)
		}

		defaultDriver := "postgres"
		// Check driver option presence
		if drv, ok := u.Options["driver"]; ok {
			switch drv {
			case "postgres", "pgx":
				defaultDriver = drv
			default:
				return nil, nil, xerrors.New("testing: invalid 'driver' option value, 'postgres' or 'pgx' supported")
			}
		}

		// Try to connect
		log.Println("Found postgresql test database config, skipping dockertest...")
		db, err := sqlx.Open(defaultDriver, u.String())
		if err != nil {
			return nil, nil, xerrors.Errorf("testing: unable to connect to postgresql: %w", err)
		}

		// Return connection
		return db, &Configuration{
			ConnectionString: u.String(),
			DatabaseName:     u.Database,
			DatabaseUser:     u.User,
			Password:         u.Password,
		}, nil
	}

	// Initialize a docker container
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, xerrors.Errorf("testing: unable to initialize docker connection: %w", err)
	}

	// Build postgres container
	container := newPostgresContainer(pool)

	var db *sqlx.DB

	// Wait for connection
	if err = pool.Retry(func() error {
		var err error

		// Try to connect using connection string
		db, err = sqlx.Open("postgres", container.config.ConnectionString)
		if err != nil {
			return xerrors.Errorf("testing: unable to initialize postgresql driver: %w", err)
		}

		// Ping database
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// Everything is ready
	log.Printf("Postgres (%v): up", container.Name)

	// Add container to resources
	resources = append(resources, container)

	// Return connection
	return db, container.Configuration(), nil
}
