package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/dropwhile/refid/v2"
	"github.com/go-chi/chi/v5"
	"github.com/samber/mo"
	"go.uber.org/mock/gomock"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/internal/errs"
	"github.com/dropwhile/icbt/internal/middleware/auth"
	"github.com/dropwhile/icbt/internal/util"
)

func TestHandler_Event_Create(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		ID:           1,
		RefID:        refid.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      ts,
		LastModified: ts,
	}
	event := &model.Event{
		ID:           1,
		RefID:        refid.Must(model.NewEventRefID()),
		UserID:       user.ID,
		Name:         "event",
		Description:  "description",
		Archived:     false,
		StartTime:    ts,
		StartTimeTz:  util.Must(service.ParseTimeZone("Etc/UTC")),
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
			CreateEvent(
				ctx, user, event.Name, event.Description,
				gomock.AssignableToTypeOf(time.Time{}),
				event.StartTimeTz.String(),
			).
			Return(event, nil)

		data := url.Values{
			"name":        {event.Name},
			"description": {event.Description},
			"when":        {event.StartTime.Format("2006-01-02T15:04")},
			"timezone":    {event.StartTimeTz.String()},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/event", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.CreateEvent(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Assert(t,
			strings.HasPrefix(rr.Header().Get("location"), "/events/"),
			"handler returned wrong redirect: expected prefix %s didnt match %s",
			"/events/", rr.Header().Get("location"),
		)
		// we make sure that all expectations were met
	})

	t.Run("create missing form value name", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		data := url.Values{
			"description": {event.Description},
			"when":        {event.StartTime.Format("2006-01-02T15:04")},
			"timezone":    {event.StartTimeTz.String()},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/event", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.CreateEvent(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
	})

	t.Run("create missing form value description", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		data := url.Values{
			"name":     {event.Name},
			"when":     {event.StartTime.Format("2006-01-02T15:04")},
			"timezone": {event.StartTimeTz.String()},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/event", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.CreateEvent(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
	})

	t.Run("create missing form value when", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		data := url.Values{
			"name":        {event.Name},
			"description": {event.Description},
			"timezone":    {event.StartTimeTz.String()},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/event", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.CreateEvent(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
	})

	t.Run("create missing form value timezone", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		data := url.Values{
			"name":        {event.Name},
			"description": {event.Description},
			"when":        {event.StartTime.Format("2006-01-02T15:04")},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/event", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.CreateEvent(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
	})

	t.Run("create bad timezone with default to utc should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		mock.EXPECT().
			CreateEvent(
				ctx, user, event.Name, event.Description,
				gomock.AssignableToTypeOf(time.Time{}),
				event.StartTimeTz.String(),
			).
			Return(event, nil)

		data := url.Values{
			"name":        {event.Name},
			"description": {event.Description},
			"when":        {event.StartTime.Format("2006-01-02T15:04")},
			"timezone":    {"morbin/time"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/event", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.CreateEvent(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Assert(t,
			strings.HasPrefix(rr.Header().Get("location"), "/events/"),
			"handler returned wrong redirect: expected prefix %s didnt match %s",
			"/events/", rr.Header().Get("location"),
		)
		// we make sure that all expectations were met
	})

	t.Run("create bad time", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		data := url.Values{
			"name":        {event.Name},
			"description": {event.Description},
			"when":        {"It's ho-ho-ho time!"},
			"timezone":    {event.StartTimeTz.String()},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/event", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.CreateEvent(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
	})

	t.Run("create user not verified should fail", func(t *testing.T) {
		t.Parallel()

		user := &model.User{
			ID:           1,
			RefID:        refid.Must(model.NewUserRefID()),
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

		mock.EXPECT().
			CreateEvent(ctx, user, event.Name, event.Description,
				gomock.AssignableToTypeOf(time.Time{}),
				event.StartTimeTz.String(),
			).
			Return(nil, errs.PermissionDenied.Error(
				"Account must be verified before event creation is allowed.",
			))

		data := url.Values{
			"name":        {event.Name},
			"description": {event.Description},
			"when":        {event.StartTime.Format("2006-01-02T15:04")},
			"timezone":    {event.StartTimeTz.String()},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/event", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.CreateEvent(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
	})
}

func TestHandler_Event_Update(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		ID:           1,
		RefID:        refid.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      ts,
		LastModified: ts,
	}
	event := &model.Event{
		ID:            1,
		RefID:         refid.Must(model.NewEventRefID()),
		UserID:        user.ID,
		Name:          "event",
		Description:   "description",
		Archived:      false,
		ItemSortOrder: []int{},
		StartTime:     ts,
		StartTimeTz:   util.Must(service.ParseTimeZone("Etc/UTC")),
		Created:       ts,
		LastModified:  ts,
	}

	t.Run("update event should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		euvs := &service.EventUpdateValues{
			StartTime: mo.Some(event.StartTime.Truncate(time.Minute).In(
				event.StartTimeTz.Location,
			)),
			Name:        mo.Some(event.Name),
			Description: mo.Some(event.Description),
			Tz:          mo.Some(event.StartTimeTz.String()),
		}

		mock.EXPECT().
			UpdateEvent(ctx, user.ID, event.RefID, euvs).
			Return(nil)

		data := url.Values{
			"name":        {event.Name},
			"description": {event.Description},
			"when":        {event.StartTime.Format("2006-01-02T15:04")},
			"timezone":    {event.StartTimeTz.String()},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/event", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateEvent(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Assert(t,
			strings.HasPrefix(rr.Header().Get("location"), "/events/"),
			"handler returned wrong redirect: expected prefix %s didnt match %s",
			"/events/", rr.Header().Get("location"),
		)
		// we make sure that all expectations were met
	})

	t.Run("update event bad refid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", refid.Must(model.NewEarmarkRefID()).String())

		data := url.Values{
			"name": {event.Name},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/event", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateEvent(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})

	t.Run("update missing event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		euvs := &service.EventUpdateValues{
			Name: mo.Some(event.Name),
		}

		mock.EXPECT().
			UpdateEvent(ctx, user.ID, event.RefID, euvs).
			Return(errs.NotFound.Error("event not found"))

		data := url.Values{
			"name": {event.Name},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/event", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateEvent(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})

	t.Run("update user.id not match event.userid should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		euvs := &service.EventUpdateValues{
			Name: mo.Some(event.Name),
		}

		mock.EXPECT().
			UpdateEvent(ctx, user.ID, event.RefID, euvs).
			Return(errs.PermissionDenied.Error("permission denied"))

		data := url.Values{
			"name": {event.Name},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/event", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateEvent(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
	})

	t.Run("update event archived should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		euvs := &service.EventUpdateValues{
			Name: mo.Some(event.Name),
		}

		mock.EXPECT().
			UpdateEvent(ctx, user.ID, event.RefID, euvs).
			Return(errs.PermissionDenied.Error("event is archived"))

		data := url.Values{
			"name": {event.Name},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/event", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateEvent(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
	})

	t.Run("update only name should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		euvs := &service.EventUpdateValues{
			Name: mo.Some(event.Name + "x"),
		}

		mock.EXPECT().
			UpdateEvent(ctx, user.ID, event.RefID, euvs).
			Return(nil)

		data := url.Values{
			"name": {event.Name + "x"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/event", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateEvent(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Assert(t,
			strings.HasPrefix(rr.Header().Get("location"), "/events/"),
			"handler returned wrong redirect: expected prefix %s didnt match %s",
			"/events/", rr.Header().Get("location"),
		)
		// we make sure that all expectations were met
	})

	t.Run("update only update description should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		euvs := &service.EventUpdateValues{
			Description: mo.Some(event.Description + "x"),
		}

		mock.EXPECT().
			UpdateEvent(ctx, user.ID, event.RefID, euvs).
			Return(nil)

		data := url.Values{
			"description": {event.Description + "x"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/event", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateEvent(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Assert(t,
			strings.HasPrefix(rr.Header().Get("location"), "/events/"),
			"handler returned wrong redirect: expected prefix %s didnt match %s",
			"/events/", rr.Header().Get("location"),
		)
		// we make sure that all expectations were met
	})

	t.Run("update with when and tz should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		euvs := &service.EventUpdateValues{
			StartTime: mo.Some(event.StartTime.Truncate(time.Minute).In(
				event.StartTimeTz.Location,
			)),
			Tz: mo.Some(event.StartTimeTz.String()),
		}

		mock.EXPECT().
			UpdateEvent(ctx, user.ID, event.RefID, euvs).
			Return(nil)

		data := url.Values{
			"when":     {event.StartTime.Format("2006-01-02T15:04")},
			"timezone": {event.StartTimeTz.String()},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/event", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateEvent(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Assert(t,
			strings.HasPrefix(rr.Header().Get("location"), "/events/"),
			"handler returned wrong redirect: expected prefix %s didnt match %s",
			"/events/", rr.Header().Get("location"),
		)
		// we make sure that all expectations were met
	})

	t.Run("update with only when should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		data := url.Values{
			"when": {event.StartTime.Format("2006-01-02T15:04")},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/event", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateEvent(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
	})

	t.Run("update with only tz should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		data := url.Values{
			"timezone": {event.StartTimeTz.String()},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/event", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateEvent(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
	})

	t.Run("update nothing should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		euvs := &service.EventUpdateValues{}

		mock.EXPECT().
			UpdateEvent(ctx, user.ID, event.RefID, euvs).
			Return(errs.InvalidArgument.Error("missing fields"))

		data := url.Values{
			"namexxx": {event.Name},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/event", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateEvent(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
	})

	t.Run("update bad tz should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		data := url.Values{
			"when":     {event.StartTime.Format("2006-01-02T15:04")},
			"timezone": {"morbin/time"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/event", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateEvent(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
	})

	t.Run("update bad time", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		data := url.Values{
			"when":     {"It's ho-ho-ho time!"},
			"timezone": {event.StartTimeTz.String()},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/event", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateEvent(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
	})
}

