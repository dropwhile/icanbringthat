package rpc

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/twitchtv/twirp"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/rpc/dto/converter"
	pb "github.com/dropwhile/icbt/rpc"
)

var dtoConverter = converter.DTOConverter{}

func (s *Server) ListNotifications(ctx context.Context,
	request *pb.ListNotificationsRequest,
) (*pb.ListNotificationsResponse, error) {
	// get user from auth in context
	user, err := model.GetUserByID(ctx, s.Db, 1)
	if err != nil || user == nil {
		log.Debug().Err(err).Msg("invalid credentials: no user match")
		// reutrn nil, twirp.a Unauthenticated
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	limit := 20
	offset := 0
	if request.Pagination != nil {
		limit = int(request.Pagination.Limit)
		offset = int(request.Pagination.Offset)
	}

	notifCount, err := model.GetNotificationCountByUser(ctx, s.Db, 1)
	if err != nil {
		return nil, twirp.InternalError("db error")
	}

	notifications := make([]*model.Notification, 0)
	if notifCount > 0 {
		notifications, err = model.GetNotificationsByUserPaginated(ctx, s.Db, 1, limit, offset)
		if err != nil {
			return nil, twirp.InternalError("db error")
		}
	}

	dtoNotificationes := dtoConverter.ConvertNotifications(notifications)
	response := &pb.ListNotificationsResponse{
		Notifications: dtoNotificationes,
		Pagination: &pb.PaginationResult{
			Limit:  uint32(limit),
			Offset: uint32(offset),
			Count:  uint32(notifCount),
		},
	}
	return response, nil
}
