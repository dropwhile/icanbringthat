package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/someerr"
)

type Pagination struct {
	Limit  uint32
	Offset uint32
	Count  uint32
}

func GetNotifcationsPaginated(
	ctx context.Context, db model.PgxHandle, userID int,
	limit, offset int,
) ([]*model.Notification, *Pagination, someerr.Error) {
	notifCount, errx := GetNotificationCount(ctx, db, userID)
	if errx != nil {
		return nil, nil, someerr.Internal.Error("db error")
	}

	notifications := []*model.Notification{}
	if notifCount > 0 {
		notifs, err := model.GetNotificationsByUserPaginated(
			ctx, db, userID, limit, offset)
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			notifs = []*model.Notification{}
		case err != nil:
			return nil, nil, someerr.Internal.Error("db error")
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

func GetNotificationCount(
	ctx context.Context, db model.PgxHandle, userID int,
) (int, someerr.Error) {
	notifCount, err := model.GetNotificationCountByUser(ctx, db, userID)
	if err != nil {
		return 0, someerr.Internal.Error("db error")
	}
	return notifCount, nil
}

func GetNotifications(
	ctx context.Context, db model.PgxHandle, userID int,
) ([]*model.Notification, someerr.Error) {
	notifications, err := model.GetNotificationsByUser(ctx, db, userID)
	if err != nil {
		return nil, someerr.Internal.Error("db error")
	}
	return notifications, nil
}

func DeleteNotification(
	ctx context.Context, db model.PgxHandle, userID int,
	refID model.NotificationRefID,
) someerr.Error {
	if userID == 0 {
		return someerr.Unauthenticated.Error("invalid credentials")
	}

	notification, err := model.GetNotificationByRefID(ctx, db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return someerr.NotFound.
			Error("notification not found").
			Wrap(err)
	case err != nil:
		return someerr.Internal.
			Error("db error").
			Wrap(err)
	}

	if userID != notification.UserID {
		return someerr.PermissionDenied.Error("permission denied")
	}

	err = model.DeleteNotification(ctx, db, notification.ID)
	if err != nil {
		return someerr.Internal.
			Error("db error").
			Wrap(err)
	}
	return nil
}

func DeleteAllNotifications(
	ctx context.Context, db model.PgxHandle, userID int,
) someerr.Error {
	if userID == 0 {
		return someerr.Unauthenticated.Error("invalid credentials")
	}

	err := model.DeleteNotificationsByUser(ctx, db, userID)
	if err != nil {
		return someerr.Internal.Error("db error")
	}

	return nil
}
