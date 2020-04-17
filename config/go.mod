module go.zenithar.org/pkg/config

go 1.14

replace go.zenithar.org/pkg/log => ../log

replace go.zenithar.org/pkg/flags => ../flags

require (
	github.com/mcuadros/go-defaults v1.2.0
	github.com/pelletier/go-toml v1.7.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.6.3
	go.uber.org/zap v1.14.1
	go.zenithar.org/pkg/flags v0.0.0-00010101000000-000000000000
	go.zenithar.org/pkg/log v0.1.1
)
