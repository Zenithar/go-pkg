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
)

var defaultFactory LoggerFactory

// -----------------------------------------------------------------------------

// SetLogger defines the default package logger
func SetLogger(instance LoggerFactory) {
	if defaultFactory != nil {
		defaultFactory.Bg().Debug("Replacing logger factory", String("old", defaultFactory.Name()), String("new", instance.Name()))
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

// DefaultFactory returns the logger factory
func DefaultFactory() LoggerFactory {
	return checkFactory(defaultFactory)
}

// -----------------------------------------------------------------------------

func checkFactory(defaultFactory LoggerFactory) LoggerFactory {
	if defaultFactory == nil {
		panic("Unable to create logger instance, you have to register an adapter first.")
	}
	return defaultFactory
}
