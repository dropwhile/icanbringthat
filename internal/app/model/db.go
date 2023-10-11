package model

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type DB struct {
	*pgxpool.Pool
}

type PgxHandle interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

func (db *DB) GetPool() *pgxpool.Pool {
	return db.Pool
}

func Get[T any](ctx context.Context, db PgxHandle, query string, args ...interface{}) (T, error) {
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		log.Info().Err(err).Msg("db query error")
		return *new(T), fmt.Errorf("failed query: %w", err)
	}
	// note: collectonerow closes rows
	return pgx.CollectOneRow(rows, pgx.RowTo[T])
}

func QueryOne[T any](ctx context.Context, db PgxHandle, query string, args ...interface{}) (*T, error) {
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		log.Info().Err(err).Msg("db query error")
		return nil, fmt.Errorf("failed query: %w", err)
	}
	// note: collectonerow closes rows
	return pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[T])
}

func Query[T any](ctx context.Context, db PgxHandle, query string, args ...interface{}) ([]*T, error) {
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		log.Info().Err(err).Msg("db query error")
		return nil, fmt.Errorf("failed query: %w", err)
	}
	// note: collectrows closes rows
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[T])
}

func QueryOneTx[T any](ctx context.Context, db PgxHandle, query string, args ...interface{}) (*T, error) {
	var t *T
	err := pgx.BeginFunc(ctx, db, func(tx pgx.Tx) error {
		tt, err := QueryOne[T](ctx, tx, query, args...)
		t = tt
		return err
	})
	if err != nil {
		log.Info().Err(err).Msg("DB Error")
	}
	return t, err
}

func QueryTx[T any](ctx context.Context, db PgxHandle, query string, args ...interface{}) ([]*T, error) {
	var t []*T
	err := pgx.BeginFunc(ctx, db, func(tx pgx.Tx) error {
		tt, err := Query[T](ctx, tx, query, args...)
		t = tt
		return err
	})
	if err != nil {
		return nil, err
	}
	return t, nil
}

func Exec[T any](ctx context.Context, db PgxHandle, query string, args ...interface{}) error {
	commandTag, err := db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows affected")
	}
	return nil
}

func ExecTx[T any](ctx context.Context, db PgxHandle, query string, args ...interface{}) error {
	err := pgx.BeginFunc(ctx, db, func(tx pgx.Tx) error {
		commandTag, err := tx.Exec(ctx, query, args...)
		if err != nil {
			log.Info().Err(err).Msg("DB Error")
			return err
		}
		if commandTag.RowsAffected() != 1 {
			return errors.New("no rows affected")
		}
		return nil
	})
	return err
}
