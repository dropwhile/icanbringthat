package rpc

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/rpc/dto"
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

		favCounts, err := model.GetFavoriteCountByUser(ctx, s.Db, user.ID)
		if err != nil {
			return nil, twirp.InternalError("db error")
		}

		count := favCounts.Current
		if showArchived {
			count = favCounts.Archived
		}

		if count > 0 {
			events, err = model.GetFavoriteEventsByUserPaginatedFiltered(
				ctx, s.Db, user.ID, limit, offset, showArchived)
			switch {
			case errors.Is(err, pgx.ErrNoRows):
				events = []*model.Event{}
			case err != nil:
				return nil, twirp.InternalError("db error")
			}
		}
		paginationResult = &pb.PaginationResult{
			Limit:  uint32(limit),
			Offset: uint32(offset),
			Count:  uint32(count),
		}
	} else {
		events, err = model.GetFavoriteEventsByUserFiltered(
			ctx, s.Db, user.ID, showArchived)
		if err != nil {
			return nil, twirp.InternalError("db error")
		}

	}

	response := &pb.ListFavoriteEventsResponse{
		Events:     dto.ToPbList(dto.ToPbEvent, events),
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

	event, err := model.GetEventByRefID(ctx, s.Db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, twirp.NotFoundError("event not found")
	case err != nil:
		return nil, twirp.InternalError("db error")
	}

	favorite, err := model.GetFavoriteByUserEvent(ctx, s.Db, user.ID, event.ID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, twirp.NotFoundError("favorite not found")
	case err != nil:
		return nil, twirp.InternalError("db error")
	}

	// superfluous check, but fine to leave in
	if user.ID != favorite.UserID {
		return nil, twirp.PermissionDenied.Error("permission denied")
	}

	err = model.DeleteFavorite(ctx, s.Db, favorite.ID)
	if err != nil {
		return nil, twirp.InternalError("db error")
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

	event, err := model.GetEventByRefID(ctx, s.Db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, twirp.NotFoundError("event not found")
	case err != nil:
		return nil, twirp.InternalError("db error")
	}

	// check if favorite already exists
	_, err = model.GetFavoriteByUserEvent(ctx, s.Db, user.ID, event.ID)
	if err == nil {
		return nil, twirp.AlreadyExists.Error("favorite already exists")
	}

	favorite, err := model.CreateFavorite(ctx, s.Db, user.ID, event.ID)
	if err != nil {
		return nil, twirp.InternalError("db error")
	}

	response := &pb.CreateFavoriteResponse{
		Favorite: &pb.Favorite{
			EventRefId: r.EventRefId,
			Created:    timestamppb.New(favorite.Created),
		},
	}
	return response, nil
}
