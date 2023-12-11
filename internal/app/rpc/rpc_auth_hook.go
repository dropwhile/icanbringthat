package rpc

import (
	"context"

	"github.com/twitchtv/twirp"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/internal/errs"
	"github.com/dropwhile/icbt/internal/middleware/auth"
)

func AuthHook(db model.PgxHandle) func(context.Context) (context.Context, error) {
	return func(ctx context.Context) (context.Context, error) {
		// get apikey from context
		apiKey, ok := auth.ContextGet[string](ctx, "api-key")
		if !ok {
			// auth missing: either middleware wasn't run, or it wasnt in the req
			return ctx, twirp.Unauthenticated.Error("invalid auth")
		}

		// lookup user
		user, errx := service.GetUserByApiKey(ctx, db, apiKey)
		switch {
		case errx != nil && errx.Code() != errs.NotFound:
			return ctx, twirp.Internal.Error("db error")
		case errx != nil || user == nil:
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

		ctx = auth.ContextSet(ctx, "user", user)
		// do any authorization here if needed in the future...
		// in a request routed
		// ref: https://github.com/twitchtv/twirp/issues/90#issuecomment-373108190
		return ctx, nil
	}
}
