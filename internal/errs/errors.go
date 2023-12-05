// Copyright 2018 Twitch Interactive, Inc.  All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may not
// use this file except in compliance with the License. A copy of the License is
// located at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// or in the "license" file accompanying this file. This file is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.
//
// originally from: https://github.com/twitchtv/twirp/blob/3c51d65f753a1049c77bc51e0e8c7922b1fb7e4d/errors.go
// modified 2023
package errs

import (
	"fmt"
)

// Error represents an error in a Twirp service call.
type Error interface {
	// Code is of the valid error codes.
	Code() ErrorCode

	// Msg returns a human-readable, unstructured messages describing the error.
	Msg() string

	// WithMeta returns a copy of the Error with the given key-value pair attached
	// as metadata. If the key is already set, it is overwritten.
	WithMeta(key string, val string) Error

	// Meta returns the stored value for the given key. If the key has no set
	// value, Meta returns an empty string. There is no way to distinguish between
	// an unset value and an explicit empty string.
	Meta(key string) string

	// MetaMap returns the complete key-value metadata map stored on the error.
	MetaMap() map[string]string

	// Error returns a string of the form "twirp error <Code>: <Msg>"
	Error() string

	// wrap another error
	Wrap(err error) Error
}

// ErrorCode represents a Twirp error type.
//
//go:generate stringer -type=ErrorCode
type ErrorCode byte

// Valid error types. Most error types are equivalent to gRPC status codes
// and follow similar semantics.
// ref: https://grpc.github.io/grpc/core/md_doc_statuscodes.html
const (
	// NoError is the zero-value, is considered an empty error and should not be
	// used.
	NoError ErrorCode = 0x00

	// Canceled indicates the operation was cancelled (typically by the caller).
	Canceled ErrorCode = 0x01

	// Unknown error. For example when handling errors raised by APIs that do not
	// return enough error information.
	Unknown ErrorCode = 0x02

	// InvalidArgument indicates client specified an invalid argument. It
	// indicates arguments that are problematic regardless of the state of the
	// system (i.e. a malformed file name, required argument, number out of range,
	// etc.).
	InvalidArgument ErrorCode = 0x03

	// DeadlineExceeded means operation expired before completion. For operations
	// that change the state of the system, this error may be returned even if the
	// operation has completed successfully (timeout).
	DeadlineExceeded ErrorCode = 0x04

	// NotFound means some requested entity was not found.
	NotFound ErrorCode = 0x05

	// AlreadyExists means an attempt to create an entity failed because one
	// already exists.
	AlreadyExists ErrorCode = 0x06

	// PermissionDenied indicates the caller does not have permission to execute
	// the specified operation. It must not be used if the caller cannot be
	// identified (Unauthenticated).
	PermissionDenied ErrorCode = 0x07

	// ResourceExhausted indicates some resource has been exhausted or rate-limited,
	// perhaps a per-user quota, or perhaps the entire file system is out of space.
	ResourceExhausted ErrorCode = 0x08

	// FailedPrecondition indicates operation was rejected because the system is
	// not in a state required for the operation's execution. For example, doing
	// an rmdir operation on a directory that is non-empty, or on a non-directory
	// object, or when having conflicting read-modify-write on the same resource.
	FailedPrecondition ErrorCode = 0x09

	// Aborted indicates the operation was aborted, typically due to a concurrency
	// issue like sequencer check failures, transaction aborts, etc.
	Aborted ErrorCode = 0x0a

	// OutOfRange means operation was attempted past the valid range. For example,
	// seeking or reading past end of a paginated collection.
	//
	// Unlike InvalidArgument, this error indicates a problem that may be fixed if
	// the system state changes (i.e. adding more items to the collection).
	//
	// There is a fair bit of overlap between FailedPrecondition and OutOfRange.
	// We recommend using OutOfRange (the more specific error) when it applies so
	// that callers who are iterating through a space can easily look for an
	// OutOfRange error to detect when they are done.
	OutOfRange ErrorCode = 0x0b

	// Unimplemented indicates operation is not implemented or not
	// supported/enabled in this service.
	Unimplemented ErrorCode = 0x0c

	// Internal errors. When some invariants expected by the underlying system
	// have been broken. In other words, something bad happened in the library or
	// backend service. Do not confuse with HTTP Internal Server Error; an
	// Internal error could also happen on the client code, i.e. when parsing a
	// server response.
	Internal ErrorCode = 0x0d

	// Unavailable indicates the service is currently unavailable. This is a most
	// likely a transient condition and may be corrected by retrying with a
	// backoff.
	Unavailable ErrorCode = 0x0e

	// DataLoss indicates unrecoverable data loss or corruption.
	DataLoss ErrorCode = 0x0f

	// Unauthenticated indicates the request does not have valid authentication
	// credentials for the operation.
	Unauthenticated ErrorCode = 0x10
)

