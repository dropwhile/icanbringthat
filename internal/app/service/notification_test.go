package service

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/errs"
	"github.com/dropwhile/icbt/internal/util"
)

func TestService_GetNotificationsCount(t *testing.T) {
	t.Parallel()

	t.Run("count with results should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		userID := 4
		count := 1

		mock.ExpectQuery("^SELECT (.+) FROM notification_").
			WithArgs(userID).
			WillReturnRows(pgxmock.NewRows(
				[]string{"count"}).
				AddRow(count),
			)

		result, err := svc.GetNotificationsCount(ctx, userID)
		assert.NilError(t, err)
		assert.Equal(t, result, count)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_GetNotificationsPaginated(t *testing.T) {
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

	t.Run("get with results should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		userID := 4
		limit := 5
		offset := 0
		count := 2

		mock.ExpectQuery("SELECT count(.+) FROM notification_").
			WithArgs(userID).
			WillReturnRows(pgxmock.NewRows(
				[]string{"count"}).
				AddRow(count),
			)
		mock.ExpectQuery(`SELECT [*] FROM notification_`).
			WithArgs(pgx.NamedArgs{
				"userID": userID,
				"limit":  limit,
				"offset": offset,
			}).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "read", "message",
				}).
				AddRow(
					1, util.Must(model.NewNotificationRefID()), user.ID,
					false, "some message 1",
				).
				AddRow(
					2, util.Must(model.NewNotificationRefID()), user.ID,
					false, "some message 2",
				),
			)

		notifications, pagination, err := svc.GetNotificationsPaginated(ctx, userID, limit, offset)
		assert.NilError(t, err)
		assert.Equal(t, len(notifications), count)
		assert.Equal(t, pagination.Limit, uint32(limit))
		assert.Equal(t, pagination.Offset, uint32(offset))
		assert.Equal(t, pagination.Count, uint32(count))
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get with no results should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		userID := 4
		limit := 5
		offset := 0

		mock.ExpectQuery("SELECT count(.+) FROM notification_").
			WithArgs(userID).
			WillReturnRows(pgxmock.NewRows(
				[]string{"count"}).
				AddRow(0),
			)

		notifications, pagination, err := svc.GetNotificationsPaginated(ctx, userID, limit, offset)
		assert.NilError(t, err)
		assert.Equal(t, len(notifications), 0)
		assert.Equal(t, pagination.Limit, uint32(limit))
		assert.Equal(t, pagination.Offset, uint32(offset))
		assert.Equal(t, pagination.Count, uint32(0))
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_GetNotifications(t *testing.T) {
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

	t.Run("get with results should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		count := 2

		mock.ExpectQuery(`SELECT [*] FROM notification_`).
			WithArgs(pgx.NamedArgs{
				"userID": user.ID,
			}).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "read", "message",
				}).
				AddRow(
					1, util.Must(model.NewNotificationRefID()), user.ID,
					false, "some message 1",
				).
				AddRow(
					2, util.Must(model.NewNotificationRefID()), user.ID,
					false, "some message 2",
				),
			)

		notifications, err := svc.GetNotifications(ctx, user.ID)
		assert.NilError(t, err)
		assert.Equal(t, len(notifications), count)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get with no results should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery(`SELECT [*] FROM notification_`).
			WithArgs(pgx.NamedArgs{
				"userID": user.ID,
			}).
			WillReturnError(pgx.ErrNoRows)

		notifications, err := svc.GetNotifications(ctx, user.ID)
		assert.NilError(t, err)
		assert.Equal(t, len(notifications), 0)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_DeleteNotification(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		ID:           1,
		RefID:        util.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      ts,
		LastModified: ts,
	}
	notification := &model.Notification{
		ID:      2,
		RefID:   util.Must(model.NewNotificationRefID()),
		UserID:  user.ID,
		Message: "message",
		Read:    false,
	}

	t.Run("delete should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery(`SELECT [*] FROM notification_`).
			WithArgs(notification.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "read", "message",
				}).
				AddRow(
					notification.ID, notification.RefID,
					notification.UserID, notification.Read,
					notification.Message,
				),
			)
		mock.ExpectBegin()
		mock.ExpectExec("^DELETE FROM notification_").
			WithArgs(notification.ID).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		err := svc.DeleteNotification(ctx, user.ID, notification.RefID)
		assert.NilError(t, err)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete with different user owner should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery(`SELECT [*] FROM notification_`).
			WithArgs(notification.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "read", "message",
				}).
				AddRow(
					notification.ID, notification.RefID,
					notification.UserID, notification.Read,
					notification.Message,
				),
			)

		err := svc.DeleteNotification(ctx, user.ID+1, notification.RefID)
		errs.AssertError(t, err, errs.PermissionDenied, "permission denied")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete with missing notification should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery(`SELECT [*] FROM notification_`).
			WithArgs(notification.RefID).
			WillReturnError(pgx.ErrNoRows)

		err := svc.DeleteNotification(ctx, user.ID, notification.RefID)
		errs.AssertError(t, err, errs.NotFound, "notification not found")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_DeleteAllNotifications(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		ID:           1,
		RefID:        util.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      ts,
		LastModified: ts,
	}

	t.Run("delete should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectBegin()
		mock.ExpectExec("^DELETE FROM notification_").
			WithArgs(user.ID).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		err := svc.DeleteAllNotifications(ctx, user.ID)
		assert.NilError(t, err)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_NewNotification(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		ID:           1,
		RefID:        util.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      ts,
		LastModified: ts,
	}
	notification := &model.Notification{
		ID:      2,
		RefID:   util.Must(model.NewNotificationRefID()),
		UserID:  user.ID,
		Message: "message",
		Read:    false,
	}

	t.Run("add notification should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectBegin()
		mock.ExpectQuery("^INSERT INTO notification_").
			WithArgs(pgx.NamedArgs{
				"refID":   NotificationRefIDMatcher,
				"userID":  user.ID,
				"message": notification.Message,
			}).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "read", "message",
				}).
				AddRow(
					notification.ID,
					notification.RefID,
					notification.UserID,
					notification.Read,
					notification.Message,
				),
			)
		mock.ExpectCommit()
		mock.ExpectRollback()

		result, err := svc.NewNotification(ctx, user.ID, notification.Message)
		assert.NilError(t, err)
		assert.Equal(t, result.RefID, notification.RefID)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("add notification with bad message should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		_, err := svc.NewNotification(ctx, user.ID, "")
		errs.AssertError(t, err, errs.InvalidArgument, "message bad value")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}
