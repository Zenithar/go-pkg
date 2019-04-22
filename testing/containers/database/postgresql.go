package database

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
	"github.com/dchest/uniuri"
	dockertest "gopkg.in/ory-am/dockertest.v3"

	"go.zenithar.org/pkg/testing/containers"
	"go.zenithar.org/pkg/db/adapter/postgresql"
)

var (
	// PostgreSQLVersion defines version to use
	PostgreSQLVersion = "10"
)

// PostgreSQLContainer represents database container handler
type PostgreSQLContainer struct {
	Name             string
	ConnectionString string
	Password         string
	DatabaseName     string
	DatabaseUser     string
	pool             *dockertest.Pool
	resource         *dockertest.Resource
}

// NewPostgresContainer initialize a PostgreSQL server in a docker container
func NewPostgresContainer(pool *dockertest.Pool) *PostgreSQLContainer {

	var (
		db *sqlx.DB
		databaseName = fmt.Sprintf("test-%s", uniuri.NewLen(8))
		databaseUser = fmt.Sprintf("user-%s", uniuri.NewLen(8))
		password = uniuri.NewLen(32)
	)

	// Initialize a PostgreSQL server
	resource, err := pool.Run("postgres", PostgreSQLVersion, []string{
		fmt.Sprintf("POSTGRES_PASSWORD=%s", password),
		fmt.Sprintf("POSTGRES_DB=%s", databaseName),
		fmt.Sprintf("POSTGRES_USER=%s", databaseUser),
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// Prepare connection string
	connectionString := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", databaseUser, password, resource.GetPort("5432/tcp"), databaseName)
	
	// Retrieve container name
	containerName := containers.GetName(resource)

	// Wait for connection
	if err = pool.Retry(func() error {
		var err error

		// Try to connect using connection string
		db, err = sqlx.Open("postgres", connectionString)
		if err != nil {
			return err
		}

		// Ping database
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// Everything is ready
	log.Printf("Postgres (%v): up", containerName)

	// Return container information
	return &PostgreSQLContainer{
		Name:             containerName,
		ConnectionString: connectionString,
		Password:         password,
		DatabaseName:     databaseName,
		DatabaseUser:     databaseUser,
		pool:             pool,
		resource:         resource,
	}
}

// -------------------------------------------------------------------

// Configuration returns an adapter configuraion object
func (postgres *PostgreSQLContainer) Configuration() *postgresql.Configuration {
	return &postgresql.Configuration{
		ConnectionString : postgres.ConnectionString,
		Username         : postgres.DatabaseUser,
		Password         : postgres.Password,
	}
}

// Stop the container
func (postgres *PostgreSQLContainer) Stop() error {
	log.Printf("Postgres (%v): shutting down", postgres.Name)
	return postgres.pool.Purge(postgres.resource)
}