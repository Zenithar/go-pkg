package database

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	mongowrapper "github.com/opencensus-integrations/gomongowrapper"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.zenithar.org/pkg/db/adapter/postgresql"
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
	resources = []*dockertest.Resource{}
}

// -----------------------------------------------------------------------------

// ConnectToPostgreSQL returns a PostgreSQL connection form a container or a running instance
func ConnectToPostgreSQL(_ context.Context) (*sqlx.DB, *Configuration, error) {

	// Check environment variable first
	if url := os.Getenv("TEST_DATABASE_POSTGRESQL"); url != "" {

		// Check URL syntax
		u, err := postgresql.ParseURL(url)
		if err != nil {
			return nil, nil, errors.Wrap(err, "unable to parse PostgreSQL DSN")
		}

		// Try to connect
		log.Println("Found postgresql test database config, skipping dockertest...")
		db, err := sqlx.Open("postgres", u.String())
		if err != nil {
			return nil, nil, errors.Wrap(err, "unable to initialize PostgreSQL connection")
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
		return nil, nil, err
	}

	// Build postgres container
	container := newPostgresContainer(pool)

	var db sqlx.DB

	// Wait for connection
	if err = pool.Retry(func() error {
		var err error

		// Try to connect using connection string
		db, err = sqlx.Open("postgres", container.config.ConnectionString)
		if err != nil {
			return err
		}

		// Ping database
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// Everything is ready
	log.Printf("Postgres (%v): up", container.Name)

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
			return nil, errors.Wrap(err, "unable to connect to MongoDB")
		}

		// Return connection
		return client, nil
	}

	// Initialize a docker container
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, err
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
			return errors.Wrap(err, "unable to connect to MongoDB")
		}

		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// Everything is ready
	log.Printf("MongoDB (%v): up", container.Name)

	// Return connection
	return db, nil
}
