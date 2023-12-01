package rpc

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
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

func (s *Server) CreateEvent(ctx context.Context,
	r *pb.CreateEventRequest,
) (*pb.CreateEventResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	name := r.Name
	description := r.Description
	when := r.When.Ts
	tz := r.When.Tz

	if name == "" {
		return nil, twirp.RequiredArgumentError("name")
	}
	if description == "" {
		return nil, twirp.RequiredArgumentError("description")
	}
	if when == nil {
		return nil, twirp.RequiredArgumentError("when")
	}
	if tz == "" {
		return nil, twirp.RequiredArgumentError("tz")
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		return nil, twirp.InvalidArgument.Error("unrecognized tz")
	}

	startTime := when.AsTime()

	event, err := model.NewEvent(ctx, s.Db,
		user.ID, name, description, startTime, &model.TimeZone{Location: loc})
	if err != nil {
		return nil, twirp.InternalError("db error")
	}

	response := &pb.CreateEventResponse{
		Event: dto.ToPbEvent(event),
	}
	return response, nil
}

func (s *Server) UpdateEvent(ctx context.Context,
	r *pb.UpdateEventRequest,
) (*pb.UpdateEventResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	refID, err := model.ParseEventRefID(r.RefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "bad event ref-id")
	}

	event, err := model.GetEventByRefID(ctx, s.Db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, twirp.NotFoundError("event not found")
	case err != nil:
		return nil, twirp.InternalError("db error")
	}

	if user.ID != event.UserID {
		return nil, twirp.PermissionDenied.Error("permission denied")
	}

	if r.Name == nil && r.Description == nil && r.When == nil {
		return nil, twirp.InvalidArgument.Error("missing fields")
	}

	changes := false
	if r.Name != nil && *r.Description != event.Description {
		if *r.Name == "" {
			return nil, twirp.RequiredArgumentError("name")
		}
		event.Name = *r.Name
		changes = true
	}
	if r.Description != nil && *r.Description != event.Description {
		if *r.Description == "" {
			return nil, twirp.RequiredArgumentError("description")
		}
		event.Description = *r.Description
		changes = true
	}
	if r.When != nil {
		loc, err := time.LoadLocation(r.When.Tz)
		if err != nil {
			return nil, twirp.InvalidArgument.Error("unrecognized tz")
		}

		event.StartTime = r.When.Ts.AsTime()
		event.StartTimeTz = &model.TimeZone{Location: loc}
		changes = true
	}

	if !changes {
		return nil, twirp.FailedPrecondition.Error("no changes")
	}

	if err := model.UpdateEvent(
		ctx, s.Db, event.ID,
		event.Name, event.Description, event.ItemSortOrder,
		event.StartTime, event.StartTimeTz,
	); err != nil {
		return nil, twirp.InternalError("db error")
	}

	response := &pb.UpdateEventResponse{
		Event: dto.ToPbEvent(event),
	}
	return response, nil
}

func (s *Server) GetEventDetails(ctx context.Context,
	r *pb.GetEventDetailsRequest,
) (*pb.GetEventDetailsResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	refID, err := model.ParseEventRefID(r.RefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "bad event ref-id")
	}

	event, err := model.GetEventByRefID(ctx, s.Db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, twirp.NotFoundError("event not found")
	case err != nil:
		return nil, twirp.InternalError("db error")
	}

	eventItems, err := model.GetEventItemsByEvent(ctx, s.Db, event.ID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Debug().Err(err).Msg("no rows for event items")
		eventItems = []*model.EventItem{}
	case err != nil:
		return nil, twirp.InternalError("db error")
	}
	pbEventItems := dto.ToPbList(dto.ToPbEventItem, eventItems)

	earmarks, err := model.GetEarmarksByEvent(ctx, s.Db, event.ID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, twirp.NotFoundError("event not found")
	case err != nil:
		return nil, twirp.InternalError("db error")
	}

	pbEarmarks, err := dto.ToPbListWithDb(dto.ToPbEarmark, s.Db, earmarks)
	if err != nil {
		return nil, twirp.InternalError("db error")
	}

	response := &pb.GetEventDetailsResponse{
		Event:    dto.ToPbEvent(event),
		Items:    pbEventItems,
		Earmarks: pbEarmarks,
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
		return nil, twirp.InvalidArgumentError("ref_id", "bad event ref-id")
	}

	event, err := model.GetEventByRefID(ctx, s.Db, refID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, twirp.NotFoundError("event not found")
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
