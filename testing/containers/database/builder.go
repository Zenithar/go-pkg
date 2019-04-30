package database

import (
	"context"
	"io"
	"log"
	"os"

	// Load postgresql drivers
	_ "github.com/jackc/pgx"
	_ "github.com/jackc/pgx/pgtype"
	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	mongowrapper "github.com/opencensus-integrations/gomongowrapper"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.zenithar.org/pkg/db/adapter/postgresql"
	"golang.org/x/xerrors"
	dockertest "gopkg.in/ory-am/dockertest.v3"
)

var (
	resources []io.Closer
)

// KillAll resources allocated via docker
func KillAll(ctx context.Context) {
	log.Println("Killing all docker resources ...")

	for _, resource := range resources {
		if err := resource.Close(); err != nil {
			log.Printf("Go error while closing the resource, %v", err)
		}
	}

	// Clean the resource array
	resources = []io.Closer{}
}

// -----------------------------------------------------------------------------

// ConnectToPostgreSQL returns a PostgreSQL connection form a container or a running instance
func ConnectToPostgreSQL(_ context.Context) (*sqlx.DB, *Configuration, error) {

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

// ConnectToMongoDB returns a MongoDB connection
func ConnectToMongoDB(ctx context.Context) (*mongowrapper.WrappedClient, error) {
	// Check environment variable first
	if url := os.Getenv("TEST_DATABASE_MONGODB"); url != "" {
		log.Println("Found mongodb test database config, skipping dockertest...")

		// Extract database name from connection string
		client, err := mongowrapper.Connect(ctx, options.Client().ApplyURI(url))
		if err != nil {
			return nil, xerrors.Errorf("testing: unable to connect to mongodb: %w", err)
		}

		// Return connection
		return client, nil
	}

	// Initialize a docker container
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, xerrors.Errorf("testing: unable to initialize docker connection: %w", err)
	}

	// Build mongo container
	container := newMongoDBContainer(pool)

	var db *mongowrapper.WrappedClient

	// Wait for connection
	if err = pool.Retry(func() error {
		var err error

		// Extract database name from connection string
		db, err = mongowrapper.Connect(ctx, options.Client().ApplyURI(container.ConnectionString))
		if err != nil {
			return xerrors.Errorf("testing: unable to connect to mongodb: %w", err)
		}

		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// Everything is ready
	log.Printf("MongoDB (%v): up", container.Name)

	// Add container to resources
	resources = append(resources, container)

	// Return connection
	return db, nil
}
