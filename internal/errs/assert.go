// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package errs

import (
	"testing"
)

type Errorer[T ~byte | ~string] interface {
	Code() T
	Msg() string
	Meta(key string) string
}

func AssertError[T ~byte | ~string](t *testing.T, err error, code T, msg string, meta ...map[string]string) {
	t.Helper()
	terr, ok := err.(Errorer[T])
	if !ok {
		t.Errorf("not a twirp error type")
		return
	}
	if terr.Code() != code {
		t.Errorf("wrong code. have=%q, want=%q", terr.Code(), code)
	}
	if terr.Msg() != msg {
		t.Errorf("wrong msg. have=%q, want=%q", terr.Msg(), msg)
	}
	for _, m := range meta {
		for k, v := range m {
			x := terr.Meta(k)
			if x != v {
				t.Errorf("meta value %q mismatch. have=%q, want=%q",
					k, x, v)
			}
		}
	}
}
