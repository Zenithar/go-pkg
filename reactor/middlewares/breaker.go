package middlewares

import (
	"context"

	"github.com/sony/gobreaker"

	"go.zenithar.org/pkg/reactor"
	"go.zenithar.org/pkg/reactor/chain"
)

// Tmeout implements timeout pattern for invocation.
func Breaker(cb *gobreaker.CircuitBreaker) chain.Constructor {
	return func(fn reactor.Handler) reactor.Handler {
		return reactor.HandlerFunc(func(ctx context.Context, req interface{}) (interface{}, error) {
			return cb.Execute(func() (interface{}, error) { return fn.Handle(ctx, req) })
		})
	}
}
