package rpc

import (
	"context"
	"testing"

	"github.com/dropwhile/refid/v2"
	"github.com/twitchtv/twirp"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/errs"
	"github.com/dropwhile/icbt/internal/middleware/auth"
)

func TestRpc_AuthHook(t *testing.T) {
	t.Parallel()

	t.Run("auth hook with no api-key should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		_, mock := NewTestServer(t)

		_, err := AuthHook(mock)(ctx)
		errs.AssertError(t, err, twirp.Unauthenticated, "invalid auth")
	})

	t.Run("auth hook with api-key not finding user should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		_, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "api-key", "user-123")

		mock.EXPECT().
			GetUserByApiKey(ctx, "user-123").
			Return(nil, errs.NotFound.Error("user not found"))

		_, err := AuthHook(mock)(ctx)
		errs.AssertError(t, err, twirp.Unauthenticated, "invalid auth")
	})

	t.Run("auth hook with user ApiAccess disabled should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		_, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "api-key", "user-123")
		refID := refid.Must(model.NewUserRefID())

		mock.EXPECT().
			GetUserByApiKey(ctx, "user-123").
			Return(
				&model.User{
					ID:        1,
					RefID:     refID,
					ApiAccess: false,
					Verified:  true,
				}, nil,
			)

		_, err := AuthHook(mock)(ctx)
		errs.AssertError(t, err, twirp.Unauthenticated, "invalid auth")
	})

	t.Run("auth hook with user not verified disabled should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		_, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "api-key", "user-123")
		refID := refid.Must(model.NewUserRefID())

		mock.EXPECT().
			GetUserByApiKey(ctx, "user-123").
			Return(
				&model.User{
					ID:        1,
					RefID:     refID,
					ApiAccess: true,
					Verified:  false,
				}, nil,
			)

		_, err := AuthHook(mock)(ctx)
		errs.AssertError(t, err, twirp.Unauthenticated, "account not verified")
	})

	t.Run("auth hook should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		_, mock := NewTestServer(t)
		ctx = auth.ContextSet(ctx, "api-key", "user-123")
		refID := refid.Must(model.NewUserRefID())

		mock.EXPECT().
			GetUserByApiKey(ctx, "user-123").
			Return(
				&model.User{
					ID:        1,
					RefID:     refID,
					ApiAccess: true,
					Verified:  true,
				}, nil,
			)

		ctx, err := AuthHook(mock)(ctx)
		assert.NilError(t, err)
		_, ok := auth.ContextGet[*model.User](ctx, "user")
		assert.Check(t, ok, "user is a *mode.user")
	})
}
