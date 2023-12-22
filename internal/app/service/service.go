package service

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"github.com/jackc/pgx/v5"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/errs"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
	_ = validate.RegisterValidation("notblank", validators.NotBlank)
}

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
