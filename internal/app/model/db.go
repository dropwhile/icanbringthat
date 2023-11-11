package model

import (
	"context"

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
		return *new(T), err
	}
	// note: collectonerow closes rows
	return pgx.CollectOneRow(rows, pgx.RowTo[T])
}

func QueryOne[T any](ctx context.Context, db PgxHandle, query string, args ...interface{}) (*T, error) {
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		log.Info().Err(err).Msg("db query error")
		return nil, err
	}
	// note: collectonerow closes rows
	return pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[T])
}

func Query[T any](ctx context.Context, db PgxHandle, query string, args ...interface{}) ([]*T, error) {
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		log.Info().Err(err).Msg("db query error")
		return nil, err
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
		log.Info().Err(err).Msg("db tx error")
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
		log.Info().Err(err).Msg("db tx error")
	}
	return t, err
}

func Exec[T any](ctx context.Context, db PgxHandle, query string, args ...interface{}) error {
	commandTag, err := db.Exec(ctx, query, args...)
	if err != nil {
		log.Info().Err(err).Msg("db exec error")
		return err
	}
	if commandTag.RowsAffected() == 0 {
		log.Debug().Msg("query affected zero rows!")
		// return errors.New("no rows affected")
	}
	return nil
}

func ExecTx[T any](ctx context.Context, db PgxHandle, query string, args ...interface{}) error {
	err := pgx.BeginFunc(ctx, db, func(tx pgx.Tx) error {
		return Exec[T](ctx, tx, query, args...)
	})
	if err != nil {
		log.Info().Err(err).Msg("db tx error")
	}
	return err
}
