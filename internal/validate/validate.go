package validate

import (
	"database/sql/driver"
	"log/slog"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"github.com/samber/mo"

	"github.com/dropwhile/icbt/internal/app/model"
)

// use a single instance of Validate, it caches struct info
var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
	_ = Validate.RegisterValidation("notblank", validators.NotBlank)
	Validate.RegisterCustomTypeFunc(OptionValuer,
		mo.Option[string]{},
		mo.Option[bool]{},
		mo.Option[[]byte]{},
		mo.Option[[]int]{},
		mo.Option[time.Time]{},
		mo.Option[*model.TimeZone]{},
	)
	// validate.RegisterCustomMatcherFunc(OptionValuer)
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
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
		val, err := valuer.Value()
		if err == nil {
			// ref: https://github.com/go-playground/validator/issues/1209
			// return pointer here, so omitnil checks work for validator.
			return &val
		}
		slog.
			With("error", err).
			Info("error unwrapping valuer type for validation")
	}
	return nil
}

func GetErrorField(err error) string {
	if vlderr, ok := err.(validator.ValidationErrors); ok {
		return vlderr[0].Field()
	}
	return ""
}
