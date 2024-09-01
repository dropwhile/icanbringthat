// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rpc

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"gotest.tools/v3/assert"
)

func dummyHandler(w http.ResponseWriter, r *http.Request) {}

func TestRpc_RequireApiKey(t *testing.T) {
	t.Parallel()

	t.Run("auth hook with no api-key should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, mock := NewTestServer(t)

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/rpc", &bytes.Buffer{})
		rr := httptest.NewRecorder()
		RequireApiKey(mock)(http.HandlerFunc(dummyHandler)).ServeHTTP(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		assert.Equal(t, rr.Code, http.StatusUnauthorized,
			"handler returned wrong status code: got %d expected %d",
			rr.Code, http.StatusUnauthorized)
	})
	/*
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
			refID := util.Must(model.NewUserRefID())

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
			refID := util.Must(model.NewUserRefID())

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
			refID := util.Must(model.NewUserRefID())

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
	*/
}
