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

		e.info = info
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

func GetInfo(err error) string {
	var details []string
	for err != nil {
		if subErr, ok := err.(*codeErr); ok {
			if subErr.info != "" {
				details = append(details, subErr.info)
			}
		}

		err = errors.Unwrap(err)
	}
	return strings.Join(details, " ")
}

// ArgumentError is a convenience constructor for InvalidArgument errors.
// The argument name is included on the "argument" metadata for convenience.
func ArgumentError(argument string, msg string) Error {
	e := InvalidArgument.
		Error(strings.Join([]string{argument, msg}, " ")).
		WithMeta("argument", argument)
	return e
}
