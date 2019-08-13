package errors

import (
	"fmt"

	"golang.org/x/xerrors"
)

// An Error describes a Go CDK error.
type Error struct {
	Code  ErrorCode
	msg   string
	frame xerrors.Frame
	err   error
}

func (e *Error) Error() string {
	return fmt.Sprint(e)
}

// Format error with xerrors
func (e *Error) Format(s fmt.State, c rune) {
	xerrors.FormatError(e, s, c)
}

// FormatError ...
func (e *Error) FormatError(p xerrors.Printer) (next error) {
	if e.msg == "" {
		p.Printf("code=%v", e.Code)
	} else {
		p.Printf("%s (code=%v)", e.msg, e.Code)
	}
	e.frame.Format(p)
	return e.err
}

// Unwrap returns the error underlying the receiver, which may be nil.
func (e *Error) Unwrap() error {
	return e.err
}

// New returns a new error with the given code, underlying error and message. Pass 1
// for the call depth if New is called from the function raising the error; pass 2 if
// it is called from a helper function that was invoked by the original function; and
// so on.
func New(c ErrorCode, err error, callDepth int, msg string) *Error {
	return &Error{
		Code:  c,
		msg:   msg,
		frame: xerrors.Caller(callDepth),
		err:   err,
	}
}

// Newf uses format and args to format a message, then calls New.
func Newf(c ErrorCode, err error, format string, args ...interface{}) *Error {
	return New(c, err, 2, fmt.Sprintf(format, args...))
}
