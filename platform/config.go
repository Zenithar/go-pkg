package platform

import (
	"go.zenithar.org/pkg/platform/telemetry"
)

// InstrumentationConfig holds all platform instrumentation settings
type InstrumentationConfig struct {
	Network string `toml:"network" default:"tcp" comment:"Network class used for listen (tcp, tcp4, tcp6, unixsocket)"`
	Listen  string `toml:"listen" default:":5556" comment:"Listen address for instrumentation server"`
	Log     struct {
		Level     string `toml:"level" default:"warn" comment:"Log level: debug, info, warn, error, dpanic, panic, and fatal"`
		SentryDSN string `toml:"sentryDSN" comment:"Sentry DSN"`
	} `toml:"Log" comment:"###############################\n Log Settings \n##############################"`
	Telemetry telemetry.Configuration `toml:"Telemetry" comment:"###############################\n Telemetry Settings \n##############################"`
}
