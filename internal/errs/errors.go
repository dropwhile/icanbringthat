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
// modified 2023-2024
package errs

import (
	"errors"
	"fmt"
)

// Error represents an error in a service call.
type Error interface {
	// Code is of the valid error codes.
	Code() Code

	// Msg returns a human-readable, unstructured messages describing the error.
	Msg() string

	// Meta returns the stored value for the given key. If the key has no set
	// value, Meta returns an empty string. There is no way to distinguish between
	// an unset value and an explicit empty string.
	Meta(key string) string

	// MetaMap returns the complete key-value metadata map stored on the error.
	MetaMap() map[string]string

	// Error returns a string of the form "error <Code>: <Text>"
	Error() string

	// Unwrap returns the underlying error
	Unwrap() error

	// Is reports whether any error in err's tree matches target.
	Is(error) bool

	// WithMsg returns the Error with the given message set.
	WithMsg(info string) Error

	// WithMeta returns the Error with the key-value pair provided
	// as metadata. If the key is already set, it is overwritten.
	WithMeta(key, value string) Error

	// WithMetaVals returns the Error with the Meta values provided
	// as metadata. If a Meta.Key is already present, it is overwritten.
	WithMetaVals(vals map[string]string) Error
}

// newError builds a errs.Error. The code must be one of the valid predefined constants.
// To add metadata, use .WithMeta(key, value) method after building the error.
func newError(code Code, text string) Error {
	return &codeErr{
		code: code,
		err:  errors.New(text),
	}
}

// newErrorf builds a errs.Error with a formatted text.
// The format may include "%w" to wrap other errors. Examples:
//
//	errs.newErrorf(errs.Internal, "oops: %w", originalErr)
//	errs.newErrorf(errs.NotFound, "resource not found with id: %q", resourceID)
//
// To add metadata, use .WithMeta(key, value) method after building the error.
func newErrorf(code Code, format string, a ...any) Error {
	return &codeErr{
		code: code,
		err:  fmt.Errorf(format, a...),
	}
}

type codeErr struct {
	code Code              // error code
	err  error             // underlying error
	msg  string            // friendly messages
	meta map[string]string // metadata
}

// Code returns the svcErr ErrorCode.
func (e *codeErr) Code() Code {
	if e == nil {
		return NoError
	}
	return e.code
}

// Msg returns a human-readable, unstructured message describing the error from
// the error stack. If no human-readable message is present in the error stack,
// return the underlying top level error.Error() value.
func (e *codeErr) Msg() string {
	msg := GetMsg(e)
	if msg != "" {
		return msg
	}
	if e.err != nil {
		return e.err.Error()
	}
	return ""
}

// Meta returns the stored value for the given key. If the key has no set
// value, Meta returns an empty string. There is no way to distinguish between
// an unset value and an explicit empty string.
func (e *codeErr) Meta(key string) string {
	// note: return "" (zero value of string) if
	// * e is nil
	// * key is not in meta map
	// * if map is nil
	if e == nil {
		return ""
	}
	return e.meta[key]
}

// MetaMap returns the complete key-value metadata map stored on the error.
func (e *codeErr) MetaMap() map[string]string {
	if e == nil {
		return nil
	}
	return e.meta
}

// Error returns a string of the form "error <Code>: <Text>"
func (e *codeErr) Error() string {
	if e == nil {
		return ""
	}

	code := e.code.String()

	etxt := ""
	if e.err != nil {
		etxt = e.err.Error()
	}

	if etxt == "" {
		return code
	}
	return code + ": " + etxt
}

// Unwrap returns the underlying error
func (e *codeErr) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.err
}

// Is reports whether any error in err's tree matches target.
func (e *codeErr) Is(err error) bool {
	if u, ok := err.(interface{ Code() Code }); ok {
		return u.Code() == e.Code()
	}
	return false
}

// WithMsg returns the Error with the given message set.
func (e *codeErr) WithMsg(info string) Error {
	return WithInfo(info)(e)
}

// WithMeta returns the Error with the key-value pair provided
// as metadata. If the key is already set, it is overwritten.
func (e *codeErr) WithMeta(key, value string) Error {
	return WithMeta(key, value)(e)
}

// WithMetaVals returns the Error with the Meta values provided
// as metadata. If a Meta.Key is already present, it is overwritten.
func (e *codeErr) WithMetaVals(vals map[string]string) Error {
	return WithMetaVals(vals)(e)
}
