package rpc

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/twitchtv/twirp"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
)

func AuthHook(db model.PgxHandle) func(context.Context) (context.Context, error) {
	return func(ctx context.Context) (context.Context, error) {
		// get apikey from context
		apiKey, ok := auth.ContextGet[string](ctx, "apikey")
		if !ok {
			// auth missing: either middleware wasn't run, or it wasnt in the req
			return ctx, twirp.Unauthenticated.Error("invalid auth")
		}

		// lookup user
		user, err := model.GetUserByApiKey(ctx, db, apiKey)
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return ctx, twirp.Unauthenticated.Error("invalid auth")
		case err != nil:
			return ctx, twirp.InternalError("db error")
		case user == nil:
			return ctx, twirp.Unauthenticated.Error("invalid auth")
		}

		ctx = auth.ContextSet(ctx, "auth", true)
		ctx = auth.ContextSet(ctx, "user", user)
		// do any authorization here if needed in the future...
		// in a request routed
		// ref: https://github.com/twitchtv/twirp/issues/90#issuecomment-373108190
		return ctx, nil
	}
}