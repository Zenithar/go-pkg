module go.zenithar.org/pkg/db/adapter/mongodb

go 1.12

replace github.com/opencensus-integrations/gomongowrapper => github.com/Zenithar/gomongowrapper v0.0.2

require (
	github.com/opencensus-integrations/gomongowrapper v0.0.0-00010101000000-000000000000
	github.com/pkg/errors v0.8.1
	go.mongodb.org/mongo-driver v1.0.1
	go.uber.org/zap v1.9.1
	go.zenithar.org/pkg/db v0.0.1
	go.zenithar.org/pkg/log v0.0.1
)
