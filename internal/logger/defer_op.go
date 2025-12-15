// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package logger

import (
	"log/slog"
)

type deferredOp[T any] struct {
	emitter func() T
}

func (e *deferredOp[V]) LogValue() slog.Value {
	if e == nil {
		return slog.Value{}
	}
	return slog.AnyValue(e.emitter())
}

// DeferOperation is used as a log attr to defer a slow logging operation, and
// to avoid processing if level would result in the log not emitting.
// The `value` param is used as input to the `f` function param.
//
// Example: (note: this example isn't especially slow)
//
//		slog.DebugContext(ctx, "something is happening",
//			"magic number",
//			logger.DeferOperation(func() string {
//				return fmt.Sprintf("%d", 42)
//			})
//	 	)
func DeferOperation[T any](f func() T) *deferredOp[T] {
	return &deferredOp[T]{f}
}
