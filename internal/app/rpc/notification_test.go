package rpc

import (
	"context"
	"testing"

	"github.com/dropwhile/refid/v2"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/twitchtv/twirp"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/app/model"
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
		mock := SetupDBMock(t, ctx)
		server := NewTestServer(mock)
		ctx = auth.ContextSet(ctx, "user", user)

		mock.ExpectQuery("SELECT count(.+) FROM notification_").
			WithArgs(user.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{"count"}).
				AddRow(1),
			)
		mock.ExpectQuery("SELECT (.+) FROM notification_").
			WithArgs(pgx.NamedArgs{
				"userID": user.ID,
				"limit":  pgxmock.AnyArg(),
				"offset": pgxmock.AnyArg(),
			}).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id",
					"user_id", "message",
					"read", "created", "last_modified",
				}).
				AddRow(
					notification.ID, notification.RefID,
					user.ID, notification.Message,
					notification.Read, tstTs, tstTs,
				),
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
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("list notifications non-paginated should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := NewTestServer(mock)
		ctx = auth.ContextSet(ctx, "user", user)

		mock.ExpectQuery("SELECT (.+) FROM notification_").
			WithArgs(pgx.NamedArgs{
				"userID": user.ID,
			}).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id",
					"user_id", "message",
					"read", "created", "last_modified",
				}).
				AddRow(
					notification.ID, notification.RefID,
					user.ID, notification.Message,
					notification.Read, tstTs, tstTs,
				),
			)

		request := &icbt.ListNotificationsRequest{}
		response, err := server.ListNotifications(ctx, request)
		assert.NilError(t, err)

		assert.Equal(t, len(response.Notifications), 1)
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
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
		mock := SetupDBMock(t, ctx)
		server := NewTestServer(mock)
		ctx = auth.ContextSet(ctx, "user", user)

		request := &icbt.DeleteNotificationRequest{
			RefId: "hodor",
		}
		_, err := server.DeleteNotification(ctx, request)
		assertTwirpError(t, err, twirp.InvalidArgument, "ref_id incorrect value type")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete notification with no matching refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := NewTestServer(mock)
		ctx = auth.ContextSet(ctx, "user", user)

		mock.ExpectQuery("SELECT (.+) FROM notification_").
			WithArgs(notification.RefID).
			WillReturnError(pgx.ErrNoRows)

		request := &icbt.DeleteNotificationRequest{
			RefId: notification.RefID.String(),
		}
		_, err := server.DeleteNotification(ctx, request)
		assertTwirpError(t, err, twirp.NotFound, "notification not found")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete notification for different user should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := NewTestServer(mock)
		ctx = auth.ContextSet(ctx, "user", user)

		mock.ExpectQuery("SELECT (.+) FROM notification_").
			WithArgs(notification.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id",
					"user_id", "message",
					"read", "created", "last_modified",
				}).
				AddRow(
					notification.ID, notification.RefID,
					33, notification.Message,
					notification.Read, tstTs, tstTs,
				),
			)

		request := &icbt.DeleteNotificationRequest{
			RefId: notification.RefID.String(),
		}
		_, err := server.DeleteNotification(ctx, request)
		assertTwirpError(t, err, twirp.PermissionDenied, "permission denied")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete notification should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		server := NewTestServer(mock)
		ctx = auth.ContextSet(ctx, "user", user)

		mock.ExpectQuery("SELECT (.+) FROM notification_").
			WithArgs(notification.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id",
					"user_id", "message",
					"read", "created", "last_modified",
				}).
				AddRow(
					notification.ID, notification.RefID,
					user.ID, notification.Message,
					notification.Read, tstTs, tstTs,
				),
			)
		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM notification_").
			WithArgs(notification.ID).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		request := &icbt.DeleteNotificationRequest{
			RefId: notification.RefID.String(),
		}
		_, err := server.DeleteNotification(ctx, request)
		assert.NilError(t, err)
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
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
		mock := SetupDBMock(t, ctx)
		server := NewTestServer(mock)
		ctx = auth.ContextSet(ctx, "user", user)

		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM notification_").
			WithArgs(user.ID).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		request := &icbt.DeleteAllNotificationsRequest{}
		_, err := server.DeleteAllNotifications(ctx, request)
		assert.NilError(t, err)
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}
