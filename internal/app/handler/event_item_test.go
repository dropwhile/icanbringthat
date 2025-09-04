// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-chi/chi/v5"
	"go.uber.org/mock/gomock"
	"github.com/dropwhile/assert"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/app/service"
	"github.com/dropwhile/icanbringthat/internal/errs"
	"github.com/dropwhile/icanbringthat/internal/middleware/auth"
	"github.com/dropwhile/icanbringthat/internal/util"
)

func TestHandler_EventItem_Create(t *testing.T) {
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

	t.Run("create", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		mock.EXPECT().
			AddEventItem(ctx, user.ID, event.RefID, eventItem.Description).
			Return(eventItem, nil)

		data := url.Values{"description": {eventItem.Description}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/eventItem", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", event.RefID.String())
		rr := httptest.NewRecorder()
		handler.EventItemCreate(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"),
			fmt.Sprintf("/events/%s", event.RefID),
			"handler returned wrong redirect")
		// we make sure that all expectations were met
	})

	t.Run("create bad event refid", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		data := url.Values{"description": {"some description"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/eventItem", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", eventItem.RefID.String())
		rr := httptest.NewRecorder()
		handler.EventItemCreate(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})

	t.Run("create missing event", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		mock.EXPECT().
			AddEventItem(ctx, user.ID, event.RefID, eventItem.Description).
			Return(nil, errs.NotFound.Error("event not found"))

		data := url.Values{"description": {eventItem.Description}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/eventItem", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", event.RefID.String())
		rr := httptest.NewRecorder()
		handler.EventItemCreate(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})

	t.Run("create user not owner", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		mock.EXPECT().
			AddEventItem(ctx, user.ID, event.RefID, eventItem.Description).
			Return(nil, errs.PermissionDenied.Error("permission denied"))

		data := url.Values{"description": {eventItem.Description}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/eventItem", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", event.RefID.String())
		rr := httptest.NewRecorder()
		handler.EventItemCreate(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
	})

	t.Run("create with missing description param should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		data := url.Values{"descriptionxxx": {"some description"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/eventItem", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", event.RefID.String())
		rr := httptest.NewRecorder()
		handler.EventItemCreate(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
	})
}

