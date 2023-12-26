package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/dropwhile/refid/v2"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/internal/middleware/auth"
	"github.com/dropwhile/icbt/internal/util"
)

func TestHandler_EventItem_Create(t *testing.T) {
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
		Archived:     false,
		StartTime:    ts,
		StartTimeTz:  util.Must(service.ParseTimeZone("Etc/UTC")),
		Created:      ts,
		LastModified: ts,
	}
	eventItem := &model.EventItem{
		ID:           2,
		RefID:        refid.Must(model.NewEventItemRefID()),
		EventID:      event.ID,
		Description:  "eventitem",
		Created:      ts,
		LastModified: ts,
	}

	eventColumns := []string{
		"id", "ref_id", "user_id", "name", "description", "archived",
		"start_time", "start_time_tz", "created", "last_modified",
	}
	eventItemColumns := []string{
		"id", "ref_id", "event_id", "description", "created", "last_modified",
	}

	t.Run("create", func(t *testing.T) {
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
		mock.ExpectQuery("^INSERT INTO event_item_").
			WithArgs(pgx.NamedArgs{
				"refID":       service.EventItemRefIDMatcher,
				"eventID":     eventItem.EventID,
				"description": "some description",
			}).
			WillReturnRows(pgxmock.NewRows(eventItemColumns).
				AddRow(
					eventItem.ID, eventItem.RefID, eventItem.EventID, eventItem.Description,
					ts, ts,
				))
		mock.ExpectCommit()
		mock.ExpectRollback()

		data := url.Values{"description": {"some description"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/eventItem", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.CreateEventItem(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"),
			fmt.Sprintf("/events/%s", event.RefID),
			"handler returned wrong redirect")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create bad event refid", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", eventItem.RefID.String())

		data := url.Values{"description": {"some description"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/eventItem", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.CreateEventItem(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create missing event", func(t *testing.T) {
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

		data := url.Values{"description": {"some description"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/eventItem", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.CreateEventItem(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create user not owner", func(t *testing.T) {
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

		data := url.Values{"description": {"some description"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/eventItem", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.CreateEventItem(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create event archived", func(t *testing.T) {
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

		data := url.Values{"description": {"some description"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/eventItem", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.CreateEventItem(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create with missing description param should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		data := url.Values{"descriptionxxx": {"some description"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/eventItem", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.CreateEventItem(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestHandler_EventItem_Update(t *testing.T) {
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
		Archived:     false,
		StartTime:    ts,
		StartTimeTz:  util.Must(service.ParseTimeZone("Etc/UTC")),
		Created:      ts,
		LastModified: ts,
	}
	eventItem := &model.EventItem{
		ID:           2,
		RefID:        refid.Must(model.NewEventItemRefID()),
		EventID:      event.ID,
		Description:  "eventitem",
		Created:      ts,
		LastModified: ts,
	}
	earmark := &model.Earmark{
		ID:           3,
		RefID:        refid.Must(model.NewEarmarkRefID()),
		EventItemID:  eventItem.ID,
		UserID:       user.ID,
		Note:         "nothing",
		Created:      ts,
		LastModified: ts,
	}

	eventColumns := []string{
		"id", "ref_id", "user_id", "name", "description", "archived",
		"start_time", "start_time_tz", "created", "last_modified",
	}
	eventItemColumns := []string{
		"id", "ref_id", "event_id", "description", "created", "last_modified",
	}
	earmarkColumns := []string{
		"id", "ref_id", "event_item_id", "user_id", "note", "created", "last_modified",
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
		rctx.URLParams.Add("iRefID", eventItem.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name, event.Description,
					event.Archived, event.StartTime, event.StartTimeTz, ts, ts,
				))
		mock.ExpectQuery("^SELECT (.+) FROM event_item_ (.+)").
			WithArgs(eventItem.RefID).
			WillReturnRows(pgxmock.NewRows(eventItemColumns).
				AddRow(
					eventItem.ID, eventItem.RefID, eventItem.EventID, eventItem.Description,
					ts, ts,
				))
		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.ID).
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name, event.Description,
					event.Archived, event.StartTime, event.StartTimeTz, ts, ts,
				))
		mock.ExpectQuery("^SELECT (.+) FROM earmark_ (.+)").
			WithArgs(eventItem.ID).
			WillReturnRows(pgxmock.NewRows(earmarkColumns).
				AddRow(
					earmark.ID, earmark.RefID, earmark.EventItemID, earmark.UserID,
					earmark.Note, ts, ts,
				))
		mock.ExpectBegin()
		// refid as anyarg because new refid is created on call to create
		mock.ExpectExec("^UPDATE event_item_").
			WithArgs(pgx.NamedArgs{
				"description": "new description",
				"eventItemID": eventItem.ID,
			}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		data := url.Values{"description": {"new description"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/eventItem", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateEventItem(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusOK)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update bad event refid", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", eventItem.RefID.String())
		rctx.URLParams.Add("iRefID", eventItem.RefID.String())

		data := url.Values{"description": {"new description"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/eventItem", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateEventItem(rr, req)

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
		rctx.URLParams.Add("iRefID", eventItem.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnError(pgx.ErrNoRows)

		data := url.Values{"description": {"new description"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/eventItem", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateEventItem(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update missing event item", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())
		rctx.URLParams.Add("iRefID", eventItem.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name, event.Description,
					event.Archived, event.StartTime, event.StartTimeTz, ts, ts,
				))
		mock.ExpectQuery("^SELECT (.+) FROM event_item_ (.+)").
			WithArgs(eventItem.RefID).
			WillReturnError(pgx.ErrNoRows)

		data := url.Values{"description": {"new description"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/eventItem", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateEventItem(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update event owner not match", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())
		rctx.URLParams.Add("iRefID", eventItem.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, 33, event.Name, event.Description,
					event.Archived, event.StartTime, event.StartTimeTz, ts, ts,
				))
		mock.ExpectQuery("^SELECT (.+) FROM event_item_ ").
			WithArgs(eventItem.RefID).
			WillReturnRows(pgxmock.NewRows(eventItemColumns).
				AddRow(
					eventItem.ID, eventItem.RefID, eventItem.EventID, eventItem.Description,
					ts, ts,
				))
		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.ID).
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, 33, event.Name, event.Description,
					event.Archived, event.StartTime, event.StartTimeTz, ts, ts,
				))

		data := url.Values{"description": {"new description"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/eventItem", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateEventItem(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update event archived", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())
		rctx.URLParams.Add("iRefID", eventItem.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name, event.Description,
					true, event.StartTime, event.StartTimeTz, ts, ts,
				))
		mock.ExpectQuery("^SELECT (.+) FROM event_item_ ").
			WithArgs(eventItem.RefID).
			WillReturnRows(pgxmock.NewRows(eventItemColumns).
				AddRow(
					eventItem.ID, eventItem.RefID, eventItem.EventID, eventItem.Description,
					ts, ts,
				))
		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.ID).
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name, event.Description,
					true, event.StartTime, event.StartTimeTz, ts, ts,
				))

		data := url.Values{"description": {"new description"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/eventItem", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateEventItem(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
	t.Run("update already earmarked by other user", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())
		rctx.URLParams.Add("iRefID", eventItem.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name, event.Description,
					event.Archived, event.StartTime, event.StartTimeTz, ts, ts,
				))
		mock.ExpectQuery("^SELECT (.+) FROM event_item_ (.+)").
			WithArgs(eventItem.RefID).
			WillReturnRows(pgxmock.NewRows(eventItemColumns).
				AddRow(
					eventItem.ID, eventItem.RefID, eventItem.EventID, eventItem.Description,
					ts, ts,
				))
		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.ID).
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name, event.Description,
					event.Archived, event.StartTime, event.StartTimeTz, ts, ts,
				))
		mock.ExpectQuery("^SELECT (.+) FROM earmark_ (.+)").
			WithArgs(eventItem.ID).
			WillReturnRows(pgxmock.NewRows(earmarkColumns).
				AddRow(
					earmark.ID, earmark.RefID, earmark.EventItemID, 33,
					earmark.Note, ts, ts,
				))

		data := url.Values{"description": {"new description"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/eventItem", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateEventItem(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update missing form data", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())
		rctx.URLParams.Add("iRefID", eventItem.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name, event.Description,
					event.Archived, event.StartTime, event.StartTimeTz, ts, ts,
				))

		data := url.Values{"descriptionxxxx": {"new description"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/eventItem", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateEventItem(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update eventitem not matching event", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())
		rctx.URLParams.Add("iRefID", eventItem.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name, event.Description,
					event.Archived, event.StartTime, event.StartTimeTz, ts, ts,
				))
		mock.ExpectQuery("^SELECT (.+) FROM event_item_ (.+)").
			WithArgs(eventItem.RefID).
			WillReturnRows(pgxmock.NewRows(eventItemColumns).
				AddRow(
					eventItem.ID, eventItem.RefID, 33, eventItem.Description,
					ts, ts,
				))

		data := url.Values{"description": {"new description"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/eventItem", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateEventItem(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestHandler_EventItem_Delete(t *testing.T) {
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
		Archived:     false,
		StartTime:    ts,
		StartTimeTz:  util.Must(service.ParseTimeZone("Etc/UTC")),
		Created:      ts,
		LastModified: ts,
	}
	eventItem := &model.EventItem{
		ID:           2,
		RefID:        refid.Must(model.NewEventItemRefID()),
		EventID:      event.ID,
		Description:  "eventitem",
		Created:      ts,
		LastModified: ts,
	}

	eventColumns := []string{
		"id", "ref_id", "user_id", "name", "description", "archived",
		"start_time", "start_time_tz", "created", "last_modified",
	}
	eventItemColumns := []string{
		"id", "ref_id", "event_id", "description", "created", "last_modified",
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
		rctx.URLParams.Add("iRefID", eventItem.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name, event.Description,
					event.Archived, event.StartTime, event.StartTimeTz, ts, ts,
				))
		mock.ExpectQuery("^SELECT (.+) FROM event_item_ (.+)").
			WithArgs(eventItem.RefID).
			WillReturnRows(pgxmock.NewRows(eventItemColumns).
				AddRow(
					eventItem.ID, eventItem.RefID, eventItem.EventID, eventItem.Description,
					ts, ts,
				))
		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.ID).
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name, event.Description,
					event.Archived, event.StartTime, event.StartTimeTz, ts, ts,
				))
		mock.ExpectBegin()
		// refid as anyarg because new refid is created on call to create
		mock.ExpectExec("^DELETE FROM event_item_").
			WithArgs(eventItem.ID).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/eventItem", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.DeleteEventItem(rr, req)

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
		rctx.URLParams.Add("eRefID", eventItem.RefID.String())
		rctx.URLParams.Add("iRefID", eventItem.RefID.String())

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/eventItem", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.DeleteEventItem(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete bad event item refid", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())
		rctx.URLParams.Add("iRefID", event.RefID.String())

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/eventItem", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.DeleteEventItem(rr, req)

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
		rctx.URLParams.Add("iRefID", eventItem.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnError(pgx.ErrNoRows)

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/eventItem", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.DeleteEventItem(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete missing event item", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())
		rctx.URLParams.Add("iRefID", eventItem.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name, event.Description,
					event.Archived, event.StartTime, event.StartTimeTz, ts, ts,
				))
		mock.ExpectQuery("^SELECT (.+) FROM event_item_ (.+)").
			WithArgs(eventItem.RefID).
			WillReturnError(pgx.ErrNoRows)

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/eventItem", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.DeleteEventItem(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete user not owner", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())
		rctx.URLParams.Add("iRefID", eventItem.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, 33, event.Name, event.Description,
					event.Archived, event.StartTime, event.StartTimeTz, ts, ts,
				))
		mock.ExpectQuery("^SELECT (.+) FROM event_item_ (.+)").
			WithArgs(eventItem.RefID).
			WillReturnRows(pgxmock.NewRows(eventItemColumns).
				AddRow(
					eventItem.ID, eventItem.RefID, eventItem.EventID, eventItem.Description,
					ts, ts,
				))
		mock.ExpectQuery("^SELECT (.+) FROM event_ ").
			WithArgs(eventItem.EventID).
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, 33, event.Name, event.Description,
					event.Archived, event.StartTime, event.StartTimeTz, ts, ts,
				))

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/eventItem", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.DeleteEventItem(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete event item not related to supplied event", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())
		rctx.URLParams.Add("iRefID", eventItem.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name, event.Description,
					event.Archived, event.StartTime, event.StartTimeTz, ts, ts,
				))
		mock.ExpectQuery("^SELECT (.+) FROM event_item_ (.+)").
			WithArgs(eventItem.RefID).
			WillReturnRows(pgxmock.NewRows(eventItemColumns).
				AddRow(
					eventItem.ID, eventItem.RefID, 33, eventItem.Description,
					ts, ts,
				))

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/eventItem", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.DeleteEventItem(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete event archived", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())
		rctx.URLParams.Add("iRefID", eventItem.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name, event.Description,
					true, event.StartTime, event.StartTimeTz, ts, ts,
				))
		mock.ExpectQuery("^SELECT (.+) FROM event_item_ (.+)").
			WithArgs(eventItem.RefID).
			WillReturnRows(pgxmock.NewRows(eventItemColumns).
				AddRow(
					eventItem.ID, eventItem.RefID, event.UserID, eventItem.Description,
					ts, ts,
				))
		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.ID).
			WillReturnRows(pgxmock.NewRows(eventColumns).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name, event.Description,
					true, event.StartTime, event.StartTimeTz, ts, ts,
				))

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/eventItem", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.DeleteEventItem(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}
