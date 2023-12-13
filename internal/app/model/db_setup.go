package model

import (
	"context"
	"log/slog"
	"runtime"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"

	"github.com/dropwhile/icbt/internal/logger"
)

func SetupDBPool(dbDSN string) (*pgxpool.Pool, error) {
	if !logger.Enabled(logger.LevelTrace) {
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
		// skip [runtime.Callers, this function, this function's caller]
		runtime.Callers(2, pcs[:])

		r := slog.NewRecord(time.Now(), logger.LevelTrace, msg, pcs[0])
		r.AddAttrs(attrs...)
		_ = slog.Default().Handler().Handle(ctx, r)
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