func TestHandler_EventItem_Update(t *testing.T) {
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

	t.Run("update", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		description := "new description"

		mock.EXPECT().
			GetEvent(ctx, event.RefID).
			Return(event, nil)
		mock.EXPECT().
			UpdateEventItem(ctx, user.ID, eventItem.RefID, description, gomock.Any()).
			Return(eventItem, nil)

		data := url.Values{"description": {description}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/eventItem", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", event.RefID.String())
		req.SetPathValue("iRefID", eventItem.RefID.String())
		rr := httptest.NewRecorder()
		handler.EventItemUpdate(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusOK)
		// we make sure that all expectations were met
	})

	t.Run("update bad event refid", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		data := url.Values{"description": {"new description"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/eventItem", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", eventItem.RefID.String())
		req.SetPathValue("iRefID", eventItem.RefID.String())
		rr := httptest.NewRecorder()
		handler.EventItemUpdate(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})

	t.Run("update missing event", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		mock.EXPECT().
			GetEvent(ctx, event.RefID).
			Return(nil, errs.NotFound.Error("event not found"))

		data := url.Values{"description": {"new description"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/eventItem", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", event.RefID.String())
		req.SetPathValue("iRefID", eventItem.RefID.String())
		rr := httptest.NewRecorder()
		handler.EventItemUpdate(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})

	t.Run("update missing event item", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		description := "new description"

		mock.EXPECT().
			GetEvent(ctx, event.RefID).
			Return(event, nil)
		mock.EXPECT().
			UpdateEventItem(ctx, user.ID, eventItem.RefID, description, gomock.Any()).
			Return(nil, errs.NotFound.Error("event-item not found"))

		data := url.Values{"description": {description}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/eventItem", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", event.RefID.String())
		req.SetPathValue("iRefID", eventItem.RefID.String())
		rr := httptest.NewRecorder()
		handler.EventItemUpdate(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})

	t.Run("update event owner not match", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		description := "new description"

		mock.EXPECT().
			GetEvent(ctx, event.RefID).
			Return(event, nil)
		mock.EXPECT().
			UpdateEventItem(ctx, user.ID, eventItem.RefID, description, gomock.Any()).
			Return(nil, errs.PermissionDenied.Error("not event owner"))

		data := url.Values{"description": {description}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/eventItem", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", event.RefID.String())
		req.SetPathValue("iRefID", eventItem.RefID.String())
		rr := httptest.NewRecorder()
		handler.EventItemUpdate(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
	})

	t.Run("update event archived", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		description := "new description"

		mock.EXPECT().
			GetEvent(ctx, event.RefID).
			Return(event, nil)
		mock.EXPECT().
			UpdateEventItem(ctx, user.ID, eventItem.RefID, description, gomock.Any()).
			Return(nil, errs.PermissionDenied.Error("event is archived"))

		data := url.Values{"description": {description}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/eventItem", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", event.RefID.String())
		req.SetPathValue("iRefID", eventItem.RefID.String())
		rr := httptest.NewRecorder()
		handler.EventItemUpdate(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
	})

	t.Run("update already earmarked by other user", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		description := "new description"

		mock.EXPECT().
			GetEvent(ctx, event.RefID).
			Return(event, nil)
		mock.EXPECT().
			UpdateEventItem(ctx, user.ID, eventItem.RefID, description, gomock.Any()).
			Return(nil, errs.PermissionDenied.Error("earmarked by other user"))

		data := url.Values{"description": {description}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/eventItem", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", event.RefID.String())
		req.SetPathValue("iRefID", eventItem.RefID.String())
		rr := httptest.NewRecorder()
		handler.EventItemUpdate(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
	})

	t.Run("update missing form data", func(t *testing.T) {
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

		data := url.Values{"descriptionxxxx": {"new description"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/eventItem", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", event.RefID.String())
		req.SetPathValue("iRefID", eventItem.RefID.String())
		rr := httptest.NewRecorder()
		handler.EventItemUpdate(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
	})

	t.Run("update eventitem not matching event", func(t *testing.T) {
		t.Parallel()

		eventItem := &model.EventItem{
			ID:          33,
			RefID:       eventItem.RefID,
			EventID:     eventItem.EventID + 1,
			Description: "some desc",
		}

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		description := "new description"

		mock.EXPECT().
			GetEvent(ctx, event.RefID).
			Return(event, nil)
		mock.EXPECT().
			UpdateEventItem(ctx, user.ID, eventItem.RefID, description, gomock.Any()).
			DoAndReturn(
				func(
					_ctx context.Context, _userID int,
					_eventItemRefID model.EventItemRefID,
					_description string,
					f func(*model.EventItem) bool,
				) (*model.EventItem, errs.Error) {
					if f(eventItem) {
						return nil, errs.FailedPrecondition.Error("extra checks failed")
					}
					return eventItem, nil
				},
			)

		data := url.Values{"description": {description}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/eventItem", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", event.RefID.String())
		req.SetPathValue("iRefID", eventItem.RefID.String())
		rr := httptest.NewRecorder()
		handler.EventItemUpdate(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})
}

func TestHandler_EventItem_Delete(t *testing.T) {
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

	t.Run("delete", func(t *testing.T) {
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
			RemoveEventItem(ctx, user.ID, eventItem.RefID, gomock.Any()).
			Return(nil)

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/eventItem", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", event.RefID.String())
		req.SetPathValue("iRefID", eventItem.RefID.String())
		rr := httptest.NewRecorder()
		handler.EventItemDelete(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusOK)
		// we make sure that all expectations were met
	})

	t.Run("delete bad event refid", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/eventItem", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", eventItem.RefID.String())
		req.SetPathValue("iRefID", eventItem.RefID.String())
		rr := httptest.NewRecorder()
		handler.EventItemDelete(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})

	t.Run("delete bad event item refid", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/eventItem", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", event.RefID.String())
		req.SetPathValue("iRefID", event.RefID.String())
		rr := httptest.NewRecorder()
		handler.EventItemDelete(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})

	t.Run("delete missing event", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		mock.EXPECT().
			GetEvent(ctx, event.RefID).
			Return(nil, errs.NotFound.Error("event not found"))

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/eventItem", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", event.RefID.String())
		req.SetPathValue("iRefID", eventItem.RefID.String())
		rr := httptest.NewRecorder()
		handler.EventItemDelete(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})

	t.Run("delete missing event item", func(t *testing.T) {
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
			RemoveEventItem(ctx, user.ID, eventItem.RefID, gomock.Any()).
			Return(errs.NotFound.Error("event-item not found"))

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/eventItem", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", event.RefID.String())
		req.SetPathValue("iRefID", eventItem.RefID.String())
		rr := httptest.NewRecorder()
		handler.EventItemDelete(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})

	t.Run("delete user not owner", func(t *testing.T) {
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
			RemoveEventItem(ctx, user.ID, eventItem.RefID, gomock.Any()).
			Return(errs.PermissionDenied.Error("not event owner"))

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/eventItem", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", event.RefID.String())
		req.SetPathValue("iRefID", eventItem.RefID.String())
		rr := httptest.NewRecorder()
		handler.EventItemDelete(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
	})

	t.Run("delete event item not related to supplied event", func(t *testing.T) {
		t.Parallel()

		eventItem := &model.EventItem{
			ID:          33,
			RefID:       eventItem.RefID,
			EventID:     eventItem.EventID + 1,
			Description: "some desc",
		}

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
			RemoveEventItem(ctx, user.ID, eventItem.RefID, gomock.Any()).
			DoAndReturn(
				func(
					_ctx context.Context,
					_userID int,
					_eventItemRefID model.EventItemRefID,
					f func(*model.EventItem) bool,
				) errs.Error {
					if f(eventItem) {
						return errs.FailedPrecondition.Error("extra checks failed")
					}
					return nil
				},
			)

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/eventItem", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", event.RefID.String())
		req.SetPathValue("iRefID", eventItem.RefID.String())
		rr := httptest.NewRecorder()
		handler.EventItemDelete(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})

	t.Run("delete event archived", func(t *testing.T) {
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
			RemoveEventItem(ctx, user.ID, eventItem.RefID, gomock.Any()).
			Return(errs.PermissionDenied.Error("event is archived"))

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/eventItem", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("eRefID", event.RefID.String())
		req.SetPathValue("iRefID", eventItem.RefID.String())
		rr := httptest.NewRecorder()
		handler.EventItemDelete(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
	})
}
