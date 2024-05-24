// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rpc

import (
	"context"

	"github.com/twitchtv/twirp"

	"github.com/dropwhile/icanbringthat/internal/app/convert"
	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/app/service"
	"github.com/dropwhile/icanbringthat/internal/middleware/auth"
	"github.com/dropwhile/icanbringthat/rpc/icbt"
)

func (s *Server) NotificationsList(ctx context.Context,
	r *icbt.NotificationsListRequest,
) (*icbt.NotificationsListResponse, error) {
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
		notifs, pagination, errx := s.svc.GetNotificationsPaginated(ctx, user.ID, limit, offset)
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
		notifs, errx := s.svc.GetNotifications(ctx, user.ID)
		if errx != nil {
			return nil, convert.ToTwirpError(errx)
		}
		notifications = notifs

	}

	response := &icbt.NotificationsListResponse{
		Notifications: convert.ToPbList(convert.ToPbNotification, notifications),
		Pagination:    paginationResult,
	}
	return response, nil
}

func (s *Server) NotificationDelete(ctx context.Context,
	r *icbt.NotificationDeleteRequest,
) (*icbt.Empty, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	refID, err := service.ParseNotificationRefID(r.RefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "incorrect value type")
	}

	errx := s.svc.DeleteNotification(ctx, user.ID, refID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	response := &icbt.Empty{}
	return response, nil
}

func (s *Server) NotificationsDeleteAll(ctx context.Context,
	r *icbt.Empty,
) (*icbt.Empty, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	errx := s.svc.DeleteAllNotifications(ctx, user.ID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	response := &icbt.Empty{}
	return response, nil
}
