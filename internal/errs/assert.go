// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package errs

import (
	"testing"
)

func AssertError[T ~byte | ~string | ~uint32](t *testing.T, err error, code T, msg string, meta ...map[string]string) {
	t.Helper()

	codeErr, ok := err.(interface{ Code() T })
	if !ok {
		t.Errorf("not an expected error type: %T", codeErr)
		return
	}
	if codeErr.Code() != code {
		t.Errorf("wrong code. have=%q, want=%q", codeErr.Code(), code)
	}

	switch terr := err.(type) {
	case interface{ Msg() string }:
		if terr.Msg() != msg {
			t.Errorf("wrong msg. have=%q, want=%q", terr.Msg(), msg)
		}
	case interface{ Message() string }:
		if terr.Message() != msg {
			t.Errorf("wrong msg. have=%q, want=%q", terr.Message(), msg)
		}
	default:
		t.Errorf("not an expected error type: %T", err)
		return
	}
}
