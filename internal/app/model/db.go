package model

import (
	"context"
	"errors"

	"github.com/cactus/mlog"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	*pgxpool.Pool
}

type pgxIface interface {
	Begin(context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Ping(context.Context) error
	Prepare(context.Context, string, string) (*pgconn.StatementDescription, error)
	PgConn() *pgconn.PgConn
}

type PgxConnIface interface {
	pgxIface
	pgx.Tx
	Close(ctx context.Context) error
	Deallocate(ctx context.Context, name string) error
}

type PgxPoolIface interface {
	pgxIface
	pgx.Tx
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
	AcquireAllIdle(ctx context.Context) []*pgxpool.Conn
	AcquireFunc(ctx context.Context, f func(*pgxpool.Conn) error) error
	AsConn() PgxConnIface
	Close()
	Stat() *pgxpool.Stat
	Reset()
}

type PgxHandle interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

func (db *DB) GetPool() *pgxpool.Pool {
	return db.Pool
}

func SetupFromDb(pool *pgxpool.Pool) *DB {
	return &DB{pool}
}

func QueryOne[T ModelType](ctx context.Context, db PgxHandle, query string, args ...interface{}) (*T, error) {
	var t T
	err := pgxscan.Get(ctx, db, &t, query, args...)
	return &t, err
}

func Query[T ModelType](ctx context.Context, db PgxHandle, query string, args ...interface{}) ([]*T, error) {
	var t []*T
	err := pgxscan.Select(ctx, db, &t, query, args...)
	return t, err
}

func Exec[T ModelType](ctx context.Context, db PgxHandle, query string, args ...interface{}) error {
	commandTag, err := db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows affected")
	}
	return nil
}

func QueryOneTx[T ModelType](ctx context.Context, db PgxHandle, query string, args ...interface{}) (*T, error) {
	var t T
	err := pgx.BeginFunc(context.Background(), db, func(tx pgx.Tx) error {
		err := pgxscan.Get(ctx, tx, &t, query, args...)
		if err != nil {
			mlog.Printx("DB Error", mlog.A("err", err))
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func QueryTx[T ModelType](ctx context.Context, db PgxHandle, query string, args ...interface{}) ([]*T, error) {
	var t []*T
	err := pgx.BeginFunc(context.Background(), db, func(tx pgx.Tx) error {
		err := pgxscan.Select(ctx, tx, &t, query, args...)
		if err != nil {
			mlog.Printx("DB Error", mlog.A("err", err))
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return t, nil
}

func ExecTx[T ModelType](ctx context.Context, db PgxHandle, query string, args ...interface{}) error {
	err := pgx.BeginFunc(context.Background(), db, func(tx pgx.Tx) error {
		commandTag, err := tx.Exec(ctx, query, args...)
		if err != nil {
			mlog.Printx("DB Error", mlog.A("err", err))
			return err
		}
		if commandTag.RowsAffected() != 1 {
			return errors.New("no rows affected")
		}
		return nil
	})
	return err
}
