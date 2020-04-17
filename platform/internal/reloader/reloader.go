package reloader

import (
	"context"
	"net"

	"github.com/oklog/run"
)

// Reloader defines socket reloader contract.
type Reloader interface {
	Listen(network, address string) (net.Listener, error)
	SetupGracefulRestart(context.Context, run.Group)
}
