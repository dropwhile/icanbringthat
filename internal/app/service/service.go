package service

import (
	"context"
	"database/sql/driver"
	"log/slog"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"github.com/jackc/pgx/v5"
	"github.com/samber/mo"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/errs"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
	_ = validate.RegisterValidation("notblank", validators.NotBlank)
	validate.RegisterCustomTypeFunc(OptionValuer,
		mo.Option[string]{},
		mo.Option[bool]{},
		mo.Option[[]byte]{},
		mo.Option[[]int]{},
		mo.Option[time.Time]{},
		mo.Option[*model.TimeZone]{},
	)
	// validate.RegisterCustomMatcherFunc(OptionValuer)
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("name"), ",", 2)[0]
		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}
		if name == "" {
			return fld.Name
		}
		return name
	})
}

func OptionValuer(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(driver.Valuer); ok {
		if val, err := valuer.Value(); err == nil {
			return val
		} else {
			slog.
				With("error", err).
				Info("error unwrapping valuer type for validation")
		}
	}
	return nil
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
