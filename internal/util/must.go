// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package util

import "io"

func MustReadAll(r io.ReadCloser) []byte {
	d, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return d
}

func Must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}
