package platform

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudflare/tableflip"
	"github.com/dchest/uniuri"
	"github.com/oklog/run"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"golang.org/x/xerrors"

	"go.zenithar.org/pkg/log"
	"go.zenithar.org/pkg/platform/diagnostic"
	"go.zenithar.org/pkg/platform/jaeger"
	"go.zenithar.org/pkg/platform/ocagent"
	"go.zenithar.org/pkg/platform/prometheus"
)

// -----------------------------------------------------------------------------

// Application represents platform application
type Application struct {
	Debug           bool
	Name            string
	Version         string
	Revision        string
	Instrumentation InstrumentationConfig
	Builder         func(upg *tableflip.Upgrader, group *run.Group)
}

// Run the dispatcher
func Run(ctx context.Context, app *Application) error {

	// Generate an instance identifier
	appID := uniuri.NewLen(64)

	// Prepare logger
	log.Setup(ctx, &log.Options{
		Debug:     app.Debug,
		AppName:   app.Name,
		AppID:     appID,
		Version:   app.Version,
		Revision:  app.Revision,
		SentryDSN: app.Instrumentation.Logs.SentryDSN,
		LogLevel:  app.Instrumentation.Logs.Level,
	})

	// Preparing instrumentation
	instrumentationRouter := http.NewServeMux()

	// Register common features
	if app.Instrumentation.Diagnostic.Enabled {
		err := diagnostic.Register(ctx, app.Instrumentation.Diagnostic.Config, instrumentationRouter)
		if err != nil {
			log.For(ctx).Fatal("Unable to register diagnostic instrumentation", zap.Error(err))
		}
	}
	if app.Instrumentation.Prometheus.Enabled {
		if _, err := prometheus.RegisterExporter(ctx, app.Instrumentation.Prometheus.Config, instrumentationRouter); err != nil {
			log.For(ctx).Fatal("Unable to register prometheus instrumentation", zap.Error(err))
		}
	}
	if app.Instrumentation.Jaeger.Enabled {
		if _, err := jaeger.RegisterExporter(ctx, app.Instrumentation.Jaeger.Config); err != nil {
			log.For(ctx).Fatal("Unable to register jaeger instrumentation", zap.Error(err))
		}
	}
	if app.Instrumentation.OCAgent.Enabled {
		if _, err := ocagent.RegisterExporter(ctx, app.Instrumentation.OCAgent.Config); err != nil {
			log.For(ctx).Fatal("Unable to register ocagent instrumentation", zap.Error(err))
		}
	}

	// Trace everything when debugging is enabled
	if app.Debug {
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	}

	// Configure graceful restart
	upg, err := tableflip.New(tableflip.Options{})
	if err != nil {
		return xerrors.Errorf("platform: unable to register graceful restart handler: %w", err)
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
		ln, err := upg.Fds.Listen(app.Instrumentation.Network, app.Instrumentation.Listen)
		if err != nil {
			return xerrors.Errorf("platform: unable to start instrumentation server: %w", err)
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
	app.Builder(upg, &group)

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
