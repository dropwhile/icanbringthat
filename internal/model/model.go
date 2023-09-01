package model

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var _ sql.NullString

type DB struct {
	pool *sqlx.DB
}

func (db *DB) Close() error {
	return db.pool.Close()
}

func NewDatabase(dbDSN string) (*DB, error) {
	dbConn, err := sqlx.Connect("postgres", dbDSN)
	if err != nil {
		return nil, err
	}
	return &DB{pool: dbConn}, nil
}
