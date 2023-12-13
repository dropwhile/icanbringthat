package logger

import (
	"context"
	"log/slog"
	"strings"

	"github.com/pkg/errors"
)

// default info
var logLevel = new(slog.LevelVar)

const (
	LevelTrace = slog.Level(-8)
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelError = slog.LevelError
	LevelFatal = slog.Level(12)
)

var (
	LevelTraceStr = slog.StringValue("TRC")
	LevelDebugStr = slog.StringValue("DBG")
	LevelInfoStr  = slog.StringValue("INF")
	LevelErrorStr = slog.StringValue("ERR")
	LevelFatalStr = slog.StringValue("FTL")
)

// text unmarshaller for envconfig
type LoggerLevel slog.Level

func (l *LoggerLevel) UnmarshalText(text []byte) error {
	var t slog.Level
	switch strings.ToLower(string(text)) {
	case "trace":
		t = LevelTrace
	case "debug":
		t = LevelDebug
	case "info":
		t = LevelInfo
	case "error":
		t = LevelError
	default:
		return errors.New("unknown log level")
	}
	*l = LoggerLevel(t)
	return nil
}

func (l LoggerLevel) String() string {
	switch slog.Level(l) {
	case LevelTrace:
		return "trace"
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelError:
		return "error"
	default:
		return "unknown"
	}
}

func (l LoggerLevel) Level() slog.Level {
	return slog.Level(l)
}

func (l LoggerLevel) Covers(lvl slog.Level) bool {
	return lvl >= slog.Level(l)
}

func (l LoggerLevel) LogValuer() slog.Value {
	return slog.StringValue(l.String())
}

func SetLevel(level slog.Leveler) {
	logLevel.Set(level.Level())
}

func Enabled[T ~int](level T) bool {
	return slog.Default().Enabled(context.Background(), slog.Level(level))
}
