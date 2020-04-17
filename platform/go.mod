module go.zenithar.org/pkg/platform

go 1.14

replace go.zenithar.org/pkg/log => ../log

require (
	github.com/cloudflare/tableflip v1.0.0
	github.com/dchest/uniuri v0.0.0-20200228104902-7aecb25e1fe5
	github.com/go-ozzo/ozzo-validation/v4 v4.1.0
	github.com/oklog/run v1.1.0
	go.opencensus.io v0.22.3
	go.uber.org/zap v1.14.1
	go.zenithar.org/pkg/log v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.28.1
)
