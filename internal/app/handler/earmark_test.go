// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package handler

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-chi/chi/v5"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/app/service"
	"github.com/dropwhile/icanbringthat/internal/errs"
	"github.com/dropwhile/icanbringthat/internal/middleware/auth"
	"github.com/dropwhile/icanbringthat/internal/util"
)

func TestHandler_Earmark_Create(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		ID:           1,
		RefID:        util.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      ts,
		LastModified: ts,
	}
	event := &model.Event{
		ID:           1,
		RefID:        util.Must(model.NewEventRefID()),
		UserID:       user.ID,
		Name:         "event",
		Description:  "description",
		Archived:     false,
		StartTime:    ts,
		StartTimeTz:  util.Must(service.ParseTimeZone("Etc/UTC")),
		Created:      ts,
		LastModified: ts,
	}
	eventItem := &model.EventItem{
		ID:           2,
		RefID:        util.Must(model.NewEventItemRefID()),
		EventID:      event.ID,
		Description:  "eventitem",
		Created:      ts,
		LastModified: ts,
	}
	earmark := &model.Earmark{
		ID:           3,
		RefID:        util.Must(model.NewEarmarkRefID()),
		EventItemID:  eventItem.ID,
		UserID:       user.ID,
		Note:         "nothing",
		Created:      ts,
		LastModified: ts,
	}

	t.Run("create earmark should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		note := "some note"

		mock.EXPECT().
			GetEvent(ctx, event.RefID).
			Return(event, nil)
		mock.EXPECT().
			GetEventItem(ctx, eventItem.RefID).
			Return(eventItem, nil)
		mock.EXPECT().
			NewEarmark(ctx, user, eventItem.ID, note).
			Return(earmark, nil)

		data := url.Values{"note": {note}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/earmark", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", event.RefID.String())
		req.SetPathValue("iRefID", eventItem.RefID.String())
		rr := httptest.NewRecorder()
		handler.EarmarkCreate(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusOK)
		// we make sure that all expectations were met
	})

	t.Run("create earmark bad event refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		data := url.Values{"note": {"some note"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/earmark", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", "hodor")
		req.SetPathValue("iRefID", eventItem.RefID.String())
		rr := httptest.NewRecorder()
		handler.EarmarkCreate(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})

	t.Run("create earmark wrong type event refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		data := url.Values{"note": {"some note"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/earmark", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", eventItem.RefID.String())
		req.SetPathValue("iRefID", eventItem.RefID.String())
		rr := httptest.NewRecorder()
		handler.EarmarkCreate(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})

	t.Run("create earmark bad event item id should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		data := url.Values{"note": {"some note"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/earmark", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", event.RefID.String())
		req.SetPathValue("iRefID", "hodor")
		rr := httptest.NewRecorder()
		handler.EarmarkCreate(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})

	t.Run("create earmark wrong type event item id should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		data := url.Values{"note": {"some note"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/earmark", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", event.RefID.String())
		req.SetPathValue("iRefID", event.RefID.String())
		rr := httptest.NewRecorder()
		handler.EarmarkCreate(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})

	t.Run("create earmark missing event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		mock.EXPECT().
			GetEvent(ctx, event.RefID).
			Return(nil, errs.NotFound.Error("not found"))

		data := url.Values{"note": {"some note"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/earmark", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", event.RefID.String())
		req.SetPathValue("iRefID", eventItem.RefID.String())
		rr := httptest.NewRecorder()
		handler.EarmarkCreate(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})

	t.Run("create earmark missing event item should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		mock.EXPECT().
			GetEvent(ctx, event.RefID).
			Return(event, nil)
		mock.EXPECT().
			GetEventItem(ctx, eventItem.RefID).
			Return(nil, errs.NotFound.Error("not found"))

		data := url.Values{"note": {"some note"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/earmark", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", event.RefID.String())
		req.SetPathValue("iRefID", eventItem.RefID.String())
		rr := httptest.NewRecorder()
		handler.EarmarkCreate(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})

	t.Run("create earmark eventitem not matching event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		mock.EXPECT().
			GetEvent(ctx, event.RefID).
			Return(event, nil)
		mock.EXPECT().
			GetEventItem(ctx, eventItem.RefID).
			Return(
				&model.EventItem{
					ID:           eventItem.ID,
					RefID:        eventItem.RefID,
					EventID:      33,
					Description:  eventItem.Description,
					Created:      ts,
					LastModified: ts,
				},
				nil,
			)

		data := url.Values{"note": {"some note"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/earmark", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", event.RefID.String())
		req.SetPathValue("iRefID", eventItem.RefID.String())
		rr := httptest.NewRecorder()
		handler.EarmarkCreate(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})

	t.Run("create earmark user not verified should fail", func(t *testing.T) {
		t.Parallel()

		user := &model.User{
			ID:           33,
			RefID:        util.Must(model.NewUserRefID()),
			Email:        "user@example.com",
			Name:         "user",
			PWHash:       []byte("00x00"),
			Verified:     false,
			Created:      ts,
			LastModified: ts,
		}

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		note := "some note"

		mock.EXPECT().
			GetEvent(ctx, event.RefID).
			Return(event, nil)
		mock.EXPECT().
			GetEventItem(ctx, eventItem.RefID).
			Return(eventItem, nil)
		mock.EXPECT().
			NewEarmark(ctx, user, eventItem.ID, note).
			Return(nil, errs.PermissionDenied.Error("user not verified"))

		data := url.Values{"note": {note}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/earmark", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", event.RefID.String())
		req.SetPathValue("iRefID", eventItem.RefID.String())
		rr := httptest.NewRecorder()
		handler.EarmarkCreate(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
	})

	t.Run("create earmark already exists should fail", func(t *testing.T) {
		t.Parallel()

		user := &model.User{
			ID:           33,
			RefID:        util.Must(model.NewUserRefID()),
			Email:        "user@example.com",
			Name:         "user",
			PWHash:       []byte("00x00"),
			Verified:     false,
			Created:      ts,
			LastModified: ts,
		}

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		note := "some note"

		mock.EXPECT().
			GetEvent(ctx, event.RefID).
			Return(event, nil)
		mock.EXPECT().
			GetEventItem(ctx, eventItem.RefID).
			Return(eventItem, nil)
		mock.EXPECT().
			NewEarmark(ctx, user, eventItem.ID, note).
			Return(nil, errs.AlreadyExists.Error("already earmarked - access denied"))

		data := url.Values{"note": {note}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/earmark", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", event.RefID.String())
		req.SetPathValue("iRefID", eventItem.RefID.String())
		rr := httptest.NewRecorder()
		handler.EarmarkCreate(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
	})
}

func TestHandler_Earmark_Delete(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		ID:           1,
		RefID:        util.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      ts,
		LastModified: ts,
	}

	t.Run("delete earmark should succeed", func(t *testing.T) {
		t.Parallel()

		earmark := &model.Earmark{
			ID:           1,
			RefID:        util.Must(model.NewEarmarkRefID()),
			EventItemID:  1,
			UserID:       user.ID,
			Note:         "nothing",
			Created:      ts,
			LastModified: ts,
		}

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		mock.EXPECT().
			DeleteEarmarkByRefID(ctx, user.ID, earmark.RefID).
			Return(nil)

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/earmark", nil)
		req.SetPathValue("mRefID", earmark.RefID.String())
		rr := httptest.NewRecorder()
		handler.EarmarkDelete(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusOK)
		// we make sure that all expectations were met
	})

	t.Run("delete earmark permission denied should fail", func(t *testing.T) {
		t.Parallel()

		earmark := &model.Earmark{
			ID:           1,
			RefID:        util.Must(model.NewEarmarkRefID()),
			EventItemID:  1,
			UserID:       user.ID,
			Note:         "nothing",
			Created:      ts,
			LastModified: ts,
		}

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		mock.EXPECT().
			DeleteEarmarkByRefID(ctx, user.ID, earmark.RefID).
			Return(errs.PermissionDenied.Error("event is archived"))

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/earmark", nil)
		req.SetPathValue("mRefID", earmark.RefID.String())
		rr := httptest.NewRecorder()
		handler.EarmarkDelete(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
	})

	t.Run("delete earmark missing refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/earmark", nil)
		rr := httptest.NewRecorder()
		handler.EarmarkDelete(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})

	t.Run("delete earmark bad refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/earmark", nil)
		req.SetPathValue("mRefID", "hodor")
		rr := httptest.NewRecorder()
		handler.EarmarkDelete(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})

	t.Run("delete earmark not found should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		refID := util.Must(model.NewEarmarkRefID())

		mock.EXPECT().
			DeleteEarmarkByRefID(ctx, user.ID, refID).
			Return(errs.NotFound.Error("earmark not found"))

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/earmark", nil)
		req.SetPathValue("mRefID", refID.String())
		rr := httptest.NewRecorder()
		handler.EarmarkDelete(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})

	t.Run("delete earmark refid wrong type should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/earmark", nil)
		req.SetPathValue("mRefID", util.Must(model.NewEventRefID()).String())
		rr := httptest.NewRecorder()
		handler.EarmarkDelete(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})
}
