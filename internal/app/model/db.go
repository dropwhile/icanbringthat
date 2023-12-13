package model

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dropwhile/icbt/internal/logger"
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
	logger.Debug(ctx, "db query",
		slog.String("query", query),
		slog.Any("args", args),
	)
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		logger.Error(ctx, "db query error",
			logger.Err(err))
		return *new(T), err
	}
	// note: collectonerow closes rows
	return pgx.CollectOneRow(rows, pgx.RowTo[T])
}

func QueryOne[T any](ctx context.Context, db PgxHandle, query string, args ...interface{}) (*T, error) {
	logger.Debug(ctx, "db query",
		slog.String("query", query),
		slog.Any("args", args),
	)
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		logger.Error(ctx, "db query error",
			logger.Err(err))
		return nil, err
	}
	// note: collectonerow closes rows
	return pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByNameLax[T])
}

func Query[T any](ctx context.Context, db PgxHandle, query string, args ...interface{}) ([]*T, error) {
	logger.Debug(ctx, "db query",
		slog.String("query", query),
		slog.Any("args", args),
	)
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		logger.Error(ctx, "db query error",
			logger.Err(err))
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
			logger.Error(ctx, "inner db tx error",
				logger.Err(errIn))
		}
		return errIn
	})
	if err != nil {
		logger.Error(ctx, "outer db tx error",
			logger.Err(err))
	}
	return t, err
}

func QueryTx[T any](ctx context.Context, db PgxHandle, query string, args ...interface{}) ([]*T, error) {
	var t []*T
	err := pgx.BeginFunc(ctx, db, func(tx pgx.Tx) (errIn error) {
		t, errIn = Query[T](ctx, tx, query, args...)
		if errIn != nil {
			logger.Error(ctx, "inner db tx error",
				logger.Err(errIn))
		}
		return errIn
	})
	if err != nil {
		logger.Error(ctx, "outer db tx error",
			logger.Err(err))
	}
	return t, err
}

func Exec[T any](ctx context.Context, db PgxHandle, query string, args ...interface{}) error {
	logger.Debug(ctx, "db exec query",
		slog.String("query", query),
		slog.Any("args", args),
	)
	commandTag, err := db.Exec(ctx, query, args...)
	if err != nil {
		logger.Error(ctx, "db exec error",
			logger.Err(err))
		return err
	}
	if commandTag.RowsAffected() == 0 {
		logger.Debug(ctx, "query affected zero rows!")
		// return errors.New("no rows affected")
	}
	return nil
}

func ExecTx[T any](ctx context.Context, db PgxHandle, query string, args ...interface{}) error {
	err := pgx.BeginFunc(ctx, db, func(tx pgx.Tx) (errIn error) {
		errIn = Exec[T](ctx, tx, query, args...)
		if errIn != nil {
			logger.Error(ctx, "inner db tx error",
				logger.Err(errIn))
		}
		return errIn
	})
	if err != nil {
		logger.Error(ctx, "outer db tx error",
			logger.Err(err))
	}
	return err
}
