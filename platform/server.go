package platform

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"go.zenithar.org/pkg/log"
	"go.zenithar.org/pkg/platform/actors"
	"go.zenithar.org/pkg/platform/internal/reloader"

	"github.com/dchest/uniuri"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oklog/run"
	"go.opencensus.io/trace"
)

// -----------------------------------------------------------------------------

// Server represents platform server
type Server struct {
	Debug           bool
	Name            string
	Version         string
	Revision        string
	Network         string
	Address         string
	Instrumentation InstrumentationConfig
	Builder         func(ctx context.Context, ln net.Listener, group *run.Group)
}

// Validate server settings
func (s Server) Validate() error {
	return validation.ValidateStruct(s,
		validation.Field(&s.Name, validation.Required),
		validation.Field(&s.Version, validation.Required),
		validation.Field(&s.Revision, validation.Required),
		validation.Field(&s.Network, validation.Required),
		validation.Field(&s.Address, validation.Required),
	)
}

// Serve starts the server listening process
func Serve(ctx context.Context, srv Server) error {

	// Validate server settings first
	if err := srv.Validate(); err != nil {
		return fmt.Errorf("unable to validat server settings: %w", err)
	}

	// Generate an instance identifier
	appID := uniuri.NewLen(64)

	// Prepare logger
	log.Setup(ctx, log.Options{
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

	// Trace everything when debugging is enabled
	if srv.Debug {
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	}

	// Configure graceful restart
	upg := reloader.Create(ctx)

	var group run.Group

	// Instrumentation server
	{
		ln, err := upg.Listen(srv.Instrumentation.Network, srv.Instrumentation.Listen)
		if err != nil {
			return fmt.Errorf("platform: unable to start instrumentation server: %w", err)
		}

		server := &http.Server{
			Handler: instrumentationRouter,
		}

		// Register HTTP actor
		actors.HTTP(server, ln)(ctx, &group)
	}

	// Initialiaze network listener
	ln, err := upg.Listen(srv.Network, srv.Address)
	if err != nil {
		return fmt.Errorf("unable to start server listener: %w", err)
	}

	// Initialize the component
	srv.Builder(ctx, ln, &group)

	// Setup signal handler
	actors.Signal(ctx, &group)

	// Register graceful restart handler
	upg.SetupGracefulRestart(ctx, group)

	// Run goroutine group
	return group.Run()
}
