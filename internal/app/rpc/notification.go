package rpc

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/twitchtv/twirp"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/rpc/dto"
	pb "github.com/dropwhile/icbt/rpc"
)

func (s *Server) ListNotifications(ctx context.Context,
	r *pb.ListNotificationsRequest,
) (*pb.ListNotificationsResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	var paginationResult *pb.PaginationResult
	var notifications []*model.Notification
	if r.Pagination != nil {
		limit := int(r.Pagination.Limit)
		offset := int(r.Pagination.Offset)

		notifCount, err := model.GetNotificationCountByUser(ctx, s.Db, user.ID)
		if err != nil {
			return nil, twirp.InternalError("db error")
		}

		if notifCount > 0 {
			notifications, err = model.GetNotificationsByUserPaginated(ctx, s.Db, user.ID, limit, offset)
			switch {
			case errors.Is(err, pgx.ErrNoRows):
				notifications = []*model.Notification{}
			case err != nil:
				return nil, twirp.InternalError("db error")
			}
		}
		paginationResult = &pb.PaginationResult{
			Limit:  uint32(limit),
			Offset: uint32(offset),
			Count:  uint32(notifCount),
		}
	} else {
		notifications, err = model.GetNotificationsByUser(ctx, s.Db, user.ID)
		if err != nil {
			return nil, twirp.InternalError("db error")
		}

	}

	response := &pb.ListNotificationsResponse{
		Notifications: dto.ToPbList(dto.ToPbNotification, notifications),
		Pagination:    paginationResult,
	}
	return response, nil
}

func (s *Server) DeleteNotification(ctx context.Context,
	r *pb.DeleteNotificationRequest,
) (*pb.DeleteNotificationResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	refID, err := model.ParseNotificationRefID(r.RefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "bad notification ref-id")
	}

	notification, err := model.GetNotificationByRefID(ctx, s.Db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, twirp.NotFoundError("notification not found")
	case err != nil:
		return nil, twirp.InternalError("db error")
	}

	if user.ID != notification.UserID {
		return nil, twirp.PermissionDenied.Error("permission denied")
	}

	err = model.DeleteNotification(ctx, s.Db, notification.ID)
	if err != nil {
		return nil, twirp.InternalError("db error")
	}

	response := &pb.DeleteNotificationResponse{}
	return response, nil
}

func (s *Server) DeleteAllNotifications(ctx context.Context,
	r *pb.DeleteAllNotificationsRequest,
) (*pb.DeleteAllNotificationsResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	err = model.DeleteNotificationsByUser(ctx, s.Db, user.ID)
	if err != nil {
		return nil, twirp.InternalError("db error")
	}

	response := &pb.DeleteAllNotificationsResponse{}
	return response, nil
}
