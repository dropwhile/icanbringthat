// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package handler

import (
	"context"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dropwhile/assert"
	"github.com/go-chi/chi/v5"
	"go.uber.org/mock/gomock"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/app/resources"
	"github.com/dropwhile/icanbringthat/internal/app/service"
	"github.com/dropwhile/icanbringthat/internal/encoder"
	"github.com/dropwhile/icanbringthat/internal/errs"
	"github.com/dropwhile/icanbringthat/internal/mail"
	"github.com/dropwhile/icanbringthat/internal/middleware/auth"
	"github.com/dropwhile/icanbringthat/internal/util"
)

func TestHandler_SendVerificationEmail(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		ID:           1,
		RefID:        util.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     false,
		PWAuth:       true,
		WebAuthn:     false,
		Created:      ts,
		LastModified: ts,
	}

	uv := &model.UserVerify{
		RefID:   util.Must(model.NewUserVerifyRefID()),
		UserID:  user.ID,
		Created: ts,
	}

	verifyTpl := util.Must(template.New("").Parse(
		`{{.Subject}}: <a href="{{.VerificationUrl}}">{{.VerificationUrl}}</a>`))

	t.Run("send verification email", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		mockmailer := SetupMailerMock(t)
		handler.mailer = mockmailer
		ctx, _ = handler.sessMgr.Load(ctx, "")
		// copy user to avoid context user being modified
		// impacting future tests
		u := *user
		user := &u
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		handler.templates = &resources.TemplateMap{
			"mail_account_email_verify.gohtml": verifyTpl,
		}

		mock.EXPECT().
			NewUserVerify(ctx, user.ID).
			Return(uv, nil)
		mockmailer.EXPECT().
			SendAsync("",
				[]string{user.Email},
				"Account Verification",
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
				after, found := strings.CutPrefix(msgPlain, "Account Verification: http://example.com/verify/")
				assert.True(t, found)
				refParts := strings.Split(after, "-")
				rID := util.Must(service.ParseUserVerifyRefID(refParts[0]))
				hmacBytes, err := encoder.Base32DecodeString(refParts[1])
				assert.Nil(t, err)
				assert.True(t, handler.cMAC.Validate([]byte(rID.String()), hmacBytes))
			},
		)

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/send-verification", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.VerifySendEmail(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/settings",
			"handler returned wrong redirect")
		// we make sure that all expectations were met
	})
}

func TestHandler_VerifyEmail(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		ID:           1,
		RefID:        util.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     false,
		PWAuth:       true,
		WebAuthn:     false,
		Created:      ts,
		LastModified: ts,
	}

	uv := &model.UserVerify{
		RefID:   util.Must(model.NewUserVerifyRefID()),
		UserID:  user.ID,
		Created: ts,
	}

	t.Run("verify", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		// copy user to avoid context user being modified
		// impacting future tests
		u := *user
		user := &u
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		mock.EXPECT().
			GetUserVerifyByRefID(ctx, uv.RefID).
			Return(uv, nil)
		mock.EXPECT().
			SetUserVerified(ctx, user, uv).
			Return(nil)

		req, _ := http.NewRequestWithContext(ctx, "GET", "http://example.com/verify", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("uvRefID", uv.RefID.String())
		// generate/add hmac
		req.SetPathValue("hmac", encoder.Base32EncodeToString(
			handler.cMAC.Generate([]byte(uv.RefID.String()))))
		rr := httptest.NewRecorder()
		handler.VerifyEmail(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/settings",
			"handler returned wrong redirect")
		// we make sure that all expectations were met
	})

	t.Run("verify with bad hmac", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		// copy user to avoid context user being modified
		// impacting future tests
		u := *user
		user := &u
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		// generate hmac
		macBytes := handler.cMAC.Generate([]byte(uv.RefID.String()))
		macBytes[0] += 1
		// base32 encode hmac
		macStr := encoder.Base32EncodeToString(macBytes)

		req, _ := http.NewRequestWithContext(ctx, "GET", "http://example.com/verify", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("uvRefID", uv.RefID.String())
		req.SetPathValue("hmac", macStr)
		rr := httptest.NewRecorder()
		handler.VerifyEmail(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
	})

	t.Run("verify with bad refid", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		// copy user to avoid context user being modified
		// impacting future tests
		u := *user
		user := &u
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		refID := util.Must(model.NewEventItemRefID())

		// generate hmac
		macBytes := handler.cMAC.Generate([]byte(refID.String()))
		// base32 encode hmac
		macStr := encoder.Base32EncodeToString(macBytes)

		req, _ := http.NewRequestWithContext(ctx, "GET", "http://example.com/verify", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("uvRefID", refID.String())
		req.SetPathValue("hmac", macStr)
		rr := httptest.NewRecorder()
		handler.VerifyEmail(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})

	t.Run("verify not in db", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		// copy user to avoid context user being modified
		// impacting future tests
		u := *user
		user := &u
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		// generate hmac
		macBytes := handler.cMAC.Generate([]byte(uv.RefID.String()))
		// base32 encode hmac
		macStr := encoder.Base32EncodeToString(macBytes)

		mock.EXPECT().
			GetUserVerifyByRefID(ctx, uv.RefID).
			Return(nil, errs.NotFound.Error("verify not found"))

		req, _ := http.NewRequestWithContext(ctx, "GET", "http://example.com/verify", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("uvRefID", uv.RefID.String())
		req.SetPathValue("hmac", macStr)
		rr := httptest.NewRecorder()
		handler.VerifyEmail(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})

	t.Run("verify is expired", func(t *testing.T) {
		t.Parallel()

		refID := util.Must(model.NewUserVerifyRefID())
		rfts, _ := time.Parse(time.RFC3339, "2023-01-14T18:29:00Z")
		refID.SetTime(rfts)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		// copy user to avoid context user being modified
		// impacting future tests
		u := *user
		user := &u
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		// generate hmac
		macBytes := handler.cMAC.Generate([]byte(refID.String()))
		// base32 encode hmac
		macStr := encoder.Base32EncodeToString(macBytes)

		mock.EXPECT().
			GetUserVerifyByRefID(ctx, refID).
			Return(
				&model.UserVerify{
					UserID:  user.ID,
					RefID:   refID,
					Created: rfts,
				}, nil,
			)

		req, _ := http.NewRequestWithContext(ctx, "GET", "http://example.com/verify", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("uvRefID", refID.String())
		req.SetPathValue("hmac", macStr)
		rr := httptest.NewRecorder()
		handler.VerifyEmail(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
	})
}
