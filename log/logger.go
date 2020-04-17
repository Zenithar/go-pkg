package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// logger delegates all calls to the underlying zap.Logger
type logger struct {
	logger *zap.Logger
}

// Debug logs an debug msg with fields
func (l logger) Debug(msg string, fields ...zapcore.Field) {
	l.logger.Debug(msg, fields...)
}

// Info logs an info msg with fields
func (l logger) Info(msg string, fields ...zapcore.Field) {
	l.logger.Info(msg, fields...)
}

// Error logs an error msg with fields
func (l logger) Error(msg string, fields ...zapcore.Field) {
	l.logger.Error(msg, fields...)
}

// Warn logs a warning with fields
func (l logger) Warn(msg string, fields ...zapcore.Field) {
	l.logger.Warn(msg, fields...)
}

// Fatal logs a fatal error msg with fields
func (l logger) Fatal(msg string, fields ...zapcore.Field) {
	l.logger.Fatal(msg, fields...)
}

// With creates a child logger, and optionally adds some context fields to that logger.
func (l logger) With(fields ...zapcore.Field) Logger {
	return &logger{logger: l.logger.With(fields...)}
}
