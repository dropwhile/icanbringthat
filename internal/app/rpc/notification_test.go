// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rpc

import (
	"context"
	"testing"

	"github.com/twitchtv/twirp"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/app/service"
	"github.com/dropwhile/icanbringthat/internal/errs"
	"github.com/dropwhile/icanbringthat/internal/middleware/auth"
	"github.com/dropwhile/icanbringthat/internal/util"
	"github.com/dropwhile/icanbringthat/rpc/icbt"
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
					Limit:  uint32(limit),
					Offset: uint32(offset),
					Count:  1,
				}, nil,
			)

		request := &icbt.NotificationsListRequest{
			Pagination: &icbt.PaginationRequest{
				Limit:  10,
				Offset: 0,
			},
		}
		response, err := server.NotificationsList(ctx, request)
		assert.NilError(t, err)
		assert.Equal(t, len(response.Notifications), 1)
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
		response, err := server.NotificationsList(ctx, request)
		assert.NilError(t, err)

		assert.Equal(t, len(response.Notifications), 1)
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

		request := &icbt.NotificationDeleteRequest{
			RefId: "hodor",
		}
		_, err := server.NotificationDelete(ctx, request)
		errs.AssertError(t, err, twirp.InvalidArgument, "ref_id incorrect value type")
	})

	t.Run("delete notification with no matching refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		mock.EXPECT().
			DeleteNotification(ctx, user.ID, notification.RefID).
			Return(errs.NotFound.Error("notification not found"))

		request := &icbt.NotificationDeleteRequest{
			RefId: notification.RefID.String(),
		}
		_, err := server.NotificationDelete(ctx, request)
		errs.AssertError(t, err, twirp.NotFound, "notification not found")
	})

	t.Run("delete notification for different user should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		mock.EXPECT().
			DeleteNotification(ctx, user.ID, notification.RefID).
			Return(errs.PermissionDenied.Error("permission denied"))

		request := &icbt.NotificationDeleteRequest{
			RefId: notification.RefID.String(),
		}
		_, err := server.NotificationDelete(ctx, request)
		errs.AssertError(t, err, twirp.PermissionDenied, "permission denied")
	})

	t.Run("delete notification should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		server, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "user", user)

		mock.EXPECT().
			DeleteNotification(ctx, user.ID, notification.RefID).
			Return(nil)

		request := &icbt.NotificationDeleteRequest{
			RefId: notification.RefID.String(),
		}
		_, err := server.NotificationDelete(ctx, request)
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

		request := &icbt.Empty{}
		_, err := server.NotificationsDeleteAll(ctx, request)
		assert.NilError(t, err)
	})
}
