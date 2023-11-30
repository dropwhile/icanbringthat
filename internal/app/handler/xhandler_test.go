package handler

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dropwhile/refid"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/pashagolub/pgxmock/v3"
	pg_query "github.com/pganalyze/pg_query_go/v4"
	"github.com/rs/zerolog/log"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/crypto"
	"github.com/dropwhile/icbt/internal/logger"
	"github.com/dropwhile/icbt/internal/mail"
	"github.com/dropwhile/icbt/internal/session"
	"github.com/dropwhile/icbt/resources"
)

var tstTs time.Time = MustParseTime(time.RFC3339, "2023-01-01T03:04:05Z")

func MustParseTime(layout, value string) time.Time {
	ts, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return ts
}

type CloseTimeMatcher struct {
	value  time.Time
	within time.Duration
}

// Match satisfies sqlmock.Argument interface
func (a CloseTimeMatcher) Match(v interface{}) bool {
	ts, ok := v.(time.Time)
	if !ok {
		return false
	}
	if ts.Equal(a.value) {
		return true
	}
	if ts.Before(a.value) {
		return !ts.Add(a.within).Before(a.value)
	}
	if ts.After(a.value) {
		return !ts.Add(-a.within).After(a.value)
	}
	return true
}

func setCookie(r *http.Request, cookie string) {
	r.Header.Set("Cookie", cookie)
}

type TestMailer struct {
	Sent []*mail.Mail
}

func (tm *TestMailer) SendRaw(mail *mail.Mail) error {
	tm.Sent = append(tm.Sent, mail)
	return nil
}

func (tm *TestMailer) Send(from string, to []string, subject, bodyPlain, bodyHtml string, extraHeaders mail.MailHeader) error {
	if from == "" {
		from = "user@example.com"
	}
	mail := &mail.Mail{
		Sender:       from,
		To:           to,
		ExtraHeaders: extraHeaders,
		Subject:      subject,
		BodyPlain:    bodyPlain,
		BodyHtml:     bodyHtml,
	}
	return tm.SendRaw(mail)
}

func (tm *TestMailer) SendAsync(from string, to []string, subject, bodyPlain, bodyHtml string, extraHeaders mail.MailHeader) {
	tm.Send(from, to, subject, bodyPlain, bodyHtml, extraHeaders)
}

func SetupHandler(t *testing.T, ctx context.Context) (pgxmock.PgxConnIface, *chi.Mux, *Handler) {
	t.Helper()

	var queryMatcher pgxmock.QueryMatcher = pgxmock.QueryMatcherFunc(func(expectedSQL, actualSQL string) error {
		err := pgxmock.QueryMatcherRegexp.Match(expectedSQL, actualSQL)
		if err != nil {
			return err
		}
		_, err = pg_query.Parse(actualSQL)
		if err != nil {
			return fmt.Errorf("error parsing sql '%s': %w", actualSQL, err)
		}

		return nil
	})

	mock, err := pgxmock.NewConn(
		pgxmock.QueryMatcherOption(queryMatcher),
	)
	assert.NilError(t, err)
	t.Cleanup(func() { mock.Close(ctx) })
	tpl := template.Must(template.New("error-page.gohtml").Parse(`{{.ErrorCode}}-{{.ErrorStatus}}`))
	h := &Handler{
		Db:          mock,
		TemplateMap: resources.TemplateMap{"error-page.gohtml": tpl},
		SessMgr:     session.NewTestSessionManager(),
		Mailer:      &TestMailer{make([]*mail.Mail, 0)},
		MAC:         crypto.NewMAC([]byte("test-hmac-key")),
		BaseURL:     "http://example.com",
	}
	mux := chi.NewMux()
	mux.Use(h.SessMgr.LoadAndSave)
	mux.Use(auth.Load(h.Db, h.SessMgr))
	return mock, mux, h
}

func SetupUserSession(t *testing.T, mux *chi.Mux, mock pgxmock.PgxConnIface, x *Handler) string {
	t.Helper()

	userID := 1
	ts := tstTs

	mux.Get("/dummy", func(w http.ResponseWriter, r *http.Request) {
		err := x.SessMgr.RenewToken(r.Context())
		assert.NilError(t, err)
		// add data to session
		x.SessMgr.Put(r.Context(), "user-id", userID)
		w.WriteHeader(http.StatusOK)
	})

	refID := refid.Must(model.NewUserRefID())
	rows := pgxmock.NewRows(
		[]string{"id", "ref_id", "email", "name", "pwhash", "created", "last_modified"}).
		AddRow(userID, refID, "user@example.com", "user", []byte("00x00"), ts, ts)

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

func FormData(values url.Values) *strings.Reader {
	return strings.NewReader(values.Encode())
}

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	flag.Parse()
	log.Logger = logger.NewTestLogger(os.Stderr)
	os.Exit(m.Run())
}