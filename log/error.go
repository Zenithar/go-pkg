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
	"io"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// CheckErr handles error correctly
func CheckErr(msg string, err error, fields ...zapcore.Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
		defaultFactory.Bg().Error(msg, fields...)
	}
}

// CheckErrCtx handles error correctly
func CheckErrCtx(ctx context.Context, msg string, err error, fields ...zapcore.Field) {
	if err != nil {
		fields = append(fields, zap.Error(errors.WithStack(err)))
		defaultFactory.For(ctx).Error(msg, fields...)
	}
}

// SafeClose handles the closer error
func SafeClose(c io.Closer, msg string, fields ...zapcore.Field) {
	if cerr := c.Close(); cerr != nil {
		fields = append(fields, zap.Error(errors.WithStack(cerr)))
		defaultFactory.Bg().Error(msg, fields...)
	}
}
