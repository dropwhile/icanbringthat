// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package handler

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/gorilla/csrf"
	"github.com/jaytaylor/html2text"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/app/service"
	"github.com/dropwhile/icanbringthat/internal/encoder"
	"github.com/dropwhile/icanbringthat/internal/errs"
	"github.com/dropwhile/icanbringthat/internal/mail"
	"github.com/dropwhile/icanbringthat/internal/middleware/auth"
)

func (x *Handler) ForgotPasswordShowForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	_, err := auth.UserFromContext(ctx)
	// already a logged in user, redirect to
	if err == nil {
		http.Redirect(w, r, "/settings", http.StatusSeeOther)
		return
	}

	tplVars := MapSA{
		"title":          "Forgot Password",
		"next":           r.FormValue("next"),
		"flashes":        x.sessMgr.FlashPopAll(ctx),
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
	}
	// render user profile view
	w.Header().Set("content-type", "text/html")
	err = x.TemplateExecute(w, "forgot-password-form.gohtml", tplVars)
	if err != nil {
		x.TemplateError(w)
		return
	}
}

func (x *Handler) PasswordResetShowForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	_, err := auth.UserFromContext(ctx)
	// already a logged in user, redirect to
	if err == nil {
		http.Redirect(w, r, "/settings", http.StatusSeeOther)
		return
	}

	hmacStr := r.PathValue("hmac")
	refIDStr := r.PathValue("upwRefID")
	if hmacStr == "" || refIDStr == "" {
		x.NotFoundError(w)
		return
	}

	// decode hmac
	hmacBytes, err := encoder.Base32DecodeString(hmacStr)
	if err != nil {
		slog.DebugContext(ctx, "error decoding hmac data", "error", err)
		x.BadRequestError(w, "Bad Request Data")
		return
	}
	// check hmac
	if !x.cMAC.Validate([]byte(refIDStr), hmacBytes) {
		slog.DebugContext(ctx, "invalid hmac!")
		x.BadRequestError(w, "Bad Request Data")
		return
	}

	// hmac checks out. ok to parse refid now.
	refID, err := service.ParseUserPWResetRefID(refIDStr)
	if err != nil {
		slog.DebugContext(ctx, "bad refid", "error", err)
		x.BadRequestError(w, "Bad Request Data")
		x.BadRefIDError(w, "verify", err)
		return
	}

	upw, errx := x.svc.GetUserPWResetByRefID(ctx, refID)
	if errx != nil {
		slog.DebugContext(ctx, "no upw match", "error", errx)
		x.BadRequestError(w, "Bad Request Data")
		return
	}

	if service.IsTimerExpired(upw.RefID, model.UserPWResetExpiry) {
		slog.DebugContext(ctx, "token expired")
		x.NotFoundError(w)
		return
	}

	_, err = x.svc.GetUserByID(ctx, upw.UserID)
	if err != nil {
		slog.DebugContext(ctx, "no user match", "error", err)
		x.BadRequestError(w, "Bad Request Data")
		return
	}

	tplVars := MapSA{
		"title":          "Reset Password",
		"next":           r.FormValue("next"),
		"flashes":        x.sessMgr.FlashPopAll(ctx),
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
		"refID":          refIDStr,
		"hmac":           hmacStr,
	}
	// render user profile view
	w.Header().Set("content-type", "text/html")
	err = x.TemplateExecute(w, "password-reset-form.gohtml", tplVars)
	if err != nil {
		x.TemplateError(w)
		return
	}
}

