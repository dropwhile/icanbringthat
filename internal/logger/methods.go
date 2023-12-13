package logger

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"time"
)

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
