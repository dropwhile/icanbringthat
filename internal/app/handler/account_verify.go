package handler

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
	"github.com/dropwhile/icbt/internal/encoder"
	"github.com/dropwhile/icbt/internal/htmx"
	"github.com/dropwhile/icbt/internal/mail"
)

func (x *Handler) SendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// attempt to get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	// generate a verifier
	uv, err := model.NewUserVerify(ctx, x.Db, user)
	if err != nil {
		x.DBError(w, err)
		return
	}
	uvRefIDStr := uv.RefID.String()

	// generate hmac
	macBytes := x.MAC.Generate([]byte(uvRefIDStr))
	// base32 encode hmac
	macStr := encoder.Base32EncodeToString(macBytes)

	verificationUrl, err := url.JoinPath(x.BaseURL, fmt.Sprintf("/verify/%s-%s", uvRefIDStr, macStr))
	if err != nil {
		x.InternalServerError(w, "processing error")
		return
	}

	// construct email
	subject := "Account Verification"
	var buf bytes.Buffer
	err = x.TemplateExecute(&buf, "mail_account_email_verify.gotxt",
		MapSA{
			"Subject":         subject,
			"VerificationUrl": verificationUrl,
		},
	)
	if err != nil {
		x.TemplateError(w)
		return
	}
	messagePlain := buf.String()

	buf.Reset()
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
	x.SessMgr.FlashAppend(ctx, "success", "Account verification email sent.")
	if htmx.Hx(r).Request() {
		w.Header().Add("HX-Refresh", "true")
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

	hmacStr := chi.URLParam(r, "hmac")
	refIDStr := chi.URLParam(r, "uvRefID")
	if hmacStr == "" || refIDStr == "" {
		log.Debug().Msg("missing url query data")
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
	verifyRefID, err := model.ParseUserVerifyRefID(refIDStr)
	if err != nil {
		x.BadRefIDError(w, "verify", err)
		return
	}

	verifier, err := model.GetUserVerifyByRefID(ctx, x.Db, verifyRefID)
	if err != nil {
		log.Debug().Err(err).Msg("no verifier match")
		x.NotFoundError(w)
		return
	}

	if model.IsExpired(verifier.RefID, model.UserVerifyExpiry) {
		log.Debug().Err(err).Msg("verifier is expired")
		x.NotFoundError(w)
		return
	}

	user.Verified = true
	err = pgx.BeginFunc(ctx, x.Db, func(tx pgx.Tx) error {
		innerErr := model.UpdateUser(ctx, tx,
			user.Email, user.Name, user.PWHash,
			user.Verified, user.PWAuth, user.ApiAccess,
			user.WebAuthn, user.ID,
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
		log.Debug().Err(err).Msg("error saving verification")
		x.DBError(w, err)
		return
	}

	x.SessMgr.FlashAppend(ctx, "success", "Email verification successfull")
	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}
