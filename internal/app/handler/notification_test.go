package handler

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dropwhile/refid/v2"
	"github.com/go-chi/chi/v5"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/errs"
	"github.com/dropwhile/icbt/internal/middleware/auth"
)

func TestHandler_Notification_Delete(t *testing.T) {
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
	notification := &model.Notification{
		ID:           2,
		RefID:        refid.Must(model.NewNotificationRefID()),
		UserID:       user.ID,
		Message:      "",
		Read:         false,
		Created:      ts,
		LastModified: ts,
	}

	t.Run("delete notification", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("nRefID", notification.RefID.String())

		mock.EXPECT().
			DeleteNotification(ctx, user.ID, notification.RefID).
			Return(nil)

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/notification", nil)
		rr := httptest.NewRecorder()
		handler.DeleteNotification(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusOK)
		// we make sure that all expectations were met
		mock.AssertExpectations(t)
	})

	t.Run("delete notification missing refid", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/notification", nil)
		rr := httptest.NewRecorder()
		handler.DeleteNotification(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
		mock.AssertExpectations(t)
	})

	t.Run("delete notification bad refid", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("nRefID", "hodor")

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/notification", nil)
		rr := httptest.NewRecorder()
		handler.DeleteNotification(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
		mock.AssertExpectations(t)
	})

	t.Run("delete notification not found", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		refID := refid.Must(model.NewNotificationRefID())
		rctx.URLParams.Add("nRefID", refID.String())

		mock.EXPECT().
			DeleteNotification(ctx, user.ID, refID).
			Return(errs.NotFound.Error("notification not found"))

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/notification", nil)
		rr := httptest.NewRecorder()
		handler.DeleteNotification(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
		mock.AssertExpectations(t)
	})

	t.Run("delete notification refid wrong type", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("nRefID", refid.Must(model.NewEventRefID()).String())

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/notification", nil)
		rr := httptest.NewRecorder()
		handler.DeleteNotification(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
		mock.AssertExpectations(t)
	})

	t.Run("delete notification wrong user", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("nRefID", notification.RefID.String())

		mock.EXPECT().
			DeleteNotification(ctx, user.ID, notification.RefID).
			Return(errs.PermissionDenied.Error("permission denied"))

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/notification", nil)
		rr := httptest.NewRecorder()
		handler.DeleteNotification(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
		mock.AssertExpectations(t)
	})
}

func TestHandler_Notification_DeleteAll(t *testing.T) {
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
	notification := &model.Notification{
		ID:           2,
		RefID:        refid.Must(model.NewNotificationRefID()),
		UserID:       user.ID,
		Message:      "",
		Read:         false,
		Created:      ts,
		LastModified: ts,
	}

	t.Run("delete all notifications", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("nRefID", notification.RefID.String())

		mock.EXPECT().
			DeleteAllNotifications(ctx, user.ID).
			Return(nil)

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/notification", nil)
		rr := httptest.NewRecorder()
		handler.DeleteAllNotifications(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusOK)
		// we make sure that all expectations were met
		mock.AssertExpectations(t)
	})
}
