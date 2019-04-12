module go.zenithar.org/pkg

go 1.12

require (
	github.com/Masterminds/squirrel v1.1.0
	github.com/TheZeroSlave/zapsentry v0.0.0-20180112122240-410ad1e37c78
	github.com/bitly/go-hostpool v0.0.0-20171023180738-a3a6125de932 // indirect
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/certifi/gocertifi v0.0.0-20190410005359-59a85de7f35e // indirect
	github.com/cheekybits/is v0.0.0-20150225183255-68e9c0620927 // indirect
	github.com/fatih/structs v1.1.0
	github.com/getsentry/raven-go v0.2.0 // indirect
	github.com/gorilla/schema v1.1.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.0.0
	github.com/matryer/try v0.0.0-20161228173917-9ac251b645a2 // indirect
	github.com/mongodb/mongo-go-driver v1.0.0
	github.com/opencensus-integrations/gomongowrapper v0.0.1
	github.com/opencensus-integrations/ocsql v0.1.4
	github.com/pkg/errors v0.8.1
	github.com/smartystreets/goconvey v0.0.0-20190330032615-68dc04aab96a
	github.com/spf13/pflag v1.0.3
	github.com/stretchr/testify v1.3.0
	go.mongodb.org/mongo-driver v1.0.0
	go.opencensus.io v0.20.2
	go.uber.org/atomic v1.3.2 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.9.1
	gopkg.in/matryer/try.v1 v1.0.0-20150601225556-312d2599e12e
	gopkg.in/rethinkdb/rethinkdb-go.v5 v5.0.1
)

replace github.com/opencensus-integrations/gomongowrapper => github.com/Zenithar/gomongowrapper v0.0.2
