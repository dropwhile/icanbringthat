// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"github.com/dropwhile/icanbringthat/internal/errs"
	"github.com/dropwhile/icanbringthat/internal/validate"
)

func (s *Service) DisableRemindersWithNotification(
	ctx context.Context,
	email string, suppressionReason string,
) errs.Error {
	err := validate.Validate.VarCtx(ctx, email, "required,notblank,email")
	if err != nil {
		slog.
			With("field", "email").
			With("error", err).
			Info("bad field value")
		return errs.InvalidArgumentError("email", "bad value")
	}

	err = validate.Validate.VarCtx(ctx, suppressionReason, "required,notblank")
	if err != nil {
		slog.
			With("field", "suppressionReason").
			With("error", err).
			Info("bad field value")
		return errs.InvalidArgumentError("suppressionReason", "bad value")
	}

	user, errx := s.GetUserByEmail(ctx, email)
	if errx != nil {
		return errx
	}

	// if already disabled, no need to disable again
	if !user.Settings.EnableReminders {
		return errs.FailedPrecondition.Error("reminders already disabled")
	}

	// bounced email, marked spam, unsubscribed...etc
	// so... disable reminders
	user.Settings.EnableReminders = false
	errx = TxnFunc(ctx, s.Db, func(tx pgx.Tx) error {
		innerErr := s.updateUserSettings(ctx, tx, user.ID, &user.Settings)
		if innerErr != nil {
			return innerErr
		}
		_, innerErr = s.newNotification(ctx, tx, user.ID,
			fmt.Sprintf("email notifications disabled due to '%s'", suppressionReason),
		)
		if innerErr != nil {
			return innerErr
		}
		return nil
	})
	return errx
}
