package handler

import (
	"context"
	"io"
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

func TestHandler_Earmark_Delete(t *testing.T) {
	t.Parallel()

	refId := model.EarmarkRefIdT.MustNew()
	ts := tstTs
	user := &model.User{
		Id:           1,
		RefId:        refId,
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Created:      ts,
		LastModified: ts,
	}
	earmark := &model.Earmark{
		Id:           1,
		RefId:        refId,
		EventItemId:  1,
		UserId:       user.Id,
		Note:         "nothing",
		Created:      ts,
		LastModified: ts,
	}

	t.Run("delete earmark", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("mRefId", earmark.RefId.String())

		rows := pgxmock.NewRows(
			[]string{"id", "ref_id", "event_item_id", "user_id", "note", "created", "last_modified"}).
			AddRow(earmark.Id, earmark.RefId, earmark.EventItemId, user.Id, earmark.Note, ts, ts)

		mock.ExpectQuery("^SELECT (.+) FROM earmark_").
			WithArgs(earmark.RefId).
			WillReturnRows(rows)
		mock.ExpectBegin()
		mock.ExpectExec("^DELETE FROM earmark_").
			WithArgs(earmark.Id).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/earmark", nil)
		rr := httptest.NewRecorder()
		handler.DeleteEarmark(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusOK)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete earmark missing refid", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/earmark", nil)
		rr := httptest.NewRecorder()
		handler.DeleteEarmark(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete earmark bad refid", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("mRefId", "hodor")

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/earmark", nil)
		rr := httptest.NewRecorder()
		handler.DeleteEarmark(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete earmark not found", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		refId = model.EarmarkRefIdT.MustNew()
		rctx.URLParams.Add("mRefId", refId.String())

		mock.ExpectQuery("^SELECT (.+) FROM earmark_").
			WithArgs(refId).
			WillReturnError(pgx.ErrNoRows)

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/earmark", nil)
		rr := httptest.NewRecorder()
		handler.DeleteEarmark(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete earmark refid wrong type", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		refId = model.EventRefIdT.MustNew()
		rctx.URLParams.Add("mRefId", refId.String())

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/earmark", nil)
		rr := httptest.NewRecorder()
		handler.DeleteEarmark(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete earmark wrong user", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("mRefId", earmark.RefId.String())

		rows := pgxmock.NewRows(
			[]string{"id", "ref_id", "event_item_id", "user_id", "note", "created", "last_modified"}).
			AddRow(earmark.Id, earmark.RefId, earmark.EventItemId, user.Id+1, earmark.Note, ts, ts)

		mock.ExpectQuery("^SELECT (.+) FROM earmark_").
			WithArgs(earmark.RefId).
			WillReturnRows(rows)

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/earmark", nil)
		rr := httptest.NewRecorder()
		handler.DeleteEarmark(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestHandler_Earmark_Create(t *testing.T) {
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

	eventRows := pgxmock.NewRows(
		[]string{
			"id", "ref_id", "user_id", "name", "description",
			"start_time", "start_time_tz", "created", "last_modified",
		}).
		AddRow(
			event.Id, event.RefId, event.UserId, event.Name, event.Description,
			event.StartTime, event.StartTimeTZ, ts, ts,
		)
	eventItemRows := pgxmock.NewRows(
		[]string{
			"id", "ref_id", "event_id", "description", "created", "last_modified",
		}).
		AddRow(
			eventItem.Id, eventItem.RefId, eventItem.EventId, eventItem.Description,
			ts, ts,
		)
	earmarkRows := pgxmock.NewRows(
		[]string{
			"id", "ref_id", "event_item_id", "user_id", "note", "created", "last_modified",
		}).
		AddRow(
			earmark.Id, earmark.RefId, earmark.EventItemId, earmark.UserId,
			earmark.Note, ts, ts,
		)

	t.Run("create earmark", func(t *testing.T) {
		t.Parallel()

		eventRows := pgxmock.NewRows(
			[]string{
				"id", "ref_id", "user_id", "name", "description",
				"start_time", "start_time_tz", "created", "last_modified",
			}).
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
		mock.ExpectQuery("^SELECT (.+) FROM event_item_").
			WithArgs(eventItem.RefId).
			WillReturnRows(eventItemRows)
		mock.ExpectQuery("^SELECT (.+) FROM earmark_").
			WithArgs(eventItem.Id).
			WillReturnError(pgx.ErrNoRows)
		mock.ExpectBegin()
		// refid as anyarg because new refid is created on call to create
		mock.ExpectQuery("^INSERT INTO earmark_").
			WithArgs(model.EarmarkRefIdT.AnyMatcher(), earmark.EventItemId, earmark.UserId, "some note").
			WillReturnRows(earmarkRows)
		mock.ExpectCommit()
		mock.ExpectRollback()

		data := url.Values{"note": {"some note"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/earmark", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.CreateEarmark(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusOK)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create earmark bad event refid", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefId", "hodor")
		rctx.URLParams.Add("iRefId", eventItem.RefId.String())

		data := url.Values{"note": {"some note"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/earmark", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.CreateEarmark(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create earmark wrong type event refid", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefId", eventItem.RefId.String())
		rctx.URLParams.Add("iRefId", eventItem.RefId.String())

		data := url.Values{"note": {"some note"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/earmark", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.CreateEarmark(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create earmark bad event item id", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefId", event.RefId.String())
		rctx.URLParams.Add("iRefId", "hodor")

		data := url.Values{"note": {"some note"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/earmark", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.CreateEarmark(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create earmark wrong type event item id", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefId", event.RefId.String())
		rctx.URLParams.Add("iRefId", event.RefId.String())

		data := url.Values{"note": {"some note"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/earmark", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.CreateEarmark(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create earmark missing event", func(t *testing.T) {
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

		data := url.Values{"note": {"some note"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/earmark", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.CreateEarmark(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create earmark missing event item", func(t *testing.T) {
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
			WillReturnRows(eventRows)
		mock.ExpectQuery("^SELECT (.+) FROM event_item_").
			WithArgs(eventItem.RefId).
			WillReturnError(pgx.ErrNoRows)

		data := url.Values{"note": {"some note"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/earmark", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.CreateEarmark(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create earmark eventitem not matching event", func(t *testing.T) {
		t.Parallel()

		eventRows := pgxmock.NewRows(
			[]string{
				"id", "ref_id", "user_id", "name", "description",
				"start_time", "start_time_tz", "created", "last_modified",
			}).
			AddRow(
				event.Id, event.RefId, event.UserId, event.Name, event.Description,
				event.StartTime, event.StartTimeTZ, ts, ts,
			)
		eventItemRows := pgxmock.NewRows(
			[]string{
				"id", "ref_id", "event_id", "description", "created", "last_modified",
			}).
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
		mock.ExpectQuery("^SELECT (.+) FROM event_item_").
			WithArgs(eventItem.RefId).
			WillReturnRows(eventItemRows)

		data := url.Values{"note": {"some note"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/earmark", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.CreateEarmark(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}
