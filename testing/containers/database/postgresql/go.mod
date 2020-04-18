module go.zenithar.org/pkg/testing/database/postgresql

go 1.14

replace go.zenithar.org/pkg/testing => ../../../../testing

require (
	github.com/dchest/uniuri v0.0.0-20200228104902-7aecb25e1fe5
	github.com/jackc/pgx/v4 v4.6.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.3.0
	github.com/ory/dockertest/v3 v3.6.0
	go.zenithar.org/pkg/testing v0.0.0-00010101000000-000000000000
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543
	google.golang.org/appengine v1.6.5 // indirect
)
