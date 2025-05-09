// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package service

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/errs"
)

type FailIfCheckFunc[T any] func(T) bool

type Pagination struct {
	Limit  int
	Offset int
	Count  int
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

//go:generate tool ifacemaker -f "*.go" -s Service -i Servicer -p service -o servicer_iface.go
//go:generate tool mockgen -source servicer_iface.go -destination mockservice/servicer_mock.go -package mockservice
type Service struct {
	Db model.PgxHandle
}

type Options struct {
	Db model.PgxHandle
}

func New(opts Options) *Service {
	return &Service{Db: opts.Db}
}
