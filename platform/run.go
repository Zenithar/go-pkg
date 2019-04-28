package platform

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudflare/tableflip"
	"github.com/oklog/run"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"go.zenithar.org/pkg/log"
	"go.zenithar.org/pkg/platform/diagnostic"
	"go.zenithar.org/pkg/platform/jaeger"
	"go.zenithar.org/pkg/platform/prometheus"
)

// -----------------------------------------------------------------------------

// Run the dispatcher
func Run(ctx context.Context, debug bool, conf InstrumentationConfig, builder func(upg *tableflip.Upgrader, group run.Group)) error {
	// Preparing instrumentation
	instrumentationRouter := http.NewServeMux()

	// Register common features
	if conf.Diagnostic.Enabled {
		if err := diagnostic.Register(ctx, conf.Diagnostic.Config, instrumentationRouter); err != nil {
			log.For(ctx).Fatal("Unable to register diagnostic instrumentation", zap.Error(err))
		}
	}
	if conf.Prometheus.Enabled {
		if err := prometheus.RegisterExporter(ctx, conf.Prometheus.Config, instrumentationRouter); err != nil {
			log.For(ctx).Fatal("Unable to register prometheus instrumentation", zap.Error(err))
		}
	}
	if conf.Jaeger.Enabled {
		if err := jaeger.RegisterExporter(ctx, debug, conf.Jaeger.Config); err != nil {
			log.For(ctx).Fatal("Unable to register jaeger instrumentation", zap.Error(err))
		}
	}

	// Configure graceful restart
	upg, err := tableflip.New(tableflip.Options{})
	if err != nil {
		return errors.Wrap(err, "Unable to register graceful restart handler")
	}

	// Do an upgrade on SIGHUP
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGHUP)
		for range ch {
			log.For(ctx).Info("Graceful reloading")
			_ = upg.Upgrade()
		}
	}()

	var group run.Group

	// Instrumentation server
	{
		ln, err := upg.Fds.Listen(conf.Network, conf.Listen)
		if err != nil {
			return errors.Wrap(err, "Unable to start instrumentation server")
		}

		server := &http.Server{
			Handler: instrumentationRouter,
		}

		group.Add(
			func() error {
				log.For(ctx).Info("Starting instrumentation server", zap.Stringer("address", ln.Addr()))
				return server.Serve(ln)
			},
			func(e error) {
				log.For(ctx).Info("Shutting instrumentation server down")

				ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
				defer cancel()

				log.CheckErrCtx(ctx, "Error raised while shutting down the server", server.Shutdown(ctx))
				log.SafeClose(server, "Unable to close instrumentation server")
			},
		)
	}

	// Initialize the component
	builder(upg, group)

	// Setup signal handler
	{
		var (
			cancelInterrupt = make(chan struct{})
			ch              = make(chan os.Signal, 2)
		)
		defer close(ch)

		group.Add(
			func() error {
				signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

				select {
				case sig := <-ch:
					log.For(ctx).Info("Captured signal", zap.Any("signal", sig))
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

	// Final handler
	{
		group.Add(
			func() error {
				// Tell the parent we are ready
				_ = upg.Ready()

				// Wait for children to be ready
				// (or application shutdown)
				<-upg.Exit()

				return nil
			},
			func(e error) {
				upg.Stop()
			},
		)
	}

	// Run goroutine group
	return group.Run()
}
