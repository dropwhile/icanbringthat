package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/errs"
)

func DisableRemindersWithNotification(
	ctx context.Context, db model.PgxHandle,
	email string, suppressionReason string,
) errs.Error {
	user, errx := GetUserByEmail(ctx, db, email)
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
	errx = TxnFunc(ctx, db, func(tx pgx.Tx) error {
		innerErr := UpdateUserSettings(ctx, tx, user.ID, &user.Settings)
		if innerErr != nil {
			return innerErr
		}
		_, innerErr = NewNotification(ctx, tx, user.ID,
			fmt.Sprintf("email notifications disabled due to '%s'", suppressionReason),
		)
		if innerErr != nil {
			return innerErr
		}
		return nil
	})
	return errx
}
