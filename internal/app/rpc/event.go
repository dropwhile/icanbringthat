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

func (s *Server) ListEvents(ctx context.Context,
	r *pb.ListEventsRequest,
) (*pb.ListEventsResponse, error) {
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

		eventCounts, err := model.GetEventCountsByUser(ctx, s.Db, user.ID)
		if err != nil {
			return nil, twirp.InternalError("db error")
		}

		count := eventCounts.Current
		if showArchived {
			count = eventCounts.Archived
		}

		if count > 0 {
			events, err = model.GetEventsByUserPaginatedFiltered(
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
		events, err = model.GetEventsByUserFiltered(
			ctx, s.Db, user.ID, showArchived)
		if err != nil {
			return nil, twirp.InternalError("db error")
		}

	}

	response := &pb.ListEventsResponse{
		Events:     dto.ToPbList(dto.ToPbEvent, events),
		Pagination: paginationResult,
	}
	return response, nil
}

func (s *Server) DeleteEvent(ctx context.Context,
	r *pb.DeleteEventRequest,
) (*pb.DeleteEventResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	refID, err := model.ParseEventRefID(r.RefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "bad notification ref-id")
	}

	event, err := model.GetEventByRefID(ctx, s.Db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, twirp.NotFoundError("notification not found")
	case err != nil:
		return nil, twirp.InternalError("db error")
	}

	if user.ID != event.UserID {
		return nil, twirp.PermissionDenied.Error("permission denied")
	}

	err = model.DeleteEvent(ctx, s.Db, event.ID)
	if err != nil {
		return nil, twirp.InternalError("db error")
	}

	response := &pb.DeleteEventResponse{}
	return response, nil
}
