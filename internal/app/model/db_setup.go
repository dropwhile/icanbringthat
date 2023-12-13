package model

import (
	"context"
	"log/slog"
	"runtime"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
)

func SetupDBPool(dbDSN string, tracing bool) (*pgxpool.Pool, error) {
	if !tracing {
		return pgxpool.New(context.Background(), dbDSN)
	}

	config, err := pgxpool.ParseConfig(dbDSN)
	if err != nil {
		return nil, err
	}

	traceLoggerFunc := func(
		ctx context.Context, level tracelog.LogLevel,
		msg string, data map[string]interface{},
	) {
		if ctx == nil {
			ctx = context.Background()
		}

		attrs := make([]slog.Attr, 0, len(data))
		for k, v := range data {
			attrs = append(attrs, slog.Any(k, v))
		}

		var pcs [1]uintptr
		pcz := make([]uintptr, 12)

		// skip [runtime.Callers, this function, this function's caller]
		// how many frames to get out of pgx??
		// either 8 or 9?
		runtime.Callers(2, pcz[:])

		for ctr, pc := range pcz {
			if pc == 0 {
				continue
			}
			// fn := runtime.FuncForPC(pc)
			// funcName := fn.Name()
			// file, line := fn.FileLine(pc - 1)
			// fmt.Printf("%s:%d %s\n", file, line, funcName)
			file, _ := runtime.FuncForPC(pc).FileLine(pc - 1)
			if strings.HasPrefix(file, "github.com/jackc/pgx") {
				continue
			}
			if strings.HasPrefix(file, "github.com/dropwhile") {
				pcs[0] = pcz[ctr]
				if strings.HasSuffix(file, "model/db.go") {
					continue
				}
				break
			}
		}

		r := slog.NewRecord(time.Now(), slog.LevelDebug, msg, pcs[0])
		r.AddAttrs(attrs...)
		_ = slog.Default().Handler().WithGroup("sql").Handle(ctx, r)
	}

	config.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   tracelog.LoggerFunc(traceLoggerFunc),
		LogLevel: tracelog.LogLevelTrace,
	}
	return pgxpool.NewWithConfig(context.Background(), config)
}

func SetupFromDbPool(pool *pgxpool.Pool) *DB {
	return &DB{pool}
}
