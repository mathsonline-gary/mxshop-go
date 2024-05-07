package errors

import (
	"errors"
	"fmt"
)

const (
	UnexpectedCode int32 = 500
)

var (
	_ error = (*Error)(nil)
)

type Error struct {
	Code int32  `json:"code"`
	Msg  string `json:"message"`
	Err  error  `json:"error"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("error: code = %d desc = %s", e.Code, e.Msg)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func New(code int32, message string) *Error {
	return &Error{
		Code: code,
		Msg:  message,
	}
}

func FromError(err error) *Error {
	if err == nil {
		return nil
	}
	if e := new(Error); errors.As(err, &e) {
		return e
	}

	return New(UnexpectedCode, err.Error())
}
