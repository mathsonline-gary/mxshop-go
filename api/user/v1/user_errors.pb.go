// Code generated by protoc-gen-go-errors. DO NOT EDIT.

package v1

import (
	fmt "fmt"

	errors "github.com/zycgary/mxshop-go/gmicro/errors"
)

type ErrorCode int32

const (
	CodeUserNotFound ErrorCode = 0
	CodeUserExists   ErrorCode = 1
)

func IsErrorUserNotFound(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Code == int32(CodeUserNotFound)
}

func ErrorUserNotFound(format string, args ...interface{}) *errors.Error {
	return errors.New(int32(CodeUserNotFound), fmt.Sprintf(format, args...))
}

func IsErrorUserExists(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Code == int32(CodeUserExists)
}

func ErrorUserExists(format string, args ...interface{}) *errors.Error {
	return errors.New(int32(CodeUserExists), fmt.Sprintf(format, args...))
}
