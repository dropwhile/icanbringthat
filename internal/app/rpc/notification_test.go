// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rpc

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/app/service"
	"github.com/dropwhile/icanbringthat/internal/errs"
	"github.com/dropwhile/icanbringthat/internal/middleware/auth"
	"github.com/dropwhile/icanbringthat/internal/util"
	icbt "github.com/dropwhile/icanbringthat/rpc/icbt/rpc/v1"
)

func TestRpc_ListNotifications(t *testing.T) {
	t.Parallel()

	user := &model.User{
		ID:           1,
		RefID:        util.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      tstTs,
		LastModified: tstTs,
	}
	notification := &model.Notification{
		ID:           2,
		RefID:        util.Must(model.NewNotificationRefID()),
		UserID:       user.ID,
		Message:      "",
		Read:         false,
		Created:      tstTs,
		LastModified: tstTs,
	}

	t.Run("list notifications paginated should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		limit := 10
		offset := 0

		mock.EXPECT().
			GetNotificationsPaginated(ctx, user.ID, limit, offset).
			Return(
				[]*model.Notification{notification},
				&service.Pagination{
					Limit:  limit,
					Offset: offset,
					Count:  1,
				}, nil,
			)

		request := icbt.NotificationsListRequest_builder{
			Pagination: icbt.PaginationRequest_builder{
				Limit:  10,
				Offset: 0,
			}.Build(),
		}.Build()
		response, err := server.NotificationsList(ctx, connect.NewRequest(request))
		assert.NilError(t, err)
		assert.Equal(t, len(response.Msg.GetNotifications()), 1)
	})

	t.Run("list notifications non-paginated should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		mock.EXPECT().
			GetNotifications(ctx, user.ID).
			Return([]*model.Notification{notification}, nil)

		request := &icbt.NotificationsListRequest{}
		response, err := server.NotificationsList(ctx, connect.NewRequest(request))
		assert.NilError(t, err)

		assert.Equal(t, len(response.Msg.GetNotifications()), 1)
	})
}

func TestRpc_DeleteNotification(t *testing.T) {
	t.Parallel()

	user := &model.User{
		ID:           1,
		RefID:        util.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      tstTs,
		LastModified: tstTs,
	}
	notification := &model.Notification{
		ID:           2,
		RefID:        util.Must(model.NewNotificationRefID()),
		UserID:       user.ID,
		Message:      "",
		Read:         false,
		Created:      tstTs,
		LastModified: tstTs,
	}

	t.Run("delete notification with bad refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, _ := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		request := icbt.NotificationDeleteRequest_builder{
			RefId: "hodor",
		}.Build()
		_, err := server.NotificationDelete(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodeInvalidArgument, "bad notification ref-id")
	})

	t.Run("delete notification with no matching refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		mock.EXPECT().
			DeleteNotification(ctx, user.ID, notification.RefID).
			Return(errs.NotFound.Error("notification not found"))

		request := icbt.NotificationDeleteRequest_builder{
			RefId: notification.RefID.String(),
		}.Build()
		_, err := server.NotificationDelete(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodeNotFound, "notification not found")
	})

	t.Run("delete notification for different user should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		mock.EXPECT().
			DeleteNotification(ctx, user.ID, notification.RefID).
			Return(errs.PermissionDenied.Error("permission denied"))

		request := icbt.NotificationDeleteRequest_builder{
			RefId: notification.RefID.String(),
		}.Build()
		_, err := server.NotificationDelete(ctx, connect.NewRequest(request))
		rpcErr := AsConnectError(t, err)
		errs.AssertError(t, rpcErr, connect.CodePermissionDenied, "permission denied")
	})

	t.Run("delete notification should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		mock.EXPECT().
			DeleteNotification(ctx, user.ID, notification.RefID).
			Return(nil)

		request := icbt.NotificationDeleteRequest_builder{
			RefId: notification.RefID.String(),
		}.Build()
		_, err := server.NotificationDelete(ctx, connect.NewRequest(request))
		assert.NilError(t, err)
	})
}

func TestRpc_DeleteAllNotifications(t *testing.T) {
	t.Parallel()

	user := &model.User{
		ID:           1,
		RefID:        util.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      tstTs,
		LastModified: tstTs,
	}

	t.Run("delete all notification should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		mock.EXPECT().
			DeleteAllNotifications(ctx, user.ID).
			Return(nil)

		request := &icbt.NotificationsDeleteAllRequest{}
		_, err := server.NotificationsDeleteAll(ctx, connect.NewRequest(request))
		assert.NilError(t, err)
	})
}
