package platform

import (
	"go.zenithar.org/pkg/platform/diagnostic"
	"go.zenithar.org/pkg/platform/jaeger"
	"go.zenithar.org/pkg/platform/ocagent"
	"go.zenithar.org/pkg/platform/prometheus"
)

// InstrumentationConfig holds all platform instrumentation settings
type InstrumentationConfig struct {
	Network    string `toml:"network" default:"tcp" comment:"Network class used for listen (tcp, tcp4, tcp6, unixsocket)"`
	Listen     string `toml:"listen" default:":5556" comment:"Listen address for instrumentation server"`
	Diagnostic struct {
		Enabled bool              `toml:"enabled" default:"true" comment:"Enable diagnostic handlers"`
		Config  diagnostic.Config `toml:"Config" comment:"Diagnostic settings"`
	} `toml:"Diagnostic" comment:"###############################\n Diagnotic Settings \n##############################"`
	Logs struct {
		Level     string `toml:"level" default:"warn" comment:"Log level: debug, info, warn, error, dpanic, panic, and fatal"`
		SentryDSN string `toml:"sentryDSN" comment:"Sentry DSN"`
	} `toml:"Logs" comment:"###############################\n Logs Settings \n##############################"`
	Prometheus struct {
		Enabled bool              `toml:"enabled" default:"true" comment:"Enable metric instrumentation"`
		Config  prometheus.Config `toml:"Config" comment:"Prometheus settings"`
	} `toml:"Prometheus" comment:"###############################\n Prometheus exporter \n##############################"`
	Jaeger struct {
		Enabled bool          `toml:"enabled" default:"true" comment:"Enable trace instrumentation"`
		Config  jaeger.Config `toml:"Config" comment:"Jaeger settings"`
	} `toml:"Jaeger" comment:"###############################\n Jaeger exporter \n##############################"`
	OCAgent struct {
		Enabled bool           `toml:"enabled" default:"false" comment:"Enable trace instrumentation"`
		Config  ocagent.Config `toml:"Config" comment:"OpenCensus agent settings"`
	} `toml:"OCAgent" comment:"###############################\n OpenCensus Agent exporter \n##############################"`
}
