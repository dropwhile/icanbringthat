package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/util"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"gotest.tools/v3/assert"
)

func TestHandler_EventItem_Create(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		Id:           1,
		RefId:        model.UserRefIdT.MustNew(),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Created:      ts,
		LastModified: ts,
	}
	event := &model.Event{
		Id:           1,
		RefId:        model.EventRefIdT.MustNew(),
		UserId:       user.Id,
		Name:         "event",
		Description:  "description",
		StartTime:    ts,
		StartTimeTZ:  "Etc/UTC",
		Created:      ts,
		LastModified: ts,
	}
	eventItem := &model.EventItem{
		Id:           2,
		RefId:        model.EventItemRefIdT.MustNew(),
		EventId:      event.Id,
		Description:  "eventitem",
		Created:      ts,
		LastModified: ts,
	}

	eventColumns := []string{
		"id", "ref_id", "user_id", "name", "description",
		"start_time", "start_time_tz", "created", "last_modified",
	}
	eventItemColumns := []string{
		"id", "ref_id", "event_id", "description", "created", "last_modified",
	}

	t.Run("create", func(t *testing.T) {
		t.Parallel()

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.Id, event.RefId, event.UserId, event.Name, event.Description,
				event.StartTime, event.StartTimeTZ, ts, ts,
			)
		eventItemRows := pgxmock.NewRows(eventItemColumns).
			AddRow(
				eventItem.Id, eventItem.RefId, eventItem.EventId, eventItem.Description,
				ts, ts,
			)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefId", event.RefId.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefId).
			WillReturnRows(eventRows)
		mock.ExpectBegin()
		// refid as anyarg because new refid is created on call to create
		mock.ExpectQuery("^INSERT INTO event_item_").
			WithArgs(model.EventItemRefIdT.AnyMatcher(), eventItem.EventId, "some description").
			WillReturnRows(eventItemRows)
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
			fmt.Sprintf("/events/%s", event.RefId),
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
		rctx.URLParams.Add("eRefId", eventItem.RefId.String())

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
		rctx.URLParams.Add("eRefId", event.RefId.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefId).
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

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.Id, event.RefId, 33, event.Name, event.Description,
				event.StartTime, event.StartTimeTZ, ts, ts,
			)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefId", event.RefId.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefId).
			WillReturnRows(eventRows)

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

	t.Run("create missing description", func(t *testing.T) {
		t.Parallel()

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.Id, event.RefId, event.UserId, event.Name, event.Description,
				event.StartTime, event.StartTimeTZ, ts, ts,
			)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefId", event.RefId.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefId).
			WillReturnRows(eventRows)

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
		Id:           1,
		RefId:        model.UserRefIdT.MustNew(),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Created:      ts,
		LastModified: ts,
	}
	event := &model.Event{
		Id:           1,
		RefId:        model.EventRefIdT.MustNew(),
		UserId:       user.Id,
		Name:         "event",
		Description:  "description",
		StartTime:    ts,
		StartTimeTZ:  "Etc/UTC",
		Created:      ts,
		LastModified: ts,
	}
	eventItem := &model.EventItem{
		Id:           2,
		RefId:        model.EventItemRefIdT.MustNew(),
		EventId:      event.Id,
		Description:  "eventitem",
		Created:      ts,
		LastModified: ts,
	}
	earmark := &model.Earmark{
		Id:           3,
		RefId:        model.EarmarkRefIdT.MustNew(),
		EventItemId:  eventItem.Id,
		UserId:       user.Id,
		Note:         "nothing",
		Created:      ts,
		LastModified: ts,
	}

	eventColumns := []string{
		"id", "ref_id", "user_id", "name", "description",
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

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.Id, event.RefId, event.UserId, event.Name, event.Description,
				event.StartTime, event.StartTimeTZ, ts, ts,
			)
		eventItemRows := pgxmock.NewRows(eventItemColumns).
			AddRow(
				eventItem.Id, eventItem.RefId, eventItem.EventId, eventItem.Description,
				ts, ts,
			)
		earmarkRows := pgxmock.NewRows(earmarkColumns).
			AddRow(
				earmark.Id, earmark.RefId, earmark.EventItemId, earmark.UserId,
				earmark.Note, ts, ts,
			)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefId", event.RefId.String())
		rctx.URLParams.Add("iRefId", eventItem.RefId.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefId).
			WillReturnRows(eventRows)
		mock.ExpectQuery("^SELECT (.+) FROM event_item_ (.+)").
			WithArgs(eventItem.RefId).
			WillReturnRows(eventItemRows)
		mock.ExpectQuery("^SELECT (.+) FROM earmark_ (.+)").
			WithArgs(eventItem.Id).
			WillReturnRows(earmarkRows)
		mock.ExpectBegin()
		// refid as anyarg because new refid is created on call to create
		mock.ExpectExec("^UPDATE event_item_").
			WithArgs("new description", eventItem.Id).
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
		rctx.URLParams.Add("eRefId", eventItem.RefId.String())
		rctx.URLParams.Add("iRefId", eventItem.RefId.String())

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
		rctx.URLParams.Add("eRefId", event.RefId.String())
		rctx.URLParams.Add("iRefId", eventItem.RefId.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefId).
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

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.Id, event.RefId, event.UserId, event.Name, event.Description,
				event.StartTime, event.StartTimeTZ, ts, ts,
			)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefId", event.RefId.String())
		rctx.URLParams.Add("iRefId", eventItem.RefId.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefId).
			WillReturnRows(eventRows)
		mock.ExpectQuery("^SELECT (.+) FROM event_item_ (.+)").
			WithArgs(eventItem.RefId).
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

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.Id, event.RefId, 33, event.Name, event.Description,
				event.StartTime, event.StartTimeTZ, ts, ts,
			)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefId", event.RefId.String())
		rctx.URLParams.Add("iRefId", eventItem.RefId.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefId).
			WillReturnRows(eventRows)

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

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.Id, event.RefId, event.UserId, event.Name, event.Description,
				event.StartTime, event.StartTimeTZ, ts, ts,
			)
		eventItemRows := pgxmock.NewRows(eventItemColumns).
			AddRow(
				eventItem.Id, eventItem.RefId, eventItem.EventId, eventItem.Description,
				ts, ts,
			)
		earmarkRows := pgxmock.NewRows(earmarkColumns).
			AddRow(
				earmark.Id, earmark.RefId, earmark.EventItemId, 33,
				earmark.Note, ts, ts,
			)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefId", event.RefId.String())
		rctx.URLParams.Add("iRefId", eventItem.RefId.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefId).
			WillReturnRows(eventRows)
		mock.ExpectQuery("^SELECT (.+) FROM event_item_ (.+)").
			WithArgs(eventItem.RefId).
			WillReturnRows(eventItemRows)
		mock.ExpectQuery("^SELECT (.+) FROM earmark_ (.+)").
			WithArgs(eventItem.Id).
			WillReturnRows(earmarkRows)

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

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.Id, event.RefId, event.UserId, event.Name, event.Description,
				event.StartTime, event.StartTimeTZ, ts, ts,
			)
		eventItemRows := pgxmock.NewRows(eventItemColumns).
			AddRow(
				eventItem.Id, eventItem.RefId, eventItem.EventId, eventItem.Description,
				ts, ts,
			)
		earmarkRows := pgxmock.NewRows(earmarkColumns).
			AddRow(
				earmark.Id, earmark.RefId, earmark.EventItemId, earmark.UserId,
				earmark.Note, ts, ts,
			)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefId", event.RefId.String())
		rctx.URLParams.Add("iRefId", eventItem.RefId.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefId).
			WillReturnRows(eventRows)
		mock.ExpectQuery("^SELECT (.+) FROM event_item_ (.+)").
			WithArgs(eventItem.RefId).
			WillReturnRows(eventItemRows)
		mock.ExpectQuery("^SELECT (.+) FROM earmark_ (.+)").
			WithArgs(eventItem.Id).
			WillReturnRows(earmarkRows)

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

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.Id, event.RefId, event.UserId, event.Name, event.Description,
				event.StartTime, event.StartTimeTZ, ts, ts,
			)
		eventItemRows := pgxmock.NewRows(eventItemColumns).
			AddRow(
				eventItem.Id, eventItem.RefId, 33, eventItem.Description,
				ts, ts,
			)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefId", event.RefId.String())
		rctx.URLParams.Add("iRefId", eventItem.RefId.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefId).
			WillReturnRows(eventRows)
		mock.ExpectQuery("^SELECT (.+) FROM event_item_ (.+)").
			WithArgs(eventItem.RefId).
			WillReturnRows(eventItemRows)

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
		Id:           1,
		RefId:        model.UserRefIdT.MustNew(),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Created:      ts,
		LastModified: ts,
	}
	event := &model.Event{
		Id:           1,
		RefId:        model.EventRefIdT.MustNew(),
		UserId:       user.Id,
		Name:         "event",
		Description:  "description",
		StartTime:    ts,
		StartTimeTZ:  "Etc/UTC",
		Created:      ts,
		LastModified: ts,
	}
	eventItem := &model.EventItem{
		Id:           2,
		RefId:        model.EventItemRefIdT.MustNew(),
		EventId:      event.Id,
		Description:  "eventitem",
		Created:      ts,
		LastModified: ts,
	}

	eventColumns := []string{
		"id", "ref_id", "user_id", "name", "description",
		"start_time", "start_time_tz", "created", "last_modified",
	}
	eventItemColumns := []string{
		"id", "ref_id", "event_id", "description", "created", "last_modified",
	}

	t.Run("delete", func(t *testing.T) {
		t.Parallel()

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.Id, event.RefId, event.UserId, event.Name, event.Description,
				event.StartTime, event.StartTimeTZ, ts, ts,
			)
		eventItemRows := pgxmock.NewRows(eventItemColumns).
			AddRow(
				eventItem.Id, eventItem.RefId, eventItem.EventId, eventItem.Description,
				ts, ts,
			)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefId", event.RefId.String())
		rctx.URLParams.Add("iRefId", eventItem.RefId.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefId).
			WillReturnRows(eventRows)
		mock.ExpectQuery("^SELECT (.+) FROM event_item_ (.+)").
			WithArgs(eventItem.RefId).
			WillReturnRows(eventItemRows)
		mock.ExpectBegin()
		// refid as anyarg because new refid is created on call to create
		mock.ExpectExec("^DELETE FROM event_item_").
			WithArgs(eventItem.Id).
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
		rctx.URLParams.Add("eRefId", eventItem.RefId.String())
		rctx.URLParams.Add("iRefId", eventItem.RefId.String())

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
		rctx.URLParams.Add("eRefId", event.RefId.String())
		rctx.URLParams.Add("iRefId", event.RefId.String())

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
		rctx.URLParams.Add("eRefId", event.RefId.String())
		rctx.URLParams.Add("iRefId", eventItem.RefId.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefId).
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

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.Id, event.RefId, event.UserId, event.Name, event.Description,
				event.StartTime, event.StartTimeTZ, ts, ts,
			)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefId", event.RefId.String())
		rctx.URLParams.Add("iRefId", eventItem.RefId.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefId).
			WillReturnRows(eventRows)
		mock.ExpectQuery("^SELECT (.+) FROM event_item_ (.+)").
			WithArgs(eventItem.RefId).
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

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.Id, event.RefId, 33, event.Name, event.Description,
				event.StartTime, event.StartTimeTZ, ts, ts,
			)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefId", event.RefId.String())
		rctx.URLParams.Add("iRefId", eventItem.RefId.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefId).
			WillReturnRows(eventRows)

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

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.Id, event.RefId, event.UserId, event.Name, event.Description,
				event.StartTime, event.StartTimeTZ, ts, ts,
			)
		eventItemRows := pgxmock.NewRows(eventItemColumns).
			AddRow(
				eventItem.Id, eventItem.RefId, 33, eventItem.Description,
				ts, ts,
			)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefId", event.RefId.String())
		rctx.URLParams.Add("iRefId", eventItem.RefId.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefId).
			WillReturnRows(eventRows)
		mock.ExpectQuery("^SELECT (.+) FROM event_item_ (.+)").
			WithArgs(eventItem.RefId).
			WillReturnRows(eventItemRows)

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
}
