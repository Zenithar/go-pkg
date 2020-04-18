package log

import (
	"context"

	"go.uber.org/zap/zapcore"
)

// Logger is a simplified abstraction of the zap.Logger
type Logger interface {
	Debug(msg string, fields ...zapcore.Field)
	Info(msg string, fields ...zapcore.Field)
	Warn(msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
	Fatal(msg string, fields ...zapcore.Field)
	With(fields ...zapcore.Field) Logger
}

// LoggerFactory defines logger factory contract
type LoggerFactory interface {
	Name() string
	Bg() Logger
	For(context.Context) Logger
	With(...zapcore.Field) LoggerFactory
}
