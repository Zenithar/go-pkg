package database

import (
	"fmt"
	"log"

	"github.com/dchest/uniuri"
	dockertest "gopkg.in/ory-am/dockertest.v3"

	"go.zenithar.org/pkg/testing/containers"
)

var (
	// MongoDBVersion defines version to use
	MongoDBVersion = "latest"
)

// PostgreSQLContainer represents database container handler
type mongoDBContainer struct {
	Name             string
	ConnectionString string
	Password         string
	DatabaseName     string
	DatabaseUser     string
	pool             *dockertest.Pool
	resource         *dockertest.Resource
}

// NewPostgresContainer initialize a PostgreSQL server in a docker container
func newMongoDBContainer(pool *dockertest.Pool) *mongoDBContainer {

	var (
		databaseName = fmt.Sprintf("test-%s", uniuri.NewLen(8))
		databaseUser = fmt.Sprintf("user-%s", uniuri.NewLen(8))
		password     = uniuri.NewLen(32)
	)

	// Initialize a PostgreSQL server
	resource, err := pool.Run("mongo", MongoDBVersion, []string{
		fmt.Sprintf("MONGO_INITDB_ROOT_USERNAME=%s", databaseUser),
		fmt.Sprintf("MONGO_INITDB_ROOT_PASSWORD=%s", password),
		fmt.Sprintf("MONGO_INITDB_DATABASE=%s", databaseName),
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// Prepare connection string
	connectionString := fmt.Sprintf("localhost:%s", resource.GetPort("27017/tcp"))

	// Retrieve container name
	containerName := containers.GetName(resource)

	// Return container information
	return &mongoDBContainer{
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

// Close the container
func (container *mongoDBContainer) Close() error {
	log.Printf("Postgres (%v): shutting down", container.Name)
	return container.pool.Purge(container.resource)
}
