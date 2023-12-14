package logger

import (
	"io"
	"log"
	"log/slog"
	"os"
	"strings"
)

// defaults to info level
var logLevel = &slog.LevelVar{}

type Options struct {
	UseLocalTime bool
	OmitTime     bool
	OmitSource   bool
	Sink         io.Writer
	Prependers   []AttrExtractor
	Appenders    []AttrExtractor
}

func newContextHandler(next slog.Handler, opts Options) *slog.Logger {
	// add defaults
	prependers := []AttrExtractor{ContextExtractor}
	// add custom additions
	prependers = append(prependers, opts.Prependers...)
	h := &ContextHandler{
		next,
		prependers,
		opts.Appenders,
	}
	return slog.New(h)
}

func NewConsoleLogger(opts Options) *slog.Logger {
	if opts.Sink == nil {
		opts.Sink = os.Stderr
	}
	logHandler := slog.NewTextHandler(
		opts.Sink,
		&slog.HandlerOptions{
			Level:       logLevel,
			AddSource:   !opts.OmitSource,
			ReplaceAttr: replaceAttr(opts),
		},
	)
	return newContextHandler(logHandler, opts)
}

func NewJsonLogger(opts Options) *slog.Logger {
	if opts.Sink == nil {
		opts.Sink = os.Stderr
	}
	logHandler := slog.NewJSONHandler(
		opts.Sink,
		&slog.HandlerOptions{
			Level:       logLevel,
			AddSource:   !opts.OmitSource,
			ReplaceAttr: replaceAttr(opts),
		},
	)
	return newContextHandler(logHandler, opts)
}

func NewTestLogger(opts Options) *slog.Logger {
	if opts.Sink == nil {
		opts.Sink = os.Stderr
	}
	// always omit time for test logs,
	// to enable log matching if desired.
	opts.OmitTime = true
	logHandler := slog.NewTextHandler(
		opts.Sink,
		&slog.HandlerOptions{
			Level:       logLevel,
			AddSource:   !opts.OmitSource,
			ReplaceAttr: replaceAttr(opts),
		},
	)

	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		logLevel.Set(slog.LevelDebug)
	default:
		logLevel.Set(slog.LevelInfo)
	}

	return newContextHandler(logHandler, opts)
}

func SetupLogging(mkLogger func(Options) *slog.Logger, opts *Options) {
	if opts == nil {
		opts = &Options{}
	}
	if opts.Sink == nil {
		opts.Sink = os.Stderr
	}
	logger := mkLogger(*opts)
	slog.SetDefault(logger)
	log.SetOutput(&logWriter{opts.Sink})
	log.SetFlags(log.Lshortfile)
}

func SetLevel(level slog.Level) {
	logLevel.Set(level)
}
