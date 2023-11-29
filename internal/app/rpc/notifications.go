package rpc

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/twitchtv/twirp"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/rpc/dto"
	pb "github.com/dropwhile/icbt/rpc"
)

func (s *Server) ListNotifications(ctx context.Context,
	r *pb.ListNotificationsRequest,
) (*pb.ListNotificationsResponse, error) {
	// get user from auth in context
	user, err := model.GetUserByID(ctx, s.Db, 1)
	if err != nil || user == nil {
		log.Debug().Err(err).Msg("invalid credentials: no user match")
		// reutrn nil, twirp.a Unauthenticated
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	var paginationResult *pb.PaginationResult
	var notifications []*model.Notification
	if r.Pagination != nil {
		limit := int(r.Pagination.Limit)
		offset := int(r.Pagination.Offset)

		notifCount, err := model.GetNotificationCountByUser(ctx, s.Db, 1)
		if err != nil {
			return nil, twirp.InternalError("db error")
		}

		if notifCount > 0 {
			notifications, err = model.GetNotificationsByUserPaginated(ctx, s.Db, 1, limit, offset)
			if err != nil {
				return nil, twirp.InternalError("db error")
			}
		}
		paginationResult = &pb.PaginationResult{
			Limit:  uint32(limit),
			Offset: uint32(offset),
			Count:  uint32(notifCount),
		}
	} else {
		notifications, err = model.GetNotificationsByUser(ctx, s.Db, 1)
		if err != nil {
			return nil, twirp.InternalError("db error")
		}

	}

	response := &pb.ListNotificationsResponse{
		Notifications: dto.ToPbNotificationList(notifications),
		Pagination:    paginationResult,
	}
	return response, nil
}
