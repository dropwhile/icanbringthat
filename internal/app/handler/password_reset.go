package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/internal/encoder"
	"github.com/dropwhile/icbt/internal/errs"
	"github.com/dropwhile/icbt/internal/mail"
)

func (x *Handler) ShowForgotPasswordForm(w http.ResponseWriter, r *http.Request) {
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
		"flashes":        x.SessMgr.FlashPopAll(ctx),
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

func (x *Handler) ShowPasswordResetForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	_, err := auth.UserFromContext(ctx)
	// already a logged in user, redirect to
	if err == nil {
		http.Redirect(w, r, "/settings", http.StatusSeeOther)
		return
	}

	hmacStr := chi.URLParam(r, "hmac")
	refIDStr := chi.URLParam(r, "upwRefID")
	if hmacStr == "" || refIDStr == "" {
		x.NotFoundError(w)
		return
	}

	// decode hmac
	hmacBytes, err := encoder.Base32DecodeString(hmacStr)
	if err != nil {
		log.Debug().Err(err).Msg("error decoding hmac data")
		x.BadRequestError(w, "Bad Request Data")
		return
	}
	// check hmac
	if !x.MAC.Validate([]byte(refIDStr), hmacBytes) {
		log.Debug().Msg("invalid hmac!")
		x.BadRequestError(w, "Bad Request Data")
		return
	}

	// hmac checks out. ok to parse refid now.
	refID, err := model.ParseUserPWResetRefID(refIDStr)
	if err != nil {
		log.Debug().Err(err).Msg("bad refid")
		x.BadRequestError(w, "Bad Request Data")
		x.BadRefIDError(w, "verify", err)
		return
	}

	upw, errx := service.GetUserPWResetByRefID(ctx, x.Db, refID)
	if errx != nil {
		log.Debug().Err(err).Msg("no upw match")
		x.BadRequestError(w, "Bad Request Data")
		return
	}

	if model.IsExpired(upw.RefID, model.UserPWResetExpiry) {
		log.Debug().Err(err).Msg("token expired")
		x.NotFoundError(w)
		return
	}

	_, err = service.GetUserByID(ctx, x.Db, upw.UserID)
	if err != nil {
		log.Debug().Err(err).Msg("no user match")
		x.BadRequestError(w, "Bad Request Data")
		return
	}

	tplVars := MapSA{
		"title":          "Reset Password",
		"next":           r.FormValue("next"),
		"flashes":        x.SessMgr.FlashPopAll(ctx),
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

func (x *Handler) SendResetPasswordEmail(w http.ResponseWriter, r *http.Request) {
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
	user, errx := service.GetUserByEmail(ctx, x.Db, email)
	if errx != nil {
		switch errx.Code() {
		case errs.NotFound:
			log.Info().Err(errx).Msg("no user found")
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
		log.Info().
			Str("email", email).
			Msg("pretending to sent password reset email")
	}

	if !doFake {
		// generate a upw
		upw, errx := service.NewUserPWReset(ctx, x.Db, user.ID)
		if errx != nil {
			x.InternalServerError(w, errx.Msg())
			return
		}
		upwRefIDStr := upw.RefID.String()

		// generate hmac
		macBytes := x.MAC.Generate([]byte(upwRefIDStr))
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
		err := x.TemplateExecute(&buf, "mail_password_reset.gotxt",
			MapSA{
				"Subject":          subject,
				"PasswordResetUrl": u.String(),
			},
		)
		if err != nil {
			x.TemplateError(w)
			return
		}
		messagePlain := buf.String()

		buf.Reset()
		err = x.TemplateExecute(&buf, "mail_password_reset.gohtml",
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

		log.Debug().
			Str("plain", messagePlain).
			Str("html", messageHtml).
			Msg("email content")

		_ = user
		x.Mailer.SendAsync("", []string{user.Email},
			subject, messagePlain, messageHtml,
			mail.MailHeader{
				"X-PM-Message-Stream": "outbound",
			},
		)
	}

	x.SessMgr.FlashAppend(ctx, "success", "Password reset email sent.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (x *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// attempt to get user from session
	if _, err := auth.UserFromContext(ctx); err == nil {
		// already a logged in user, reject password reset
		x.AccessDeniedError(w)
		return
	}

	hmacStr := chi.URLParam(r, "hmac")
	refIDStr := chi.URLParam(r, "upwRefID")
	if hmacStr == "" || refIDStr == "" {
		log.Debug().Msg("missing url query data")
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
		log.Debug().Err(err).Msg("error decoding hmac data")
		x.BadRequestError(w, "Bad Request Data")
		return
	}
	// check hmac
	if !x.MAC.Validate([]byte(refIDStr), hmacBytes) {
		log.Debug().Msg("invalid hmac!")
		x.BadRequestError(w, "Bad Request Data")
		return
	}

	// hmac checks out. ok to parse refid now.
	refID, err := model.ParseUserPWResetRefID(refIDStr)
	if err != nil {
		x.BadRefIDError(w, "upw", err)
		return
	}

	upw, errx := service.GetUserPWResetByRefID(ctx, x.Db, refID)
	if errx != nil {
		log.Debug().Err(err).Msg("no upw match")
		x.BadRequestError(w, "Bad Request Data")
		return
	}

	if model.IsExpired(upw.RefID, model.UserPWResetExpiry) {
		log.Debug().Err(err).Msg("token expired")
		x.NotFoundError(w)
		return
	}

	user, errx := service.GetUserByID(ctx, x.Db, upw.UserID)
	if errx != nil {
		log.Debug().Err(err).Msg("no user match")
		x.BadRequestError(w, "Bad Request Data")
		return
	}

	// if pw auth is disabled, do not send email either
	if !user.PWAuth {
		log.Info().Msg("pw reset attempt but pw auth disabled")
		x.AccessDeniedError(w)
		return
	}

	pwHash, err := model.HashPass(ctx, []byte(newPasswd))
	if err != nil {
		log.Debug().Err(err).Msg("error updating password")
		x.InternalServerError(w, "error updating user password")
		return
	}
	user.PWHash = pwHash

	errx = service.UpdateUserPWReset(ctx, x.Db, user, upw)
	if errx != nil {
		x.DBError(w, errx)
		return
	}

	// renew sesmgr token to help prevent session fixation. ref:
	//   https://github.com/OWASP/CheatSheetSeries/blob/master/cheatsheets/Session_Management_Cheat_Sheet.md
	//   #renew-the-session-id-after-any-privilege-level-change
	err = x.SessMgr.RenewToken(ctx)
	if err != nil {
		x.InternalServerError(w, "Session Error")
		return
	}
	// Then make the privilege-level change.
	x.SessMgr.Put(r.Context(), "user-id", user.ID)
	target := "/dashboard"
	http.Redirect(w, r, target, http.StatusSeeOther)
}