func TestHandler_Event_UpdateSorting(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		ID:           1,
		RefID:        refid.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      ts,
		LastModified: ts,
	}
	event := &model.Event{
		ID:            1,
		RefID:         refid.Must(model.NewEventRefID()),
		UserID:        user.ID,
		Name:          "event",
		Description:   "description",
		Archived:      false,
		ItemSortOrder: []int{1, 2, 3},
		StartTime:     ts,
		StartTimeTz:   util.Must(service.ParseTimeZone("Etc/UTC")),
		Created:       ts,
		LastModified:  ts,
	}

	t.Run("update", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		mock.EXPECT().
			UpdateEventItemSorting(ctx, user.ID, event.RefID, []int{1, 3, 2}).
			Return(event, nil)

		data := url.Values{
			"sortOrder": {"1", "3", "2"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/event", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateEventItemSorting(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusOK)
		// we make sure that all expectations were met
	})

	t.Run("update bad refid", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", refid.Must(model.NewEarmarkRefID()).String())

		data := url.Values{
			"sortOrder": {"1", "3", "2"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/event", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateEventItemSorting(rr, req)

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
		rctx.URLParams.Add("eRefID", event.RefID.String())

		mock.EXPECT().
			UpdateEventItemSorting(ctx, user.ID, event.RefID, []int{1, 3, 2}).
			Return(nil, errs.NotFound.Error("event not found"))

		data := url.Values{
			"sortOrder": {"1", "3", "2"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/event", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateEventItemSorting(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})

	t.Run("update user.id not match event.userid", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		mock.EXPECT().
			UpdateEventItemSorting(ctx, user.ID, event.RefID, []int{1, 3, 2}).
			Return(nil, errs.PermissionDenied.Error("permission denied"))

		data := url.Values{
			"sortOrder": {"1", "3", "2"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/event", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateEventItemSorting(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
	})

	t.Run("update archived event should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		mock.EXPECT().
			UpdateEventItemSorting(ctx, user.ID, event.RefID, []int{1, 3, 2}).
			Return(nil, errs.PermissionDenied.Error("event is archived"))

		data := url.Values{
			"sortOrder": {"1", "3", "2"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/event", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateEventItemSorting(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
	})

	t.Run("update update nothing", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		data := url.Values{
			"namexxx": {event.Name},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/event", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateEventItemSorting(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
	})

	t.Run("update bad values", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		data := url.Values{
			"sortOrder": {"a", "3", "2"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/event", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateEventItemSorting(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
	})

	t.Run("update bad values 2", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		data := url.Values{
			"sortOrder": {"a123"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/event", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateEventItemSorting(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
	})
}

func TestHandler_Event_Delete(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		ID:           1,
		RefID:        refid.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      ts,
		LastModified: ts,
	}
	event := &model.Event{
		ID:           1,
		RefID:        refid.Must(model.NewEventRefID()),
		UserID:       user.ID,
		Name:         "event",
		Description:  "description",
		Archived:     false,
		StartTime:    ts,
		StartTimeTz:  util.Must(service.ParseTimeZone("Etc/UTC")),
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
		rctx.URLParams.Add("eRefID", event.RefID.String())

		mock.EXPECT().
			DeleteEvent(ctx, user.ID, event.RefID).
			Return(nil)

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/event", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.DeleteEvent(rr, req)

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
		rctx.URLParams.Add("eRefID", refid.Must(model.NewEventItemRefID()).String())

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/event", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.DeleteEvent(rr, req)

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
		rctx.URLParams.Add("eRefID", event.RefID.String())

		mock.EXPECT().
			DeleteEvent(ctx, user.ID, event.RefID).
			Return(errs.NotFound.Error("event not found"))

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/event", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.DeleteEvent(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})

	t.Run("delete mismatch user id", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		mock.EXPECT().
			DeleteEvent(ctx, user.ID, event.RefID).
			Return(errs.PermissionDenied.Error("permission denied"))

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/event", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.DeleteEvent(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
	})
}
