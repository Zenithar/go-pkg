package mongodb

import (
	"context"
	"log"
	"os"

	dockertest "github.com/ory/dockertest/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/xerrors"
)

// Connect returns a MongoDB connection
func Connect(ctx context.Context) (*mongo.Client, error) {
	// Check environment variable first
	if url := os.Getenv("TEST_DATABASE_MONGODB"); url != "" {
		log.Println("Found mongodb test database config, skipping dockertest...")

		// Extract database name from connection string
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
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

	var db *mongo.Client

	// Wait for connection
	if err = pool.Retry(func() error {
		var err error

		// Extract database name from connection string
		db, err = mongo.Connect(ctx, options.Client().ApplyURI(container.ConnectionString))
		if err != nil {
			return xerrors.Errorf("testing: unable to connect to mongodb: %w", err)
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
