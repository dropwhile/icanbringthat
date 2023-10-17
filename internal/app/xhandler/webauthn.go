package xhandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/service"
)

func getAuthnInstance(r *http.Request, isProd bool, baseURL string) (*webauthn.WebAuthn, error) {
	if isProd {
		siteURL, err := url.Parse(baseURL)
		if err != nil {
			return nil, err
		}
		wconfig := &webauthn.Config{
			// Display Name for site
			RPDisplayName: "ICanBringThat",
			// Generally the FQDN for site
			RPID: siteURL.Hostname(),
			// The origin URLs allowed for WebAuthn requests
			RPOrigins: []string{siteURL.String()},
		}
		return webauthn.New(wconfig)
	} else {
		protocol := "http"
		if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
			protocol = "https"
		}
		host := r.Host
		if r.Header.Get("X-Forwarded-Host") != "" {
			host = r.Header.Get("X-Forwarded-Host")
		}
		baseURL := fmt.Sprintf("%s://%s", protocol, host)

		siteURL, err := url.Parse(baseURL)
		if err != nil {
			return nil, err
		}
		wconfig := &webauthn.Config{
			// Display Name for site
			RPDisplayName: "ICanBringThat",
			// Generally the FQDN for site
			RPID: siteURL.Hostname(),
			// The origin URLs allowed for WebAuthn requests
			RPOrigins: []string{siteURL.String()},
		}
		return webauthn.New(wconfig)
	}
}

func (x *XHandler) WebAuthnBeginRegistration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	authnInstance, err := getAuthnInstance(r, x.IsProd, x.BaseURL)
	if err != nil {
		log.Info().Err(err).Msg("webauthn error")
		x.Error(w, "webauthn error", http.StatusInternalServerError)
		return
	}

	authNUser := service.WebAuthnUserFrom(x.Db, user)

	options, sessionData, err := authnInstance.BeginRegistration(authNUser)
	if err != nil {
		log.Info().Err(err).Msg("webauthn error")
		x.Error(w, "webauthn error", http.StatusInternalServerError)
		return
	}

	val, err := json.Marshal(sessionData)
	if err != nil {
		log.Info().Err(err).Msg("webauthn error")
		x.Error(w, "webauthn error", http.StatusInternalServerError)
		return
	}
	x.SessMgr.Put(ctx, "webauthn-session-register", val)
	x.Json(w, http.StatusOK, options)
}

func (x *XHandler) WebAuthnFinishRegistration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	authnInstance, err := getAuthnInstance(r, x.IsProd, x.BaseURL)
	if err != nil {
		log.Info().Err(err).Msg("webauthn error")
		x.Error(w, "webauthn error", http.StatusInternalServerError)
		return
	}

	authNUser := service.WebAuthnUserFrom(x.Db, user)

	var sessionData webauthn.SessionData
	sessionBytes := x.SessMgr.Pop(ctx, "webauthn-session-register").([]byte)
	if err = json.Unmarshal(sessionBytes, &sessionData); err != nil {
		log.Info().Err(err).Msg("error decoding json webauthn session")
		x.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	credential, err := authnInstance.FinishRegistration(authNUser, sessionData, r)
	if err != nil {
		log.Info().Err(err).Msg("error finishing webauthn registration")
		x.Error(w, "webauthn registration error", http.StatusInternalServerError)
		return
	}
	if err := authNUser.AddCredential(credential); err != nil {
		log.Info().Err(err).Msg("error finishing webauthn registration")
		x.Error(w, "webauthn registration error", http.StatusInternalServerError)
		return
	}
}

func (x *XHandler) WebAuthnBeginLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	_, err := auth.UserFromContext(ctx)
	// already a logged in user
	if err == nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	if email == "" {
		log.Debug().Msg("missing param data")
		x.Error(w, "Missing param data", http.StatusBadRequest)
		return
	}

	// find user...
	user, err := model.GetUserByEmail(ctx, x.Db, email)
	if err != nil || user == nil {
		log.Debug().Err(err).Msg("invalid credentials: no user match")
		x.SessMgr.FlashAppend(ctx, "error", "Invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	authnInstance, err := getAuthnInstance(r, x.IsProd, x.BaseURL)
	if err != nil {
		log.Info().Err(err).Msg("webauthn error")
		x.Error(w, "webauthn error", http.StatusInternalServerError)
		return
	}

	authNUser := service.WebAuthnUserFrom(x.Db, user)

	options, sessionData, err := authnInstance.BeginLogin(authNUser)
	if err != nil {
		log.Info().Err(err).Msg("webauthn error")
		x.Error(w, "webauthn error", http.StatusBadRequest)
		return
	}

	val, err := json.Marshal(sessionData)
	if err != nil {
		log.Info().Err(err).Msg("webauthn error")
		x.Error(w, "webauthn error", http.StatusInternalServerError)
		return
	}
	x.SessMgr.Put(ctx, "webauthn-session-login", val)
	x.Json(w, http.StatusOK, options)
}

func (x *XHandler) WebAuthnFinishLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	_, err := auth.UserFromContext(ctx)
	// already a logged in user
	if err == nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	email := r.PostFormValue("email")
	if email == "" {
		log.Debug().Msg("missing param data")
		x.Error(w, "Missing param data", http.StatusBadRequest)
		return
	}

	user, err := model.GetUserByEmail(ctx, x.Db, email)
	if err != nil || user == nil {
		log.Debug().Err(err).Msg("invalid credentials: no user match")
		x.SessMgr.FlashAppend(ctx, "error", "Invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	authnInstance, err := getAuthnInstance(r, x.IsProd, x.BaseURL)
	if err != nil {
		log.Info().Err(err).Msg("webauthn error")
		x.Error(w, "webauthn error", http.StatusInternalServerError)
		return
	}

	authNUser := service.WebAuthnUserFrom(x.Db, user)

	var sessionData webauthn.SessionData
	sessionBytes := x.SessMgr.Pop(ctx, "webauthn-session-login").([]byte)
	if err = json.Unmarshal(sessionBytes, &sessionData); err != nil {
		log.Info().Err(err).Msg("error decoding json webauthn session")
		x.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	_, err = authnInstance.FinishLogin(authNUser, sessionData, r)
	if err != nil {
		log.Info().Err(err).Msg("error finishing webauthn login")
		x.Error(w, "webauthn login error", http.StatusForbidden)
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
	if r.PostFormValue("next") != "" {
		target = r.FormValue("next")
	}
	x.SessMgr.FlashAppend(ctx, "success", "Login successful")
	http.Redirect(w, r, target, http.StatusSeeOther)
}
