package database

import (
	"fmt"
	"log"
	"time"

	"github.com/dchest/uniuri"
	_ "github.com/lib/pq"
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

	// Hard killing resource timeout
	resource.Expire(15 * time.Minute)

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
func (postgres *postgreSQLContainer) Close() error {
	log.Printf("Postgres (%v): shutting down", postgres.Name)
	return postgres.pool.Purge(postgres.resource)
}

// Configuration return database settings
func (postgres *postgreSQLContainer) Configuration() *Configuration {
	return postgres.config
}
