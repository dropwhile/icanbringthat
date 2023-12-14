package logger

import (
	"context"
	"io"
	"log"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"
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
	return NewContextHandler(logHandler, opts)
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
	return NewContextHandler(logHandler, opts)
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

	return NewContextHandler(logHandler, opts)
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

func Fatal(msg string, args ...any) {
	LogSkip(slog.Default(), 1, slog.LevelError,
		context.Background(), msg, args...)
	os.Exit(1)
}

func LogSkip(logger *slog.Logger, skip int, level slog.Level,
	ctx context.Context, msg string, args ...any,
) {
	if ctx == nil {
		ctx = context.Background()
	}
	if !logger.Enabled(ctx, level) {
		return
	}

	var pcs [1]uintptr
	// skip [runtime.Callers, this function, this function's caller]
	runtime.Callers(2+skip, pcs[:]) // skip [Callers, log, wrapper]

	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.Add(args...)
	_ = logger.Handler().Handle(ctx, r)
}

func LogAttrsSkip(logger *slog.Logger, skip int, level slog.Level,
	ctx context.Context, msg string, attrs ...slog.Attr,
) {
	if ctx == nil {
		ctx = context.Background()
	}
	if !logger.Enabled(ctx, level) {
		return
	}

	var pcs [1]uintptr
	// skip [runtime.Callers, this function, this function's caller]
	runtime.Callers(2+skip, pcs[:])

	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.AddAttrs(attrs...)
	_ = logger.Handler().Handle(ctx, r)
}
