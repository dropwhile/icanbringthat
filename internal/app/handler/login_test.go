package handler

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/util"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/rs/zerolog/log"
	"gotest.tools/v3/assert"
)

func TestHandler_Login_InvalidCredentials(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	mock, _, handler := SetupHandler(t, ctx)

	refId, _ := model.UserRefIdT.New()
	ts := tstTs
	pwhash, _ := util.HashPW([]byte("00x00"))
	rows := pgxmock.NewRows(
		[]string{"id", "ref_id", "email", "name", "pwhash", "created", "last_modified"}).
		AddRow(1, refId, "user@example.com", "user", pwhash, ts, ts)

	// mock.ExpectBegin()
	mock.ExpectQuery("^SELECT (.+) FROM user_").
		WithArgs("user@example.com").
		WillReturnRows(rows)

	data := url.Values{
		"email":    []string{"user@example.com"},
		"password": []string{"00x01"},
	}

	ctx, _ = handler.SessMgr.Load(ctx, "")
	req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/login", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler.Login(rr, req)

	response := rr.Result()
	out, err := io.ReadAll(response.Body)
	log.Debug().Msg(string(out))
	assert.NilError(t, err)

	// Check the status code is what we expect.
	AssertStatusEqual(t, rr, http.StatusSeeOther)
	assert.Equal(t, rr.Header().Get("location"), "/login",
		"handler returned wrong redirect")
	// we make sure that all expectations were met
	assert.Assert(t, mock.ExpectationsWereMet(),
		"there were unfulfilled expectations")
}

func TestHandler_Login_ValidCredentials(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	mock, _, handler := SetupHandler(t, ctx)

	refId, _ := model.UserRefIdT.New()
	ts := tstTs
	pwhash, _ := util.HashPW([]byte("00x00"))
	rows := pgxmock.NewRows(
		[]string{"id", "ref_id", "email", "name", "pwhash", "created", "last_modified"}).
		AddRow(1, refId, "user@example.com", "user", pwhash, ts, ts)

	// mock.ExpectBegin()
	mock.ExpectQuery("^SELECT (.+) FROM user_").
		WithArgs("user@example.com").
		WillReturnRows(rows)

	data := url.Values{
		"email":    []string{"user@example.com"},
		"password": []string{"00x00"},
	}

	// inject session into context
	ctx, _ = handler.SessMgr.Load(ctx, "")
	req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/login", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler.Login(rr, req)

	response := rr.Result()
	out, err := io.ReadAll(response.Body)
	log.Debug().Msg(string(out))
	assert.NilError(t, err)

	// Check the status code is what we expect.
	AssertStatusEqual(t, rr, http.StatusSeeOther)
	assert.Equal(t, rr.Header().Get("location"), "/dashboard",
		"handler returned wrong redirect")
	// we make sure that all expectations were met
	assert.Assert(t, mock.ExpectationsWereMet(),
		"there were unfulfilled expectations")
}
