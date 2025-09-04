// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package service

import (
	"context"
	"testing"

	"github.com/dropwhile/assert"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/errs"
	"github.com/dropwhile/icanbringthat/internal/util"
)

func TestService_DisableRemindersWithNotification(t *testing.T) {
	t.Parallel()

	user := &model.User{
		ID:     1,
		RefID:  util.Must(model.NewUserRefID()),
		Email:  "user@example.com",
		Name:   "user",
		PWHash: []byte("00x00"),
		Settings: model.UserSettings{
			ReminderThresholdHours: 60,
			EnableReminders:        true,
		},
		Verified:     true,
		Created:      tstTs,
		LastModified: tstTs,
	}

	t.Run("disable reminders should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		email := "user@example.com"
		reason := "just-because"
		msg := "email notifications disabled due to 'just-because'"

		mock.ExpectQuery("^SELECT (.+) FROM user_").
			WithArgs(user.Email).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "email", "name", "pwhash", "pwauth",
					"created", "last_modified", "settings",
				}).
				AddRow(
					user.ID, user.RefID, user.Email, user.Name,
					user.PWHash, user.PWAuth, user.Created, user.LastModified,
					user.Settings,
				),
			)
		// outer tx begin
		mock.ExpectBegin()
		// inner tx 1 begin
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE user_ ").
			WithArgs(pgx.NamedArgs{
				"userID": user.ID,
				"settings": &model.UserSettings{
					ReminderThresholdHours: user.Settings.ReminderThresholdHours,
					EnableReminders:        false,
				},
			}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()
		// inner tx 2 begin
		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO notification_ ").
			WithArgs(pgx.NamedArgs{
				"refID":   pgxmock.AnyArg(),
				"userID":  user.ID,
				"message": msg,
			}).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "message", "read",
				}).
				AddRow(1, util.Must(model.NewNotificationRefID()),
					msg, false,
				),
			)
		mock.ExpectCommit()
		mock.ExpectRollback()
		// outer tx end
		mock.ExpectCommit()
		mock.ExpectRollback()

		err := svc.DisableRemindersWithNotification(ctx, email, reason)
		assert.Nil(t, err)
		// we make sure that all expectations were met
		assert.Nil(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("disable reminders with user not found should fail", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		email := "user@example.com"
		reason := "just-because"

		mock.ExpectQuery("^SELECT (.+) FROM user_").
			WithArgs(user.Email).
			WillReturnError(pgx.ErrNoRows)

		err := svc.DisableRemindersWithNotification(ctx, email, reason)
		errs.AssertError(t, err, errs.NotFound, "user not found")
		// we make sure that all expectations were met
		assert.Nil(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("disable reminders already disabled should fail", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		email := "user@example.com"
		reason := "just-because"

		settings := model.UserSettings{
			ReminderThresholdHours: 60,
			EnableReminders:        false,
		}

		mock.ExpectQuery("^SELECT (.+) FROM user_").
			WithArgs(user.Email).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "email", "name", "pwhash", "pwauth",
					"created", "last_modified", "settings",
				}).
				AddRow(
					user.ID, user.RefID, user.Email, user.Name,
					user.PWHash, user.PWAuth, user.Created, user.LastModified,
					settings,
				),
			)

		err := svc.DisableRemindersWithNotification(ctx, email, reason)
		errs.AssertError(t, err, errs.FailedPrecondition, "reminders already disabled")
		// we make sure that all expectations were met
		assert.Nil(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("disable reminders with bad email should fail", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		email := "userexample.com"
		reason := "just-because"

		err := svc.DisableRemindersWithNotification(ctx, email, reason)
		errs.AssertError(t, err, errs.InvalidArgument, "email bad value")
		// we make sure that all expectations were met
		assert.Nil(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("disable reminders with empty reason should fail", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		email := "user@example.com"
		reason := ""

		err := svc.DisableRemindersWithNotification(ctx, email, reason)
		errs.AssertError(t, err, errs.InvalidArgument, "suppressionReason bad value")
		// we make sure that all expectations were met
		assert.Nil(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}
