package xhandler

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/util"
	"github.com/dropwhile/icbt/internal/util/htmx"
)

func (x *XHandler) SendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// attempt to get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		log.Debug().Err(err).Msg("bad session data")
		x.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	// generate a verifier
	uv, err := model.NewUserVerify(ctx, x.Db, user)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	uvRefIDStr := uv.RefID.String()

	// generate hmac
	macBytes := x.Hmac.Generate([]byte(uvRefIDStr))
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
	u = u.JoinPath(fmt.Sprintf("/verify/%s-%s", uvRefIDStr, macStr))

	// construct email
	subject := "Account Verification"
	var buf bytes.Buffer
	err = x.TemplateExecute(&buf, "mail_account_email_verify.gohtml",
		map[string]any{
			"Subject":         subject,
			"VerificationUrl": u.String(),
		},
	)
	if err != nil {
		x.Error(w, "template error", http.StatusInternalServerError)
		return
	}
	messagePlain := fmt.Sprintf("Account Verification url: %s", u.String())
	messageHtml := buf.String()
	log.Debug().
		Str("plain", messagePlain).
		Str("html", messageHtml).
		Msg("email content")

	_ = user
	x.Mailer.SendAsync("", []string{user.Email}, subject, messagePlain, messageHtml)
	x.SessMgr.FlashAppend(ctx, "success", "Account verification email sent.")
	if htmx.Hx(r).Request() {
		w.Header().Add("HX-Refresh", "true")
		w.WriteHeader(200)
		return
	}
	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}

func (x *XHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	hmacStr := chi.URLParam(r, "hmac")
	refIDStr := chi.URLParam(r, "uvRefID")
	if hmacStr == "" || refIDStr == "" {
		log.Debug().Msg("missing url query data")
		x.Error(w, "not found", http.StatusNotFound)
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
		x.Error(w, "bad data", http.StatusNotFound)
		return
	}

	// hmac checks out. ok to parse refid now.
	verifyRefID, err := model.ParseUserVerifyRefID(refIDStr)
	if err != nil {
		x.Error(w, "bad verify-ref-id", http.StatusNotFound)
		return
	}

	verifier, err := model.GetUserVerifyByRefID(ctx, x.Db, verifyRefID)
	if err != nil {
		log.Debug().Err(err).Msg("no verifier match")
		x.Error(w, "bad data", http.StatusNotFound)
		return
	}

	if model.IsExpired(verifier.RefID, model.UserVerifyExpiry) {
		log.Debug().Err(err).Msg("verifier is expired")
		x.Error(w, "bad data", http.StatusNotFound)
		return
	}

	user.Verified = true
	err = pgx.BeginFunc(ctx, x.Db, func(tx pgx.Tx) error {
		innerErr := model.UpdateUser(ctx, tx,
			user.Email, user.Name, user.PWHash,
			user.Verified, user.PWAuth, user.WebAuthn, user.ID,
		)
		if innerErr != nil {
			log.Debug().Err(innerErr).Msg("inner db error saving user")
			return innerErr
		}

		innerErr = model.DeleteUserVerify(ctx, x.Db, verifier.RefID)
		if innerErr != nil {
			log.Debug().Err(innerErr).Msg("inner db error cleaning up verifier token")
			return innerErr
		}
		return nil
	})
	if err != nil {
		log.Debug().Err(err).Msg("db error")
		x.Error(w, "error saving verification", http.StatusInternalServerError)
		return
	}

	x.SessMgr.FlashAppend(ctx, "success", "Email verification successfull")
	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}
