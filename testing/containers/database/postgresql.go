package database

import (
	"fmt"
	"log"

	// Load driver if not already done
	_ "github.com/lib/pq"

	"github.com/dchest/uniuri"
	dockertest "gopkg.in/ory-am/dockertest.v3"

	"go.zenithar.org/pkg/testing/containers"
)

var (
	// PostgreSQLVersion defines version to use
	PostgreSQLVersion = "10"
)

// PostgreSQLContainer represents database container handler
type postgreSQLContainer struct {
	Name     string
	pool     *dockertest.Pool
	resource *dockertest.Resource
	config   *Configuration
}

// NewPostgresContainer initialize a PostgreSQL server in a docker container
func newPostgresContainer(pool *dockertest.Pool) *postgreSQLContainer {

	var (
		databaseName = fmt.Sprintf("test-%s", uniuri.NewLen(8))
		databaseUser = fmt.Sprintf("user-%s", uniuri.NewLen(8))
		password     = uniuri.NewLen(32)
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

	// Return container information
	return &postgreSQLContainer{
		Name:     containerName,
		pool:     pool,
		resource: resource,
		config: &Configuration{
			ConnectionString: connectionString,
			Password:         password,
			DatabaseName:     databaseName,
			DatabaseUser:     databaseUser,
		},
	}
}

// -------------------------------------------------------------------

// Close the container
func (container *postgreSQLContainer) Close() error {
	log.Printf("Postgres (%v): shutting down", container.Name)
	return container.pool.Purge(container.resource)
}

// Configuration return database settings
func (container *postgreSQLContainer) Configuration() *Configuration {
	return container.config
}
