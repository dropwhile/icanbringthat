// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package logger

import (
	"log/slog"
)

type defOp[T any, V any] struct {
	elem T
	emit func(T) V
}

func (e *defOp[T, V]) LogValue() slog.Value {
	if e == nil {
		return slog.Value{}
	}
	return slog.AnyValue(e.emit(e.elem))
}

func DeferOperation[T any, V any](value T, f func(T) V) *defOp[T, V] {
	return &defOp[T, V]{value, f}
}
