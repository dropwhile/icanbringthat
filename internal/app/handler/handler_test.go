package handler

import (
	"context"
	"flag"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/session"
	"github.com/dropwhile/icbt/internal/util"
	"github.com/dropwhile/icbt/resources"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/rs/zerolog/log"
	"gotest.tools/v3/assert"
)

var tstTs time.Time

func init() {
	tstTs, _ = time.Parse(time.RFC3339, "2023-01-01T03:04:05Z")
}

func setCookie(r *http.Request, cookie string) {
	r.Header.Set("Cookie", cookie)
}

func SetupHandler(t *testing.T, ctx context.Context) (pgxmock.PgxConnIface, *chi.Mux, *Handler) {
	t.Helper()

	mock, err := pgxmock.NewConn()
	assert.NilError(t, err)
	t.Cleanup(func() { mock.Close(ctx) })
	h := &Handler{
		Db:      mock,
		Tpl:     resources.TemplateMap{},
		SessMgr: session.NewMemorySessionManager(),
	}
	mux := chi.NewMux()
	mux.Use(h.SessMgr.LoadAndSave)
	mux.Use(auth.Load(h.Db, h.SessMgr))
	return mock, mux, h
}

func SetupUserSession(t *testing.T, mux *chi.Mux, mock pgxmock.PgxConnIface, h *Handler) string {
	t.Helper()

	userId := 1
	ts := tstTs

	mux.Get("/dummy", func(w http.ResponseWriter, r *http.Request) {
		err := h.SessMgr.RenewToken(r.Context())
		assert.NilError(t, err)
		// add data to session
		h.SessMgr.Put(r.Context(), "user-id", userId)
		w.WriteHeader(http.StatusOK)
	})

	refId, _ := model.UserRefIdT.New()
	rows := pgxmock.NewRows(
		[]string{"id", "ref_id", "email", "name", "pwhash", "created", "last_modified"}).
		AddRow(userId, refId, "user@example.com", "user", []byte("00x00"), ts, ts)

	// mock.ExpectBegin()
	mock.ExpectQuery("^SELECT (.+) FROM user_").
		WithArgs(1).
		WillReturnRows(rows)

	// create request to set up session/cookies
	req, err := http.NewRequest("GET", "/dummy", nil)
	assert.NilError(t, err)
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	response := rr.Result()
	_, err = io.ReadAll(response.Body)
	assert.NilError(t, err)

	return rr.Header().Get("Set-Cookie")
}

func StatusEqual(rr *httptest.ResponseRecorder, status int) bool {
	return rr.Code == status
}

func AssertStatusEqual(t *testing.T, rr *httptest.ResponseRecorder, status int) {
	t.Helper()
	assert.Assert(t, rr.Code == status,
		"handler returned wrong status code: got %d expected %d", rr.Code, status)
}

func GetTokenViaRequest(mux *chi.Mux) (string, string) {
	var csrfToken string
	mux.Get("/get_token", func(_ http.ResponseWriter, r *http.Request) {
		csrfToken = csrf.Token(r)
	})

	req, _ := http.NewRequest("GET", "http://example.com/get_token", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return csrfToken, rr.Header().Get("Set-Cookie")
}

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	flag.Parse()
	log.Logger = util.NewTestLogger(os.Stderr)
	os.Exit(m.Run())
}
