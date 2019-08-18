package reactor

import (
	"context"
)

// Handler describes a command handler
type Handler interface {
	Handle(ctx context.Context, req interface{}) (interface{}, error)
}

// -----------------------------------------------------------------------------

// HandlerFunc describes a function implementation.
type HandlerFunc func(context.Context, interface{}) (interface{}, error)

// Handle call the wrapped function
func (f HandlerFunc) Handle(ctx context.Context, req interface{}) (interface{}, error) {
	return f(ctx, req)
}

// -----------------------------------------------------------------------------

// Callback function for asynchronous event handling.
type Callback func(context.Context, interface{}, error)

// Reactor defines reactor contract.
type Reactor interface {
	// Send the reques to the reactor as an asynchronous call.
	Send(ctx context.Context, req interface{}, cb Callback) error
	// Do the request as a synchronous call.
	Do(ctx context.Context, req interface{}) (interface{}, error)
	// Register a message type handler
	RegisterHandler(msg interface{}, fn Handler)
}
