module go.zenithar.org/pkg/db/adapter/postgresql

go 1.14

replace go.zenithar.org/pkg/db => ../../../db

replace go.zenithar.org/pkg/log => ../../../log

require (
	github.com/Masterminds/squirrel v1.2.0
	github.com/cheekybits/is v0.0.0-20150225183255-68e9c0620927 // indirect
	github.com/jackc/pgx/v4 v4.6.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.3.0
	github.com/matryer/try v0.0.0-20161228173917-9ac251b645a2 // indirect
	github.com/opencensus-integrations/ocsql v0.1.5
	go.opencensus.io v0.22.3 // indirect
	go.zenithar.org/pkg/db v0.0.0-00010101000000-000000000000
	go.zenithar.org/pkg/log v0.0.0-00010101000000-000000000000
	google.golang.org/appengine v1.6.5 // indirect
	gopkg.in/matryer/try.v1 v1.0.0-20150601225556-312d2599e12e
)
