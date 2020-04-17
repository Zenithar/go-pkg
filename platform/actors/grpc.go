package actors

import (
	"context"
	"net"

	"go.zenithar.org/pkg/log"

	"github.com/oklog/run"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// GRPC registers an gRPC listener actor.
func GRPC(server *grpc.Server, ln net.Listener) func(context.Context, *run.Group) error {
	return func(ctx context.Context, group *run.Group) error {
		// Register grpc actor
		group.Add(
			func() error {
				log.For(ctx).Info("Starting gRPC server", zap.String("address", ln.Addr().String()))
				return server.Serve(ln)
			},
			func(e error) {
				log.For(ctx).Info("Shutting gRPC server down")
				server.GracefulStop()
			},
		)
		return nil
	}
}
