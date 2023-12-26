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
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/samber/mo"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/service"
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

	eventColumns := []string{
		"id", "ref_id", "user_id", "name", "description", "archived",
		"start_time", "start_time_tz", "created", "last_modified",
	}

	t.Run("create", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		mock.ExpectBegin()
		// refid as anyarg because new refid is created on call to create
		mock.ExpectQuery("INSERT INTO event_ ").
			WithArgs(pgx.NamedArgs{
				"refID":       service.EventRefIDMatcher,
				"userID":      event.UserID,
				"name":        event.Name,
				"description": event.Description,
				"startTime":   util.CloseTimeMatcher{Value: event.StartTime, Within: time.Minute},
				"startTimeTz": event.StartTimeTz,
			}).
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name, event.Description,
					event.Archived, event.StartTime, event.StartTimeTz, ts, ts,
				))
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

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		mock.ExpectBegin()
		// refid as anyarg because new refid is created on call to create
		mock.ExpectQuery("INSERT INTO event_ ").
			WithArgs(pgx.NamedArgs{
				"refID":       service.EventRefIDMatcher,
				"userID":      event.UserID,
				"name":        event.Name,
				"description": event.Description,
				"startTime":   util.CloseTimeMatcher{Value: event.StartTime, Within: time.Minute},
				"startTimeTz": event.StartTimeTz,
			}).
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name, event.Description,
					event.Archived, event.StartTime, event.StartTimeTz, ts, ts,
				))
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
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

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
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create user not validated", func(t *testing.T) {
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

	eventColumns := []string{
		"id", "ref_id", "user_id", "name", "description", "archived",
		"start_time", "start_time_tz", "created", "last_modified",
	}

	t.Run("update event should succeed", func(t *testing.T) {
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
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name, event.Description,
					event.Archived, event.StartTime, event.StartTimeTz, ts, ts,
				))
		mock.ExpectBegin()
		// begin inner tx
		mock.ExpectBegin()
		// refid as anyarg because new refid is created on call to create
		mock.ExpectExec("^UPDATE event_ ").
			WithArgs(pgx.NamedArgs{
				"eventID":       event.ID,
				"name":          mo.Some(event.Name),
				"description":   mo.Some(event.Description),
				"itemSortOrder": pgxmock.AnyArg(),
				"startTime": util.OptionMatcher[time.Time](
					util.CloseTimeMatcher{
						Value: event.StartTime, Within: time.Minute,
					},
				),
				"startTimeTz": mo.Some(event.StartTimeTz),
			}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()
		// end inner tx
		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(event.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id",
					"user_id", "archived",
					"name", "description",
					"start_time", "start_time_tz",
					"created", "last_modified",
				}).
				AddRow(
					event.ID, event.RefID,
					user.ID, false,
					event.Name, event.Description,
					event.StartTime, event.StartTimeTz,
					tstTs, tstTs,
				),
			)
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

	t.Run("update event bad refid should fail", func(t *testing.T) {
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

	t.Run("update missing event should fail", func(t *testing.T) {
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

	t.Run("update user.id not match event.userid should fail", func(t *testing.T) {
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
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, 33, event.Name, event.Description,
					event.Archived, event.StartTime, event.StartTimeTz, ts, ts,
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
		handler.UpdateEvent(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update event archived should fail", func(t *testing.T) {
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
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name, event.Description,
					true, event.StartTime, event.StartTimeTz, ts, ts,
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
		handler.UpdateEvent(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update only update name should succeed", func(t *testing.T) {
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
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name, event.Description,
					event.Archived, event.StartTime, event.StartTimeTz, ts, ts,
				))
		mock.ExpectBegin()
		// inner tx
		mock.ExpectBegin()
		// refid as anyarg because new refid is created on call to create
		mock.ExpectExec("^UPDATE event_ ").
			WithArgs(pgx.NamedArgs{
				"name":          mo.Some(event.Name + "x"),
				"description":   mo.None[string](),
				"itemSortOrder": pgxmock.AnyArg(),
				"startTime":     mo.None[time.Time](),
				"startTimeTz":   mo.None[*model.TimeZone](),
				"eventID":       event.ID,
			}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()
		// end inner tx
		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(event.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id",
					"user_id", "archived",
					"name", "description",
					"start_time", "start_time_tz",
					"created", "last_modified",
				}).
				AddRow(
					event.ID, event.RefID,
					user.ID, false,
					event.Name, event.Description,
					event.StartTime, event.StartTimeTz,
					tstTs, tstTs,
				),
			)
		mock.ExpectCommit()
		mock.ExpectRollback()

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
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update only update description should succeed", func(t *testing.T) {
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
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name, event.Description,
					event.Archived, event.StartTime, event.StartTimeTz, ts, ts,
				))
		mock.ExpectBegin()
		// inner tx
		mock.ExpectBegin()
		// refid as anyarg because new refid is created on call to create
		mock.ExpectExec("^UPDATE event_ ").
			WithArgs(pgx.NamedArgs{
				"eventID":       event.ID,
				"name":          mo.None[string](),
				"description":   mo.Some(event.Description + "x"),
				"itemSortOrder": pgxmock.AnyArg(),
				"startTime":     mo.None[time.Time](),
				"startTimeTz":   mo.None[*model.TimeZone](),
			}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()
		// end inner tx
		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(event.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id",
					"user_id", "archived",
					"name", "description",
					"start_time", "start_time_tz",
					"created", "last_modified",
				}).
				AddRow(
					event.ID, event.RefID,
					user.ID, false,
					event.Name, event.Description,
					event.StartTime, event.StartTimeTz,
					tstTs, tstTs,
				),
			)
		mock.ExpectCommit()
		mock.ExpectRollback()

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
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update with when and tz should succeed", func(t *testing.T) {
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
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name, event.Description,
					event.Archived, event.StartTime, event.StartTimeTz, ts, ts,
				))
		mock.ExpectBegin()
		// inner tx
		mock.ExpectBegin()
		// refid as anyarg because new refid is created on call to create
		mock.ExpectExec("^UPDATE event_ ").
			WithArgs(pgx.NamedArgs{
				"name":          mo.None[string](),
				"description":   mo.None[string](),
				"itemSortOrder": pgxmock.AnyArg(),
				"startTime": util.OptionMatcher[time.Time](
					util.CloseTimeMatcher{
						Value: event.StartTime, Within: time.Minute,
					},
				),
				"startTimeTz": mo.Some(event.StartTimeTz),
				"eventID":     event.ID,
			}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()
		// end inner tx
		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(event.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id",
					"user_id", "archived",
					"name", "description",
					"start_time", "start_time_tz",
					"created", "last_modified",
				}).
				AddRow(
					event.ID, event.RefID,
					user.ID, false,
					event.Name, event.Description,
					event.StartTime, event.StartTimeTz,
					tstTs, tstTs,
				),
			)
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

	t.Run("update with only when should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
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
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update with only tz should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
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
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update nothing should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
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
		handler.UpdateEvent(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update bad tz default to utc should succeed", func(t *testing.T) {
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
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name, event.Description,
					event.Archived, event.StartTime, event.StartTimeTz, ts, ts,
				))
		mock.ExpectBegin()
		// start inner tx
		mock.ExpectBegin()
		// refid as anyarg because new refid is created on call to create
		mock.ExpectExec("^UPDATE event_ ").
			WithArgs(pgx.NamedArgs{
				"name":          mo.None[string](),
				"description":   mo.None[string](),
				"itemSortOrder": pgxmock.AnyArg(),
				"startTime": util.OptionMatcher[time.Time](
					util.CloseTimeMatcher{
						Value: event.StartTime, Within: time.Minute,
					},
				),
				"startTimeTz": mo.Some(event.StartTimeTz),
				"eventID":     event.ID,
			}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()
		// end innertx
		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(event.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id",
					"user_id", "archived",
					"name", "description",
					"start_time", "start_time_tz",
					"created", "last_modified",
				}).
				AddRow(
					event.ID, event.RefID,
					user.ID, false,
					event.Name, event.Description,
					event.StartTime, event.StartTimeTz,
					tstTs, tstTs,
				),
			)
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

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
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

	eventColumns := []string{
		"id", "ref_id", "user_id", "name", "description", "archived",
		"item_sort_order", "start_time", "start_time_tz", "created", "last_modified",
	}

	t.Run("update", func(t *testing.T) {
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
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name, event.Description,
					event.Archived, event.ItemSortOrder, event.StartTime,
					event.StartTimeTz, ts, ts,
				))
		mock.ExpectBegin()
		// refid as anyarg because new refid is created on call to create
		mock.ExpectExec("^UPDATE event_ ").
			WithArgs(pgx.NamedArgs{
				"eventID":       event.ID,
				"name":          mo.None[string](),
				"description":   mo.None[string](),
				"itemSortOrder": mo.Some([]int{1, 3, 2}),
				"startTime":     mo.None[time.Time](),
				"startTimeTz":   mo.None[*model.TimeZone](),
			}).
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

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, 33, event.Name, event.Description,
					event.Archived, event.ItemSortOrder, event.StartTime,
					event.StartTimeTz, ts, ts,
				))

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

	t.Run("update archived event should fail", func(t *testing.T) {
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
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name, event.Description,
					true, event.ItemSortOrder, event.StartTime,
					event.StartTimeTz, ts, ts,
				))

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

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
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
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update bad values", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
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
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update bad values 2", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
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

	eventColumns := []string{
		"id", "ref_id", "user_id", "name", "description", "archived",
		"start_time", "start_time_tz", "created", "last_modified",
	}

	t.Run("delete", func(t *testing.T) {
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
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name, event.Description,
					event.Archived, event.StartTime, event.StartTimeTz, ts, ts,
				))
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

	t.Run("delete event archived should succeed", func(t *testing.T) {
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
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name, event.Description,
					true, event.StartTime, event.StartTimeTz, ts, ts,
				))
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

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, 33, event.Name, event.Description,
					event.Archived, event.StartTime, event.StartTimeTz, ts, ts,
				))

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
