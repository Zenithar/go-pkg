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
	IntervalStr string `toml:"interval" default:"20s" comment:"Refresh interval for runtime metrics update"`
	Interval time.Duration `toml:"-"`
}

// Validate rntime config parameters
func (c *Config) Validate() error {

	interval, err := time.ParseDuration(c.IntervalStr)
	if err != nil {
		return xerrors.Errorf("invalid interval duration value for runtime metrics: %w", err)
	}
	
	// Assign parsed interval
	c.Interval = interval

	return nil
}