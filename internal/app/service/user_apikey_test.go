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

func TestService_GetApiKeyByUser(t *testing.T) {
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

	t.Run("get user apikey should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		token := "some-token"

		mock.ExpectQuery("^SELECT (.+) FROM api_key_").
			WithArgs(user.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{"user_id", "token"}).
				AddRow(user.ID, token),
			)

		result, err := svc.GetApiKeyByUser(ctx, user.ID)
		assert.NilError(t, err)
		assert.Equal(t, result.UserID, user.ID)
		assert.Equal(t, result.Token, token)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get user apikey not found should fail", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		mock.ExpectQuery("^SELECT (.+) FROM api_key_").
			WithArgs(user.ID).
			WillReturnError(pgx.ErrNoRows)

		_, err := svc.GetApiKeyByUser(ctx, user.ID)
		errs.AssertError(t, err, errs.NotFound, "user-api-key not found")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_GetUserByApiKey(t *testing.T) {
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

	t.Run("get user by apikey should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		token := "some-token"

		mock.ExpectQuery("^SELECT (.+) FROM user_").
			WithArgs(token).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "email", "name", "verified",
				}).
				AddRow(
					user.ID, user.RefID, user.Email, user.Name,
					user.Verified,
				),
			)

		result, err := svc.GetUserByApiKey(ctx, token)
		assert.NilError(t, err)
		assert.Equal(t, result.RefID, user.RefID)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("get user by apikey not found should fail", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		token := "some-token"

		mock.ExpectQuery("^SELECT (.+) FROM user_").
			WithArgs(token).
			WillReturnError(pgx.ErrNoRows)

		_, err := svc.GetUserByApiKey(ctx, token)
		errs.AssertError(t, err, errs.NotFound, "user not found")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_NewApiKey(t *testing.T) {
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

	t.Run("add user apikey should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		token := "some-token"

		mock.ExpectBegin()
		mock.ExpectQuery("^INSERT INTO api_key_").
			WithArgs(pgx.NamedArgs{
				"userID": user.ID,
				"token":  pgxmock.AnyArg(),
			}).
			WillReturnRows(pgxmock.NewRows(
				[]string{"user_id", "token"}).
				AddRow(user.ID, token),
			)
		mock.ExpectCommit()
		mock.ExpectRollback()

		result, err := svc.NewApiKey(ctx, user.ID)
		assert.NilError(t, err)
		assert.Equal(t, result.Token, token)
		assert.Equal(t, result.UserID, user.ID)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestService_NewApiKeyIfNotExists(t *testing.T) {
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

	t.Run("add user apikey should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		token := "some-token"

		mock.ExpectQuery("^SELECT (.+) FROM api_key_").
			WithArgs(user.ID).
			WillReturnError(pgx.ErrNoRows)
		mock.ExpectBegin()
		mock.ExpectQuery("^INSERT INTO api_key_").
			WithArgs(pgx.NamedArgs{
				"userID": user.ID,
				"token":  pgxmock.AnyArg(),
			}).
			WillReturnRows(pgxmock.NewRows(
				[]string{"user_id", "token"}).
				AddRow(user.ID, token),
			)
		mock.ExpectCommit()
		mock.ExpectRollback()

		result, err := svc.NewApiKeyIfNotExists(ctx, user.ID)
		assert.NilError(t, err)
		assert.Equal(t, result.Token, token)
		assert.Equal(t, result.UserID, user.ID)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("add user apikey already exists should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})

		token := "some-token"

		mock.ExpectQuery("^SELECT (.+) FROM api_key_").
			WithArgs(user.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{"user_id", "token"}).
				AddRow(user.ID, token),
			)

		result, err := svc.NewApiKeyIfNotExists(ctx, user.ID)
		assert.NilError(t, err)
		assert.Equal(t, result.Token, token)
		assert.Equal(t, result.UserID, user.ID)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}
