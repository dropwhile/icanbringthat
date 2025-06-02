package errs

import (
	"errors"
	"maps"
	"strings"
)

type withFunc func(err *codeErr) *codeErr

func WithInfo(info string) withFunc {
	return func(e *codeErr) *codeErr {
		if e == nil {
			return nil
		}

		e.msg = info
		return e
	}
}

func WithMeta(key, value string) withFunc {
	return func(e *codeErr) *codeErr {
		if e == nil {
			return nil
		}

		if e.meta == nil {
			e.meta = make(map[string]string, 0)
		}

		e.meta[key] = value
		return e
	}
}

func WithMetaVals(vals map[string]string) withFunc {
	return func(e *codeErr) *codeErr {
		if e == nil {
			return nil
		}

		if e.meta == nil {
			e.meta = make(map[string]string, 0)
		}

		maps.Copy(e.meta, vals)
		return e
	}
}

// GetMsgs collects any human-readable messages in the error stack,
// and returns it as a space separated string.
func GetMsg(err error) string {
	details := CollectMsgs(err)

	if len(details) > 0 {
		return strings.Join(details, " ")
	}
	return ""
}

// CollectMsgs collects any human-readable messages in the error stack.
func CollectMsgs(err error) []string {
	var details []string

	for err != nil {
		if e, ok := err.(*codeErr); ok {
			if e.msg != "" {
				details = append(details, e.msg)
			}
		}

		err = errors.Unwrap(err)
	}

	return details
}

// ArgumentError is a convenience constructor for InvalidArgument errors.
// The argument name is included on the "argument" metadata for convenience.
func ArgumentError(argument string, text string) Error {
	e := InvalidArgument.
		Error(strings.Join([]string{argument, text}, " ")).
		WithMeta("argument", argument)
	return e
}
