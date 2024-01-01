package service

import (
	"context"
	"time"

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

type Timer interface {
	Time() time.Time
}

func IsTimerExpired(tm Timer, expiry time.Duration) bool {
	return tm.Time().Add(expiry).Before(time.Now())
}

func ParseTimeZone(tz string) (*model.TimeZone, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return nil, err
	}
	return &model.TimeZone{Location: loc}, nil
}

//go:generate ifacemaker -f "*.go" -s Service -i Servicer -p service -o servicer_iface.go
type Service struct {
	Db model.PgxHandle
}

type Options struct {
	Db model.PgxHandle
}

func New(opts Options) *Service {
	return &Service{Db: opts.Db}
}
