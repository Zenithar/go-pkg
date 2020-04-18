package log

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Factory is the default logging wrapper that can create
// logger instances either for a given Context or context-less.
type factory struct {
	logger *zap.Logger
}

// NewFactory creates a new Factory.
func NewFactory(logger *zap.Logger) LoggerFactory {
	return &factory{logger: logger}
}

// -----------------------------------------------------------------------------

// Name returns the logger adapter name
func (b factory) Name() string {
	return "zap"
}

// Bg creates a context-unaware logger.
func (b factory) Bg() Logger {
	return &logger{logger: b.logger}
}

// For returns a context-aware Logger.
func (b factory) For(ctx context.Context) Logger {
	return b.Bg()
}

// With creates a child logger, and optionally adds some context fields to that logger.
func (b factory) With(fields ...zapcore.Field) LoggerFactory {
	return &factory{logger: b.logger.With(fields...)}
}
