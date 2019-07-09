module go.zenithar.org/pkg/db/adapter/mongodb

go 1.12

replace github.com/opencensus-integrations/gomongowrapper => github.com/eug48/gomongowrapper v0.0.2

require (
	github.com/opencensus-integrations/gomongowrapper v0.0.1
	github.com/stretchr/testify v1.3.0 // indirect
	github.com/tidwall/pretty v1.0.0 // indirect
	go.mongodb.org/mongo-driver v1.0.4
	go.uber.org/zap v1.10.0
	go.zenithar.org/pkg/db v0.0.3
	go.zenithar.org/pkg/log v0.0.3
	golang.org/x/xerrors v0.0.0-20190513163551-3ee3066db522
)
