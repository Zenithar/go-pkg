// MIT License
//
// Copyright (c) 2019 Thibault NORMAND
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package log

import (
	"context"
	"io"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultFactory LoggerFactory

// -----------------------------------------------------------------------------

func init() {
	SetLoggerFactory(NewFactory(zap.L()))
}

// SetLoggerFactory defines the default package logger
func SetLoggerFactory(instance LoggerFactory) {
	if defaultFactory != nil {
		defaultFactory.Bg().Debug("Replacing logger factory", zap.String("old", defaultFactory.Name()), zap.String("new", instance.Name()))
	} else {
		instance.Bg().Debug("Initializing logger factory", zap.String("factory", instance.Name()))
	}
	defaultFactory = instance
}

// -----------------------------------------------------------------------------

// Bg delegates a no-context logger
func Bg() Logger {
	return checkFactory(defaultFactory).Bg()
}

// For delegates a context logger
func For(ctx context.Context) Logger {
	return checkFactory(defaultFactory).For(ctx)
}

// Default returns the logger factory
func Default() LoggerFactory {
	return checkFactory(defaultFactory)
}

// CheckErr handles error correctly
func CheckErr(msg string, err error, fields ...zapcore.Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
		Default().Bg().Error(msg, fields...)
	}
}

// CheckErrCtx handles error correctly
func CheckErrCtx(ctx context.Context, msg string, err error, fields ...zapcore.Field) {
	if err != nil {
		fields = append(fields, zap.Error(errors.WithStack(err)))
		Default().For(ctx).Error(msg, fields...)
	}
}

// SafeClose handles the closer error
func SafeClose(c io.Closer, msg string, fields ...zapcore.Field) {
	if cerr := c.Close(); cerr != nil {
		fields = append(fields, zap.Error(errors.WithStack(cerr)))
		Default().Bg().Error(msg, fields...)
	}
}

// SafeCloseCtx handles the closer error
func SafeCloseCtx(ctx context.Context, c io.Closer, msg string, fields ...zapcore.Field) {
	if cerr := c.Close(); cerr != nil {
		fields = append(fields, zap.Error(errors.WithStack(cerr)))
		Default().For(ctx).Error(msg, fields...)
	}
}

// -----------------------------------------------------------------------------

func checkFactory(defaultFactory LoggerFactory) LoggerFactory {
	if defaultFactory == nil {
		panic("Unable to create logger instance, you have to register an adapter first.")
	}
	return defaultFactory
}
