package errs

// Code represents a error type.
//
//go:generate tool stringer -type=Code
type Code byte

// Valid error types. Most error types are equivalent to gRPC status codes
// and follow similar semantics.
// ref: https://grpc.github.io/grpc/core/md_doc_statuscodes.html
const (
	// NoError is the zero-value, is considered an empty error and should not be
	// used.
	NoError Code = 0x00

	// Canceled indicates the operation was cancelled (typically by the caller).
	Canceled Code = 0x01

	// Unknown error. For example when handling errors raised by APIs that do not
	// return enough error information.
	Unknown Code = 0x02

	// InvalidArgument indicates client specified an invalid argument. It
	// indicates arguments that are problematic regardless of the state of the
	// system (i.e. a malformed file name, required argument, number out of range,
	// etc.).
	InvalidArgument Code = 0x03

	// DeadlineExceeded means operation expired before completion. For operations
	// that change the state of the system, this error may be returned even if the
	// operation has completed successfully (timeout).
	DeadlineExceeded Code = 0x04

	// NotFound means some requested entity was not found.
	NotFound Code = 0x05

	// AlreadyExists means an attempt to create an entity failed because one
	// already exists.
	AlreadyExists Code = 0x06

	// PermissionDenied indicates the caller does not have permission to execute
	// the specified operation. It must not be used if the caller cannot be
	// identified (Unauthenticated).
	PermissionDenied Code = 0x07

	// ResourceExhausted indicates some resource has been exhausted or rate-limited,
	// perhaps a per-user quota, or perhaps the entire file system is out of space.
	ResourceExhausted Code = 0x08

	// FailedPrecondition indicates operation was rejected because the system is
	// not in a state required for the operation's execution. For example, doing
	// an rmdir operation on a directory that is non-empty, or on a non-directory
	// object, or when having conflicting read-modify-write on the same resource.
	FailedPrecondition Code = 0x09

	// Aborted indicates the operation was aborted, typically due to a concurrency
	// issue like sequencer check failures, transaction aborts, etc.
	Aborted Code = 0x0a

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
	OutOfRange Code = 0x0b

	// Unimplemented indicates operation is not implemented or not
	// supported/enabled in this service.
	Unimplemented Code = 0x0c

	// Internal errors. When some invariants expected by the underlying system
	// have been broken. In other words, something bad happened in the library or
	// backend service. Do not confuse with HTTP Internal Server Error; an
	// Internal error could also happen on the client code, i.e. when parsing a
	// server response.
	Internal Code = 0x0d

	// Unavailable indicates the service is currently unavailable. This is a most
	// likely a transient condition and may be corrected by retrying with a
	// backoff.
	Unavailable Code = 0x0e

	// DataLoss indicates unrecoverable data loss or corruption.
	DataLoss Code = 0x0f

	// Unauthenticated indicates the request does not have valid authentication
	// credentials for the operation.
	Unauthenticated Code = 0x10
)

// code.Wrap(error) wraps an error with code, and optional modifiers.
// Example:
//
//	errs.NotFound.Wrap(errors.New("resource not found"))
//	errs.Internal.Error(errors.New("oops"))
//	errs.Internal.Error(
//	  errors.New("resource not found"),
//	  WithInfo("Internal server error."))
func (code Code) Wrap(err error, wfs ...withFunc) Error {
	e := &codeErr{code: code, err: err}
	for _, w := range wfs {
		e = w(e)
	}
	return e
}

// code.Error(text) builds a new error with code and text.
// Example:
//
//	errs.NotFound.Error("resource not found")
//	errs.Internal.Error("oops")
func (code Code) Error(text string) Error {
	return newError(code, text)
}

// code.Errorf(format, args...) builds a new error with code and formatted text.
// The format may include "%w" to wrap other errors. Examples:
//
//	errs.Internal.Error("oops: %w", originalErr)
//	errs.NotFound.Error("resource not found with id: %q", resourceID)
func (code Code) Errorf(format string, a ...any) Error {
	return newErrorf(code, format, a...)
}
