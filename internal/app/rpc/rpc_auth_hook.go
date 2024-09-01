// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rpc

import (
	"context"
	"errors"
	"io"
	"net/http"

	"connectrpc.com/connect"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/errs"
	"github.com/dropwhile/icanbringthat/internal/middleware/auth"
)

type GetUserProvider interface {
	GetUserByApiKey(context.Context, string) (*model.User, errs.Error)
}

/*
func NewAuthInterceptor(up GetUserProvider) connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			// ensure context isn't cancelled
			if err := ctx.Err(); err != nil {
				return nil, err
			}

			// get apikey from context
			apiKey, ok := auth.ContextGet[string](ctx, "api-key")
			if !ok {
				// auth missing: either middleware wasn't run, or it wasnt in the req
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid auth"))
			}

			// lookup user
			user, errx := up.GetUserByApiKey(ctx, apiKey)
			switch {
			case errx != nil && errx.Code() != errs.NotFound:
				return nil, connect.NewError(connect.CodeInternal, errors.New("db error"))
			case errx != nil || user == nil:
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid auth"))
			}

			// the above query checks to ensure user.apikey is true as well,
			// but double check to make sure (in case sql above changes), as
			// this is a cheap local comparison anyway
			if !user.ApiAccess {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid auth"))
			}

			if !user.Verified {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("account not verified"))
			}

			ctx = auth.ContextSet(ctx, "user", user)
			// do any authorization here if needed in the future...
			// in a request routed
			// ref: https://github.com/twitchtv/twirp/issues/90#issuecomment-373108190
			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
*/

func AuthError(msg string) *connect.Error {
	return connect.NewError(connect.CodeUnauthenticated, errors.New(msg))
}

func InternalError(msg string) *connect.Error {
	return connect.NewError(connect.CodeInternal, errors.New(msg))
}

func RequireApiKey(up GetUserProvider, opts ...connect.HandlerOption) func(http.Handler) http.Handler {
	errW := connect.NewErrorWriter(opts...)

	// send a protocol appropriate error response
	writeErr := func(w http.ResponseWriter, r *http.Request, err *connect.Error) {
		defer r.Body.Close()
		defer io.Copy(io.Discard, r.Body)
		if errW.IsSupported(r) {
			// Send a protocol-appropriate error to RPC clients, so that they receive
			// the right code, message, and any metadata or error details.
			_ = errW.Write(w, r, err)
		} else {
			// Send an error to non-RPC clients.
			if err.Code() == connect.CodeInternal {
				http.Error(w, err.Message(), http.StatusInternalServerError)
			} else {
				http.Error(w, err.Message(), http.StatusUnauthorized)
			}
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// get apikey from context
			apiKey, ok := auth.ContextGet[string](ctx, "api-key")
			if !ok {
				// auth missing: either middleware wasn't run, or it wasnt in the req
				writeErr(w, r, AuthError("invalid auth"))
				return
			}

			// lookup user
			user, errx := up.GetUserByApiKey(ctx, apiKey)
			switch {
			case errx != nil && errx.Code() != errs.NotFound:
				writeErr(w, r, InternalError("db error"))
				return
			case errx != nil || user == nil:
				writeErr(w, r, AuthError("invalid auth"))
				return
			}

			// the above query checks to ensure user.apikey is true as well,
			// but double check to make sure (in case sql above changes), as
			// this is a cheap local comparison anyway
			if !user.ApiAccess {
				writeErr(w, r, AuthError("invalid auth"))
				return
			}

			if !user.Verified {
				writeErr(w, r, AuthError("invalid auth"))
				return
			}

			ctx = auth.ContextSet(ctx, "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
