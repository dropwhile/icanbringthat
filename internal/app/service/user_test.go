package service

import (
	"context"
	"testing"

	"github.com/dropwhile/refid/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/samber/mo"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/errs"
)

func TestService_GetUser(t *testing.T) {
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

	t.Run("get user should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM user_").
			WithArgs(user.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "email", "name", "verified",
				}).
				AddRow(
					user.ID, user.RefID, user.Email, user.Name,
					user.Verified,
				),
			)

		result, err := svc.GetUser(ctx, user.RefID)
		assert.NilError(t, err)
		assert.Equal(t, result.RefID, user.RefID)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get user with no result should fail", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM user_").
			WithArgs(user.RefID).
			WillReturnError(pgx.ErrNoRows)

		_, err := svc.GetUser(ctx, user.RefID)
		errs.AssertError(t, err, errs.NotFound, "user not found")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_GetUserByEmail(t *testing.T) {
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

	t.Run("get user should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM user_").
			WithArgs(user.Email).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "email", "name", "verified",
				}).
				AddRow(
					user.ID, user.RefID, user.Email, user.Name,
					user.Verified,
				),
			)

		result, err := svc.GetUserByEmail(ctx, user.Email)
		assert.NilError(t, err)
		assert.Equal(t, result.RefID, user.RefID)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get user with no result should fail", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM user_").
			WithArgs(user.Email).
			WillReturnError(pgx.ErrNoRows)

		_, err := svc.GetUserByEmail(ctx, user.Email)
		errs.AssertError(t, err, errs.NotFound, "user not found")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_GetUserByID(t *testing.T) {
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

	t.Run("get user should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM user_").
			WithArgs(user.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "email", "name", "verified",
				}).
				AddRow(
					user.ID, user.RefID, user.Email, user.Name,
					user.Verified,
				),
			)

		result, err := svc.GetUserByID(ctx, user.ID)
		assert.NilError(t, err)
		assert.Equal(t, result.RefID, user.RefID)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get user with no result should fail", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM user_").
			WithArgs(user.ID).
			WillReturnError(pgx.ErrNoRows)

		_, err := svc.GetUserByID(ctx, user.ID)
		errs.AssertError(t, err, errs.NotFound, "user not found")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_GetUsersByIDs(t *testing.T) {
	t.Parallel()

	user1 := &model.User{
		ID:           1,
		RefID:        refid.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      tstTs,
		LastModified: tstTs,
	}
	user2 := &model.User{
		ID:           2,
		RefID:        refid.Must(model.NewUserRefID()),
		Email:        "user2@example.com",
		Name:         "user2",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      tstTs,
		LastModified: tstTs,
	}

	t.Run("get users should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM user_").
			WithArgs([]int{user1.ID, user2.ID}).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "email", "name", "verified",
				}).
				AddRow(
					user1.ID, user1.RefID, user1.Email, user1.Name,
					user1.Verified,
				).
				AddRow(
					user2.ID, user2.RefID, user2.Email, user2.Name,
					user2.Verified,
				),
			)

		result, err := svc.GetUsersByIDs(ctx, []int{user1.ID, user2.ID})
		assert.NilError(t, err)
		assert.Equal(t, len(result), 2)
		assert.Equal(t, result[0].RefID, user1.RefID)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get users with no result should fail", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM user_").
			WithArgs([]int{user1.ID, user2.ID}).
			WillReturnError(pgx.ErrNoRows)

		result, err := svc.GetUsersByIDs(ctx, []int{user1.ID, user2.ID})
		assert.NilError(t, err)
		assert.Equal(t, len(result), 0)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_NewUser(t *testing.T) {
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

	t.Run("add new user should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectBegin()
		mock.ExpectQuery("^INSERT INTO user_").
			WithArgs(pgx.NamedArgs{
				"refID":    pgxmock.AnyArg(),
				"email":    user.Email,
				"name":     user.Name,
				"pwHash":   pgxmock.AnyArg(),
				"pwAuth":   true,
				"settings": model.NewUserPropertyMap(),
			}).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "email", "name", "verified",
				}).
				AddRow(
					user.ID, user.RefID, user.Email, user.Name,
					user.Verified,
				),
			)
		mock.ExpectCommit()
		mock.ExpectRollback()

		result, err := svc.NewUser(ctx, user.Email, user.Name, []byte("00x00"))
		assert.NilError(t, err)
		assert.Equal(t, result.RefID, user.RefID)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("add new user with empty name should fail", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		_, err := svc.NewUser(ctx, user.Email, "", []byte("00x00"))
		errs.AssertError(t, err, errs.InvalidArgument, "name bad value")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("add new user with bad email should fail", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		_, err := svc.NewUser(ctx, "junky", user.Name, []byte("00x00"))
		errs.AssertError(t, err, errs.InvalidArgument, "email bad value")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("add new user with empty email should fail", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		_, err := svc.NewUser(ctx, "", user.Name, []byte("00x00"))
		errs.AssertError(t, err, errs.InvalidArgument, "email bad value")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("add new user that already exists should fail", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectBegin()
		mock.ExpectQuery("^INSERT INTO user_").
			WithArgs(pgx.NamedArgs{
				"refID":    pgxmock.AnyArg(),
				"email":    user.Email,
				"name":     user.Name,
				"pwHash":   pgxmock.AnyArg(),
				"pwAuth":   true,
				"settings": model.NewUserPropertyMap(),
			}).
			WillReturnError(&pgconn.PgError{
				ConstraintName: "user_email_idx",
			})
		mock.ExpectRollback()
		mock.ExpectRollback()

		_, err := svc.NewUser(ctx, user.Email, user.Name, []byte("00x00"))
		errs.AssertError(t, err, errs.AlreadyExists, "user already exists")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_UpdateUser(t *testing.T) {
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

	t.Run("update user should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		vals := &UserUpdateValues{
			Name:  mo.Some(user.Name),
			Email: mo.Some(user.Email),
		}

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE user_").
			WithArgs(pgx.NamedArgs{
				"userID":    user.ID,
				"email":     vals.Email,
				"name":      vals.Name,
				"pwHash":    vals.PWHash,
				"verified":  vals.Verified,
				"pwAuth":    vals.PWAuth,
				"apiAccess": vals.ApiAccess,
				"webAuthn":  vals.WebAuthn,
			}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		err := svc.UpdateUser(ctx, user.ID, vals)
		assert.NilError(t, err)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update user empty name should fail", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		vals := &UserUpdateValues{
			Name: mo.Some(""),
		}

		err := svc.UpdateUser(ctx, user.ID, vals)
		errs.AssertError(t, err, errs.InvalidArgument, "Name bad value")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update user bad email should fail", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		vals := &UserUpdateValues{
			Email: mo.Some("hodor"),
		}

		err := svc.UpdateUser(ctx, user.ID, vals)
		errs.AssertError(t, err, errs.InvalidArgument, "Email bad value")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update user empty email should fail", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		vals := &UserUpdateValues{
			Email: mo.Some(""),
		}

		err := svc.UpdateUser(ctx, user.ID, vals)
		errs.AssertError(t, err, errs.InvalidArgument, "Email bad value")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update user bad pwhash should fail", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		vals := &UserUpdateValues{
			PWHash: mo.Some([]byte{}),
		}

		err := svc.UpdateUser(ctx, user.ID, vals)
		errs.AssertError(t, err, errs.InvalidArgument, "PWHash bad value")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_UpdateUserSettings(t *testing.T) {
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

	t.Run("update user settings should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		settings := &model.UserSettings{
			ReminderThresholdHours: 0,
			EnableReminders:        false,
		}

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE user_").
			WithArgs(pgx.NamedArgs{
				"settings": settings,
				"userID":   user.ID,
			}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		err := svc.UpdateUserSettings(ctx, user.ID, settings)
		assert.NilError(t, err)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_DeleteUser(t *testing.T) {
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

	t.Run("delete user should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM user_").
			WithArgs(user.ID).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		err := svc.DeleteUser(ctx, user.ID)
		assert.NilError(t, err)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}
