package runtime

import (
	"time"
)

// Config holds the runtime metric exporter informations
type Config struct {
	Name     string        `toml:"-"`
	ID       string        `toml:"-"`
	Version  string        `toml:"-"`
	Revision string        `toml:"-"`
	Interval time.Duration `toml:"interval" default:"20s" comment:"Refresh interval for runtime metrics update"`
}
