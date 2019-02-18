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
	l "log"

	"go.uber.org/zap"
)

var defaultFactory LoggerFactory

// -----------------------------------------------------------------------------

func init() {
	defaultLogger, err := zap.NewProduction()
	if err != nil {
		l.Fatalln(err)
	}

	defaultFactory = NewFactory(defaultLogger)
}

// SetLogger defines the default package logger
func SetLogger(instance LoggerFactory) {
	defaultFactory = instance
}

// -----------------------------------------------------------------------------

// Bg delegates a no-context logger
func Bg() Logger {
	return defaultFactory.Bg()
}

// For delegates a context logger
func For(ctx context.Context) Logger {
	return defaultFactory.For(ctx)
}

// Default returns the logger factory
func Default() LoggerFactory {
	return defaultFactory
}
