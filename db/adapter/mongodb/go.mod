module go.zenithar.org/pkg/db/adapter/mongodb

go 1.14

replace go.zenithar.org/pkg/db => ../../../db

replace go.zenithar.org/pkg/log => ../../../log

require (
	github.com/google/go-cmp v0.3.0 // indirect
	github.com/xdg/stringprep v1.0.0 // indirect
	go.mongodb.org/mongo-driver v1.3.2
	go.uber.org/zap v1.14.1
	go.zenithar.org/pkg/db v0.0.0-00010101000000-000000000000
	go.zenithar.org/pkg/log v0.0.0-00010101000000-000000000000
)
