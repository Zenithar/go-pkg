module go.zenithar.org/pkg/db/adapter/mongodb

go 1.12

replace github.com/opencensus-integrations/gomongowrapper => github.com/Zenithar/gomongowrapper v0.0.2

require (
	github.com/opencensus-integrations/gomongowrapper v0.0.0-00010101000000-000000000000
	go.mongodb.org/mongo-driver v1.0.1-0.20190712184055-9ec4480161a7
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.10.0
	go.zenithar.org/pkg/db v0.0.3
	go.zenithar.org/pkg/log v0.1.1
	golang.org/x/xerrors v0.0.0-20190513163551-3ee3066db522
)
