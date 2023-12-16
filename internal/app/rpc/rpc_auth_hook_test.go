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
)

func TestRpc_AuthHook(t *testing.T) {
	t.Parallel()

	t.Run("auth hook with no api-key should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)

		_, err := AuthHook(mock)(ctx)
		assertTwirpError(t, err, twirp.Unauthenticated, "invalid auth")
	})

	t.Run("auth hook with api-key not finding user should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		ctx = auth.ContextSet(ctx, "api-key", "user-123")

		mock.ExpectQuery("^SELECT (.+) FROM user_").
			WithArgs("user-123").
			WillReturnError(pgx.ErrNoRows)

		_, err := AuthHook(mock)(ctx)
		assertTwirpError(t, err, twirp.Unauthenticated, "invalid auth")
	})

	t.Run("auth hook with user ApiAccess disabled should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		ctx = auth.ContextSet(ctx, "api-key", "user-123")
		refID := refid.Must(model.NewUserRefID())

		rows := pgxmock.NewRows(
			[]string{"id", "ref_id", "email", "name", "api_access", "verified", "created", "last_modified"}).
			AddRow(1, refID, "user@example.com", "user", false, true, tstTs, tstTs)
		mock.ExpectQuery("SELECT (.+) FROM user_").
			WithArgs("user-123").
			WillReturnRows(rows)

		_, err := AuthHook(mock)(ctx)
		assertTwirpError(t, err, twirp.Unauthenticated, "invalid auth")
	})

	t.Run("auth hook with user not verified disabled should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		ctx = auth.ContextSet(ctx, "api-key", "user-123")
		refID := refid.Must(model.NewUserRefID())

		rows := pgxmock.NewRows(
			[]string{"id", "ref_id", "email", "name", "api_access", "verified", "created", "last_modified"}).
			AddRow(1, refID, "user@example.com", "user", true, false, tstTs, tstTs)
		mock.ExpectQuery("SELECT (.+) FROM user_").
			WithArgs("user-123").
			WillReturnRows(rows)

		_, err := AuthHook(mock)(ctx)
		assertTwirpError(t, err, twirp.Unauthenticated, "account not verified")
	})

	t.Run("auth hook should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		ctx = auth.ContextSet(ctx, "api-key", "user-123")
		refID := refid.Must(model.NewUserRefID())

		rows := pgxmock.NewRows(
			[]string{"id", "ref_id", "email", "name", "api_access", "verified", "created", "last_modified"}).
			AddRow(1, refID, "user@example.com", "user", true, true, tstTs, tstTs)
		mock.ExpectQuery("SELECT (.+) FROM user_").
			WithArgs("user-123").
			WillReturnRows(rows)

		ctx, err := AuthHook(mock)(ctx)
		assert.NilError(t, err)
		_, ok := auth.ContextGet[*model.User](ctx, "user")
		assert.Check(t, ok, "user is a *mode.user")
	})
}
