// +build !windows

package reloader

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.zenithar.org/pkg/log"

	"github.com/cloudflare/tableflip"
	"github.com/oklog/run"
)

// TableflipReloader deleagtes socket reloading to tableflip library which his
// not windows compatible.
type TableflipReloader struct {
	*tableflip.Upgrader
}

// Create a descriptor reload based on tableflip.
func Create(ctx context.Context) Reloader {
	upg, _ := tableflip.New(tableflip.Options{})

	// Do an upgrade on SIGHUP
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGHUP)
		for range ch {
			log.For(ctx).Warn("Graceful reloading socket descriptor")
			_ = upg.Upgrade()
		}
	}()

	return &TableflipReloader{upg}
}

// SetupGracefulRestart arms the graceful restart handler.
func (t *TableflipReloader) SetupGracefulRestart(ctx context.Context, group run.Group) {
	ctx, cancel := context.WithCancel(ctx)

	// Register an actor, i.e. an execute and interrupt func, that
	// terminates when graceful restart is initiated and the child process
	// signals to be ready, or the parent context is canceled.
	group.Add(func() error {
		// Tell the parent we are ready
		err := t.Ready()
		if err != nil {
			return err
		}

		select {
		case <-t.Exit(): // Wait for child to be ready (or application shutdown)
			return nil

		case <-ctx.Done():
			return ctx.Err()
		}
	}, func(error) {
		cancel()
		t.Stop()
	})
}
