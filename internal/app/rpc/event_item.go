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
	"github.com/dropwhile/icanbringthat/internal/app/service"
	"github.com/dropwhile/icanbringthat/internal/middleware/auth"

	icbt "github.com/dropwhile/icanbringthat/rpc/icbt/rpc/v1"
)

func (s *Server) EventListItems(ctx context.Context,
	req *connect.Request[icbt.EventListItemsRequest],
) (*connect.Response[icbt.EventListItemsResponse], error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid credentials"))
	}

	refID, err := service.ParseEventRefID(req.Msg.GetRefId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("bad event ref-id"))
	}

	items, errx := s.svc.GetEventItemsByEvent(ctx, refID)
	if errx != nil {
		return nil, convert.ToConnectRpcError(errx)
	}

	response := icbt.EventListItemsResponse_builder{
		Items: convert.ToPbList(convert.ToPbEventItem, items),
	}.Build()
	return connect.NewResponse(response), nil
}

func (s *Server) EventRemoveItem(ctx context.Context,
	req *connect.Request[icbt.EventRemoveItemRequest],
) (*connect.Response[emptypb.Empty], error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid credentials"))
	}

	refID, err := service.ParseEventItemRefID(req.Msg.GetRefId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("bad event-item ref-id"))
	}

	errx := s.svc.RemoveEventItem(ctx, user.ID, refID, nil)
	if errx != nil {
		return nil, convert.ToConnectRpcError(errx)
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (s *Server) EventAddItem(ctx context.Context,
	req *connect.Request[icbt.EventAddItemRequest],
) (*connect.Response[icbt.EventAddItemResponse], error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid credentials"))
	}

	refID, err := service.ParseEventRefID(req.Msg.GetEventRefId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("bad event ref-id"))
	}

	eventItem, errx := s.svc.AddEventItem(
		ctx, user.ID, refID, req.Msg.GetDescription(),
	)
	if errx != nil {
		return nil, convert.ToConnectRpcError(errx)
	}

	response := icbt.EventAddItemResponse_builder{
		EventItem: convert.ToPbEventItem(eventItem),
	}.Build()
	return connect.NewResponse(response), nil
}

func (s *Server) EventUpdateItem(ctx context.Context,
	req *connect.Request[icbt.EventUpdateItemRequest],
) (*connect.Response[icbt.EventUpdateItemResponse], error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid credentials"))
	}

	refID, err := service.ParseEventItemRefID(req.Msg.GetRefId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("bad event-item ref-id"))
	}

	eventItem, errx := s.svc.UpdateEventItem(
		ctx, user.ID, refID, req.Msg.GetDescription(), nil,
	)
	if errx != nil {
		return nil, convert.ToConnectRpcError(errx)
	}

	response := icbt.EventUpdateItemResponse_builder{
		EventItem: convert.ToPbEventItem(eventItem),
	}.Build()
	return connect.NewResponse(response), nil
}
