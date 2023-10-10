package modelx

import (
	"context"

	pgxz "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type PgxHandle interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

func DsnToPool(dbDSN string) (*pgxpool.Pool, error) {
	if zerolog.GlobalLevel() != zerolog.TraceLevel {
		return pgxpool.New(context.Background(), dbDSN)
	}

	config, err := pgxpool.ParseConfig(dbDSN)
	if err != nil {
		return nil, err
	}
	config.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   pgxz.NewLogger(log.Logger),
		LogLevel: tracelog.LogLevelTrace,
	}
	return pgxpool.NewWithConfig(context.Background(), config)
}
