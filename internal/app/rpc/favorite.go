// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rpc

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/dropwhile/icanbringthat/internal/app/convert"
	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/app/service"
	"github.com/dropwhile/icanbringthat/internal/middleware/auth"

	icbt "github.com/dropwhile/icanbringthat/rpc/icbt/rpc/v1"
)

func (s *Server) FavoriteListEvents(ctx context.Context,
	req *connect.Request[icbt.FavoriteListEventsRequest],
) (*connect.Response[icbt.FavoriteListEventsResponse], error) {
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

		favs, pagination, errx := s.svc.GetFavoriteEventsPaginated(
			ctx, user.ID, limit, offset, showArchived)
		if errx != nil {
			return nil, convert.ToConnectRpcError(errx)
		}

		events = favs
		paginationResult = convert.ToPbPagination(pagination)
	} else {
		favs, errx := s.svc.GetFavoriteEvents(
			ctx, user.ID, showArchived)
		if errx != nil {
			return nil, convert.ToConnectRpcError(errx)
		}
		events = favs
	}

	response := icbt.FavoriteListEventsResponse_builder{
		Events:     convert.ToPbList(convert.ToPbEvent, events),
		Pagination: paginationResult,
	}.Build()
	return connect.NewResponse(response), nil
}

func (s *Server) FavoriteRemove(ctx context.Context,
	req *connect.Request[icbt.FavoriteRemoveRequest],
) (*connect.Response[emptypb.Empty], error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid credentials"))
	}

	refID, err := service.ParseEventRefID(req.Msg.GetEventRefId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("bad event ref-id"))
	}

	errx := s.svc.RemoveFavorite(ctx, user.ID, refID)
	if errx != nil {
		return nil, convert.ToConnectRpcError(errx)
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (s *Server) FavoriteAdd(ctx context.Context,
	req *connect.Request[icbt.FavoriteAddRequest],
) (*connect.Response[icbt.FavoriteAddResponse], error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid credentials"))
	}

	refID, err := service.ParseEventRefID(req.Msg.GetEventRefId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("bad event ref-id"))
	}

	favorite, errx := s.svc.AddFavorite(ctx, user.ID, refID)
	if errx != nil {
		return nil, convert.ToConnectRpcError(errx)
	}

	response := icbt.FavoriteAddResponse_builder{
		Favorite: icbt.Favorite_builder{
			EventRefId: req.Msg.GetEventRefId(),
			Created:    timestamppb.New(favorite.Created),
		}.Build(),
	}.Build()
	return connect.NewResponse(response), nil
}
