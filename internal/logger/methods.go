package logger

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"time"
)

func Trace(msg string, args ...any) {
	logx(context.Background(), slog.Default(), LevelTrace, msg, args...)
}

func Debug(msg string, args ...any) {
	logx(context.Background(), slog.Default(), LevelDebug, msg, args...)
}

func Info(msg string, args ...any) {
	logx(context.Background(), slog.Default(), LevelInfo, msg, args...)
}

func Error(msg string, args ...any) {
	logx(context.Background(), slog.Default(), LevelError, msg, args...)
}

func Fatal(msg string, args ...any) {
	logx(context.Background(), slog.Default(), LevelFatal, msg, args...)
	os.Exit(1)
}

func TraceCtx(ctx context.Context, msg string, attrs ...slog.Attr) {
	logxAttrs(ctx, slog.Default(), LevelTrace, msg, attrs...)
}

func DebugCtx(ctx context.Context, msg string, attrs ...slog.Attr) {
	logxAttrs(ctx, slog.Default(), LevelDebug, msg, attrs...)
}

func InfoCtx(ctx context.Context, msg string, attrs ...slog.Attr) {
	logxAttrs(ctx, slog.Default(), LevelInfo, msg, attrs...)
}

func ErrorCtx(ctx context.Context, msg string, attrs ...slog.Attr) {
	logxAttrs(ctx, slog.Default(), LevelError, msg, attrs...)
}

func FatalCtx(ctx context.Context, msg string, attrs ...slog.Attr) {
	logxAttrs(ctx, slog.Default(), LevelFatal, msg, attrs...)
	os.Exit(1)
}

func Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	logx(ctx, slog.Default(), level, msg, args...)
}

func LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	logxAttrs(ctx, slog.Default(), level, msg, attrs...)
}

func With(args ...any) *slog.Logger {
	return slog.Default().With(args...)
}

func WithGroup(name string) *slog.Logger {
	return slog.Default().WithGroup(name)
}

func logx(ctx context.Context, logger *slog.Logger, level slog.Level, msg string, args ...any) {
	if ctx == nil {
		ctx = context.Background()
	}

	if !logger.Enabled(ctx, level) {
		return
	}

	var pcs [1]uintptr
	// skip [runtime.Callers, this function, this function's caller]
	runtime.Callers(3, pcs[:]) // skip [Callers, log, wrapper]

	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.Add(args...)
	_ = logger.Handler().Handle(ctx, r)
}

func logxAttrs(ctx context.Context, logger *slog.Logger, level slog.Level, msg string, attrs ...slog.Attr) {
	if ctx == nil {
		ctx = context.Background()
	}

	if !logger.Enabled(ctx, level) {
		return
	}

	var pcs [1]uintptr
	// skip [runtime.Callers, this function, this function's caller]
	runtime.Callers(3, pcs[:])

	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.AddAttrs(attrs...)
	_ = logger.Handler().Handle(ctx, r)
}
