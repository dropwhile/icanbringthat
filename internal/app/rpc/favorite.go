package rpc

import (
	"context"

	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/dropwhile/icbt/internal/app/convert"
	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/service"
	pb "github.com/dropwhile/icbt/rpc"
)

func (s *Server) ListFavoriteEvents(ctx context.Context,
	r *pb.ListFavoriteEventsRequest,
) (*pb.ListFavoriteEventsResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	showArchived := false
	if r.Archived != nil && *r.Archived {
		showArchived = true
	}

	var paginationResult *pb.PaginationResult
	var events []*model.Event
	if r.Pagination != nil {
		limit := int(r.Pagination.Limit)
		offset := int(r.Pagination.Offset)

		favs, pagination, errx := service.GetFavoriteEventsPaginated(
			ctx, s.Db, user.ID, limit, offset, showArchived)
		if errx != nil {
			return nil, convert.ToTwirpError(errx)
		}

		events = favs
		paginationResult = &pb.PaginationResult{
			Limit:  uint32(limit),
			Offset: uint32(offset),
			Count:  uint32(pagination.Count),
		}
	} else {
		favs, errx := service.GetFavoriteEvents(
			ctx, s.Db, user.ID, showArchived)
		if errx != nil {
			return nil, convert.ToTwirpError(errx)
		}
		events = favs
	}

	response := &pb.ListFavoriteEventsResponse{
		Events:     convert.ToPbList(convert.ToPbEvent, events),
		Pagination: paginationResult,
	}
	return response, nil
}

func (s *Server) RemoveFavorite(ctx context.Context,
	r *pb.RemoveFavoriteRequest,
) (*pb.RemoveFavoriteResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	refID, err := model.ParseEventRefID(r.EventRefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "bad notification ref-id")
	}

	errx := service.RemoveFavorite(ctx, s.Db, user.ID, refID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	return &pb.RemoveFavoriteResponse{}, nil
}

func (s *Server) AddFavorite(ctx context.Context,
	r *pb.CreateFavoriteRequest,
) (*pb.CreateFavoriteResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	refID, err := model.ParseEventRefID(r.EventRefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "bad notification ref-id")
	}

	favorite, errx := service.AddFavorite(ctx, s.Db, user.ID, refID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	response := &pb.CreateFavoriteResponse{
		Favorite: &pb.Favorite{
			EventRefId: r.EventRefId,
			Created:    timestamppb.New(favorite.Created),
		},
	}
	return response, nil
}
