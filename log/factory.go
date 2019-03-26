// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggerFactory defines logger factory contract
type LoggerFactory interface {
	Bg() Logger
	For(context.Context) Logger
	With(...zapcore.Field) LoggerFactory
}

// Factory is the default logging wrapper that can create
// logger instances either for a given Context or context-less.
type factory struct {
	logger *zap.Logger
}

// NewFactory creates a new Factory.
func NewFactory(logger *zap.Logger) LoggerFactory {
	return &factory{logger: logger}
}

// Bg creates a context-unaware logger.
func (b factory) Bg() Logger {
	return &logger{logger: b.logger}
}

// For returns a context-aware Logger.
// TODO: OpenCensus implementation
func (b factory) For(ctx context.Context) Logger {
	return b.Bg()
}

// With creates a child logger, and optionally adds some context fields to that logger.
func (b factory) With(fields ...zapcore.Field) LoggerFactory {
	return &factory{logger: b.logger.With(fields...)}
}
