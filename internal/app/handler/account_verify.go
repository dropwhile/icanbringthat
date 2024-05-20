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

	"github.com/jaytaylor/html2text"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/app/service"
	"github.com/dropwhile/icanbringthat/internal/encoder"
	"github.com/dropwhile/icanbringthat/internal/htmx"
	"github.com/dropwhile/icanbringthat/internal/logger"
	"github.com/dropwhile/icanbringthat/internal/mail"
	"github.com/dropwhile/icanbringthat/internal/middleware/auth"
)

func (x *Handler) VerifySendEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// attempt to get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	// generate a verifier
	uv, errx := x.svc.NewUserVerify(ctx, user.ID)
	if errx != nil {
		x.InternalServerError(w, errx.Msg())
		return
	}
	uvRefIDStr := uv.RefID.String()

	// generate hmac
	macBytes := x.cMAC.Generate([]byte(uvRefIDStr))
	// base32 encode hmac
	macStr := encoder.Base32EncodeToString(macBytes)

	verificationUrl, err := url.JoinPath(x.baseURL, fmt.Sprintf("/verify/%s-%s", uvRefIDStr, macStr))
	if err != nil {
		x.InternalServerError(w, "processing error")
		return
	}

	// construct email
	subject := "Account Verification"
	var buf bytes.Buffer
	err = x.TemplateExecute(&buf, "mail_account_email_verify.gohtml",
		MapSA{
			"Subject":         subject,
			"VerificationUrl": verificationUrl,
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
	x.sessMgr.FlashAppend(ctx, "success", "Account verification email sent.")
	if htmx.Request(r).IsRequest() {
		htmx.Response(w).HxLocation(htmx.Request(r).CurrentUrl().Path)
		w.WriteHeader(200)
		return
	}
	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}

func (x *Handler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	hmacStr := r.PathValue("hmac")
	refIDStr := r.PathValue("uvRefID")
	if hmacStr == "" || refIDStr == "" {
		slog.DebugContext(ctx, "missing url query data")
		x.NotFoundError(w)
		return
	}

	// decode hmac
	hmacBytes, err := encoder.Base32DecodeString(hmacStr)
	if err != nil {
		slog.DebugContext(ctx, "error decoding hmac data", logger.Err(err))
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
	verifyRefID, err := service.ParseUserVerifyRefID(refIDStr)
	if err != nil {
		x.BadRefIDError(w, "verify", err)
		return
	}

	verifier, errx := x.svc.GetUserVerifyByRefID(ctx, verifyRefID)
	if errx != nil {
		slog.DebugContext(ctx, "no verifier match", logger.Err(errx))
		x.NotFoundError(w)
		return
	}

	if service.IsTimerExpired(verifier.RefID, model.UserVerifyExpiry) {
		slog.DebugContext(ctx, "verifier is expired")
		x.NotFoundError(w)
		return
	}

	user.Verified = true
	errx = x.svc.SetUserVerified(ctx, user, verifier)
	if errx != nil {
		slog.DebugContext(ctx, "error saving verification", logger.Err(errx))
		x.InternalServerError(w, errx.Msg())
		return
	}

	x.sessMgr.FlashAppend(ctx, "success", "Email verification successfull")
	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}
