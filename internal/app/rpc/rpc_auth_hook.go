package rpc

import (
	"context"

	"github.com/twitchtv/twirp"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/rpc/dto"
	"github.com/dropwhile/icbt/internal/app/service"
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
		user, errx := service.GetUserByApiKey(ctx, db, apiKey)
		if errx != nil {
			return nil, dto.ToTwirpError(errx)
		}
		if user == nil {
			return ctx, twirp.Unauthenticated.Error("invalid auth")
		}

		// the above query checks to ensure user.apikey is true as well,
		// but double check to make sure (in case sql above changes), as
		// this is a cheap local comparison anyway
		if !user.ApiAccess {
			return ctx, twirp.Unauthenticated.Error("invalid auth")
		}

		if !user.Verified {
			return ctx, twirp.Unauthenticated.Error("account not verified")
		}

		ctx = auth.ContextSet(ctx, "auth", true)
		ctx = auth.ContextSet(ctx, "user", user)
		// do any authorization here if needed in the future...
		// in a request routed
		// ref: https://github.com/twitchtv/twirp/issues/90#issuecomment-373108190
		return ctx, nil
	}
}
