package rpc

import (
	"context"
	"testing"

	"github.com/dropwhile/refid/v2"
	"github.com/twitchtv/twirp"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/internal/errs"
	"github.com/dropwhile/icbt/internal/middleware/auth"
	"github.com/dropwhile/icbt/rpc/icbt"
)

func TestRpc_ListNotifications(t *testing.T) {
	t.Parallel()

	user := &model.User{
		ID:           1,
		RefID:        refid.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      tstTs,
		LastModified: tstTs,
	}
	notification := &model.Notification{
		ID:           2,
		RefID:        refid.Must(model.NewNotificationRefID()),
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

		request := &icbt.ListNotificationsRequest{
			Pagination: &icbt.PaginationRequest{
				Limit:  10,
				Offset: 0,
			},
		}
		response, err := server.ListNotifications(ctx, request)
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

		request := &icbt.ListNotificationsRequest{}
		response, err := server.ListNotifications(ctx, request)
		assert.NilError(t, err)

		assert.Equal(t, len(response.Notifications), 1)
	})
}

func TestRpc_DeleteNotification(t *testing.T) {
	t.Parallel()

	user := &model.User{
		ID:           1,
		RefID:        refid.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      tstTs,
		LastModified: tstTs,
	}
	notification := &model.Notification{
		ID:           2,
		RefID:        refid.Must(model.NewNotificationRefID()),
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

		request := &icbt.DeleteNotificationRequest{
			RefId: "hodor",
		}
		_, err := server.DeleteNotification(ctx, request)
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

		request := &icbt.DeleteNotificationRequest{
			RefId: notification.RefID.String(),
		}
		_, err := server.DeleteNotification(ctx, request)
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

		request := &icbt.DeleteNotificationRequest{
			RefId: notification.RefID.String(),
		}
		_, err := server.DeleteNotification(ctx, request)
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

		request := &icbt.DeleteNotificationRequest{
			RefId: notification.RefID.String(),
		}
		_, err := server.DeleteNotification(ctx, request)
		assert.NilError(t, err)
	})
}

func TestRpc_DeleteAllNotifications(t *testing.T) {
	t.Parallel()

	user := &model.User{
		ID:           1,
		RefID:        refid.Must(model.NewUserRefID()),
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
		_, err := server.DeleteAllNotifications(ctx, request)
		assert.NilError(t, err)
	})
}
