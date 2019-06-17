package runtime

import (
	"context"
	"time"
)

// Monitor will start the runtime monitoring process to publish metrics
func Monitor(ctx context.Context, cfg Config) error {

	// Validate config
	if err := cfg.Validate(); err != nil {
		return err
	}

	// Initialize runtime stats
	rstats := &Stats{}

	// Collect first
	rstats.Collect()

	// Initialize publisher
	publisher := OpenCensus(cfg)
	publisher.Publish(rstats)

	// Fork monitoring routine
	go func() {
		tick := time.NewTicker(cfg.Interval)
		defer tick.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				// Update metrics
				rstats.Collect()
				// Publish to metric protocol
				publisher.Publish(rstats)
			}
		}
	}()

	return nil
}
