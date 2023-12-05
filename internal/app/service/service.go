package service

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/errs"
)

type Pagination struct {
	Limit  uint32
	Offset uint32
	Count  uint32
}

func TxnFunc(ctx context.Context, db model.PgxHandle,
	dbfn func(pgx.Tx) error,
) errs.Error {
	err := pgx.BeginFunc(ctx, db, dbfn)
	if err != nil {
		return errs.Internal.Errorf("db error: %w", err)
	}
	return nil
}
