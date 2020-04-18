module go.zenithar.org/pkg/reactor

go 1.14

replace go.zenithar.org/pkg/errors => ../errors

replace go.zenithar.org/pkg/log => ../log

replace go.zenithar.org/pkg/types => ../types

require (
	github.com/google/go-cmp v0.4.0
	github.com/sony/gobreaker v0.4.1
	go.uber.org/zap v1.14.1
	go.zenithar.org/pkg/errors v0.0.0-00010101000000-000000000000
	go.zenithar.org/pkg/log v0.0.0-00010101000000-000000000000
	go.zenithar.org/pkg/types v0.0.0-00010101000000-000000000000
)
