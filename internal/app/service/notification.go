package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/dropwhile/refid/v2/reftag"
	"github.com/jackc/pgx/v5"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/errs"
)

var (
	NotificationRefIDMatcher   = reftag.NewMatcher[model.NotificationRefID]()
	NotificationRefIDFromBytes = reftag.FromBytes[model.NotificationRefID]
	ParseNotificationRefID     = reftag.Parse[model.NotificationRefID]
)

func GetNotifcationsPaginated(
	ctx context.Context, db model.PgxHandle, userID int,
	limit, offset int,
) ([]*model.Notification, *Pagination, errs.Error) {
	notifCount, errx := GetNotificationsCount(ctx, db, userID)
	if errx != nil {
		return nil, nil, errs.Internal.Error("db error")
	}

	notifications := []*model.Notification{}
	if notifCount > 0 {
		notifs, err := model.GetNotificationsByUserPaginated(
			ctx, db, userID, limit, offset)
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			notifs = []*model.Notification{}
		case err != nil:
			return nil, nil, errs.Internal.Error("db error")
		}
		notifications = notifs
	}
	pagination := &Pagination{
		Limit:  uint32(limit),
		Offset: uint32(offset),
		Count:  uint32(notifCount),
	}
	return notifications, pagination, nil
}

func GetNotificationsCount(
	ctx context.Context, db model.PgxHandle, userID int,
) (int, errs.Error) {
	notifCount, err := model.GetNotificationCountByUser(ctx, db, userID)
	if err != nil {
		return 0, errs.Internal.Error("db error")
	}
	return notifCount, nil
}

func GetNotifications(
	ctx context.Context, db model.PgxHandle, userID int,
) ([]*model.Notification, errs.Error) {
	notifications, err := model.GetNotificationsByUser(ctx, db, userID)
	if err != nil {
		return nil, errs.Internal.Error("db error")
	}
	return notifications, nil
}

func DeleteNotification(
	ctx context.Context, db model.PgxHandle, userID int,
	refID model.NotificationRefID,
) errs.Error {
	if userID == 0 {
		return errs.Unauthenticated.Error("invalid credentials")
	}

	notification, err := model.GetNotificationByRefID(ctx, db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return errs.NotFound.
			Error("notification not found").
			Wrap(err)
	case err != nil:
		return errs.Internal.
			Error("db error").
			Wrap(err)
	}

	if userID != notification.UserID {
		return errs.PermissionDenied.Error("permission denied")
	}

	err = model.DeleteNotification(ctx, db, notification.ID)
	if err != nil {
		return errs.Internal.
			Error("db error").
			Wrap(err)
	}
	return nil
}

func DeleteAllNotifications(
	ctx context.Context, db model.PgxHandle, userID int,
) errs.Error {
	if userID == 0 {
		return errs.Unauthenticated.Error("invalid credentials")
	}

	err := model.DeleteNotificationsByUser(ctx, db, userID)
	if err != nil {
		return errs.Internal.Error("db error")
	}

	return nil
}

func NewNotification(
	ctx context.Context, db model.PgxHandle, userID int,
	message string,
) (*model.Notification, errs.Error) {
	err := validate.VarCtx(ctx, message, "required,notblank")
	if err != nil {
		slog.
			With("field", "message").
			With("error", err).
			Info("bad field value")
		return nil, errs.InvalidArgumentError("message", "bad value")
	}

	notification, err := model.NewNotification(ctx, db, userID, message)
	if err != nil {
		return nil, errs.Internal.Errorf("db error: %w", err)
	}
	return notification, nil
}
