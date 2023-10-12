package xhandler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/dropwhile/refid"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
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
		Created:      ts,
		LastModified: ts,
	}
	event := &model.Event{
		ID:           1,
		RefID:        refid.Must(model.NewEventRefID()),
		UserID:       user.ID,
		Name:         "event",
		Description:  "description",
		StartTime:    ts,
		StartTimeTz:  model.Must(model.ParseTimeZone("Etc/UTC")),
		Created:      ts,
		LastModified: ts,
	}

	eventColumns := []string{
		"id", "ref_id", "user_id", "name", "description",
		"start_time", "start_time_tz", "created", "last_modified",
	}

	t.Run("create", func(t *testing.T) {
		t.Parallel()

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.ID, event.RefID, event.UserID, event.Name, event.Description,
				event.StartTime, event.StartTimeTz, ts, ts,
			)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		mock.ExpectBegin()
		// refid as anyarg because new refid is created on call to create
		mock.ExpectQuery("INSERT INTO event_ ").
			WithArgs(util.NewPgxNamedArgsMatcher(pgx.NamedArgs{
				"refID":       model.EventRefIDMatcher{},
				"userID":      event.UserID,
				"name":        event.Name,
				"description": event.Description,
				"startTime":   CloseTimeMatcher{event.StartTime, time.Minute},
				"startTimeTz": event.StartTimeTz,
			})).
			WillReturnRows(eventRows)
		mock.ExpectCommit()
		mock.ExpectRollback()

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
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create missing form value name", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
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
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create missing form value description", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
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
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create missing form value when", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
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
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create missing form value timezone", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
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
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create bad timezone", func(t *testing.T) {
		t.Parallel()

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.ID, event.RefID, event.UserID, event.Name, event.Description,
				event.StartTime, event.StartTimeTz, ts, ts,
			)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		mock.ExpectBegin()
		// refid as anyarg because new refid is created on call to create
		mock.ExpectQuery("INSERT INTO event_ ").
			WithArgs(util.NewPgxNamedArgsMatcher(pgx.NamedArgs{
				"refID":       model.EventRefIDMatcher{},
				"userID":      event.UserID,
				"name":        event.Name,
				"description": event.Description,
				"startTime":   CloseTimeMatcher{event.StartTime, time.Minute},
				"startTimeTz": event.StartTimeTz,
			})).
			WillReturnRows(eventRows)
		mock.ExpectCommit()
		mock.ExpectRollback()

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
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create bad time", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
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
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
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
		Created:      ts,
		LastModified: ts,
	}
	event := &model.Event{
		ID:            1,
		RefID:         refid.Must(model.NewEventRefID()),
		UserID:        user.ID,
		Name:          "event",
		Description:   "description",
		ItemSortOrder: []int{},
		StartTime:     ts,
		StartTimeTz:   model.Must(model.ParseTimeZone("Etc/UTC")),
		Created:       ts,
		LastModified:  ts,
	}

	eventColumns := []string{
		"id", "ref_id", "user_id", "name", "description",
		"start_time", "start_time_tz", "created", "last_modified",
	}

	t.Run("update", func(t *testing.T) {
		t.Parallel()

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.ID, event.RefID, event.UserID, event.Name, event.Description,
				event.StartTime, event.StartTimeTz, ts, ts,
			)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(eventRows)
		mock.ExpectBegin()
		// refid as anyarg because new refid is created on call to create
		mock.ExpectExec("^UPDATE event_ ").
			WithArgs(util.NewPgxNamedArgsMatcher(pgx.NamedArgs{
				"name":          event.Name,
				"description":   event.Description,
				"itemSortOrder": pgxmock.AnyArg(),
				"startTime":     CloseTimeMatcher{event.StartTime, time.Minute},
				"startTimeTz":   event.StartTimeTz,
				"eventID":       event.ID,
			})).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

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
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update bad refid", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", refid.Must(model.NewEarmarkRefID()).String())

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
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update missing event", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnError(pgx.ErrNoRows)

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
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update user.id not match event.userid", func(t *testing.T) {
		t.Parallel()

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.ID, event.RefID, 33, event.Name, event.Description,
				event.StartTime, event.StartTimeTz, ts, ts,
			)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(eventRows)

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
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update only update name", func(t *testing.T) {
		t.Parallel()

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.ID, event.RefID, event.UserID, event.Name, event.Description,
				event.StartTime, event.StartTimeTz, ts, ts,
			)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(eventRows)
		mock.ExpectBegin()
		// refid as anyarg because new refid is created on call to create
		mock.ExpectExec("^UPDATE event_ ").
			WithArgs(util.NewPgxNamedArgsMatcher(pgx.NamedArgs{
				"name":          event.Name,
				"description":   event.Description,
				"itemSortOrder": pgxmock.AnyArg(),
				"startTime":     CloseTimeMatcher{event.StartTime, time.Minute},
				"startTimeTz":   event.StartTimeTz,
				"eventID":       event.ID,
			})).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

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
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Assert(t,
			strings.HasPrefix(rr.Header().Get("location"), "/events/"),
			"handler returned wrong redirect: expected prefix %s didnt match %s",
			"/events/", rr.Header().Get("location"),
		)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update only update description", func(t *testing.T) {
		t.Parallel()

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.ID, event.RefID, event.UserID, event.Name, event.Description,
				event.StartTime, event.StartTimeTz, ts, ts,
			)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(eventRows)
		mock.ExpectBegin()
		// refid as anyarg because new refid is created on call to create
		mock.ExpectExec("^UPDATE event_ ").
			WithArgs(util.NewPgxNamedArgsMatcher(pgx.NamedArgs{
				"name":          event.Name,
				"description":   event.Description,
				"itemSortOrder": pgxmock.AnyArg(),
				"startTime":     CloseTimeMatcher{event.StartTime, time.Minute},
				"startTimeTz":   event.StartTimeTz,
				"eventID":       event.ID,
			})).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		data := url.Values{
			"description": {event.Description},
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
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update only update when and tz", func(t *testing.T) {
		t.Parallel()

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.ID, event.RefID, event.UserID, event.Name, event.Description,
				event.StartTime, event.StartTimeTz, ts, ts,
			)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(eventRows)
		mock.ExpectBegin()
		// refid as anyarg because new refid is created on call to create
		mock.ExpectExec("^UPDATE event_ ").
			WithArgs(util.NewPgxNamedArgsMatcher(pgx.NamedArgs{
				"name":          event.Name,
				"description":   event.Description,
				"itemSortOrder": pgxmock.AnyArg(),
				"startTime":     CloseTimeMatcher{event.StartTime, time.Minute},
				"startTimeTz":   event.StartTimeTz,
				"eventID":       event.ID,
			})).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

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
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update only update when", func(t *testing.T) {
		t.Parallel()

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.ID, event.RefID, event.UserID, event.Name, event.Description,
				event.StartTime, event.StartTimeTz, ts, ts,
			)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(eventRows)

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
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update only update tz", func(t *testing.T) {
		t.Parallel()

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.ID, event.RefID, event.UserID, event.Name, event.Description,
				event.StartTime, event.StartTimeTz, ts, ts,
			)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(eventRows)

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
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update update nothing", func(t *testing.T) {
		t.Parallel()

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.ID, event.RefID, event.UserID, event.Name, event.Description,
				event.StartTime, event.StartTimeTz, ts, ts,
			)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(eventRows)

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
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update bad tz default to utc", func(t *testing.T) {
		t.Parallel()

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.ID, event.RefID, event.UserID, event.Name, event.Description,
				event.StartTime, event.StartTimeTz, ts, ts,
			)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(eventRows)
		mock.ExpectBegin()
		// refid as anyarg because new refid is created on call to create
		mock.ExpectExec("^UPDATE event_ ").
			WithArgs(util.NewPgxNamedArgsMatcher(pgx.NamedArgs{
				"name":          event.Name,
				"description":   event.Description,
				"itemSortOrder": pgxmock.AnyArg(),
				"startTime":     CloseTimeMatcher{event.StartTime, time.Minute},
				"startTimeTz":   event.StartTimeTz,
				"eventID":       event.ID,
			})).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

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
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Assert(t,
			strings.HasPrefix(rr.Header().Get("location"), "/events/"),
			"handler returned wrong redirect: expected prefix %s didnt match %s",
			"/events/", rr.Header().Get("location"),
		)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update bad time", func(t *testing.T) {
		t.Parallel()

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.ID, event.RefID, event.UserID, event.Name, event.Description,
				event.StartTime, event.StartTimeTz, ts, ts,
			)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(eventRows)

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
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
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
		Created:      ts,
		LastModified: ts,
	}
	event := &model.Event{
		ID:            1,
		RefID:         refid.Must(model.NewEventRefID()),
		UserID:        user.ID,
		Name:          "event",
		Description:   "description",
		ItemSortOrder: []int{1, 2, 3},
		StartTime:     ts,
		StartTimeTz:   model.Must(model.ParseTimeZone("Etc/UTC")),
		Created:       ts,
		LastModified:  ts,
	}

	eventColumns := []string{
		"id", "ref_id", "user_id", "name", "description", "item_sort_order",
		"start_time", "start_time_tz", "created", "last_modified",
	}

	t.Run("update", func(t *testing.T) {
		t.Parallel()

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.ID, event.RefID, event.UserID, event.Name, event.Description,
				event.ItemSortOrder, event.StartTime, event.StartTimeTz, ts, ts,
			)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(eventRows)
		mock.ExpectBegin()
		// refid as anyarg because new refid is created on call to create
		mock.ExpectExec("^UPDATE event_ ").
			WithArgs(util.NewPgxNamedArgsMatcher(pgx.NamedArgs{
				"name":          event.Name,
				"description":   event.Description,
				"itemSortOrder": []int{1, 3, 2},
				"startTime":     CloseTimeMatcher{event.StartTime, time.Minute},
				"startTimeTz":   event.StartTimeTz,
				"eventID":       event.ID,
			})).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

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
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update bad refid", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
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
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update missing event", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnError(pgx.ErrNoRows)

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
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update user.id not match event.userid", func(t *testing.T) {
		t.Parallel()

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.ID, event.RefID, 33, event.Name, event.Description,
				event.ItemSortOrder, event.StartTime, event.StartTimeTz, ts, ts,
			)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(eventRows)

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
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update update nothing", func(t *testing.T) {
		t.Parallel()

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.ID, event.RefID, event.UserID, event.Name, event.Description,
				event.ItemSortOrder, event.StartTime, event.StartTimeTz, ts, ts,
			)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(eventRows)

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
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update bad values", func(t *testing.T) {
		t.Parallel()

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.ID, event.RefID, event.UserID, event.Name, event.Description,
				event.ItemSortOrder, event.StartTime, event.StartTimeTz, ts, ts,
			)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(eventRows)

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
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update bad values 2", func(t *testing.T) {
		t.Parallel()

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.ID, event.RefID, event.UserID, event.Name, event.Description,
				event.ItemSortOrder, event.StartTime, event.StartTimeTz, ts, ts,
			)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(eventRows)

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
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
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
		Created:      ts,
		LastModified: ts,
	}
	event := &model.Event{
		ID:           1,
		RefID:        refid.Must(model.NewEventRefID()),
		UserID:       user.ID,
		Name:         "event",
		Description:  "description",
		StartTime:    ts,
		StartTimeTz:  model.Must(model.ParseTimeZone("Etc/UTC")),
		Created:      ts,
		LastModified: ts,
	}

	eventColumns := []string{
		"id", "ref_id", "user_id", "name", "description",
		"start_time", "start_time_tz", "created", "last_modified",
	}

	t.Run("delete", func(t *testing.T) {
		t.Parallel()

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.ID, event.RefID, event.UserID, event.Name, event.Description,
				event.StartTime, event.StartTimeTz, ts, ts,
			)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(eventRows)
		mock.ExpectBegin()
		// refid as anyarg because new refid is created on call to create
		mock.ExpectExec("^DELETE FROM event_ ").
			WithArgs(event.ID).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/event", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.DeleteEvent(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusOK)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete bad event refid", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
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
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete missing event", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnError(pgx.ErrNoRows)

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/event", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.DeleteEvent(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete mismatch user id", func(t *testing.T) {
		t.Parallel()

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.ID, event.RefID, 33, event.Name, event.Description,
				event.StartTime, event.StartTimeTz, ts, ts,
			)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(eventRows)

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/event", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.DeleteEvent(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}