func (x *Handler) ResetPasswordSendEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// attempt to get user from session
	if _, err := auth.UserFromContext(ctx); err == nil {
		// already a logged in user, reject password reset
		x.AccessDeniedError(w)
		return
	}

	email := r.PostFormValue("email")
	if email == "" {
		x.BadFormDataError(w, nil, "email")
		return
	}

	// don't leak existence of user. if email doens't match,
	// behave like we sent a reset anyway...
	doFake := false
	user, errx := x.svc.GetUserByEmail(ctx, email)
	if errx != nil {
		switch errx.Code() {
		case errs.NotFound:
			slog.InfoContext(ctx, "no user found", "error", errx)
			doFake = true
		default:
			x.InternalServerError(w, errx.Msg())
		}
	}

	// if pw auth is disabled, behave the same as faking it
	// (do not send email either)
	if !doFake && !user.PWAuth {
		doFake = true
	}

	if doFake {
		slog.InfoContext(ctx,
			"pretending to sent password reset email",
			slog.String("email", email),
		)
	}

	if !doFake {
		// generate a upw
		upw, errx := x.svc.NewUserPWReset(ctx, user.ID)
		if errx != nil {
			x.InternalServerError(w, errx.Msg())
			return
		}
		upwRefIDStr := upw.RefID.String()

		// generate hmac
		macBytes := x.cMAC.Generate([]byte(upwRefIDStr))
		// base32 encode hmac
		macStr := encoder.Base32EncodeToString(macBytes)

		// construct url
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		u := &url.URL{
			Scheme: scheme,
			Host:   r.Host,
		}
		u = u.JoinPath(fmt.Sprintf("/forgot-password/%s-%s", upwRefIDStr, macStr))

		// construct email
		subject := "Password reset"
		var buf bytes.Buffer
		err := x.TemplateExecute(&buf, "mail_password_reset.gohtml",
			MapSA{
				"Subject":          subject,
				"PasswordResetUrl": u.String(),
			},
		)
		if err != nil {
			x.TemplateError(w)
			return
		}
		messageHtml := buf.String()

		messagePlain, err := html2text.FromString(messageHtml)
		if err != nil {
			x.TemplateError(w)
		}

		slog.DebugContext(ctx, "email content",
			slog.String("plain", messagePlain),
			slog.String("html", messageHtml),
		)

		_ = user
		x.mailer.SendAsync("", []string{user.Email},
			subject, messagePlain, messageHtml,
			mail.MailHeader{
				"X-PM-Message-Stream": "outbound",
			},
		)
	}

	x.sessMgr.FlashAppend(ctx, "success", "Password reset email sent.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (x *Handler) PasswordReset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// attempt to get user from session
	if _, err := auth.UserFromContext(ctx); err == nil {
		// already a logged in user, reject password reset
		x.AccessDeniedError(w)
		return
	}

	hmacStr := r.PathValue("hmac")
	refIDStr := r.PathValue("upwRefID")
	if hmacStr == "" || refIDStr == "" {
		slog.DebugContext(ctx, "missing url query data")
		x.NotFoundError(w)
		return
	}

	newPasswd := r.PostFormValue("password")
	confirmPassword := r.PostFormValue("confirm_password")
	if newPasswd == "" || newPasswd != confirmPassword {
		x.BadFormDataError(w, nil)
		return
	}

	// decode hmac
	hmacBytes, err := encoder.Base32DecodeString(hmacStr)
	if err != nil {
		slog.DebugContext(ctx, "error decoding hmac data", "error", err)
		x.BadRequestError(w, "Bad Request Data")
		return
	}
	// check hmac
	if !x.cMAC.Validate([]byte(refIDStr), hmacBytes) {
		slog.DebugContext(ctx, "invalid hmac!")
		x.BadRequestError(w, "Bad Request Data")
		return
	}

	// hmac checks out. ok to parse refid now.
	refID, err := service.ParseUserPWResetRefID(refIDStr)
	if err != nil {
		x.BadRefIDError(w, "upw", err)
		return
	}

	upw, errx := x.svc.GetUserPWResetByRefID(ctx, refID)
	if errx != nil {
		slog.DebugContext(ctx, "no upw match", "error", errx)
		x.BadRequestError(w, "Bad Request Data")
		return
	}

	if service.IsTimerExpired(upw.RefID, model.UserPWResetExpiry) {
		slog.DebugContext(ctx, "token expired")
		x.NotFoundError(w)
		return
	}

	user, errx := x.svc.GetUserByID(ctx, upw.UserID)
	if errx != nil {
		slog.DebugContext(ctx, "no user match", "error", errx)
		x.BadRequestError(w, "Bad Request Data")
		return
	}

	// if pw auth is disabled, do not send email either
	if !user.PWAuth {
		slog.InfoContext(ctx, "pw reset attempt but pw auth disabled")
		x.AccessDeniedError(w)
		return
	}

	pwHash, err := model.HashPass(ctx, []byte(newPasswd))
	if err != nil {
		slog.DebugContext(ctx, "error updating password", "error", err)
		x.InternalServerError(w, "error updating user password")
		return
	}
	user.PWHash = pwHash

	errx = x.svc.UpdateUserPWReset(ctx, user, upw)
	if errx != nil {
		x.DBError(w, errx)
		return
	}

	// renew sesmgr token to help prevent session fixation. ref:
	//   https://github.com/OWASP/CheatSheetSeries/blob/master/cheatsheets/Session_Management_Cheat_Sheet.md
	//   #renew-the-session-id-after-any-privilege-level-change
	err = x.sessMgr.RenewToken(ctx)
	if err != nil {
		x.InternalServerError(w, "Session Error")
		return
	}
	// Then make the privilege-level change.
	x.sessMgr.Put(r.Context(), "user-id", user.ID)
	target := "/dashboard"
	http.Redirect(w, r, target, http.StatusSeeOther)
}
