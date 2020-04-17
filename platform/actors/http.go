package actors

import (
	"context"
	"net"
	"net/http"
	"time"

	"go.zenithar.org/pkg/log"

	"github.com/oklog/run"
	"go.uber.org/zap"
)

// HTTP registers an HTTP listener actor.
func HTTP(server *http.Server, ln net.Listener) func(context.Context, *run.Group) error {
	return func(ctx context.Context, group *run.Group) error {
		// Register http actor
		group.Add(
			func() error {
				log.For(ctx).Info("Starting HTTP server", zap.String("address", ln.Addr().String()))
				return server.Serve(ln)
			},
			func(e error) {
				log.For(ctx).Info("Shutting HTTP server down")

				ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
				defer cancel()

				log.CheckErrCtx(ctx, "Error raised while shutting down the server", server.Shutdown(ctx))
				log.SafeClose(server, "Unable to close HTTP server")
			},
		)
		return nil
	}
}
