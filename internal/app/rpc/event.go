package rpc

import (
	"context"
	"time"

	"github.com/samber/mo"
	"github.com/twitchtv/twirp"

	"github.com/dropwhile/icbt/internal/app/convert"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/internal/middleware/auth"
	"github.com/dropwhile/icbt/rpc/icbt"
)

func (s *Server) ListEvents(ctx context.Context,
	r *icbt.ListEventsRequest,
) (*icbt.ListEventsResponse, error) {
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

		evts, pagination, errx := service.GetEventsPaginated(
			ctx, s.Db, user.ID, limit, offset, showArchived,
		)
		if errx != nil {
			return nil, convert.ToTwirpError(errx)
		}

		events = evts
		paginationResult = &icbt.PaginationResult{
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

	response := &icbt.ListEventsResponse{
		Events:     convert.ToPbList(convert.ToPbEvent, events),
		Pagination: paginationResult,
	}
	return response, nil
}

func (s *Server) CreateEvent(ctx context.Context,
	r *icbt.CreateEventRequest,
) (*icbt.CreateEventResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	if !r.When.Ts.IsValid() {
		return nil, twirp.InvalidArgumentError("start_time", "bad empty value")
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

	response := &icbt.CreateEventResponse{
		Event: convert.ToPbEvent(event),
	}
	return response, nil
}

func (s *Server) UpdateEvent(ctx context.Context,
	r *icbt.UpdateEventRequest,
) (*icbt.UpdateEventResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	refID, err := service.ParseEventRefID(r.RefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "bad event ref-id")
	}

	var startTime *time.Time
	var tz *string

	if r.When != nil {
		if r.When.Ts.IsValid() {
			t := r.When.Ts.AsTime()
			startTime = &t
		}
		if r.When.Tz != "" {
			tz = &r.When.Tz
		}
	}

	euvs := &service.EventUpdateValues{}
	euvs.Name = mo.PointerToOption(r.Name)
	euvs.Description = mo.PointerToOption(r.Description)
	euvs.StartTime = mo.PointerToOption(startTime)
	euvs.Tz = mo.PointerToOption(tz)
	errx := service.UpdateEvent(ctx, s.Db, user.ID, refID, euvs)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	response := &icbt.UpdateEventResponse{}
	return response, nil
}

func (s *Server) GetEventDetails(ctx context.Context,
	r *icbt.GetEventDetailsRequest,
) (*icbt.GetEventDetailsResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	refID, err := service.ParseEventRefID(r.RefId)
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

	response := &icbt.GetEventDetailsResponse{
		Event:    pbEvent,
		Items:    pbEventItems,
		Earmarks: pbEarmarks,
	}
	return response, nil
}

func (s *Server) DeleteEvent(ctx context.Context,
	r *icbt.DeleteEventRequest,
) (*icbt.DeleteEventResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	refID, err := service.ParseEventRefID(r.RefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "bad event ref-id")
	}

	errx := service.DeleteEvent(ctx, s.Db, user.ID, refID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	response := &icbt.DeleteEventResponse{}
	return response, nil
}
