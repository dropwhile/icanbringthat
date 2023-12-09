package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
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

func TestHandler_Favorite_Delete(t *testing.T) {
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
	favoriteColumns := []string{"id", "user_id", "event_id", "created"}

	t.Run("delete favorite", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.ID, event.RefID, event.UserID, event.Name, event.Description,
				event.StartTime, event.StartTimeTz, ts, ts,
			)
		favoriteRows := pgxmock.NewRows(favoriteColumns).
			AddRow(33, user.ID, event.ID, ts)

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(eventRows)
		mock.ExpectQuery("^SELECT (.+) FROM favorite_ (.+)").
			WithArgs(pgx.NamedArgs{"userID": user.ID, "eventID": event.ID}).
			WillReturnRows(favoriteRows)
		mock.ExpectBegin()
		mock.ExpectExec("^DELETE FROM favorite_").
			WithArgs(33).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))
		mock.ExpectCommit()
		mock.ExpectRollback()

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/favorite", nil)
		rr := httptest.NewRecorder()
		handler.DeleteFavorite(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusOK)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete favorite bad event refid", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", "hodor")

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/favorite", nil)
		rr := httptest.NewRecorder()
		handler.DeleteFavorite(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete favorite event not found", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		refID := refid.Must(model.NewEventRefID())
		rctx.URLParams.Add("eRefID", refID.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ ").
			WithArgs(refID).
			WillReturnError(pgx.ErrNoRows)

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/favorite", nil)
		rr := httptest.NewRecorder()
		handler.DeleteFavorite(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete favorite event refid wrong type", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		refID := refid.Must(model.NewEventItemRefID())
		rctx.URLParams.Add("eRefID", refID.String())

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/favorite", nil)
		rr := httptest.NewRecorder()
		handler.DeleteFavorite(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("delete favorite not exist", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.ID, event.RefID, event.UserID, event.Name, event.Description,
				event.StartTime, event.StartTimeTz, ts, ts,
			)

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(eventRows)
		mock.ExpectQuery("^SELECT (.+) FROM favorite_ (.+)").
			WithArgs(pgx.NamedArgs{"userID": user.ID, "eventID": event.ID}).
			WillReturnError(pgx.ErrNoRows)

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/favorite", nil)
		rr := httptest.NewRecorder()
		handler.DeleteFavorite(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestHandler_Favorite_Add(t *testing.T) {
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
		UserID:       2,
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
	favoriteColumns := []string{"id", "user_id", "event_id", "created"}

	t.Run("create favorite", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.ID, event.RefID, event.UserID, event.Name, event.Description,
				event.StartTime, event.StartTimeTz, ts, ts,
			)
		favoriteRows := pgxmock.NewRows(favoriteColumns).
			AddRow(33, user.ID, event.ID, ts)

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(eventRows)
		mock.ExpectQuery("^SELECT (.+) FROM favorite_ (.+)").
			WithArgs(pgx.NamedArgs{"userID": user.ID, "eventID": event.ID}).
			WillReturnError(pgx.ErrNoRows)
		mock.ExpectBegin()
		mock.ExpectQuery("^INSERT INTO favorite_").
			WithArgs(pgx.NamedArgs{"userID": user.ID, "eventID": event.ID}).
			WillReturnRows(favoriteRows)
		mock.ExpectCommit()
		mock.ExpectRollback()

		req, _ := http.NewRequestWithContext(ctx, "PUT", "http://example.com/favorite", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.AddFavorite(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"),
			fmt.Sprintf("/events/%s", event.RefID.String()),
			"handler returned wrong redirect")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create favorite bad event refid", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", "hodor")

		req, _ := http.NewRequestWithContext(ctx, "PUT", "http://example.com/favorite", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.AddFavorite(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create favorite wrong type event refid", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", refid.Must(model.NewEventItemRefID()).String())

		req, _ := http.NewRequestWithContext(ctx, "PUT", "http://example.com/favorite", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.AddFavorite(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create favorite same user as event", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.ID, event.RefID, user.ID, event.Name, event.Description,
				event.StartTime, event.StartTimeTz, ts, ts,
			)
		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(eventRows)

		req, _ := http.NewRequestWithContext(ctx, "PUT", "http://example.com/favorite", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.AddFavorite(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create favorite already exists", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("eRefID", event.RefID.String())

		eventRows := pgxmock.NewRows(eventColumns).
			AddRow(
				event.ID, event.RefID, event.UserID, event.Name, event.Description,
				event.StartTime, event.StartTimeTz, ts, ts,
			)
		favoriteRows := pgxmock.NewRows(favoriteColumns).
			AddRow(33, user.ID, event.ID, ts)

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(eventRows)
		mock.ExpectQuery("^SELECT (.+) FROM favorite_ (.+)").
			WithArgs(pgx.NamedArgs{"userID": user.ID, "eventID": event.ID}).
			WillReturnRows(favoriteRows)

		req, _ := http.NewRequestWithContext(ctx, "PUT", "http://example.com/favorite", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.AddFavorite(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}
