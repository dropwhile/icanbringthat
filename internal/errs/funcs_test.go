package errs

import (
	"errors"
	"fmt"
	"testing"
)

func TestGetMsg(t *testing.T) {
	t.Parallel()

	check := func(err error, expected string) {
		t.Helper()
		if got := GetInfo(err); got != expected {
			t.Errorf("GetMsg(err) = %v, want %v", got, expected)
		}
	}

	t.Run("test non-svcErr error", func(t *testing.T) {
		check(errors.New("not-an-svc-err"), "")
	})

	t.Run("test a non-svcErr wrapping an svcErr error with no info", func(t *testing.T) {
		err := NotFound.Error("oops")
		e := fmt.Errorf("not-an-svc-err: %w", err)
		check(e, "")
	})

	t.Run("test a non-svcErr wrapping an svcErr error with info", func(t *testing.T) {
		err := NotFound.Error("oops").WithInfo("i-have-info")
		e := fmt.Errorf("not-an-svc-err: %w", err)
		check(e, "i-have-info")
	})

	t.Run("test an svcErr wrapping an svcErr error with info", func(t *testing.T) {
		err1 := NotFound.Error("oops1").WithInfo("The frobnicator is stuck.")
		err2 := NotFound.Wrap(err1).WithInfo("The frobnicator will not frobnicate.")
		e := fmt.Errorf("not-an-svc-err: %w", err2)
		check(e, "The frobnicator will not frobnicate. The frobnicator is stuck.")
	})
}
