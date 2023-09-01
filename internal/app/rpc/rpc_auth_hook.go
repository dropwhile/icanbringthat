// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rpc

import (
	"context"

	"github.com/twitchtv/twirp"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/errs"
	"github.com/dropwhile/icanbringthat/internal/middleware/auth"
)

type GetUserProvider interface {
	GetUserByApiKey(context.Context, string) (*model.User, errs.Error)
}

func AuthHook(up GetUserProvider) func(context.Context) (context.Context, error) {
	return func(ctx context.Context) (context.Context, error) {
		// get apikey from context
		apiKey, ok := auth.ContextGet[string](ctx, "api-key")
		if !ok {
			// auth missing: either middleware wasn't run, or it wasnt in the req
			return ctx, twirp.Unauthenticated.Error("invalid auth")
		}

		// lookup user
		user, errx := up.GetUserByApiKey(ctx, apiKey)
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
