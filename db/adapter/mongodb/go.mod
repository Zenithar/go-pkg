module go.zenithar.org/pkg/db/adapter/mongodb

go 1.12

replace github.com/opencensus-integrations/gomongowrapper v0.0.1 => github.com/Zenithar/gomongowrapper v0.0.2

require (
	github.com/opencensus-integrations/gomongowrapper v0.0.1
	go.mongodb.org/mongo-driver v1.0.4
	go.zenithar.org/pkg/db v0.0.3
	go.zenithar.org/pkg/log v0.1.3
	golang.org/x/xerrors v0.0.0-20190717185122-a985d3407aa7
)
