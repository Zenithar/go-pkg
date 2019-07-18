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
	"go.zenithar.org/pkg/platform/runtime"
)

// -----------------------------------------------------------------------------

// Server represents platform server
type Server struct {
	Debug           bool
	Name            string
	Version         string
	Revision        string
	Instrumentation InstrumentationConfig
	Builder         func(upg *tableflip.Upgrader, group *run.Group)
}

// Serve starts the server listening process
func Serve(ctx context.Context, srv *Server) error {

	// Generate an instance identifier
	appID := uniuri.NewLen(64)

	// Prepare logger
	log.Setup(ctx, &log.Options{
		Debug:     srv.Debug,
		AppName:   srv.Name,
		AppID:     appID,
		Version:   srv.Version,
		Revision:  srv.Revision,
		SentryDSN: srv.Instrumentation.Logs.SentryDSN,
		LogLevel:  srv.Instrumentation.Logs.Level,
	})

	// Preparing instrumentation
	instrumentationRouter := http.NewServeMux()

	// Register common features
	if srv.Instrumentation.Diagnostic.Enabled {
		cancelFunc, err := diagnostic.Register(ctx, srv.Instrumentation.Diagnostic.Config, instrumentationRouter)
		if err != nil {
			log.For(ctx).Fatal("Unable to register diagnostic instrumentation", zap.Error(err))
		}
		defer cancelFunc()
	}
	if srv.Instrumentation.Prometheus.Enabled {
		if _, err := prometheus.RegisterExporter(ctx, srv.Instrumentation.Prometheus.Config, instrumentationRouter); err != nil {
			log.For(ctx).Fatal("Unable to register prometheus instrumentation", zap.Error(err))
		}
	}
	if srv.Instrumentation.Jaeger.Enabled {
		// Apply default service name in empty
		if srv.Instrumentation.Jaeger.Config.ServiceName == "" {
			log.For(ctx).Debug("No Jaeger service name given, applying server name as service name", zap.String("name", srv.Name))
			srv.Instrumentation.Jaeger.Config.ServiceName = srv.Name
		}

		// Register exporter
		cancelFunc, err := jaeger.RegisterExporter(ctx, srv.Instrumentation.Jaeger.Config)
		if err != nil {
			log.For(ctx).Fatal("Unable to register jaeger instrumentation", zap.Error(err))
		}
		defer cancelFunc()
	}
	if srv.Instrumentation.OCAgent.Enabled {
		cancelFunc, err := ocagent.RegisterExporter(ctx, srv.Instrumentation.OCAgent.Config)
		if err != nil {
			log.For(ctx).Fatal("Unable to register ocagent instrumentation", zap.Error(err))
		}
		defer cancelFunc()
	}
	if srv.Instrumentation.Runtime.Enabled {
		if err := runtime.Monitor(ctx, runtime.Config{
			Name:     srv.Name,
			ID:       appID,
			Version:  srv.Version,
			Revision: srv.Revision,
			Interval: srv.Instrumentation.Runtime.Config.Interval,
		}); err != nil {
			log.For(ctx).Fatal("Unable to start runtime monitoring", zap.Error(err))
		}
	}

	// Trace everything when debugging is enabled
	if srv.Debug {
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
		ln, err := upg.Fds.Listen(srv.Instrumentation.Network, srv.Instrumentation.Listen)
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
	srv.Builder(upg, &group)

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
