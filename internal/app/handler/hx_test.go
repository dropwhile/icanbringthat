package handler

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

	tests := []struct {
		name   string
		fields http.Header
		want   bool
	}{
		{"test missing", http.Header{}, false},
		{"present", headerWith("hx-boosted", "true"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hxx := &Hxx{
				Header: tt.fields,
			}
			if got := hxx.Boosted(); got != tt.want {
				t.Errorf("Hxx.Boosted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHxx_CurrentUrl_HasPrefix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		fields http.Header
		prefix string
		want   bool
	}{
		{"test missing", http.Header{}, "/something", false},
		{
			"present",
			headerWith("hx-current-url", "http://example.com/something"),
			"/something",
			true,
		},
		{
			"present but wrong",
			headerWith("hx-current-url", "http://example.com/somethinX"),
			"/something",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hxx := &Hxx{
				Header: tt.fields,
			}
			if got := hxx.CurrentUrl().HasPathPrefix(tt.prefix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Hxx.CurrentUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHxx_HistoryRestoreRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		fields http.Header
		want   bool
	}{
		{"test missing", http.Header{}, false},
		{"present", headerWith("hx-history-restore-request", "true"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hxx := &Hxx{
				Header: tt.fields,
			}
			if got := hxx.HistoryRestoreRequest(); got != tt.want {
				t.Errorf("Hxx.HistoryRestoreRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHxx_Prompt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		fields http.Header
		want   string
	}{
		{"test missing", http.Header{}, ""},
		{"present", headerWith("hx-prompt", "hodor"), "hodor"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hxx := &Hxx{
				Header: tt.fields,
			}
			if got := hxx.Prompt(); got != tt.want {
				t.Errorf("Hxx.Prompt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHxx_Request(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		fields http.Header
		want   bool
	}{
		{"test missing", http.Header{}, false},
		{"present", headerWith("hx-request", "true"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hxx := &Hxx{
				Header: tt.fields,
			}
			if got := hxx.Request(); got != tt.want {
				t.Errorf("Hxx.Request() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHxx_Target(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		fields http.Header
		want   string
	}{
		{"test missing", http.Header{}, ""},
		{"present", headerWith("hx-target", "#someid"), "#someid"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hxx := &Hxx{
				Header: tt.fields,
			}
			if got := hxx.Target(); got != tt.want {
				t.Errorf("Hxx.Target() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHxx_TriggerName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		fields http.Header
		want   string
	}{
		{"test missing", http.Header{}, ""},
		{"present", headerWith("hx-trigger-name", "someid"), "someid"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hxx := &Hxx{
				Header: tt.fields,
			}
			if got := hxx.TriggerName(); got != tt.want {
				t.Errorf("Hxx.TriggerName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHxx_Trigger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		fields http.Header
		want   string
	}{
		{"test missing", http.Header{}, ""},
		{"present", headerWith("hx-trigger", "someid"), "someid"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hxx := &Hxx{
				Header: tt.fields,
			}
			if got := hxx.Trigger(); got != tt.want {
				t.Errorf("Hxx.Trigger() = %v, want %v", got, tt.want)
			}
		})
	}
}
