package service

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/samber/mo"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/errs"
	"github.com/dropwhile/icbt/internal/util"
)

func TestService_GetUserVerifyByRefID(t *testing.T) {
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

	t.Run("get user verify should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		refID := util.Must(model.NewUserVerifyRefID())

		mock.ExpectQuery("^SELECT (.+) FROM user_verify_").
			WithArgs(refID).
			WillReturnRows(pgxmock.NewRows(
				[]string{"ref_id", "user_id"}).
				AddRow(refID, user.ID),
			)

		result, err := svc.GetUserVerifyByRefID(ctx, refID)
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

		refID := util.Must(model.NewUserVerifyRefID())

		mock.ExpectQuery("^SELECT (.+) FROM user_verify_").
			WithArgs(refID).
			WillReturnError(pgx.ErrNoRows)

		_, err := svc.GetUserVerifyByRefID(ctx, refID)
		errs.AssertError(t, err, errs.NotFound, "verify not found")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_NewUserVerify(t *testing.T) {
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

	t.Run("add new user verify should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		refID := util.Must(model.NewUserVerifyRefID())

		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO user_verify_").
			WithArgs(pgx.NamedArgs{
				"refID":  UserVerifyRefIDMatcher,
				"userID": user.ID,
			}).
			WillReturnRows(pgxmock.NewRows(
				[]string{"ref_id", "user_id"}).
				AddRow(refID, user.ID),
			)
		mock.ExpectCommit()
		mock.ExpectRollback()

		result, err := svc.NewUserVerify(ctx, user.ID)
		assert.NilError(t, err)
		assert.Equal(t, result.RefID, refID)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_SetUserVerified(t *testing.T) {
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

	t.Run("set user verify should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		refID := util.Must(model.NewUserVerifyRefID())
		verify := &model.UserVerify{
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
				"pwHash":    mo.None[[]byte](),
				"verified":  mo.Some(true),
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
		mock.ExpectExec("DELETE FROM user_verify_").
			WithArgs(refID).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()
		// end inner tx
		mock.ExpectCommit()
		mock.ExpectRollback()

		err := svc.SetUserVerified(ctx, user, verify)
		assert.NilError(t, err)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("set user verify not owner should fail", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		refID := util.Must(model.NewUserVerifyRefID())
		verify := &model.UserVerify{
			RefID:   refID,
			UserID:  user.ID + 1,
			Created: tstTs,
		}

		err := svc.SetUserVerified(ctx, user, verify)
		errs.AssertError(t, err, errs.PermissionDenied, "permission denied")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}
