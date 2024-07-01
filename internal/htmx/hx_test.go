// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package htmx

import (
	"net/http"
	"reflect"
	"testing"
)

func headerWith(key, value string) http.Header {
	h := http.Header{}
	h.Add(key, value)
	return h
}

func TestHxx_Boosted(t *testing.T) {
	t.Parallel()

	f := func(name string, headers http.Header, expected bool) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			t.Helper()
			t.Parallel()
			hxx := &Hxr{Header: headers}
			if got := hxx.Boosted(); got != expected {
				t.Errorf("Hxx.Boosted() = %v, want %v", got, expected)
			}
		})
	}

	f("test missing", http.Header{}, false)
	f("test presence", headerWith("hx-boosted", "true"), true)
}

func TestHxx_CurrentUrl_HasPrefix(t *testing.T) {
	t.Parallel()

	f := func(name string, headers http.Header, prefix string, want bool) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			t.Helper()
			t.Parallel()
			hxx := &Hxr{Header: headers}
			if got := hxx.CurrentUrl().HasPathPrefix(prefix); !reflect.DeepEqual(got, want) {
				t.Errorf("Hxx.CurrentUrl() = %v, want %v", got, want)
			}
		})
	}

	f("test missing", http.Header{}, "/something", false)
	f("present", headerWith("hx-current-url", "http://example.com/something"), "/something", true)
	f("present but wrong", headerWith("hx-current-url", "http://example.com/somethinX"), "/something", false)
}

func TestHxx_HistoryRestoreRequest(t *testing.T) {
	t.Parallel()

	f := func(name string, headers http.Header, want bool) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			t.Helper()
			t.Parallel()
			hxx := &Hxr{Header: headers}
			if got := hxx.HistoryRestoreRequest(); got != want {
				t.Errorf("Hxx.HistoryRestoreRequest() = %v, want %v", got, want)
			}
		})
	}

	f("test missing", http.Header{}, false)
	f("present", headerWith("hx-history-restore-request", "true"), true)
}

func TestHxx_Prompt(t *testing.T) {
	t.Parallel()

	f := func(name string, headers http.Header, want string) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			t.Helper()
			t.Parallel()
			hxx := &Hxr{Header: headers}
			if got := hxx.Prompt(); got != want {
				t.Errorf("Hxx.Prompt() = %v, want %v", got, want)
			}
		})
	}

	f("test missing", http.Header{}, "")
	f("present", headerWith("hx-prompt", "hodor"), "hodor")
}

func TestHxx_Request(t *testing.T) {
	t.Parallel()

	f := func(name string, headers http.Header, want bool) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			t.Helper()
			t.Parallel()
			hxx := &Hxr{Header: headers}
			if got := hxx.IsRequest(); got != want {
				t.Errorf("Hxx.Request() = %v, want %v", got, want)
			}
		})
	}

	f("test missing", http.Header{}, false)
	f("present", headerWith("hx-request", "true"), true)
}

func TestHxx_Target(t *testing.T) {
	t.Parallel()

	f := func(name string, headers http.Header, want string) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			t.Helper()
			t.Parallel()
			hxx := &Hxr{Header: headers}
			if got := hxx.Target(); got != want {
				t.Errorf("Hxx.Target() = %v, want %v", got, want)
			}
		})
	}

	f("test missing", http.Header{}, "")
	f("present", headerWith("hx-target", "#someid"), "#someid")
}

func TestHxx_TriggerName(t *testing.T) {
	t.Parallel()

	f := func(name string, headers http.Header, want string) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			t.Helper()
			t.Parallel()
			hxx := &Hxr{Header: headers}
			if got := hxx.TriggerName(); got != want {
				t.Errorf("Hxx.TriggerName() = %v, want %v", got, want)
			}
		})
	}

	f("test missing", http.Header{}, "")
	f("present", headerWith("hx-trigger-name", "someid"), "someid")
}

func TestHxx_Trigger(t *testing.T) {
	t.Parallel()

	f := func(name string, headers http.Header, want string) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			t.Helper()
			t.Parallel()
			hxx := &Hxr{Header: headers}
			if got := hxx.Trigger(); got != want {
				t.Errorf("Hxx.Trigger() = %v, want %v", got, want)
			}
		})
	}

	f("test missing", http.Header{}, "")
	f("present", headerWith("hx-trigger", "someid"), "someid")
}
