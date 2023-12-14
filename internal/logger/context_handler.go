package logger

import (
	"context"
	"log/slog"
	"time"
)

type AttrExtractor func(context.Context, time.Time, slog.Level, string) []slog.Attr

type ContextHandler struct {
	slog.Handler
	Prependers []AttrExtractor
	Appenders  []AttrExtractor
}

// Handle handles the Record.
// It will only be called when Enabled returns true.
// The Context argument is as for Enabled.
// It is present solely to provide Handlers access to the context's values.
// Canceling the context should not affect record processing.
// (Among other things, log messages may be necessary to debug a
// cancellation-related problem.)
//
// Handle methods that produce output should observe the following rules:
//   - If r.Time is the zero time, ignore the time.
//   - If r.PC is zero, ignore it.
//   - Attr's values should be resolved.
//   - If an Attr's key and value are both the zero value, ignore the Attr.
//     This can be tested with attr.Equal(Attr{}).
//   - If a group's key is empty, inline the group's Attrs.
//   - If a group has no Attrs (even if it has a non-empty key),
//     ignore it.
func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	// Add all attributes to new record (because old record has all the old attributes as private members)
	newR := slog.Record{
		Time:    r.Time,
		Level:   r.Level,
		Message: r.Message,
		PC:      r.PC,
	}

	// Collect all attributes from the record (which is the most recent attribute set).
	// These attributes are ordered from oldest to newest, and our collection will be too.
	finalAttrs := make([]slog.Attr, 0, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		finalAttrs = append(finalAttrs, a)
		return true
	})

	for _, f := range h.Prependers {
		attrs := f(ctx, r.Time, r.Level, r.Message)
		if len(attrs) > 0 {
			newR.AddAttrs(attrs...)
		}
	}
	newR.AddAttrs(finalAttrs...)
	for _, f := range h.Appenders {
		attrs := f(ctx, r.Time, r.Level, r.Message)
		if len(attrs) > 0 {
			newR.AddAttrs(attrs...)
		}
	}
	return h.Handler.Handle(ctx, newR)
}

// WithAttrs returns a new Handler whose attributes consist of
// both the receiver's attributes and the arguments.
// The Handler owns the slice: it may retain, modify or discard it.
func (h *ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ContextHandler{h.Handler.WithAttrs(attrs), h.Prependers, h.Appenders}
}

// WithGroup returns a new Handler with the given group appended to
// the receiver's existing groups.
// The keys of all subsequent attributes, whether added by With or in a
// Record, should be qualified by the sequence of group names.
//
// How this qualification happens is up to the Handler, so long as
// this Handler's attribute keys differ from those of another Handler
// with a different sequence of group names.
//
// A Handler should treat WithGroup as starting a Group of Attrs that ends
// at the end of the log event. That is,
//
//	logger.WithGroup("s").LogAttrs(level, msg, slog.Int("a", 1), slog.Int("b", 2))
//
// should behave like
//
//	logger.LogAttrs(level, msg, slog.Group("s", slog.Int("a", 1), slog.Int("b", 2)))
//
// If the name is empty, WithGroup returns the receiver.
func (h *ContextHandler) WithGroup(name string) slog.Handler {
	return &ContextHandler{h.Handler.WithGroup(name), h.Prependers, h.Appenders}
}

func NewContextHandler(next slog.Handler, opts Options) *slog.Logger {
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
