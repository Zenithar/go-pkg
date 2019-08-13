package errors

import (
	"context"
	"io"
	"reflect"

	"golang.org/x/xerrors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DoNotWrap reports whether an error should not be wrapped in the Error
// type from this package.
// It returns true if err is a retry error, a context error, io.EOF, or if it wraps
// one of those.
func DoNotWrap(err error) bool {
	if xerrors.Is(err, io.EOF) {
		return true
	}
	if xerrors.Is(err, context.Canceled) {
		return true
	}
	if xerrors.Is(err, context.DeadlineExceeded) {
		return true
	}
	return false
}

var (
	grpcCodeMap = map[codes.Code]ErrorCode{
		codes.AlreadyExists:      AlreadyExists,
		codes.Aborted:            Aborted,
		codes.Canceled:           Canceled,
		codes.DataLoss:           DataLoss,
		codes.DeadlineExceeded:   DeadlineExceeded,
		codes.FailedPrecondition: FailedPrecondition,
		codes.Internal:           Internal,
		codes.InvalidArgument:    InvalidArgument,
		codes.NotFound:           NotFound,
		codes.OK:                 OK,
		codes.OutOfRange:         OutOfRange,
		codes.PermissionDenied:   PermissionDenied,
		codes.Unauthenticated:    Unauthenticated,
		codes.Unavailable:        Unavailable,
		codes.Unimplemented:      Unimplemented,
		codes.Unknown:            Unknown,
	}
)

// GRPCCode extracts the gRPC status code and converts it into an ErrorCode.
// It returns Unknown if the error isn't from gRPC.
func GRPCCode(err error) ErrorCode {
	if code, ok := grpcCodeMap[status.Code(err)]; ok {
		return code
	}
	return Unknown
}

// ErrorAs is a helper for the ErrorAs method of an API's portable type.
// It performs some initial nil checks, and does a single level of unwrapping
// when err is a *gcerr.Error. Then it calls its errorAs argument, which should
// be a driver implementation of ErrorAs.
func ErrorAs(err error, target interface{}, errorAs func(error, interface{}) bool) bool {
	if err == nil {
		return false
	}
	if target == nil {
		panic("ErrorAs target cannot be nil")
	}
	val := reflect.ValueOf(target)
	if val.Type().Kind() != reflect.Ptr || val.IsNil() {
		panic("ErrorAs target must be a non-nil pointer")
	}
	if e, ok := err.(*Error); ok {
		err = e.Unwrap()
	}
	return errorAs(err, target)
}
