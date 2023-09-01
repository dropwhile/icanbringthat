// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package util

import (
	"github.com/samber/mo"
)

func OptionFlatMapConvert[T any, V any](
	input mo.Option[T], mapF func(mo.Option[T]) (mo.Option[V], error),
) (mo.Option[V], error) {
	return mapF(input)
}

func OptionMapConvert[T any, V any](
	input mo.Option[T], mapF func(T) (V, error),
) (mo.Option[V], error) {
	if v, ok := input.Get(); ok {
		nv, err := mapF(v)
		if err != nil {
			return mo.None[V](), err
		}
		return mo.Some(nv), nil
	}
	return mo.None[V](), nil
}
