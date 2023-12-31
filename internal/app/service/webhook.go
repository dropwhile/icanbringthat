package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/dropwhile/icbt/internal/errs"
)

func (s *Service) DisableRemindersWithNotification(
	ctx context.Context,
	email string, suppressionReason string,
) errs.Error {
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
