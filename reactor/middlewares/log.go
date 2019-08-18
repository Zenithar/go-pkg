package middlewares

import (
	"context"

	"go.uber.org/zap"

	"go.zenithar.org/pkg/log"
	"go.zenithar.org/pkg/reactor"
	"go.zenithar.org/pkg/reactor/chain"
)

// Log is a logger middleware.
func Log(lf log.LoggerFactory) chain.Constructor {
	return func(fn reactor.Handler) reactor.Handler {
		return reactor.HandlerFunc(func(ctx context.Context, req interface{}) (interface{}, error) {
			lf.For(ctx).Debug("Handling message ...", zap.Any(req))

			res, err := fn.Handle(ctx, req)
			if err != nil {
				lf.For(ctx).Error("Unable to process message", zap.Error(err))
			}

			// Delegate to next handler
			return res, err
		})
	}
}
