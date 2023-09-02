package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cactus/mlog"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
}

/*
func (db *DB) Close() error {
	return db.sqlx.Close()
}
*/

func SetupFromDb(sqlDb *sql.DB) (*DB, error) {
	/*
		dbConn, err := sqlx.Connect("pgx", dbDSN)
		if err != nil {
			return nil, err
		}
	*/
	sqlxDb := sqlx.NewDb(sqlDb, "pgx")
	db := &DB{sqlxDb}
	return db, nil
}

func withTx[T any](db *DB, ctx context.Context, txf func(*sqlx.Tx) (T, error)) (T, error) {
	var result T

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return result, err
	}

	result, err = txf(tx)
	if err != nil {
		if tErr := tx.Rollback(); tErr != nil && !errors.Is(tErr, sql.ErrTxDone) {
			return result, errors.Join(tErr, err)
		}
		return result, err
	}

	err = tx.Commit()
	return result, err
}

func QueryRow[T ModelType](qc sqlx.QueryerContext, ctx context.Context, query string, args ...interface{}) (*T, error) {
	if mlog.HasDebug() {
		mlog.Debugx("SQL", mlog.A("query", query), mlog.A("args", args))
	}
	var t T
	err := qc.QueryRowxContext(ctx, query, args...).StructScan(&t)
	return &t, err
}

func Query[T ModelType](qc sqlx.QueryerContext, ctx context.Context, query string, args ...interface{}) ([]*T, error) {
	if mlog.HasDebug() {
		mlog.Debugx("SQL", mlog.A("query", query), mlog.A("args", args))
	}
	result := make([]*T, 0)
	rows, err := qc.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var t T
		err = rows.StructScan(&t)
		if err != nil {
			return nil, err
		}
		result = append(result, &t)
	}

	return result, nil
}

func Exec[T ModelType](ec sqlx.ExecerContext, ctx context.Context, query string, args ...interface{}) error {
	if mlog.HasDebug() {
		mlog.Debugx("SQL", mlog.A("query", query), mlog.A("args", args))
	}
	res, err := ec.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if mlog.HasDebug() {
		mlog.Debugf("sql exec result affected rows: %d", n)
	}

	return err
}

func QueryRowTx[T ModelType](db *DB, ctx context.Context, query string, args ...interface{}) (*T, error) {
	return withTx(db, ctx, func(tx *sqlx.Tx) (*T, error) {
		fmt.Println(query, args)
		return QueryRow[T](tx, ctx, query, args...)
	})
}

func QueryTx[T ModelType](db *DB, ctx context.Context, query string, args ...interface{}) ([]*T, error) {
	return withTx(db, ctx, func(tx *sqlx.Tx) ([]*T, error) {
		return Query[T](tx, ctx, query, args...)
	})
}

func ExecTx[T ModelType](db *DB, ctx context.Context, query string, args ...interface{}) error {
	_, err := withTx(db, ctx, func(tx *sqlx.Tx) (bool, error) {
		err := Exec[T](tx, ctx, query, args...)
		if err != nil {
			return false, err
		}
		return true, nil
	})
	return err
}
