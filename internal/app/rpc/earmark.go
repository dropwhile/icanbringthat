// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rpc

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/dropwhile/icanbringthat/internal/app/convert"
	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/app/service"
	"github.com/dropwhile/icanbringthat/internal/errs"
	"github.com/dropwhile/icanbringthat/internal/middleware/auth"
	icbt "github.com/dropwhile/icanbringthat/rpc/icbt/rpc/v1"
)

func (s *Server) EventListEarmarks(ctx context.Context,
	req *connect.Request[icbt.EventListEarmarksRequest],
) (*connect.Response[icbt.EventListEarmarksResponse], error) {
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

	earmarks, errx := s.svc.GetEarmarksByEventID(ctx, event.ID)
	if errx != nil {
		return nil, convert.ToConnectRpcError(errx)
	}

	pbEarmarks, err := convert.ToPbListWithService(ctx, convert.ToPbEarmark, s.svc, earmarks)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("db error"))
	}
	response := icbt.EventListEarmarksResponse_builder{
		Earmarks: pbEarmarks,
	}.Build()
	return connect.NewResponse(response), nil
}

func (s *Server) EarmarksList(ctx context.Context,
	req *connect.Request[icbt.EarmarksListRequest],
) (*connect.Response[icbt.EarmarksListResponse], error) {
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
	var earmarks []*model.Earmark
	if req.Msg.HasPagination() {
		ems, pgResult, errx := s.svc.GetEarmarksPaginated(
			ctx, user.ID,
			int(req.Msg.GetPagination().GetLimit()),
			int(req.Msg.GetPagination().GetOffset()),
			showArchived)
		if errx != nil {
			return nil, convert.ToConnectRpcError(errx)
		}
		paginationResult = convert.ToPbPagination(pgResult)
		earmarks = ems
	} else {
		var errx errs.Error
		earmarks, errx = s.svc.GetEarmarks(ctx, user.ID, showArchived)
		if errx != nil {
			return nil, convert.ToConnectRpcError(errx)
		}
	}

	pbEarmarks, err := convert.ToPbListWithService(ctx, convert.ToPbEarmark, s.svc, earmarks)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("db error"))
	}
	response := icbt.EarmarksListResponse_builder{
		Earmarks:   pbEarmarks,
		Pagination: paginationResult,
	}.Build()
	return connect.NewResponse(response), nil
}

func (s *Server) EarmarkCreate(ctx context.Context,
	req *connect.Request[icbt.EarmarkCreateRequest],
) (*connect.Response[icbt.EarmarkCreateResponse], error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid credentials"))
	}

	eventItemRefID, err := service.ParseEventItemRefID(req.Msg.GetEventItemRefId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("bad event-item ref-id"))
	}

	eventItem, errx := s.svc.GetEventItem(ctx, eventItemRefID)
	if errx != nil {
		return nil, convert.ToConnectRpcError(errx)
	}

	earmark, errx := s.svc.NewEarmark(ctx, user, eventItem.ID, req.Msg.GetNote())
	if errx != nil {
		return nil, convert.ToConnectRpcError(errx)
	}

	pbEarmark, err := convert.ToPbEarmark(ctx, s.svc, earmark)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("db error"))
	}

	response := icbt.EarmarkCreateResponse_builder{
		Earmark: pbEarmark,
	}.Build()
	return connect.NewResponse(response), nil
}

func (s *Server) EarmarkGetDetails(ctx context.Context,
	req *connect.Request[icbt.EarmarkGetDetailsRequest],
) (*connect.Response[icbt.EarmarkGetDetailsResponse], error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid credentials"))
	}

	refID, err := service.ParseEarmarkRefID(req.Msg.GetRefId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("bad earmark ref-id"))
	}

	earmark, errx := s.svc.GetEarmark(ctx, refID)
	if errx != nil {
		return nil, convert.ToConnectRpcError(errx)
	}

	eventItem, errx := s.svc.GetEventItemByID(ctx, earmark.EventItemID)
	if errx != nil {
		return nil, convert.ToConnectRpcError(errx)
	}

	event, errx := s.svc.GetEventByID(ctx, eventItem.EventID)
	if errx != nil {
		return nil, convert.ToConnectRpcError(errx)
	}

	pbEarmark, err := convert.ToPbEarmark(ctx, s.svc, earmark)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("db error"))
	}
	response := icbt.EarmarkGetDetailsResponse_builder{
		Earmark:    pbEarmark,
		EventRefId: event.RefID.String(),
	}.Build()
	return connect.NewResponse(response), nil
}

func (s *Server) EarmarkRemove(ctx context.Context,
	req *connect.Request[icbt.EarmarkRemoveRequest],
) (*connect.Response[emptypb.Empty], error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid credentials"))
	}

	refID, err := service.ParseEarmarkRefID(req.Msg.GetRefId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("bad earmark ref-id"))
	}

	errx := s.svc.DeleteEarmarkByRefID(ctx, user.ID, refID)
	if errx != nil {
		return nil, convert.ToConnectRpcError(errx)
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}
