module go.zenithar.org/pkg/db/adapter/rethinkdb

go 1.14

replace go.zenithar.org/pkg/db => ../../../db

replace go.zenithar.org/pkg/log => ../../../log

require (
	go.zenithar.org/pkg/db v0.0.0-00010101000000-000000000000
	gopkg.in/rethinkdb/rethinkdb-go.v6 v6.2.1
)
