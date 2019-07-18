package runtime

import (
	"time"

	"golang.org/x/xerrors"
)

// Config holds the runtime metric exporter informations
type Config struct {
	Name     string        `toml:"-"`
	ID       string        `toml:"-"`
	Version  string        `toml:"-"`
	Revision string        `toml:"-"`
	Interval time.Duration `toml:"interval" default:"1m" comment:"Refresh interval for runtime metrics update"`
}

// Validate rntime config parameters
func (c *Config) Validate() error {

	if c.Interval < 1*time.Second {
		return xerrors.Errorf("invalid interval duration value for runtime metrics, refresh interval must be more than 1 second")
	}

	return nil
}
