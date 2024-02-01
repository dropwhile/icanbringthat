package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/errs"
	"github.com/dropwhile/icbt/internal/middleware/auth"
	"github.com/dropwhile/icbt/internal/util"
)

func Test_getAuthnInstance(t *testing.T) {
	t.Parallel()

	t.Run("call isProd should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://127.0.0.1:8080/event", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		baseURL := "https://some.example.com"

		webAuthn, err := getAuthnInstance(req, baseURL)
		assert.NilError(t, err)

		assert.Equal(t, webAuthn.Config.RPDisplayName, "ICanBringThat")
		assert.Equal(t, webAuthn.Config.RPID, "some.example.com")
		assert.DeepEqual(t,
			webAuthn.Config.RPOrigins, []string{baseURL},
		)
	})
}

func TestHandler_DeleteWebAuthnKey(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		ID:           1,
		RefID:        util.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		WebAuthn:     false,
		Created:      ts,
		LastModified: ts,
	}
	credential := &model.UserCredential{
		ID:         3,
		RefID:      util.Must(model.NewCredentialRefID()),
		UserID:     user.ID,
		KeyName:    "key-name",
		Credential: []byte{0x01, 0x02},
		Created:    tstTs,
	}

	t.Run("delete should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("cRefID", credential.RefID.String())

		mock.EXPECT().
			DeleteUserCredential(ctx, user, credential.RefID).
			Return(nil)

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/event", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.DeleteWebAuthnKey(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusOK)
		// we make sure that all expectations were met
	})

	t.Run("delete with failed precondition should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("cRefID", credential.RefID.String())

		mock.EXPECT().
			DeleteUserCredential(ctx, user, credential.RefID).
			Return(errs.FailedPrecondition.Error("pre-condition failed"))

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/event", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.DeleteWebAuthnKey(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
	})

	t.Run("delete with credential not found should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("cRefID", credential.RefID.String())

		mock.EXPECT().
			DeleteUserCredential(ctx, user, credential.RefID).
			Return(errs.NotFound.Error("credential not found"))

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/event", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.DeleteWebAuthnKey(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})

	t.Run("delete other user credential should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("cRefID", credential.RefID.String())

		mock.EXPECT().
			DeleteUserCredential(ctx, user, credential.RefID).
			Return(errs.PermissionDenied.Error("permission denied"))

		req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/event", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.DeleteWebAuthnKey(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
	})
}
