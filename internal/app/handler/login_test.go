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
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/crypto"
)

func TestHandler_Login_InvalidCredentials(t *testing.T) {
	t.Parallel()

	refID := refid.Must(model.NewUserRefID())
	ts := tstTs
	pwhash, _ := crypto.HashPW([]byte("00x00"))

	t.Run("bad password", func(t *testing.T) {
		t.Parallel()
		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)

		rows := pgxmock.NewRows(
			[]string{"id", "ref_id", "email", "name", "pwhash", "pwauth", "created", "last_modified"}).
			AddRow(1, refID, "user@example.com", "user", pwhash, true, ts, ts)
		mock.ExpectQuery("^SELECT (.+) FROM user_").
			WithArgs("user@example.com").
			WillReturnRows(rows)

		// bad password
		data := url.Values{
			"email":    {"user@example.com"},
			"password": {"00x01"},
		}

		ctx, _ = handler.SessMgr.Load(ctx, "")
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

		// no matching user
		mock.ExpectQuery("^SELECT (.+) FROM user_").
			WithArgs("userXYZ@example.com").
			WillReturnError(pgx.ErrNoRows)
		data := url.Values{
			"email":    {"userXYZ@example.com"},
			"password": {"00x01"},
		}
		ctx, _ = handler.SessMgr.Load(ctx, "")
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
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("missing form data", func(t *testing.T) {
		t.Parallel()
		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)

		data := url.Values{
			"email": {"userXYZ@example.com"},
		}
		ctx, _ = handler.SessMgr.Load(ctx, "")
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
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestHandler_Login_ValidCredentials(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	mock, _, handler := SetupHandler(t, ctx)

	refID := refid.Must(model.NewUserRefID())
	ts := tstTs
	pwhash, _ := crypto.HashPW([]byte("00x00"))
	rows := pgxmock.NewRows(
		[]string{"id", "ref_id", "email", "name", "pwhash", "pwauth", "created", "last_modified"}).
		AddRow(1, refID, "user@example.com", "user", pwhash, true, ts, ts)

	// mock.ExpectBegin()
	mock.ExpectQuery("^SELECT (.+) FROM user_").
		WithArgs("user@example.com").
		WillReturnRows(rows)

	data := url.Values{
		"email":    {"user@example.com"},
		"password": {"00x00"},
	}

	// inject session into context
	ctx, _ = handler.SessMgr.Load(ctx, "")
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
