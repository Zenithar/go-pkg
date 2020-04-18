package actors

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.zenithar.org/pkg/log"

	"github.com/oklog/run"
)

// Signal register an Signal handle actor
func Signal(ctx context.Context, group run.Group) {
	var (
		cancelInterrupt = make(chan struct{})
		ch              = make(chan os.Signal, 2)
	)

	// Close signal channel on exit signal
	group.Add(
		func() error {
			<-ctx.Done()
			close(ch)
			return nil
		},
		nil,
	)

	// Add to run group
	group.Add(
		func() error {
			signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

			select {
			case <-ch:
				log.For(ctx).Info("Captured signal")
			case <-cancelInterrupt:
			}

			return nil
		},
		func(e error) {
			close(cancelInterrupt)
			signal.Stop(ch)
		},
	)
}
