package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/pkg/errors"
	slogcontext "github.com/veqryn/slog-context"
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

// for text unmarshaller
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

func (l LoggerLevel) LogValuer() slog.Value {
	return slog.StringValue(l.String())
}

// fmtErr returns a slog.GroupValue with keys "msg" and "trace". If the error
// does not implement interface { StackTrace() errors.StackTrace }, the "trace"
// key is omitted.
func fmtErr(err error) slog.Value {
	var groupValues []slog.Attr

	groupValues = append(groupValues, slog.String("msg", err.Error()))

	type StackTracer interface {
		StackTrace() errors.StackTrace
	}

	// Find the trace to the location of the first errors.New,
	// errors.Wrap, or errors.WithStack call.
	var st StackTracer
	for err := err; err != nil; err = errors.Unwrap(err) {
		if x, ok := err.(StackTracer); ok {
			st = x
		}
	}

	if st != nil {
		groupValues = append(groupValues,
			slog.Any("trace", traceLines(st.StackTrace())),
		)
	}

	return slog.GroupValue(groupValues...)
}

func traceLines(frames errors.StackTrace) []string {
	traceLines := make([]string, len(frames))

	// Iterate in reverse to skip uninteresting, consecutive runtime frames at
	// the bottom of the trace.
	var skipped int
	skipping := true
	for i := len(frames) - 1; i >= 0; i-- {
		// Adapted from errors.Frame.MarshalText(), but avoiding repeated
		// calls to FuncForPC and FileLine.
		pc := uintptr(frames[i]) - 1
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			traceLines[i] = "unknown"
			skipping = false
			continue
		}

		name := fn.Name()

		if skipping && strings.HasPrefix(name, "runtime.") {
			skipped++
			continue
		} else {
			skipping = false
		}

		filename, lineNr := fn.FileLine(pc)

		traceLines[i] = fmt.Sprintf("%s %s:%d", name, filename, lineNr)
	}

	return traceLines[:len(traceLines)-skipped]
}

var re = regexp.MustCompile(`^v[0-9]+@`)

func trimFilePath(file string) string {
	short := file
	counter := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			if counter > 0 {
				short = file[i+1:]
				if re.MatchString(short) {
					continue
				}
				break
			}
			counter += 1
		}
	}

	// prune from after @ to next /
	atIdx := strings.Index(short, "@")
	if atIdx >= 0 && atIdx+7 <= len(short) {
		for i := atIdx; i < len(short); i++ {
			if short[i] == '/' {
				if i-atIdx > 7 {
					short = short[:atIdx+7] + "..." + short[i:]
				}
				break
			}
		}
	}

	return short
}

func replaceAttr(opts Options) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		// Remove time if in test mode
		if opts.OmitTime && a.Key == slog.TimeKey && len(groups) == 0 {
			return slog.Attr{}
		}

		switch a.Key {
		case slog.TimeKey:
			a.Key = "ts"
			if !opts.UseLocalTime {
				a.Value = slog.TimeValue(a.Value.Time().UTC())
			}
		// Remove the directory from the source's filename.
		case slog.SourceKey:
			a.Key = "src"
			source := a.Value.Any().(*slog.Source)
			source.File = trimFilePath(source.File)
		// Customize the name of the level key and the output string, including
		// custom level values.
		case slog.LevelKey:
			if v, ok := a.Value.Any().(slog.Level); ok {
				a.Key = "lvl"
				switch v {
				case LevelTrace:
					a.Value = LevelTraceStr
				case LevelDebug:
					a.Value = LevelDebugStr
				case LevelError:
					a.Value = LevelErrorStr
				case LevelFatal:
					a.Value = LevelFatalStr
				default:
					a.Value = LevelInfoStr
				}
			}
		}

		switch a.Value.Kind() {
		// other cases
		case slog.KindAny:
			switch v := a.Value.Any().(type) {
			case error:
				a.Value = fmtErr(v)
			}
		}

		return a
	}
}

type AttrExtractor = slogcontext.AttrExtractor

type Options struct {
	UseLocalTime bool
	OmitTime     bool
	OmitSource   bool
	Sink         io.Writer
	Prependers   []AttrExtractor
	Appenders    []AttrExtractor
}

func newContextHandler(next slog.Handler, opts Options) *slog.Logger {
	prependers := []AttrExtractor{
		slogcontext.ExtractPrepended,
	}
	prependers = append(prependers, opts.Prependers...)

	appenders := []AttrExtractor{
		slogcontext.ExtractAppended,
	}
	appenders = append(appenders, opts.Appenders...)

	h := slogcontext.NewHandler(
		next,
		&slogcontext.HandlerOptions{
			Prependers: prependers,
			Appenders:  appenders,
		},
	)
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
	case "trace":
		logLevel.Set(LevelTrace)
	case "debug":
		logLevel.Set(LevelDebug)
	default:
		logLevel.Set(LevelInfo)
	}

	return newContextHandler(logHandler, opts)
}

type logWriter struct {
	out io.Writer
}

func (lw *logWriter) Write(b []byte) (int, error) {
	_, err := io.WriteString(lw.out, time.Now().UTC().Format("2006-01-02T15:04:05.999Z")+" ")
	if err != nil {
		return 0, err
	}
	return lw.out.Write(b)
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

func SetLevel[T ~int](level T) {
	logLevel.Set(slog.Level(level))
}

func PrependAttr(ctx context.Context, args ...any) context.Context {
	return slogcontext.Prepend(ctx, args...)
}

func Err(err error) slog.Attr {
	return slog.Any("error", err)
}

func Trace(ctx context.Context, msg string, attrs ...slog.Attr) {
	logxAttrs(ctx, slog.Default(), LevelTrace, msg, attrs...)
}

func Debug(ctx context.Context, msg string, attrs ...slog.Attr) {
	logxAttrs(ctx, slog.Default(), LevelDebug, msg, attrs...)
}

func Info(ctx context.Context, msg string, attrs ...slog.Attr) {
	logxAttrs(ctx, slog.Default(), LevelInfo, msg, attrs...)
}

func Error(ctx context.Context, msg string, attrs ...slog.Attr) {
	logxAttrs(ctx, slog.Default(), LevelError, msg, attrs...)
}

func Fatal(ctx context.Context, msg string, attrs ...slog.Attr) {
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

func Enabled[T ~int](level T) bool {
	return slog.Default().Enabled(context.Background(), slog.Level(level))
}
