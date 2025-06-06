// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/dropwhile/refid/v2/reftag"
	"github.com/jackc/pgx/v5"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/errs"
	"github.com/dropwhile/icanbringthat/internal/validate"
)

var (
	NotificationRefIDMatcher   = reftag.NewMatcher[model.NotificationRefID]()
	NotificationRefIDFromBytes = reftag.FromBytes[model.NotificationRefID]
	ParseNotificationRefID     = reftag.Parse[model.NotificationRefID]
)

func (s *Service) GetNotificationsCount(
	ctx context.Context, userID int,
) (int, errs.Error) {
	notifCount, err := model.GetNotificationCountByUser(ctx, s.Db, userID)
	if err != nil {
		return 0, errs.Internal.Error("db error")
	}
	return notifCount, nil
}

func (s *Service) GetNotificationsPaginated(
	ctx context.Context, userID int, limit, offset int,
) ([]*model.Notification, *Pagination, errs.Error) {
	notifCount, errx := s.GetNotificationsCount(ctx, userID)
	if errx != nil {
		slog.
			With("error", errx).
			Info("db error")
		return nil, nil, errs.Internal.Error("db error")
	}

	notifications := []*model.Notification{}
	if notifCount > 0 {
		notifs, err := model.GetNotificationsByUserPaginated(
			ctx, s.Db, userID, limit, offset)
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			notifs = []*model.Notification{}
		case err != nil:
			slog.
				With("error", err).
				Info("db error")
			return nil, nil, errs.Internal.Error("db error")
		}
		notifications = notifs
	}
	pagination := &Pagination{
		Limit:  limit,
		Offset: offset,
		Count:  notifCount,
	}
	return notifications, pagination, nil
}

func (s *Service) GetNotifications(
	ctx context.Context, userID int,
) ([]*model.Notification, errs.Error) {
	notifications, err := model.GetNotificationsByUser(ctx, s.Db, userID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return []*model.Notification{}, nil
	case err != nil:
		return nil, errs.Internal.Error("db error")
	}
	return notifications, nil
}

func (s *Service) DeleteNotification(
	ctx context.Context, userID int, refID model.NotificationRefID,
) errs.Error {
	if userID == 0 {
		return errs.Unauthenticated.Error("invalid credentials")
	}

	notification, err := model.GetNotificationByRefID(ctx, s.Db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return errs.NotFound.Error("notification not found")
	case err != nil:
		return errs.Internal.Errorf("db error: %w", err)
	}

	if userID != notification.UserID {
		return errs.PermissionDenied.Error("permission denied")
	}

	err = model.DeleteNotification(ctx, s.Db, notification.ID)
	if err != nil {
		return errs.Internal.Errorf("db error: %w", err)
	}
	return nil
}

func (s *Service) DeleteAllNotifications(
	ctx context.Context, userID int,
) errs.Error {
	if userID == 0 {
		return errs.Unauthenticated.Error("invalid credentials")
	}

	err := model.DeleteNotificationsByUser(ctx, s.Db, userID)
	if err != nil {
		return errs.Internal.Error("db error")
	}

	return nil
}

func (s *Service) NewNotification(
	ctx context.Context, userID int, message string,
) (*model.Notification, errs.Error) {
	return s.newNotification(ctx, s.Db, userID, message)
}

func (s *Service) newNotification(
	ctx context.Context, db model.PgxHandle, userID int, message string,
) (*model.Notification, errs.Error) {
	err := validate.Validate.VarCtx(ctx, message, "required,notblank")
	if err != nil {
		slog.
			With("field", "message").
			With("error", err).
			Info("bad field value")
		return nil, errs.ArgumentError("message", "bad value")
	}

	notification, err := model.NewNotification(ctx, db, userID, message)
	if err != nil {
		slog.
			With("error", err).
			Info("error creating notification")
		return nil, errs.Internal.Errorf("db error: %w", err)
	}
	return notification, nil
}
