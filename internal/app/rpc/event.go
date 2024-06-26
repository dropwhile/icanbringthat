// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rpc

import (
	"context"
	"time"

	"github.com/samber/mo"
	"github.com/twitchtv/twirp"

	"github.com/dropwhile/icanbringthat/internal/app/convert"
	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/app/service"
	"github.com/dropwhile/icanbringthat/internal/middleware/auth"
	"github.com/dropwhile/icanbringthat/rpc/icbt"
)

func (s *Server) EventsList(ctx context.Context,
	r *icbt.EventsListRequest,
) (*icbt.EventsListResponse, error) {
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

		evts, pagination, errx := s.svc.GetEventsPaginated(
			ctx, user.ID, limit, offset, showArchived,
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
		evts, errx := s.svc.GetEvents(
			ctx, user.ID, showArchived,
		)
		if errx != nil {
			return nil, convert.ToTwirpError(errx)
		}
		events = evts
	}

	response := &icbt.EventsListResponse{
		Events:     convert.ToPbList(convert.ToPbEvent, events),
		Pagination: paginationResult,
	}
	return response, nil
}

func (s *Server) EventCreate(ctx context.Context,
	r *icbt.EventCreateRequest,
) (*icbt.EventCreateResponse, error) {
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

	event, errx := s.svc.CreateEvent(
		ctx, user, name, description, when, tz,
	)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	response := &icbt.EventCreateResponse{
		Event: convert.ToPbEvent(event),
	}
	return response, nil
}

func (s *Server) EventUpdate(ctx context.Context,
	r *icbt.EventUpdateRequest,
) (*icbt.Empty, error) {
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
	errx := s.svc.UpdateEvent(ctx, user.ID, refID, euvs)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	response := &icbt.Empty{}
	return response, nil
}

func (s *Server) EventGetDetails(ctx context.Context,
	r *icbt.EventGetDetailsRequest,
) (*icbt.EventGetDetailsResponse, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	refID, err := service.ParseEventRefID(r.RefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "bad event ref-id")
	}

	event, errx := s.svc.GetEvent(ctx, refID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}
	pbEvent := convert.ToPbEvent(event)

	eventItems, errx := s.svc.GetEventItemsByEventID(ctx, event.ID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}
	pbEventItems := convert.ToPbList(convert.ToPbEventItem, eventItems)

	earmarks, errx := s.svc.GetEarmarksByEventID(ctx, event.ID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}
	pbEarmarks, err := convert.ToPbListWithService(ctx, convert.ToPbEarmark, s.svc, earmarks)
	if err != nil {
		return nil, twirp.InternalError("db error")
	}

	response := &icbt.EventGetDetailsResponse{
		Event:    pbEvent,
		Items:    pbEventItems,
		Earmarks: pbEarmarks,
	}
	return response, nil
}

func (s *Server) EventDelete(ctx context.Context,
	r *icbt.EventDeleteRequest,
) (*icbt.Empty, error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, twirp.Unauthenticated.Error("invalid credentials")
	}

	refID, err := service.ParseEventRefID(r.RefId)
	if err != nil {
		return nil, twirp.InvalidArgumentError("ref_id", "bad event ref-id")
	}

	errx := s.svc.DeleteEvent(ctx, user.ID, refID)
	if errx != nil {
		return nil, convert.ToTwirpError(errx)
	}

	response := &icbt.Empty{}
	return response, nil
}