// code.Error(msg) builds a new Twirp error with code and msg. Example:
//
//	twirp.NotFound.Error("Resource not found")
//	twirp.Internal.Error("Oops")
func (code ErrorCode) Error(msg string) Error {
	return NewError(code, msg)
}

func (code ErrorCode) Is(err error) bool {
	if u, ok := err.(interface{ Code() ErrorCode }); ok {
		return u.Code() == code
	}
	return false
}

// code.Errorf(msg, args...) builds a new Twirp error with code and formatted msg.
// The format may include "%w" to wrap other errors. Examples:
//
//	twirp.Internal.Error("Oops: %w", originalErr)
//	twirp.NotFound.Error("Resource not found with id: %q", resourceID)
func (code ErrorCode) Errorf(msgFmt string, a ...interface{}) Error {
	return NewErrorf(code, msgFmt, a...)
}

// NewError builds a twirp.Error. The code must be one of the valid predefined constants.
// To add metadata, use .WithMeta(key, value) method after building the error.
func NewError(code ErrorCode, msg string) Error {
	return &svcerr{code: code, msg: msg}
}

// NewErrorf builds a twirp.Error with a formatted msg.
// The format may include "%w" to wrap other errors. Examples:
//
//	twirp.NewErrorf(twirp.Internal, "Oops: %w", originalErr)
//	twirp.NewErrorf(twirp.NotFound, "resource with id: %q", resourceID)
func NewErrorf(code ErrorCode, msgFmt string, a ...interface{}) Error {
	err := fmt.Errorf(msgFmt, a...)       // format error message, may include "%w" with an original error
	svcerr := NewError(code, err.Error()) // use the error as msg
	return svcerr.Wrap(err)               // wrap so the original error can be identified with errors.Is
}

// InvalidArgumentError is a convenience constructor for InvalidArgument errors.
// The argument name is included on the "argument" metadata for convenience.
func InvalidArgumentError(argument string, msg string) Error {
	err := NewError(InvalidArgument, msg)
	err = err.WithMeta("argument", argument)
	return err
}

// twirp.Error implementation
type svcerr struct {
	code ErrorCode
	msg  string
	meta map[string]string
}

func (e *svcerr) Code() ErrorCode            { return e.code }
func (e *svcerr) Msg() string                { return e.msg }
func (e *svcerr) Wrap(err error) Error       { return &wraperr{e, err} }
func (e *svcerr) MetaMap() map[string]string { return e.meta }
func (e *svcerr) Error() string              { return fmt.Sprintf("error %s: %s", e.code, e.msg) }

func (e *svcerr) Meta(key string) string {
	if e.meta != nil {
		return e.meta[key] // also returns "" if key is not in meta map
	}
	return ""
}

func (e *svcerr) WithMeta(key string, value string) Error {
	newErr := &svcerr{
		code: e.code,
		msg:  e.msg,
		meta: make(map[string]string, len(e.meta)),
	}
	for k, v := range e.meta {
		newErr.meta[k] = v
	}
	newErr.meta[key] = value
	return newErr
}

// wraperr is the error returned by twirp.InternalErrorWith(err), which is used by clients.
// Implements Unwrap() to allow go 1.13+ errors.Is/As checks,
// and Cause() to allow (github.com/pkg/errors).Unwrap.
type wraperr struct {
	wrapper Error
	cause   error
}

func (e *wraperr) Code() ErrorCode            { return e.wrapper.Code() }
func (e *wraperr) Msg() string                { return e.wrapper.Msg() }
func (e *wraperr) Meta(key string) string     { return e.wrapper.Meta(key) }
func (e *wraperr) MetaMap() map[string]string { return e.wrapper.MetaMap() }
func (e *wraperr) Error() string              { return e.wrapper.Error() }
func (e *wraperr) Unwrap() error              { return e.cause }
func (e *wraperr) Wrap(err error) Error       { return &wraperr{e, err} }
func (e *wraperr) WithMeta(key string, val string) Error {
	return &wraperr{e.wrapper.WithMeta(key, val), e.cause}
}

func (e *wraperr) Is(err error) bool {
	if u, ok := err.(interface{ Code() ErrorCode }); ok {
		return u.Code() == e.Code()
	}
	return false
}
