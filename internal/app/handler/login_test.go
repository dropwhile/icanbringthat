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

	"github.com/dropwhile/refid/v2"
	"github.com/pashagolub/pgxmock/v3"
	mox "github.com/stretchr/testify/mock"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/crypto"
	"github.com/dropwhile/icbt/internal/errs"
)

func TestHandler_Login_InvalidCredentials(t *testing.T) {
	t.Parallel()

	pwhash, _ := crypto.HashPW([]byte("00x00"))
	user := &model.User{
		ID:       1,
		RefID:    refid.Must(model.NewUserRefID()),
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
			GetUserByEmail(mox.AnythingOfType("*context.valueCtx"), user.Email).
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
		mock.AssertExpectations(t)
	})

	t.Run("no matching user", func(t *testing.T) {
		t.Parallel()
		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)

		mock.EXPECT().
			GetUserByEmail(mox.AnythingOfType("*context.valueCtx"), user.Email).
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
		mock.AssertExpectations(t)
	})

	t.Run("missing form data", func(t *testing.T) {
		t.Parallel()
		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)

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
		mock.AssertExpectations(t)
	})
}

func TestHandler_Login_ValidCredentials(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	mock, _, handler := SetupHandlerOld(t, ctx)

	refID := refid.Must(model.NewUserRefID())
	ts := tstTs
	pwhash, _ := crypto.HashPW([]byte("00x00"))

	// mock.ExpectBegin()
	mock.ExpectQuery("^SELECT (.+) FROM user_").
		WithArgs("user@example.com").
		WillReturnRows(pgxmock.NewRows(
			[]string{
				"id", "ref_id", "email", "name", "pwhash", "pwauth",
				"created", "last_modified",
			}).
			AddRow(
				1, refID, "user@example.com", "user", pwhash, true,
				ts, ts,
			),
		)

	data := url.Values{
		"email":    {"user@example.com"},
		"password": {"00x00"},
	}

	// inject session into context
	ctx, _ = handler.sessMgr.Load(ctx, "")
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
	// we make sure that all expectations were met
	assert.Assert(t, mock.ExpectationsWereMet(),
		"there were unfulfilled expectations")
}
