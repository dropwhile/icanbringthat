// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dropwhile/assert"
	"github.com/go-chi/chi/v5"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/app/service"
	"github.com/dropwhile/icanbringthat/internal/errs"
	"github.com/dropwhile/icanbringthat/internal/middleware/auth"
	"github.com/dropwhile/icanbringthat/internal/util"
)

func TestHandler_Favorite_Delete(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		ID:           1,
		RefID:        util.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Created:      ts,
		LastModified: ts,
	}
	event := &model.Event{
		ID:           1,
		RefID:        util.Must(model.NewEventRefID()),
		UserID:       user.ID,
		Name:         "event",
		Description:  "description",
		StartTime:    ts,
		StartTimeTz:  util.Must(service.ParseTimeZone("Etc/UTC")),
		Created:      ts,
		LastModified: ts,
	}

	t.Run("delete favorite", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		mock.EXPECT().
			RemoveFavorite(ctx, user.ID, event.RefID).
			Return(nil)

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/favorite", nil)
		req.SetPathValue("eRefID", event.RefID.String())
		rr := httptest.NewRecorder()
		handler.FavoriteDelete(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusOK)
		// we make sure that all expectations were met
	})

	t.Run("delete favorite bad event refid", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/favorite", nil)
		req.SetPathValue("eRefID", "hodor")
		rr := httptest.NewRecorder()
		handler.FavoriteDelete(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})

	t.Run("delete favorite event not found", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		refID := util.Must(model.NewEventRefID())

		mock.EXPECT().
			RemoveFavorite(ctx, user.ID, refID).
			Return(errs.NotFound.Error("event not found"))

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/favorite", nil)
		req.SetPathValue("eRefID", refID.String())
		rr := httptest.NewRecorder()
		handler.FavoriteDelete(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})

	t.Run("delete favorite event refid wrong type", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		refID := util.Must(model.NewEventItemRefID())

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/favorite", nil)
		req.SetPathValue("eRefID", refID.String())
		rr := httptest.NewRecorder()
		handler.FavoriteDelete(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})
}

func TestHandler_Favorite_Add(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		ID:           1,
		RefID:        util.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Created:      ts,
		LastModified: ts,
	}
	event := &model.Event{
		ID:           1,
		RefID:        util.Must(model.NewEventRefID()),
		UserID:       2,
		Name:         "event",
		Description:  "description",
		StartTime:    ts,
		StartTimeTz:  util.Must(service.ParseTimeZone("Etc/UTC")),
		Created:      ts,
		LastModified: ts,
	}

	t.Run("create favorite", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		mock.EXPECT().
			AddFavorite(ctx, user.ID, event.RefID).
			Return(event, nil)

		req, _ := http.NewRequestWithContext(ctx, "PUT", "http://example.com/favorite", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", event.RefID.String())
		rr := httptest.NewRecorder()
		handler.FavoriteAdd(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"),
			fmt.Sprintf("/events/%s", event.RefID.String()),
			"handler returned wrong redirect")
		// we make sure that all expectations were met
	})

	t.Run("create favorite bad event refid", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		req, _ := http.NewRequestWithContext(ctx, "PUT", "http://example.com/favorite", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", "hodor")
		rr := httptest.NewRecorder()
		handler.FavoriteAdd(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})

	t.Run("create favorite wrong type event refid", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		req, _ := http.NewRequestWithContext(ctx, "PUT", "http://example.com/favorite", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", util.Must(model.NewEventItemRefID()).String())
		rr := httptest.NewRecorder()
		handler.FavoriteAdd(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})

	t.Run("create favorite same user as event", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		mock.EXPECT().
			AddFavorite(ctx, user.ID, event.RefID).
			Return(nil, errs.PermissionDenied.Error("can't favorite own event"))

		req, _ := http.NewRequestWithContext(ctx, "PUT", "http://example.com/favorite", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", event.RefID.String())
		rr := httptest.NewRecorder()
		handler.FavoriteAdd(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
	})

	t.Run("create favorite already exists", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		mock.EXPECT().
			AddFavorite(ctx, user.ID, event.RefID).
			Return(nil, errs.AlreadyExists.Error("favorite already exists"))

		req, _ := http.NewRequestWithContext(ctx, "PUT", "http://example.com/favorite", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", event.RefID.String())
		rr := httptest.NewRecorder()
		handler.FavoriteAdd(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
	})
}
