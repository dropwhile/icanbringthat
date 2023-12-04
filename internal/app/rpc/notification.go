package rpc

import (
	"context"

	"github.com/twitchtv/twirp"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/rpc/dto"
	"github.com/dropwhile/icbt/internal/app/service"
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
		notifs, pagination, errx := service.GetNotifcationsPaginated(ctx, s.Db, user.ID, limit, offset)
		if errx != nil {
			return nil, dto.ToTwirpError(errx)
		}

		notifications = notifs
		paginationResult = &pb.PaginationResult{
			Limit:  uint32(limit),
			Offset: uint32(offset),
			Count:  pagination.Count,
		}
	} else {
		notifs, errx := service.GetNotifications(ctx, s.Db, user.ID)
		if err != nil {
			return nil, dto.ToTwirpError(errx)
		}
		notifications = notifs

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

	errx := service.DeleteNotification(ctx, s.Db, user.ID, refID)
	if errx != nil {
		return nil, dto.ToTwirpError(errx)
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

	errx := service.DeleteAllNotifications(ctx, s.Db, user.ID)
	if errx != nil {
		return nil, dto.ToTwirpError(errx)
	}

	response := &pb.DeleteAllNotificationsResponse{}
	return response, nil
}
