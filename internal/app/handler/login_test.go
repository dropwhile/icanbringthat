package handler

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"go.uber.org/mock/gomock"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/crypto"
	"github.com/dropwhile/icbt/internal/errs"
	"github.com/dropwhile/icbt/internal/util"
)

func TestHandler_Login_InvalidCredentials(t *testing.T) {
	t.Parallel()

	pwhash, _ := crypto.HashPW([]byte("00x00"))
	user := &model.User{
		ID:       1,
		RefID:    util.Must(model.NewUserRefID()),
		Email:    "user@example.com",
		Name:     "user",
		PWHash:   pwhash,
		Verified: true,
	}

	t.Run("bad password", func(t *testing.T) {
		t.Parallel()
		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)

		mock.EXPECT().
			GetUserByEmail(gomock.Any(), user.Email).
			Return(user, nil)

		// bad password
		data := url.Values{
			"email":    {"user@example.com"},
			"password": {"00x01"},
		}

		ctx, _ = handler.sessMgr.Load(ctx, "")
		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/login", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.Login(rr, req)

		response := rr.Result()
		out, err := io.ReadAll(response.Body)
		assert.NilError(t, err)
		slog.DebugContext(ctx, "response",
			slog.String("body", string(out)))

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/login",
			"handler returned wrong redirect")
	})

	t.Run("no matching user", func(t *testing.T) {
		t.Parallel()
		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)

		mock.EXPECT().
			GetUserByEmail(gomock.Any(), user.Email).
			Return(nil, errs.NotFound.Error("user not found"))

		data := url.Values{
			"email":    {user.Email},
			"password": {"00x01"},
		}
		ctx, _ = handler.sessMgr.Load(ctx, "")
		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/login", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.Login(rr, req)

		response := rr.Result()
		out, err := io.ReadAll(response.Body)
		assert.NilError(t, err)
		slog.DebugContext(ctx, "response",
			slog.String("body", string(out)))

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/login",
			"handler returned wrong redirect")

		// we make sure that all expectations were met
	})

	t.Run("missing form data", func(t *testing.T) {
		t.Parallel()
		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)

		data := url.Values{
			"email": {"userXYZ@example.com"},
		}
		ctx, _ = handler.sessMgr.Load(ctx, "")
		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/login", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.Login(rr, req)

		response := rr.Result()
		out, err := io.ReadAll(response.Body)
		assert.NilError(t, err)
		slog.DebugContext(ctx, "response",
			slog.String("body", string(out)))

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
	})
}

func TestHandler_Login_ValidCredentials(t *testing.T) {
	t.Parallel()

	user := &model.User{
		ID:       1,
		RefID:    util.Must(model.NewUserRefID()),
		Email:    "user@example.com",
		Name:     "user",
		PWHash:   util.Must(crypto.HashPW([]byte("00x00"))),
		PWAuth:   true,
		Verified: true,
	}

	ctx := context.TODO()
	mock, _, handler := SetupHandler(t, ctx)
	// inject session into context
	ctx, _ = handler.sessMgr.Load(ctx, "")

	mock.EXPECT().
		GetUserByEmail(ctx, user.Email).
		Return(user, nil)

	data := url.Values{
		"email":    {"user@example.com"},
		"password": {"00x00"},
	}

	req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/login", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler.Login(rr, req)

	response := rr.Result()
	out, err := io.ReadAll(response.Body)
	assert.NilError(t, err)
	slog.DebugContext(ctx, "response",
		slog.String("body", string(out)))

	// Check the status code is what we expect.
	AssertStatusEqual(t, rr, http.StatusSeeOther)
	assert.Equal(t, rr.Header().Get("location"), "/dashboard",
		"handler returned wrong redirect")
}
