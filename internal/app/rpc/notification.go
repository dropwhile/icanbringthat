package rpc

import (
	"context"

	"github.com/twitchtv/twirp"

	"github.com/dropwhile/icbt/internal/app/convert"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/internal/middleware/auth"
	"github.com/dropwhile/icbt/rpc/icbt"
)

func (s *Server) ListNotifications(ctx context.Context,
	r *icbt.ListNotificationsRequest,
) (*icbt.ListNotificationsResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	var paginationResult *icbt.PaginationResult
	var notifications []*model.Notification
	if r.Pagination != nil {
		limit := int(r.Pagination.Limit)
		offset := int(r.Pagination.Offset)
		notifs, pagination, errx := service.GetNotifcationsPaginated(ctx, s.Db, user.ID, limit, offset)
		if errx != nil {
			return nil, convert.ToTwirpError(errx)
		}

		notifications = notifs
		paginationResult = &icbt.PaginationResult{
			Limit:  uint32(limit),
			Offset: uint32(offset),
			Count:  pagination.Count,
		}
	} else {
		notifs, errx := service.GetNotifications(ctx, s.Db, user.ID)
		if errx != nil {
			return nil, convert.ToTwirpError(errx)
		}
		notifications = notifs

	}

	response := &icbt.ListNotificationsResponse{
		Notifications: convert.ToPbList(convert.ToPbNotification, notifications),
		Pagination:    paginationResult,
	}
	return response, nil
}

func (s *Server) DeleteNotification(ctx context.Context,
	r *icbt.DeleteNotificationRequest,
) (*icbt.DeleteNotificationResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	refID, err := service.ParseNotificationRefID(r.RefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "incorrect value type")
	}

	errx := service.DeleteNotification(ctx, s.Db, user.ID, refID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	response := &icbt.DeleteNotificationResponse{}
	return response, nil
}

func (s *Server) DeleteAllNotifications(ctx context.Context,
	r *icbt.DeleteAllNotificationsRequest,
) (*icbt.DeleteAllNotificationsResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	errx := service.DeleteAllNotifications(ctx, s.Db, user.ID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	response := &icbt.DeleteAllNotificationsResponse{}
	return response, nil
}
