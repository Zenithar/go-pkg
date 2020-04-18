package telemetry

// Configuration contains telemetry settings
type Configuration struct {
	Enable bool `toml:"enable" default:"false" comment:"Enable / Disable OpenTelemetry agent usages"`
}
