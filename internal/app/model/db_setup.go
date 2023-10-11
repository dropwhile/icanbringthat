package model

import (
	"context"

	pgxz "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func SetupDBPool(dbDSN string) (*pgxpool.Pool, error) {
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

func SetupFromDbPool(pool *pgxpool.Pool) *DB {
	return &DB{pool}
}
