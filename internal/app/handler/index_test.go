// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package handler

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dropwhile/assert"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/middleware/auth"
	"github.com/dropwhile/icanbringthat/internal/util"
)

func TestHandler_ShowIndex_LoggedOut(t *testing.T) {
	t.Parallel()

	t.Run("logged out", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/notification", nil)
		rr := httptest.NewRecorder()
		handler.IndexShow(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/login",
			"handler returned wrong redirect")
		// we make sure that all expectations were met
	})

	t.Run("logged in", func(t *testing.T) {
		t.Parallel()

		user := &model.User{
			ID:       1,
			RefID:    util.Must(model.NewUserRefID()),
			Email:    "user@example.com",
			Name:     "user",
			PWHash:   []byte("00x00"),
			Verified: true,
		}

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/notification", nil)
		rr := httptest.NewRecorder()
		handler.IndexShow(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/dashboard",
			"handler returned wrong redirect")
		// we make sure that all expectations were met
	})
}
