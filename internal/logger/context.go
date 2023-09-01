// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package logger

import (
	"context"
	"log/slog"
	"slices"
	"time"
)

type logCtxAttrKey struct{}

func PrependAttr(ctx context.Context, args ...any) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if v, ok := ctx.Value(logCtxAttrKey{}).([]slog.Attr); ok {
		// Clip to ensure this is a scoped copy
		return context.WithValue(ctx, logCtxAttrKey{},
			append(slices.Clip(v), argsToAttrSlice(args)...))
	}
	return context.WithValue(ctx, logCtxAttrKey{}, argsToAttrSlice(args))
}

// ExtractPrepended is an AttrExtractor that returns the prepended attributes
// stored in the context. The returned slice should not be appended to or
// modified in any way. Doing so will cause a race condition.
func ContextExtractor(ctx context.Context, _ time.Time, _ slog.Level, _ string) []slog.Attr {
	if v, ok := ctx.Value(logCtxAttrKey{}).([]slog.Attr); ok {
		return v
	}
	return nil
}

// Turn a slice of arguments, some of which pairs of primitives,
// some might be attributes already, into a slice of attributes.
// This is copied from golang sdk.
func argsToAttrSlice(args []any) []slog.Attr {
	var (
		attr  slog.Attr
		attrs []slog.Attr
	)
	for len(args) > 0 {
		attr, args = argsToAttr(args)
		attrs = append(attrs, attr)
	}
	return attrs
}

// This is copied from golang sdk.
const badKey = "!BADKEY"

// argsToAttr turns a prefix of the nonempty args slice into an Attr
// and returns the unconsumed portion of the slice.
// If args[0] is an Attr, it returns it.
// If args[0] is a string, it treats the first two elements as
// a key-value pair.
// Otherwise, it treats args[0] as a value with a missing key.
// This is copied from golang sdk.
func argsToAttr(args []any) (slog.Attr, []any) {
	switch x := args[0].(type) {
	case string:
		if len(args) == 1 {
			return slog.String(badKey, x), nil
		}
		return slog.Any(x, args[1]), args[2:]

	case slog.Attr:
		return x, args[1:]

	default:
		return slog.Any(badKey, x), args[1:]
	}
}
