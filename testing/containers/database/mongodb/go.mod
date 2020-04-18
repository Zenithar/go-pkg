module go.zenithar.org/pkg/testing/database/mongodb

go 1.14

replace go.zenithar.org/pkg/testing => ../../../../testing

require (
	github.com/dchest/uniuri v0.0.0-20200228104902-7aecb25e1fe5
	github.com/ory/dockertest/v3 v3.6.0
	go.mongodb.org/mongo-driver v1.3.2
	go.zenithar.org/pkg/testing v0.0.0-00010101000000-000000000000
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543
)
