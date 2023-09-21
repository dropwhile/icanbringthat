package handler

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/util"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func (h *Handler) ShowForgotPasswordForm(w http.ResponseWriter, r *http.Request) {
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
		"flashes":        h.SessMgr.FlashPopKey(ctx, "forgot-password"),
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
	}
	// render user profile view
	w.Header().Set("content-type", "text/html")
	err = h.TemplateExecute(w, "forgot-password-form.gohtml", tplVars)
	if err != nil {
		log.Debug().Err(err).Msg("template error")
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ShowPasswordResetForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	_, err := auth.UserFromContext(ctx)
	// already a logged in user, redirect to
	if err == nil {
		http.Redirect(w, r, "/settings", http.StatusSeeOther)
		return
	}

	hmacStr := chi.URLParam(r, "hmac")
	refIdStr := chi.URLParam(r, "upwRefId")
	if hmacStr == "" || refIdStr == "" {
		log.Debug().Msg("missing url query data")
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	// decode hmac
	hmacBytes, err := util.Base32DecodeString(hmacStr)
	if err != nil {
		log.Info().Err(err).Msg("error decoding hmac data")
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}
	// check hmac
	if !h.Hmac.Validate([]byte(refIdStr), hmacBytes) {
		log.Info().Msg("invalid hmac!")
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	// hmac checks out. ok to parse refid now.
	refId, err := model.UserPWResetRefIdT.Parse(refIdStr)
	if err != nil {
		log.Info().Err(err).Msg("bad refid")
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	if refId.Time().Add(30 * time.Minute).Before(time.Now()) {
		log.Debug().Err(err).Msg("token expired")
		http.Error(w, "token expired", http.StatusBadRequest)
		return
	}

	upw, err := model.GetUserPWResetByRefId(ctx, h.Db, refId)
	if err != nil {
		log.Debug().Err(err).Msg("no upw match")
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	_, err = model.GetUserById(ctx, h.Db, upw.UserId)
	if err != nil {
		log.Debug().Err(err).Msg("no user match")
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	tplVars := map[string]any{
		"title":          "Reset Password",
		"next":           r.FormValue("next"),
		"flashes":        h.SessMgr.FlashPopKey(ctx, "reset-password"),
		csrf.TemplateTag: csrf.TemplateField(r),
		"csrfToken":      csrf.Token(r),
		"refId":          refIdStr,
		"hmac":           hmacStr,
	}
	// render user profile view
	w.Header().Set("content-type", "text/html")
	err = h.TemplateExecute(w, "password-reset-form.gohtml", tplVars)
	if err != nil {
		log.Debug().Err(err).Msg("template error")
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) SendResetPasswordEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Debug().Msg("test")

	// attempt to get user from session
	if _, err := auth.UserFromContext(ctx); err == nil {
		// already a logged in user, reject password reset
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	email := r.PostFormValue("email")
	if email == "" {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	// don't leak existence of user. if email doens't match,
	// behave like we sent a reset anyway...
	doFake := false
	user, err := model.GetUserByEmail(ctx, h.Db, email)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Info().Err(err).Msg("no user found")
		doFake = true
	case err != nil:
		log.Info().Err(err).Msg("db error")
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if doFake {
		log.Info().Str("email", email).Msg("pretending to sent password reset email")
	}

	if !doFake {
		// generate a upw
		upw, err := model.NewUserPWReset(ctx, h.Db, user)
		if err != nil {
			log.Info().Err(err).Msg("db error")
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		upwRefIdStr := upw.RefId.String()

		// generate hmac
		macBytes := h.Hmac.Generate([]byte(upwRefIdStr))
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
		u = u.JoinPath(fmt.Sprintf("/forgot-password/%s-%s", upwRefIdStr, macStr))

		// construct email
		subject := "Password reset"
		var buf bytes.Buffer
		err = h.TemplateExecute(&buf, "mail_password_reset.gohtml",
			map[string]any{
				"Subject":          subject,
				"PasswordResetUrl": u.String(),
			},
		)
		if err != nil {
			http.Error(w, "template error", http.StatusInternalServerError)
			return
		}
		messagePlain := fmt.Sprintf("Password reset url: %s", u.String())
		messageHtml := buf.String()
		log.Debug().
			Str("plain", messagePlain).
			Str("html", messageHtml).
			Msg("email content")

		_ = user
		/*go func() {
			err := h.Mailer.Send("", []string{user.Email}, subject, messagePlain, messageHtml)
			if err != nil {
				log.Info().Err(err).Msg("error sending email")
			}
		}()
		*/
	}

	h.SessMgr.FlashAppend(ctx, "login", "Password reset email sent.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	_, err := auth.UserFromContext(ctx)
	// already a logged in user
	if err == nil {
		http.Redirect(w, r, "/settings", http.StatusSeeOther)
		return
	}

	hmacStr := chi.URLParam(r, "hmac")
	refIdStr := chi.URLParam(r, "upwRefId")
	if hmacStr == "" || refIdStr == "" {
		log.Debug().Msg("missing url query data")
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	newPasswd := r.PostFormValue("password")
	confirmPassword := r.PostFormValue("confirm_password")
	if newPasswd == "" || newPasswd != confirmPassword {
		log.Debug().Msg("bad form data")
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	// decode hmac
	hmacBytes, err := util.Base32DecodeString(hmacStr)
	if err != nil {
		log.Info().Err(err).Msg("error decoding hmac data")
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}
	// check hmac
	if !h.Hmac.Validate([]byte(refIdStr), hmacBytes) {
		log.Info().Msg("invalid hmac!")
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	// hmac checks out. ok to parse refid now.
	refId, err := model.UserPWResetRefIdT.Parse(refIdStr)
	if err != nil {
		log.Info().Err(err).Msg("bad refid")
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	if refId.Time().Add(30 * time.Minute).Before(time.Now()) {
		log.Debug().Err(err).Msg("token expired")
		http.Error(w, "token expired", http.StatusBadRequest)
		return
	}

	upw, err := model.GetUserPWResetByRefId(ctx, h.Db, refId)
	if err != nil {
		log.Debug().Err(err).Msg("no upw match")
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	user, err := model.GetUserById(ctx, h.Db, upw.UserId)
	if err != nil {
		log.Debug().Err(err).Msg("no user match")
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	err = user.SetPass(ctx, []byte(newPasswd))
	if err != nil {
		log.Debug().Err(err).Msg("error updating password")
		http.Error(w, "error updating user password", http.StatusInternalServerError)
		return
	}
	err = user.Save(ctx, h.Db)
	if err != nil {
		log.Debug().Err(err).Msg("db error")
		http.Error(w, "error updating user password", http.StatusInternalServerError)
		return
	}

	err = upw.Delete(ctx, h.Db)
	if err != nil {
		log.Info().Err(err).Msg("error cleaning up pw reset token")
		// do not throw a 500 level error here though, just carry on since
		// the user was already found..
	}

	// renew sesmgr token to help prevent session fixation. ref:
	//   https://github.com/OWASP/CheatSheetSeries/blob/master/cheatsheets/Session_Management_Cheat_Sheet.md
	//   #renew-the-session-id-after-any-privilege-level-change
	err = h.SessMgr.RenewToken(ctx)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// Then make the privilege-level change.
	h.SessMgr.Put(r.Context(), "user-id", user.Id)
	target := "/dashboard"
	http.Redirect(w, r, target, http.StatusSeeOther)
}
