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

package zap

import (
	"context"

	"go.zenithar.org/pkg/log"

	"go.uber.org/zap"
)

// Factory is the default logging wrapper that can create
// logger instances either for a given Context or context-less.
type factory struct {
	logger *zap.Logger
}

// DefaultFactory returns a default zap logger factory
func defaultFactory() log.LoggerFactory {
	defaultLogger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	return newFactory(defaultLogger)
}

// NewFactory creates a new Factory.
func newFactory(logger *zap.Logger) log.LoggerFactory {
	return &factory{logger: logger}
}

// -----------------------------------------------------------------------------

// Name returns the logger adapter name
func (b factory) Name() string {
	return "zap"
}

// Bg creates a context-unaware logger.
func (b factory) Bg() log.Logger {
	return &logger{logger: b.logger}
}

// For returns a context-aware Logger.
// TODO: OpenCensus implementation
func (b factory) For(ctx context.Context) log.Logger {
	return b.Bg()
}

// With creates a child logger, and optionally adds some context fields to that logger.
func (b factory) With(fields ...log.Field) log.LoggerFactory {
	return &factory{logger: b.logger.With(zfields(fields)...)}
}

// -----------------------------------------------------------------------------

func init() {
	log.SetLoggerFactory(defaultFactory())
}
