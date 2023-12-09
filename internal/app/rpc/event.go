package rpc

import (
	"context"
	"time"

	"github.com/twitchtv/twirp"

	"github.com/dropwhile/icbt/internal/app/convert"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/internal/middleware/auth"
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

		evts, pagination, errx := service.GetEventsPaginated(
			ctx, s.Db, user.ID, limit, offset, showArchived,
		)
		if errx != nil {
			return nil, convert.ToTwirpError(errx)
		}

		events = evts
		paginationResult = &pb.PaginationResult{
			Limit:  uint32(limit),
			Offset: uint32(offset),
			Count:  uint32(pagination.Count),
		}
	} else {
		evts, errx := service.GetEvents(
			ctx, s.Db, user.ID, showArchived,
		)
		if errx != nil {
			return nil, convert.ToTwirpError(errx)
		}
		events = evts
	}

	response := &pb.ListEventsResponse{
		Events:     convert.ToPbList(convert.ToPbEvent, events),
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
	when := r.When.Ts.AsTime()
	tz := r.When.Tz

	event, errx := service.CreateEvent(ctx, s.Db, user,
		name, description, when, tz)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	response := &pb.CreateEventResponse{
		Event: convert.ToPbEvent(event),
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

	var start_time *time.Time
	var tz *string

	if r.When != nil {
		t := r.When.Ts.AsTime()
		start_time = &t
		tz = &r.When.Tz
	}

	event, errx := service.UpdateEvent(ctx, s.Db, user.ID, refID,
		r.Name, r.Description, start_time, tz,
	)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	response := &pb.UpdateEventResponse{
		Event: convert.ToPbEvent(event),
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

	event, errx := service.GetEvent(ctx, s.Db, refID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}
	pbEvent := convert.ToPbEvent(event)

	eventItems, errx := service.GetEventItemsByEventID(ctx, s.Db, event.ID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}
	pbEventItems := convert.ToPbList(convert.ToPbEventItem, eventItems)

	earmarks, errx := service.GetEarmarksByEventID(ctx, s.Db, event.ID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}
	pbEarmarks, err := convert.ToPbListWithDb(convert.ToPbEarmark, s.Db, earmarks)
	if err != nil {
		return nil, twirp.InternalError("db error")
	}

	response := &pb.GetEventDetailsResponse{
		Event:    pbEvent,
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

	errx := service.DeleteEvent(ctx, s.Db, user.ID, refID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	response := &pb.DeleteEventResponse{}
	return response, nil
}
