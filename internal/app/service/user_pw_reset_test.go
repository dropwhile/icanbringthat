// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/samber/mo"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/errs"
	"github.com/dropwhile/icanbringthat/internal/util"
)

func TestService_GetUserPWResetByRefID(t *testing.T) {
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

	t.Run("get user pwreset should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		refID := util.Must(model.NewUserPWResetRefID())

		mock.ExpectQuery("^SELECT (.+) FROM user_pw_reset_").
			WithArgs(refID).
			WillReturnRows(pgxmock.NewRows(
				[]string{"ref_id", "user_id"}).
				AddRow(refID, user.ID),
			)

		result, err := svc.GetUserPWResetByRefID(ctx, refID)
		assert.NilError(t, err)
		assert.Equal(t, result.RefID, refID)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get user verify not found should fail", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		refID := util.Must(model.NewUserPWResetRefID())

		mock.ExpectQuery("^SELECT (.+) FROM user_pw_reset_").
			WithArgs(refID).
			WillReturnError(pgx.ErrNoRows)

		_, err := svc.GetUserPWResetByRefID(ctx, refID)
		errs.AssertError(t, err, errs.NotFound, "pwreset not found")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_NewUserPWReset(t *testing.T) {
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

	t.Run("add new user pwreset should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		refID := util.Must(model.NewUserPWResetRefID())

		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO user_pw_reset_").
			WithArgs(pgx.NamedArgs{
				"refID":  UserPWResetRefIDMatcher,
				"userID": user.ID,
			}).
			WillReturnRows(pgxmock.NewRows(
				[]string{"ref_id", "user_id"}).
				AddRow(refID, user.ID),
			)
		mock.ExpectCommit()
		mock.ExpectRollback()

		result, err := svc.NewUserPWReset(ctx, user.ID)
		assert.NilError(t, err)
		assert.Equal(t, result.RefID, refID)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_UpdateUserPWReset(t *testing.T) {
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

	t.Run("set user pwreset should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		refID := util.Must(model.NewUserPWResetRefID())
		upw := &model.UserPWReset{
			RefID:   refID,
			UserID:  user.ID,
			Created: tstTs,
		}

		mock.ExpectBegin()
		// inner tx start
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE user_").
			WithArgs(pgx.NamedArgs{
				"userID":    user.ID,
				"email":     mo.None[string](),
				"name":      mo.None[string](),
				"pwHash":    mo.Some(user.PWHash),
				"verified":  mo.None[bool](),
				"pwAuth":    mo.None[bool](),
				"apiAccess": mo.None[bool](),
				"webAuthn":  mo.None[bool](),
			}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()
		// end inner tx
		// inner tx start
		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM user_pw_reset_").
			WithArgs(refID).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()
		// end inner tx
		mock.ExpectCommit()
		mock.ExpectRollback()

		err := svc.UpdateUserPWReset(ctx, user, upw)
		assert.NilError(t, err)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("set user pwreset with user upw delete error should fail", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		refID := util.Must(model.NewUserPWResetRefID())
		upw := &model.UserPWReset{
			RefID:   refID,
			UserID:  user.ID,
			Created: tstTs,
		}

		mock.ExpectBegin()
		// inner tx start
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE user_").
			WithArgs(pgx.NamedArgs{
				"userID":    user.ID,
				"email":     mo.None[string](),
				"name":      mo.None[string](),
				"pwHash":    mo.Some(user.PWHash),
				"verified":  mo.None[bool](),
				"pwAuth":    mo.None[bool](),
				"apiAccess": mo.None[bool](),
				"webAuthn":  mo.None[bool](),
			}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()
		// end inner tx
		// inner tx start
		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM user_pw_reset_").
			WithArgs(refID).
			WillReturnError(fmt.Errorf("honk honk"))
		mock.ExpectRollback()
		mock.ExpectRollback()
		// end inner tx
		mock.ExpectRollback()
		mock.ExpectRollback()

		err := svc.UpdateUserPWReset(ctx, user, upw)
		errs.AssertError(t, err, errs.Internal, "db error: honk honk")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("set user pwreset not owner should fail", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		refID := util.Must(model.NewUserPWResetRefID())
		upw := &model.UserPWReset{
			RefID:   refID,
			UserID:  user.ID + 1,
			Created: tstTs,
		}

		err := svc.UpdateUserPWReset(ctx, user, upw)
		errs.AssertError(t, err, errs.PermissionDenied, "permission denied")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}
