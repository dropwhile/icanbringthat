package rpc

import (
	"context"

	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/dropwhile/icbt/internal/app/convert"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/internal/middleware/auth"
	"github.com/dropwhile/icbt/rpc/icbt"
)

func (s *Server) ListFavoriteEvents(ctx context.Context,
	r *icbt.ListFavoriteEventsRequest,
) (*icbt.ListFavoriteEventsResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	showArchived := false
	if r.Archived != nil && *r.Archived {
		showArchived = true
	}

	var paginationResult *icbt.PaginationResult
	var events []*model.Event
	if r.Pagination != nil {
		limit := int(r.Pagination.Limit)
		offset := int(r.Pagination.Offset)

		favs, pagination, errx := s.svc.GetFavoriteEventsPaginated(
			ctx, user.ID, limit, offset, showArchived)
		if errx != nil {
			return nil, convert.ToTwirpError(errx)
		}

		events = favs
		paginationResult = &icbt.PaginationResult{
			Limit:  uint32(limit),
			Offset: uint32(offset),
			Count:  uint32(pagination.Count),
		}
	} else {
		favs, errx := s.svc.GetFavoriteEvents(
			ctx, user.ID, showArchived)
		if errx != nil {
			return nil, convert.ToTwirpError(errx)
		}
		events = favs
	}

	response := &icbt.ListFavoriteEventsResponse{
		Events:     convert.ToPbList(convert.ToPbEvent, events),
		Pagination: paginationResult,
	}
	return response, nil
}

func (s *Server) RemoveFavorite(ctx context.Context,
	r *icbt.RemoveFavoriteRequest,
) (*icbt.Empty, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	refID, err := service.ParseEventRefID(r.EventRefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "incorrect value type")
	}

	errx := s.svc.RemoveFavorite(ctx, user.ID, refID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	return &icbt.Empty{}, nil
}

func (s *Server) AddFavorite(ctx context.Context,
	r *icbt.CreateFavoriteRequest,
) (*icbt.CreateFavoriteResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	refID, err := service.ParseEventRefID(r.EventRefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "incorrect value type")
	}

	favorite, errx := s.svc.AddFavorite(ctx, user.ID, refID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	response := &icbt.CreateFavoriteResponse{
		Favorite: &icbt.Favorite{
			EventRefId: r.EventRefId,
			Created:    timestamppb.New(favorite.Created),
		},
	}
	return response, nil
}
