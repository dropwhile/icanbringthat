package xhandler

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/util"
)

func (x *XHandler) ShowForgotPasswordForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	_, err := auth.UserFromContext(ctx)
	// already a logged in user, redirect to
	if err == nil {
		http.Redirect(w, r, "/settings", http.StatusSeeOther)
		return
	}

	tplVars := map[string]any{
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
		log.Debug().Err(err).Msg("template error")
		x.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (x *XHandler) ShowPasswordResetForm(w http.ResponseWriter, r *http.Request) {
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
		log.Debug().Msg("missing url query data")
		x.Error(w, "not found", http.StatusNotFound)
		return
	}

	// decode hmac
	hmacBytes, err := util.Base32DecodeString(hmacStr)
	if err != nil {
		log.Info().Err(err).Msg("error decoding hmac data")
		x.Error(w, "bad data", http.StatusNotFound)
		return
	}
	// check hmac
	if !x.Hmac.Validate([]byte(refIDStr), hmacBytes) {
		log.Info().Msg("invalid hmac!")
		x.Error(w, "bad data", http.StatusNotFound)
		return
	}

	// hmac checks out. ok to parse refid now.
	refID, err := model.ParseUserPWResetRefID(refIDStr)
	if err != nil {
		log.Info().Err(err).Msg("bad refid")
		x.Error(w, "bad data", http.StatusNotFound)
		return
	}

	upw, err := model.GetUserPWResetByRefID(ctx, x.Db, refID)
	if err != nil {
		log.Debug().Err(err).Msg("no upw match")
		x.Error(w, "bad data", http.StatusNotFound)
		return
	}

	if upw.IsExpired() {
		log.Debug().Err(err).Msg("token expired")
		x.Error(w, "token expired", http.StatusNotFound)
		return
	}

	_, err = model.GetUserByID(ctx, x.Db, upw.UserID)
	if err != nil {
		log.Debug().Err(err).Msg("no user match")
		x.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	tplVars := map[string]any{
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
		log.Debug().Err(err).Msg("template error")
		x.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (x *XHandler) SendResetPasswordEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// attempt to get user from session
	if _, err := auth.UserFromContext(ctx); err == nil {
		// already a logged in user, reject password reset
		x.Error(w, "access denied", http.StatusForbidden)
		return
	}

	email := r.PostFormValue("email")
	if email == "" {
		x.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	// don't leak existence of user. if email doens't match,
	// behave like we sent a reset anyway...
	doFake := false
	user, err := model.GetUserByEmail(ctx, x.Db, email)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Info().Err(err).Msg("no user found")
		doFake = true
	case err != nil:
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if doFake {
		log.Info().Str("email", email).Msg("pretending to sent password reset email")
	}

	if !doFake {
		// generate a upw
		upw, err := model.NewUserPWReset(ctx, x.Db, user.ID)
		if err != nil {
			log.Info().Err(err).Msg("db error")
			x.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		upwRefIDStr := upw.RefID.String()

		// generate hmac
		macBytes := x.Hmac.Generate([]byte(upwRefIDStr))
		// base32 encode hmac
		macStr := util.Base32EncodeToString(macBytes)

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
		err = x.TemplateExecute(&buf, "mail_password_reset.gohtml",
			map[string]any{
				"Subject":          subject,
				"PasswordResetUrl": u.String(),
			},
		)
		if err != nil {
			x.Error(w, "template error", http.StatusInternalServerError)
			return
		}
		messagePlain := fmt.Sprintf("Password reset url: %s", u.String())
		messageHtml := buf.String()
		log.Debug().
			Str("plain", messagePlain).
			Str("html", messageHtml).
			Msg("email content")

		_ = user
		x.Mailer.SendAsync("", []string{user.Email}, subject, messagePlain, messageHtml)
	}

	x.SessMgr.FlashAppend(ctx, "success", "Password reset email sent.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (x *XHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// attempt to get user from session
	if _, err := auth.UserFromContext(ctx); err == nil {
		// already a logged in user, reject password reset
		x.Error(w, "access denied", http.StatusForbidden)
		return
	}

	hmacStr := chi.URLParam(r, "hmac")
	refIDStr := chi.URLParam(r, "upwRefID")
	if hmacStr == "" || refIDStr == "" {
		log.Debug().Msg("missing url query data")
		x.Error(w, "not found", http.StatusNotFound)
		return
	}

	newPasswd := r.PostFormValue("password")
	confirmPassword := r.PostFormValue("confirm_password")
	if newPasswd == "" || newPasswd != confirmPassword {
		log.Debug().Msg("bad form data")
		x.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	// decode hmac
	hmacBytes, err := util.Base32DecodeString(hmacStr)
	if err != nil {
		log.Info().Err(err).Msg("error decoding hmac data")
		x.Error(w, "bad data", http.StatusBadRequest)
		return
	}
	// check hmac
	if !x.Hmac.Validate([]byte(refIDStr), hmacBytes) {
		log.Info().Msg("invalid hmac!")
		x.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	// hmac checks out. ok to parse refid now.
	refID, err := model.ParseUserPWResetRefID(refIDStr)
	if err != nil {
		log.Info().Err(err).Msg("bad refid")
		x.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	upw, err := model.GetUserPWResetByRefID(ctx, x.Db, refID)
	if err != nil {
		log.Debug().Err(err).Msg("no upw match")
		x.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	if upw.IsExpired() {
		log.Debug().Err(err).Msg("token expired")
		x.Error(w, "token expired", http.StatusBadRequest)
		return
	}

	user, err := model.GetUserByID(ctx, x.Db, upw.UserID)
	if err != nil {
		log.Debug().Err(err).Msg("no user match")
		x.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	pwHash, err := model.HashPass(ctx, []byte(newPasswd))
	if err != nil {
		log.Debug().Err(err).Msg("error updating password")
		x.Error(w, "error updating user password", http.StatusInternalServerError)
		return
	}
	user.PWHash = pwHash

	err = pgx.BeginFunc(ctx, x.Db, func(tx pgx.Tx) error {
		innerErr := model.UpdateUser(ctx, tx, user.Email, user.Name, user.PWHash, user.Verified, user.ID)
		if innerErr != nil {
			log.Debug().Err(innerErr).Msg("inner db error saving user")
			return innerErr
		}

		innerErr = model.DeleteUserPWReset(ctx, tx, upw.RefID)
		if innerErr != nil {
			log.Debug().Err(innerErr).Msg("inner db error cleaning up pw reset token")
			return innerErr
		}
		return nil
	})
	if err != nil {
		log.Debug().Err(err).Msg("db error")
		x.Error(w, "error updating user password", http.StatusInternalServerError)
		return
	}

	// renew sesmgr token to help prevent session fixation. ref:
	//   https://github.com/OWASP/CheatSheetSeries/blob/master/cheatsheets/Session_Management_Cheat_Sheet.md
	//   #renew-the-session-id-after-any-privilege-level-change
	err = x.SessMgr.RenewToken(ctx)
	if err != nil {
		x.Error(w, err.Error(), 500)
		return
	}
	// Then make the privilege-level change.
	x.SessMgr.Put(r.Context(), "user-id", user.ID)
	target := "/dashboard"
	http.Redirect(w, r, target, http.StatusSeeOther)
}
