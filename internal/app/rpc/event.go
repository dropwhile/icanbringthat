// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rpc

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/samber/mo"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/dropwhile/icanbringthat/internal/app/convert"
	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/app/service"
	"github.com/dropwhile/icanbringthat/internal/middleware/auth"
	icbt "github.com/dropwhile/icanbringthat/rpc/icbt/rpc/v1"
)

func (s *Server) EventsList(ctx context.Context,
	req *connect.Request[icbt.EventsListRequest],
) (*connect.Response[icbt.EventsListResponse], error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid credentials"))
	}

	showArchived := false
	if req.Msg.HasArchived() && req.Msg.GetArchived() {
		showArchived = true
	}

	var paginationResult *icbt.PaginationResult
	var events []*model.Event
	if req.Msg.HasPagination() {
		limit := int(req.Msg.GetPagination().GetLimit())
		offset := int(req.Msg.GetPagination().GetOffset())

		evts, pagination, errx := s.svc.GetEventsPaginated(
			ctx, user.ID, limit, offset, showArchived,
		)
		if errx != nil {
			return nil, convert.ToConnectRpcError(errx)
		}

		events = evts
		paginationResult = convert.ToPbPagination(pagination)
	} else {
		evts, errx := s.svc.GetEvents(
			ctx, user.ID, showArchived,
		)
		if errx != nil {
			return nil, convert.ToConnectRpcError(errx)
		}
		events = evts
	}

	response := icbt.EventsListResponse_builder{
		Events:     convert.ToPbList(convert.ToPbEvent, events),
		Pagination: paginationResult,
	}.Build()
	return connect.NewResponse(response), nil
}

func (s *Server) EventCreate(ctx context.Context,
	req *connect.Request[icbt.EventCreateRequest],
) (*connect.Response[icbt.EventCreateResponse], error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid credentials"))
	}

	if !req.Msg.GetWhen().GetTs().IsValid() {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("start_time bad empty value"))
	}

	name := req.Msg.GetName()
	description := req.Msg.GetDescription()
	when := req.Msg.GetWhen().GetTs().AsTime()
	tz := req.Msg.GetWhen().GetTz()

	event, errx := s.svc.CreateEvent(
		ctx, user, name, description, when, tz,
	)
	if errx != nil {
		return nil, convert.ToConnectRpcError(errx)
	}

	response := icbt.EventCreateResponse_builder{
		Event: convert.ToPbEvent(event),
	}.Build()
	return connect.NewResponse(response), nil
}

func (s *Server) EventUpdate(ctx context.Context,
	req *connect.Request[icbt.EventUpdateRequest],
) (*connect.Response[emptypb.Empty], error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid credentials"))
	}

	refID, err := service.ParseEventRefID(req.Msg.GetRefId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("bad event ref-id"))
	}

	euvs := &service.EventUpdateValues{}
	if req.Msg.HasName() {
		euvs.Name = mo.Some(req.Msg.GetName())
	}
	if req.Msg.HasDescription() {
		euvs.Description = mo.Some(req.Msg.GetDescription())
	}
	if req.Msg.HasWhen() {
		if req.Msg.GetWhen().GetTs().IsValid() {
			t := req.Msg.GetWhen().GetTs().AsTime()
			euvs.StartTime = mo.Some(t)
		}
		if req.Msg.GetWhen().GetTz() != "" {
			tz := req.Msg.GetWhen().GetTz()
			euvs.Tz = mo.Some(tz)
		}
	}

	errx := s.svc.UpdateEvent(ctx, user.ID, refID, euvs)
	if errx != nil {
		return nil, convert.ToConnectRpcError(errx)
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (s *Server) EventGetDetails(ctx context.Context,
	req *connect.Request[icbt.EventGetDetailsRequest],
) (*connect.Response[icbt.EventGetDetailsResponse], error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid credentials"))
	}

	refID, err := service.ParseEventRefID(req.Msg.GetRefId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("bad event ref-id"))
	}

	event, errx := s.svc.GetEvent(ctx, refID)
	if errx != nil {
		return nil, convert.ToConnectRpcError(errx)
	}
	pbEvent := convert.ToPbEvent(event)

	eventItems, errx := s.svc.GetEventItemsByEventID(ctx, event.ID)
	if errx != nil {
		return nil, convert.ToConnectRpcError(errx)
	}
	pbEventItems := convert.ToPbList(convert.ToPbEventItem, eventItems)

	earmarks, errx := s.svc.GetEarmarksByEventID(ctx, event.ID)
	if errx != nil {
		return nil, convert.ToConnectRpcError(errx)
	}
	pbEarmarks, err := convert.ToPbListWithService(ctx, convert.ToPbEarmark, s.svc, earmarks)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("db error"))
	}

	response := icbt.EventGetDetailsResponse_builder{
		Event:    pbEvent,
		Items:    pbEventItems,
		Earmarks: pbEarmarks,
	}.Build()
	return connect.NewResponse(response), nil
}

func (s *Server) EventDelete(ctx context.Context,
	req *connect.Request[icbt.EventDeleteRequest],
) (*connect.Response[emptypb.Empty], error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid credentials"))
	}

	refID, err := service.ParseEventRefID(req.Msg.GetRefId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("bad event ref-id"))
	}

	errx := s.svc.DeleteEvent(ctx, user.ID, refID)
	if errx != nil {
		return nil, convert.ToConnectRpcError(errx)
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}
