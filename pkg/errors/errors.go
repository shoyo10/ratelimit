package errors

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type _err struct {
	Code     string
	HTTPCode int
	Message  string `json:"message"`
}

// Error error string
func (e *_err) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

var (
	ErrInvalidInput          = &_err{Code: "400001", Message: "one of the input is invalid", HTTPCode: http.StatusBadRequest}
	ErrUnauthorized          = &_err{Code: "401001", Message: "one of the input is invalid", HTTPCode: http.StatusUnauthorized}
	ErrResourceNotFound      = &_err{Code: "404001", Message: "resource not found", HTTPCode: http.StatusNotFound}
	ErrResourceAlreadyExists = &_err{Code: "409001", Message: "resource already exist", HTTPCode: http.StatusConflict}
	ErrTooManyRequests       = &_err{Code: "429001", Message: "too many request", HTTPCode: http.StatusTooManyRequests}
	ErrInternalServerError   = &_err{Code: "500001", Message: "internal server error", HTTPCode: http.StatusInternalServerError}
)

// New new error with stack
func New(err error) error {
	return errors.WithStack(err)
}

// WrapfWithMessage replace original error message by the input message
func WrapfWithMessage(err error, format string, args ...interface{}) error {
	e := Cause(err)
	_e, ok := e.(*_err)
	if !ok {
		_e := ErrInternalServerError
		_e.Message = fmt.Sprintf(format, args...)
		return WithStack(_e)
	}
	_e.Message = fmt.Sprintf(format, args...)
	return WithStack(_e)
}

// WithStack annotates err with a stack trace at the point WithStack was called.
// If err is nil, WithStack returns nil.
var WithStack = errors.WithStack

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
var Wrap = errors.Wrap

// Wrapf returns an error annotating err with a stack trace
// at the point Wrapf is called, and the format specifier.
// If err is nil, Wrapf returns nil.
var Wrapf = errors.Wrapf

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//     type causer interface {
//            Cause() error
//     }
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
var Cause = errors.Cause

// Is reports whether any error in err's chain matches target.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
var Is = errors.Is
