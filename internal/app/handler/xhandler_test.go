package handler

import (
	"context"
	"flag"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"go.uber.org/mock/gomock"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/app/resources"
	"github.com/dropwhile/icbt/internal/app/service/mockservice"
	"github.com/dropwhile/icbt/internal/crypto"
	"github.com/dropwhile/icbt/internal/logger"
	"github.com/dropwhile/icbt/internal/mail"
	"github.com/dropwhile/icbt/internal/mail/mockmail"
	"github.com/dropwhile/icbt/internal/middleware/auth"
	"github.com/dropwhile/icbt/internal/session"
	"github.com/dropwhile/icbt/internal/util"
)

var tstTs time.Time = util.MustParseTime(time.RFC3339, "2030-01-01T03:04:05Z")

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

func SetupMailerMock(t *testing.T) *mockmail.MockMailSender {
	t.Helper()

	ctrl := gomock.NewController(t)
	mailer := mockmail.NewMockMailSender(ctrl)
	return mailer
}

func SetupHandler(
	t *testing.T, ctx context.Context,
) (*mockservice.MockServicer, *chi.Mux, *Handler) {
	t.Helper()

	ctrl := gomock.NewController(t)
	mock := mockservice.NewMockServicer(ctrl)
	tpl := util.Must(template.New("error-page.gohtml").Parse(`{{.ErrorCode}}-{{.ErrorStatus}}`))
	h := &Handler{
		templates: &resources.TemplateMap{"error-page.gohtml": tpl},
		sessMgr:   session.NewTestSessionManager(),
		mailer:    &TestMailer{make([]*mail.Mail, 0)},
		cMAC:      crypto.NewMAC([]byte("test-hmac-key")),
		baseURL:   "http://example.com",
		svc:       mock,
	}
	mux := chi.NewMux()
	mux.Use(h.sessMgr.LoadAndSave)
	mux.Use(auth.Load(h.svc, h.sessMgr))
	return mock, mux, h
}

func StatusEqual(rr *httptest.ResponseRecorder, status int) bool {
	return rr.Code == status
}

func AssertStatusEqual(t *testing.T, rr *httptest.ResponseRecorder, status int) {
	t.Helper()
	assert.Equal(t, rr.Code, status,
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
	logger.SetupLogging(logger.NewTestLogger, nil)
	os.Exit(m.Run())
}
