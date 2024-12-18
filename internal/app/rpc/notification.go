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
	"github.com/dropwhile/icanbringthat/internal/middleware/auth"

	icbt "github.com/dropwhile/icanbringthat/rpc/icbt/rpc/v1"
)

func (s *Server) NotificationsList(ctx context.Context,
	req *connect.Request[icbt.NotificationsListRequest],
) (*connect.Response[icbt.NotificationsListResponse], error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid credentials"))
	}

	var paginationResult *icbt.PaginationResult
	var notifications []*model.Notification
	if req.Msg.HasPagination() {
		limit := int(req.Msg.GetPagination().GetLimit())
		offset := int(req.Msg.GetPagination().GetOffset())
		notifs, pagination, errx := s.svc.GetNotificationsPaginated(ctx, user.ID, limit, offset)
		if errx != nil {
			return nil, convert.ToConnectRpcError(errx)
		}

		notifications = notifs
		paginationResult = convert.ToPbPagination(pagination)
	} else {
		notifs, errx := s.svc.GetNotifications(ctx, user.ID)
		if errx != nil {
			return nil, convert.ToConnectRpcError(errx)
		}
		notifications = notifs

	}

	response := icbt.NotificationsListResponse_builder{
		Notifications: convert.ToPbList(convert.ToPbNotification, notifications),
		Pagination:    paginationResult,
	}.Build()
	return connect.NewResponse(response), nil
}

func (s *Server) NotificationDelete(ctx context.Context,
	req *connect.Request[icbt.NotificationDeleteRequest],
) (*connect.Response[emptypb.Empty], error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid credentials"))
	}

	refID, err := service.ParseNotificationRefID(req.Msg.GetRefId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("bad notification ref-id"))
	}

	errx := s.svc.DeleteNotification(ctx, user.ID, refID)
	if errx != nil {
		return nil, convert.ToConnectRpcError(errx)
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (s *Server) NotificationsDeleteAll(ctx context.Context,
	req *connect.Request[icbt.NotificationsDeleteAllRequest],
) (*connect.Response[emptypb.Empty], error) {
	// get user from auth in context
	user, err := auth.UserFromContext(ctx)
	if err != nil || user == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid credentials"))
	}

	errx := s.svc.DeleteAllNotifications(ctx, user.ID)
	if errx != nil {
		return nil, convert.ToConnectRpcError(errx)
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}
