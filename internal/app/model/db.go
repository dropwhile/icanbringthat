package model

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/dropwhile/icbt/internal/logger"
)

type PgxHandle interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

func Get[T any](ctx context.Context, db PgxHandle, query string, args ...interface{}) (T, error) {
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		logger.LogSkip(slog.Default(), 1, slog.LevelError,
			ctx, "db query error", logger.Err(err))
		return *new(T), err
	}
	// note: collectonerow closes rows
	return pgx.CollectOneRow(rows, pgx.RowTo[T])
}

func QueryOne[T any](ctx context.Context, db PgxHandle, query string, args ...interface{}) (*T, error) {
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		logger.LogSkip(slog.Default(), 1, slog.LevelError,
			ctx, "db query error", logger.Err(err))
		return nil, err
	}
	// note: collectonerow closes rows
	return pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[T])
}

func Query[T any](ctx context.Context, db PgxHandle, query string, args ...interface{}) ([]*T, error) {
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		logger.LogSkip(slog.Default(), 1, slog.LevelError,
			ctx, "db query error", logger.Err(err))
		return nil, err
	}
	// note: collectrows closes rows
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[T])
}

func QueryOneTx[T any](ctx context.Context, db PgxHandle, query string, args ...interface{}) (*T, error) {
	var t *T
	err := pgx.BeginFunc(ctx, db, func(tx pgx.Tx) (errIn error) {
		t, errIn = QueryOne[T](ctx, tx, query, args...)
		if errIn != nil {
			logger.LogSkip(slog.Default(), 1, slog.LevelError,
				ctx, "inner db tx error", logger.Err(errIn))
		}
		return errIn
	})
	if err != nil {
		logger.LogSkip(slog.Default(), 1, slog.LevelError,
			ctx, "outer db tx error",
			logger.Err(err))
	}
	return t, err
}

func QueryTx[T any](ctx context.Context, db PgxHandle, query string, args ...interface{}) ([]*T, error) {
	var t []*T
	err := pgx.BeginFunc(ctx, db, func(tx pgx.Tx) (errIn error) {
		t, errIn = Query[T](ctx, tx, query, args...)
		if errIn != nil {
			logger.LogSkip(slog.Default(), 1, slog.LevelError,
				ctx, "inner db tx error", logger.Err(errIn))
		}
		return errIn
	})
	if err != nil {
		logger.LogSkip(slog.Default(), 1, slog.LevelError,
			ctx, "outer db tx error", logger.Err(err))
	}
	return t, err
}

func Exec[T any](ctx context.Context, db PgxHandle, query string, args ...interface{}) error {
	commandTag, err := db.Exec(ctx, query, args...)
	if err != nil {
		logger.LogSkip(slog.Default(), 1, slog.LevelError,
			ctx, "db exec error", logger.Err(err))
		return err
	}
	if commandTag.RowsAffected() == 0 {
		logger.LogSkip(slog.Default(), 1, slog.LevelDebug,
			ctx, "query affected zero rows!")
		// return errors.New("no rows affected")
	}
	return nil
}

func ExecTx[T any](ctx context.Context, db PgxHandle, query string, args ...interface{}) error {
	err := pgx.BeginFunc(ctx, db, func(tx pgx.Tx) (errIn error) {
		errIn = Exec[T](ctx, tx, query, args...)
		if errIn != nil {
			logger.LogSkip(slog.Default(), 1, slog.LevelError,
				ctx, "inner db tx error", logger.Err(errIn))
		}
		return errIn
	})
	if err != nil {
		logger.LogSkip(slog.Default(), 1, slog.LevelError,
			ctx, "outer db tx error", logger.Err(err))
	}
	return err
}
