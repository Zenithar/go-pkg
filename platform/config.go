package platform

// InstrumentationConfig holds all platform instrumentation settings
type InstrumentationConfig struct {
	Network string `toml:"network" default:"tcp" comment:"Network class used for listen (tcp, tcp4, tcp6, unixsocket)"`
	Listen  string `toml:"listen" default:":5556" comment:"Listen address for instrumentation server"`
	Logs    struct {
		Level     string `toml:"level" default:"warn" comment:"Log level: debug, info, warn, error, dpanic, panic, and fatal"`
		SentryDSN string `toml:"sentryDSN" comment:"Sentry DSN"`
	} `toml:"Logs" comment:"###############################\n Logs Settings \n##############################"`
}
