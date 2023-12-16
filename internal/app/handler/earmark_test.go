package handler

import (
	"context"
	"io"
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
	"github.com/dropwhile/icbt/internal/middleware/auth"
	"github.com/dropwhile/icbt/internal/util"
)

func TestHandler_Earmark_Create(t *testing.T) {
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
		StartTimeTz:  model.Must(model.ParseTimeZone("Etc/UTC")),
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

	t.Run("create earmark", func(t *testing.T) {
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
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id", "user_id", "name", "description",
						"start_time", "start_time_tz", "created", "last_modified",
					}).
					AddRow(
						event.ID, event.RefID, event.UserID, event.Name, event.Description,
						event.StartTime, event.StartTimeTz, ts, ts,
					),
			)
		mock.ExpectQuery("^SELECT (.+) FROM event_item_").
			WithArgs(eventItem.RefID).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id", "event_id", "description", "created", "last_modified",
					}).
					AddRow(
						eventItem.ID, eventItem.RefID, eventItem.EventID, eventItem.Description,
						ts, ts,
					),
			)
		mock.ExpectQuery("^SELECT (.+) FROM earmark_").
			WithArgs(eventItem.ID).
			WillReturnError(pgx.ErrNoRows)
		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(eventItem.ID).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id", "user_id", "name", "description",
						"start_time", "start_time_tz", "created", "last_modified",
					}).
					AddRow(
						event.ID, event.RefID, event.UserID, event.Name, event.Description,
						event.StartTime, event.StartTimeTz, ts, ts,
					),
			)
		mock.ExpectBegin()
		// refid as anyarg because new refid is created on call to create
		mock.ExpectQuery("^INSERT INTO earmark_").
			WithArgs(pgx.NamedArgs{
				"refID":       model.EarmarkRefIDMatcher,
				"eventItemID": earmark.EventItemID,
				"userID":      earmark.UserID,
				"note":        "some note",
			}).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id", "event_item_id", "user_id", "note", "created", "last_modified",
					}).
					AddRow(
						earmark.ID, earmark.RefID, earmark.EventItemID, earmark.UserID,
						earmark.Note, ts, ts,
					),
			)
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
		rctx.URLParams.Add("eRefID", "hodor")
		rctx.URLParams.Add("iRefID", eventItem.RefID.String())

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

	t.Run("create earmark wrong type event refid", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", eventItem.RefID.String())
		rctx.URLParams.Add("iRefID", eventItem.RefID.String())

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

	t.Run("create earmark bad event item id", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())
		rctx.URLParams.Add("iRefID", "hodor")

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

	t.Run("create earmark wrong type event item id", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())
		rctx.URLParams.Add("iRefID", event.RefID.String())

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

	t.Run("create earmark missing event", func(t *testing.T) {
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
		rctx.URLParams.Add("eRefID", event.RefID.String())
		rctx.URLParams.Add("iRefID", eventItem.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id", "user_id", "name", "description", "archived",
						"start_time", "start_time_tz", "created", "last_modified",
					}).
					AddRow(
						event.ID, event.RefID, event.UserID, event.Name, event.Description,
						event.Archived, event.StartTime, event.StartTimeTz, ts, ts,
					),
			)
		mock.ExpectQuery("^SELECT (.+) FROM event_item_").
			WithArgs(eventItem.RefID).
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
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id", "user_id", "name", "description", "archived",
						"start_time", "start_time_tz", "created", "last_modified",
					}).
					AddRow(
						event.ID, event.RefID, event.UserID, event.Name, event.Description,
						event.Archived, event.StartTime, event.StartTimeTz, ts, ts,
					),
			)
		mock.ExpectQuery("^SELECT (.+) FROM event_item_").
			WithArgs(eventItem.RefID).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id", "event_id", "description", "created", "last_modified",
					}).
					AddRow(
						eventItem.ID, eventItem.RefID, 33, eventItem.Description,
						ts, ts,
					),
			)

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

	t.Run("create earmark user not verified", func(t *testing.T) {
		t.Parallel()

		user := &model.User{
			ID:           33,
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
		rctx.URLParams.Add("eRefID", event.RefID.String())
		rctx.URLParams.Add("iRefID", eventItem.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id", "user_id", "name", "description", "archived",
						"start_time", "start_time_tz", "created", "last_modified",
					}).
					AddRow(
						event.ID, event.RefID, event.UserID, event.Name, event.Description,
						event.Archived, event.StartTime, event.StartTimeTz, ts, ts,
					),
			)
		mock.ExpectQuery("^SELECT (.+) FROM event_item_").
			WithArgs(eventItem.RefID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_id", "description", "created", "last_modified",
				}).
				AddRow(
					eventItem.ID, eventItem.RefID, eventItem.EventID, eventItem.Description,
					ts, ts,
				),
			)
		mock.ExpectQuery("SELECT (.+) FROM earmark_").
			WithArgs(eventItem.ID).
			WillReturnError(pgx.ErrNoRows)
		mock.ExpectQuery("SELECT (.+) FROM event_").
			WithArgs(eventItem.ID).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id", "user_id", "name", "description",
						"archived", "created", "last_modified",
					}).
					AddRow(
						event.ID, refid.Must(model.NewEventRefID()), 34,
						"event name", "event desc",
						false, tstTs, tstTs,
					),
			)

		data := url.Values{"note": {"some note"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/earmark", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.CreateEarmark(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestHandler_Earmark_Delete(t *testing.T) {
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

	t.Run("delete earmark", func(t *testing.T) {
		t.Parallel()

		event := &model.Event{
			ID:           1,
			RefID:        refid.Must(model.NewEventRefID()),
			UserID:       user.ID,
			Name:         "event",
			Description:  "description",
			Archived:     false,
			StartTime:    ts,
			StartTimeTz:  model.Must(model.ParseTimeZone("Etc/UTC")),
			Created:      ts,
			LastModified: ts,
		}
		earmark := &model.Earmark{
			ID:           1,
			RefID:        refid.Must(model.NewEarmarkRefID()),
			EventItemID:  1,
			UserID:       user.ID,
			Note:         "nothing",
			Created:      ts,
			LastModified: ts,
		}

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("mRefID", earmark.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM earmark_").
			WithArgs(earmark.RefID).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id", "event_item_id", "user_id", "note", "created", "last_modified",
					}).
					AddRow(
						earmark.ID, earmark.RefID, earmark.EventItemID, earmark.UserID,
						earmark.Note, ts, ts,
					),
			)
		mock.ExpectQuery("^SELECT (.+) FROM event_").
			WithArgs(earmark.ID).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id", "user_id", "name", "description", "archived",
						"start_time", "start_time_tz", "created", "last_modified",
					}).
					AddRow(
						event.ID, event.RefID, event.UserID, event.Name, event.Description,
						event.Archived, event.StartTime, event.StartTimeTz, ts, ts,
					),
			)
		mock.ExpectBegin()
		mock.ExpectExec("^DELETE FROM earmark_").
			WithArgs(earmark.ID).
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

	t.Run("delete earmark event archived", func(t *testing.T) {
		t.Parallel()

		earmark := &model.Earmark{
			ID:           1,
			RefID:        refid.Must(model.NewEarmarkRefID()),
			EventItemID:  1,
			UserID:       user.ID,
			Note:         "nothing",
			Created:      ts,
			LastModified: ts,
		}
		event := &model.Event{
			ID:           1,
			RefID:        refid.Must(model.NewEventRefID()),
			UserID:       user.ID,
			Name:         "event",
			Description:  "description",
			Archived:     true,
			StartTime:    ts,
			StartTimeTz:  model.Must(model.ParseTimeZone("Etc/UTC")),
			Created:      ts,
			LastModified: ts,
		}

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("mRefID", earmark.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM earmark_").
			WithArgs(earmark.RefID).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id", "event_item_id", "user_id", "note", "created", "last_modified",
					}).
					AddRow(
						earmark.ID, earmark.RefID, earmark.EventItemID, earmark.UserID,
						earmark.Note, ts, ts,
					),
			)
		mock.ExpectQuery("^SELECT (.+) FROM event_").
			WithArgs(earmark.ID).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id", "user_id", "name", "description", "archived",
						"start_time", "start_time_tz", "created", "last_modified",
					}).
					AddRow(
						event.ID, event.RefID, event.UserID, event.Name, event.Description,
						event.Archived, event.StartTime, event.StartTimeTz, ts, ts,
					),
			)

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
		AssertStatusEqual(t, rr, http.StatusNotFound)
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
		rctx.URLParams.Add("mRefID", "hodor")

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

	t.Run("delete earmark not found", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		refID := refid.Must(model.NewEarmarkRefID())
		rctx.URLParams.Add("mRefID", refID.String())

		mock.ExpectQuery("^SELECT (.+) FROM earmark_").
			WithArgs(refID).
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
		rctx.URLParams.Add("mRefID", refid.Must(model.NewEventRefID()).String())

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

	t.Run("delete earmark wrong user", func(t *testing.T) {
		t.Parallel()

		earmark := &model.Earmark{
			ID:           1,
			RefID:        refid.Must(model.NewEarmarkRefID()),
			EventItemID:  1,
			UserID:       user.ID,
			Note:         "nothing",
			Created:      ts,
			LastModified: ts,
		}

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("mRefID", earmark.RefID.String())

		mock.ExpectQuery("^SELECT (.+) FROM earmark_").
			WithArgs(earmark.RefID).
			WillReturnRows(
				pgxmock.NewRows(
					[]string{
						"id", "ref_id", "event_item_id", "user_id", "note",
						"created", "last_modified",
					}).
					AddRow(
						earmark.ID, earmark.RefID, earmark.EventItemID,
						user.ID+1, earmark.Note, ts, ts,
					),
			)

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
