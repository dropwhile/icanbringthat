package xhandler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dropwhile/refid"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/modelx"
	"github.com/dropwhile/icbt/internal/util"
)

func TestHandler_Favorite_Delete(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		Id:           1,
		RefID:        refid.Must(modelx.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Created:      ts,
		LastModified: ts,
	}
	event := &model.Event{
		Id:           1,
		RefID:        refid.Must(model.EventRefIDT.New()),
		UserId:       user.Id,
		Name:         "event",
		Description:  "description",
		StartTime:    ts,
		StartTimeTZ:  "Etc/UTC",
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
				event.Id, event.RefID, event.UserId, event.Name, event.Description,
				event.StartTime, event.StartTimeTZ, ts, ts,
			)
		favoriteRows := pgxmock.NewRows(favoriteColumns).
			AddRow(33, user.Id, event.Id, ts)

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(eventRows)
		mock.ExpectQuery("^SELECT (.+) FROM favorite_ (.+)").
			WithArgs(user.Id, event.Id).
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
		refId := refid.Must(model.EventRefIDT.New())
		rctx.URLParams.Add("eRefID", refId.String())

		mock.ExpectQuery("^SELECT (.+) FROM event_ ").
			WithArgs(refId).
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
		refId := refid.Must(model.EventItemRefIDT.New())
		rctx.URLParams.Add("eRefID", refId.String())

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
				event.Id, event.RefID, event.UserId, event.Name, event.Description,
				event.StartTime, event.StartTimeTZ, ts, ts,
			)

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(eventRows)
		mock.ExpectQuery("^SELECT (.+) FROM favorite_ (.+)").
			WithArgs(user.Id, event.Id).
			WillReturnError(pgx.ErrNoRows)

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/favorite", nil)
		rr := httptest.NewRecorder()
		handler.DeleteFavorite(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestHandler_Favorite_Add(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		Id:           1,
		RefID:        refid.Must(modelx.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Created:      ts,
		LastModified: ts,
	}
	event := &model.Event{
		Id:           1,
		RefID:        refid.Must(model.EventRefIDT.New()),
		UserId:       2,
		Name:         "event",
		Description:  "description",
		StartTime:    ts,
		StartTimeTZ:  "Etc/UTC",
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
				event.Id, event.RefID, event.UserId, event.Name, event.Description,
				event.StartTime, event.StartTimeTZ, ts, ts,
			)
		favoriteRows := pgxmock.NewRows(favoriteColumns).
			AddRow(33, user.Id, event.Id, ts)

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(eventRows)
		mock.ExpectQuery("^SELECT (.+) FROM favorite_ (.+)").
			WithArgs(user.Id, event.Id).
			WillReturnError(pgx.ErrNoRows)
		mock.ExpectBegin()
		mock.ExpectQuery("^INSERT INTO favorite_").
			WithArgs(user.Id, event.Id).
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
		rctx.URLParams.Add("eRefID", refid.Must(model.EventItemRefIDT.New()).String())

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
				event.Id, event.RefID, user.Id, event.Name, event.Description,
				event.StartTime, event.StartTimeTZ, ts, ts,
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
				event.Id, event.RefID, event.UserId, event.Name, event.Description,
				event.StartTime, event.StartTimeTZ, ts, ts,
			)
		favoriteRows := pgxmock.NewRows(favoriteColumns).
			AddRow(33, user.Id, event.Id, ts)

		mock.ExpectQuery("^SELECT (.+) FROM event_ (.+)").
			WithArgs(event.RefID).
			WillReturnRows(eventRows)
		mock.ExpectQuery("^SELECT (.+) FROM favorite_ (.+)").
			WithArgs(user.Id, event.Id).
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
