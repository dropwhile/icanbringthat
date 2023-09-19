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
	"github.com/gorilla/csrf"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/rs/zerolog/log"
	"gotest.tools/v3/assert"
)

func TestHandler_Login(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	mock, mux, handler := SetupHandler(t, ctx)
	// CSRF middleware in use for this test
	mux.Use(csrf.Protect(
		[]byte("testkey"),
		csrf.Secure(false),
		csrf.Path("/"),
		csrf.RequestHeader("X-CSRF-Token"),
	))
	mux.Post("/login", handler.Login)
	var csrfToken string
	mux.Get("/login", func(_ http.ResponseWriter, r *http.Request) {
		csrfToken = csrf.Token(r)
	})

	// get token
	req, err := http.NewRequest("GET", "http://example.com/login", nil)
	assert.NilError(t, err)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

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

	data := url.Values{}
	data.Set("email", "user@example.com")
	data.Set("password", "00x00")
	req, err = http.NewRequest("POST", "http://example.com/login", strings.NewReader(data.Encode()))
	assert.NilError(t, err)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	setCookie(req, rr.Header().Get("Set-Cookie"))
	req.Header.Set("X-CSRF-Token", csrfToken)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

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

/*
func TestHandler_ShowIndex_LoggedIn(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	mock, mux, handler := SetupHandler(t, ctx)
	cookies := SetupUserSession(t, mux, mock, handler)
	mux.Get("/", handler.ShowIndex)

	// create request
	req, err := http.NewRequest("GET", "/", nil)
	assert.NilError(t, err)
	for _, c := range cookies {
		req.AddCookie(c)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	response := rr.Result()
	_, err = io.ReadAll(response.Body)
	assert.NilError(t, err)

	// Check the status code is what we expect.
	assert.Assert(t, StatusEqual(rr, http.StatusSeeOther),
		"handler returned wrong status code")
	assert.Equal(t, rr.Header().Get("location"), "/dashboard",
		"handler returned wrong redirect")
	// we make sure that all expectations were met
	assert.Assert(t, mock.ExpectationsWereMet(),
		"there were unfulfilled expectations")
}
*/
