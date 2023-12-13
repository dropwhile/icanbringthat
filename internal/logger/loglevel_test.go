package logger

import (
	"log/slog"
	"testing"
)

func TestLoggerLevel_Covers(t *testing.T) {
	tests := []struct {
		name string
		l    LoggerLevel
		r    slog.Level
		want bool
	}{
		{
			"debug is supported by trace logging",
			LoggerLevel(LevelTrace),
			LevelDebug,
			true,
		},
		{
			"trace is supported by trace logging",
			LoggerLevel(LevelTrace),
			LevelTrace,
			true,
		},
		{
			"trace is not supported by debug logging",
			LoggerLevel(LevelDebug),
			LevelTrace,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Covers(tt.r); got != tt.want {
				t.Errorf("LoggerLevel.Covers() = %v, want %v", got, tt.want)
			}
		})
	}
}
