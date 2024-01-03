package service

import (
	"context"
	"testing"

	"github.com/dropwhile/refid/v2"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/errs"
)

func TestService_GetUserCredentialByRefID(t *testing.T) {
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

	t.Run("get user credential should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		credRefID := refid.Must(model.NewCredentialRefID())

		mock.ExpectQuery("^SELECT (.+) FROM user_webauthn_").
			WithArgs(credRefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "key_name", "credential",
				}).
				AddRow(
					1, credRefID, user.ID, "key-name", []byte{0x00, 0x01},
				),
			)

		result, err := svc.GetUserCredentialByRefID(ctx, credRefID)
		assert.NilError(t, err)
		assert.Equal(t, result.RefID, credRefID)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get user credential not found should fail", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		credRefID := refid.Must(model.NewCredentialRefID())

		mock.ExpectQuery("^SELECT (.+) FROM user_webauthn_").
			WithArgs(credRefID).
			WillReturnError(pgx.ErrNoRows)

		_, err := svc.GetUserCredentialByRefID(ctx, credRefID)
		errs.AssertError(t, err, errs.NotFound, "credential not found")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_GetUserCredentialsByUser(t *testing.T) {
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

	t.Run("get user credentials should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		credRefID := refid.Must(model.NewCredentialRefID())
		credRefID2 := refid.Must(model.NewCredentialRefID())

		mock.ExpectQuery("^SELECT (.+) FROM user_webauthn_").
			WithArgs(user.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "key_name", "credential",
				}).
				AddRow(
					1, credRefID, user.ID, "key-name", []byte{0x00, 0x01},
				).
				AddRow(
					2, credRefID2, user.ID, "key-name2", []byte{0x00, 0x01},
				),
			)

		result, err := svc.GetUserCredentialsByUser(ctx, user.ID)
		assert.NilError(t, err)
		assert.Equal(t, len(result), 2)
		assert.Equal(t, result[0].RefID, credRefID)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get user credentials not found should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM user_webauthn_").
			WithArgs(user.ID).
			WillReturnError(pgx.ErrNoRows)

		result, err := svc.GetUserCredentialsByUser(ctx, user.ID)
		assert.NilError(t, err)
		assert.Equal(t, len(result), 0)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_GetUserCredentialCountByUser(t *testing.T) {
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

	t.Run("get user credentials count should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM user_webauthn_").
			WithArgs(user.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{"count"}).
				AddRow(2),
			)

		result, err := svc.GetUserCredentialCountByUser(ctx, user.ID)
		assert.NilError(t, err)
		assert.Equal(t, result, 2)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get user credentials count not found should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM user_webauthn_").
			WithArgs(user.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{"count"}).
				AddRow(0),
			)

		result, err := svc.GetUserCredentialCountByUser(ctx, user.ID)
		assert.NilError(t, err)
		assert.Equal(t, result, 0)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_DeleteUserCredential(t *testing.T) {
	t.Parallel()

	t.Run("delete user credential should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		credentialID := 3

		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM user_webauthn_").
			WithArgs(credentialID).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		err := svc.DeleteUserCredential(ctx, credentialID)
		assert.NilError(t, err)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_NewUserCredential(t *testing.T) {
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

	t.Run("add new credential should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		keyName := "key-name"
		credential := []byte{0x99, 0x98}
		credRefID := refid.Must(model.NewCredentialRefID())

		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO user_webauthn_").
			WithArgs(pgx.NamedArgs{
				"refID":      CredentialRefIDMatcher,
				"userID":     user.ID,
				"credential": credential,
				"keyName":    keyName,
			}).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "key_name", "credential",
				}).
				AddRow(
					1, credRefID, user.ID, keyName, credential,
				),
			)
		mock.ExpectCommit()
		mock.ExpectRollback()

		result, err := svc.NewUserCredential(ctx, user.ID, keyName, credential)
		assert.NilError(t, err)
		assert.Equal(t, result.RefID, credRefID)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("add new credential bad keyname should fail", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		keyName := ""
		credential := []byte{0x99, 0x98}

		_, err := svc.NewUserCredential(ctx, user.ID, keyName, credential)
		errs.AssertError(t, err, errs.InvalidArgument, "keyName bad value")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("add new credential bad credential should fail", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		keyName := "key-name"
		credential := []byte{}

		_, err := svc.NewUserCredential(ctx, user.ID, keyName, credential)
		errs.AssertError(t, err, errs.InvalidArgument, "credential bad value")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}
