package handler

import (
	"context"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/dropwhile/refid/v2"
	"github.com/go-chi/chi/v5"
	"go.uber.org/mock/gomock"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/resources"
	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/internal/encoder"
	"github.com/dropwhile/icbt/internal/errs"
	"github.com/dropwhile/icbt/internal/mail"
	"github.com/dropwhile/icbt/internal/middleware/auth"
	"github.com/dropwhile/icbt/internal/util"
)

func TestHandler_ResetPassword(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		ID:           1,
		RefID:        refid.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     false,
		PWAuth:       true,
		WebAuthn:     false,
		Created:      ts,
		LastModified: ts,
	}

	pwr := &model.UserPWReset{
		RefID:   refid.Must(model.NewUserPWResetRefID()),
		UserID:  user.ID,
		Created: ts,
	}

	t.Run("pwreset", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("upwRefID", pwr.RefID.String())

		// generate hmac
		macBytes := handler.cMAC.Generate([]byte(pwr.RefID.String()))
		// base32 encode hmac
		macStr := encoder.Base32EncodeToString(macBytes)
		rctx.URLParams.Add("hmac", macStr)

		mock.EXPECT().
			GetUserByID(ctx, pwr.UserID).
			Return(user, nil)
		mock.EXPECT().
			GetUserPWResetByRefID(ctx, pwr.RefID).
			Return(pwr, nil)
		mock.EXPECT().
			UpdateUserPWReset(ctx, user, pwr).
			Return(nil)

		data := url.Values{
			"password":         {"newpass"},
			"confirm_password": {"newpass"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/password-reset", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.ResetPassword(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		userID := handler.sessMgr.GetInt(ctx, "user-id")
		assert.Equal(t, userID, user.ID)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/dashboard",
			"handler returned wrong redirect")
		// we make sure that all expectations were met
	})

	t.Run("pwreset with user logged in", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("upwRefID", pwr.RefID.String())

		// generate hmac
		macBytes := handler.cMAC.Generate([]byte(pwr.RefID.String()))
		// base32 encode hmac
		macStr := encoder.Base32EncodeToString(macBytes)

		rctx.URLParams.Add("hmac", macStr)

		data := url.Values{
			"password":         {"newpass"},
			"confirm_password": {"newpass"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/password-reset", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.ResetPassword(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
	})

	t.Run("pwreset with password mismatch", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("upwRefID", pwr.RefID.String())

		// generate hmac
		macBytes := handler.cMAC.Generate([]byte(pwr.RefID.String()))
		// base32 encode hmac
		macStr := encoder.Base32EncodeToString(macBytes)

		rctx.URLParams.Add("hmac", macStr)

		data := url.Values{
			"password":         {"newpass"},
			"confirm_password": {"newpassx"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/password-reset", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.ResetPassword(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
	})

	t.Run("pwreset with bad hmac", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("upwRefID", pwr.RefID.String())

		// generate hmac
		macBytes := handler.cMAC.Generate([]byte(pwr.RefID.String()))
		macBytes[0] += 1
		// base32 encode hmac
		macStr := encoder.Base32EncodeToString(macBytes)

		rctx.URLParams.Add("hmac", macStr)

		data := url.Values{
			"password":         {"newpass"},
			"confirm_password": {"newpass"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/password-reset", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.ResetPassword(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
	})

	t.Run("pwreset with bad refid", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		refID := refid.Must(model.NewEventItemRefID())
		rctx.URLParams.Add("upwRefID", refID.String())

		// generate hmac
		macBytes := handler.cMAC.Generate([]byte(refID.String()))
		// base32 encode hmac
		macStr := encoder.Base32EncodeToString(macBytes)

		rctx.URLParams.Add("hmac", macStr)

		data := url.Values{
			"password":         {"newpass"},
			"confirm_password": {"newpass"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/password-reset", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.ResetPassword(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})

	t.Run("pwreset upw not in db", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("upwRefID", pwr.RefID.String())

		// generate hmac
		macBytes := handler.cMAC.Generate([]byte(pwr.RefID.String()))
		// base32 encode hmac
		macStr := encoder.Base32EncodeToString(macBytes)

		rctx.URLParams.Add("hmac", macStr)

		mock.EXPECT().
			GetUserPWResetByRefID(ctx, pwr.RefID).
			Return(nil, errs.NotFound.Error("upw not found"))

		data := url.Values{
			"password":         {"newpass"},
			"confirm_password": {"newpass"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/password-reset", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.ResetPassword(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
	})

	t.Run("pwreset user not found", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("upwRefID", pwr.RefID.String())

		// generate hmac
		macBytes := handler.cMAC.Generate([]byte(pwr.RefID.String()))
		// base32 encode hmac
		macStr := encoder.Base32EncodeToString(macBytes)

		rctx.URLParams.Add("hmac", macStr)

		mock.EXPECT().
			GetUserPWResetByRefID(ctx, pwr.RefID).
			Return(pwr, nil)
		mock.EXPECT().
			GetUserByID(ctx, pwr.UserID).
			Return(nil, errs.NotFound.Error("user not found"))

		data := url.Values{
			"password":         {"newpass"},
			"confirm_password": {"newpass"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/password-reset", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.ResetPassword(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
	})

	t.Run("pwreset upw is expired", func(t *testing.T) {
		t.Parallel()

		refID := refid.Must(model.NewUserPWResetRefID())
		rfts, _ := time.Parse(time.RFC3339, "2023-01-14T18:29:00Z")
		refID.SetTime(rfts)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("upwRefID", refID.String())

		// generate hmac
		macBytes := handler.cMAC.Generate([]byte(refID.String()))
		// base32 encode hmac
		macStr := encoder.Base32EncodeToString(macBytes)

		rctx.URLParams.Add("hmac", macStr)

		pwr := &model.UserPWReset{
			UserID:  user.ID,
			RefID:   refID,
			Created: rfts,
		}

		mock.EXPECT().
			GetUserPWResetByRefID(ctx, pwr.RefID).
			Return(pwr, nil)

		data := url.Values{
			"password":         {"newpass"},
			"confirm_password": {"newpass"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/password-reset", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.ResetPassword(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})
}

func TestHandler_SendResetPasswordEmail(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		ID:           1,
		RefID:        refid.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     false,
		PWAuth:       true,
		WebAuthn:     false,
		Created:      ts,
		LastModified: ts,
	}

	pwr := &model.UserPWReset{
		RefID:   refid.Must(model.NewUserPWResetRefID()),
		UserID:  user.ID,
		Created: ts,
	}

	passResetTpl := template.Must(template.New("").Parse(`{{.Subject}}: {{.PasswordResetUrl}}`))

	t.Run("send pw reset email", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		mockmailer := SetupMailerMock(t)
		handler.mailer = mockmailer
		ctx, _ = handler.sessMgr.Load(ctx, "")
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		handler.templates = resources.MockTContainer(
			resources.TemplateMap{
				"mail_password_reset.gohtml": passResetTpl,
				"mail_password_reset.gotxt":  passResetTpl,
			})

		mock.EXPECT().
			GetUserByEmail(ctx, user.Email).
			Return(user, nil)
		mock.EXPECT().
			NewUserPWReset(ctx, user.ID).
			Return(pwr, nil)
		mockmailer.EXPECT().
			SendAsync("",
				[]string{user.Email},
				"Password reset",
				gomock.AssignableToTypeOf("string"),
				gomock.AssignableToTypeOf("string"),
				mail.MailHeader{
					"X-PM-Message-Stream": "outbound",
				},
			).Do(
			func(
				from string, recipients []string,
				subject, msgPlain, msgHtml string,
				headers mail.MailHeader,
			) {
				after, found := strings.CutPrefix(msgPlain, "Password reset: http://example.com/forgot-password/")
				assert.Assert(t, found)
				refParts := strings.Split(after, "-")
				rID := refid.Must(service.ParseUserPWResetRefID(refParts[0]))
				hmacBytes, err := encoder.Base32DecodeString(refParts[1])
				assert.NilError(t, err)
				assert.Assert(t, handler.cMAC.Validate([]byte(rID.String()), hmacBytes))
			},
		)

		data := url.Values{"email": {"user@example.com"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/send-password-reset", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.SendResetPasswordEmail(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/login",
			"handler returned wrong redirect")
		// we make sure that all expectations were met
	})

	t.Run("send pw reset email no user", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		handler.templates = resources.MockTContainer(
			resources.TemplateMap{
				"mail_password_reset.gohtml": passResetTpl,
				"mail_password_reset.gotxt":  passResetTpl,
			})

		mock.EXPECT().
			GetUserByEmail(ctx, user.Email).
			Return(nil, errs.NotFound.Error("user not found"))

		data := url.Values{"email": {"user@example.com"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/send-password-reset", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.SendResetPasswordEmail(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		tm := handler.mailer.(*TestMailer)
		assert.Equal(t, len(tm.Sent), 0)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/login",
			"handler returned wrong redirect")
		// we make sure that all expectations were met
	})
}
