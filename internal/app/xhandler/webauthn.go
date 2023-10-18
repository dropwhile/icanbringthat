package xhandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/jackc/pgx/v5"
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

	// exclude any existing credentials so the user can't accidentally
	// reregister the same device twice
	credentials := authNUser.WebAuthnCredentials()
	excludeList := make([]protocol.CredentialDescriptor, 0, len(credentials))
	for _, cred := range credentials {
		excludeList = append(excludeList, cred.Descriptor())
	}

	authSelect := protocol.AuthenticatorSelection{
		// We want computer or phone as authenticator device. Another
		// option would be to require an USB Security Key for example.
		AuthenticatorAttachment: protocol.AuthenticatorAttachment("platform"),
		// This feature is also referred to as "Discoverable Credential" and it
		// enables us to authenticate without a username, just by providing the
		// passkey. Pretty convenient, but we don't support that yet.
		RequireResidentKey: protocol.ResidentKeyNotRequired(),
		// This triggers the Authenticator to ask for Face ID, Touch ID or a PIN
		// whenever the new passkey is to be used. Your device decides which
		// mechanism is active. We want multi-factor authentication!
		UserVerification: protocol.VerificationPreferred,
	}
	// This determines if we want to receive so called attestation
	// information. Think of it as a certificate about the capabilities of
	// the Authenticator. You would be using that in a scenario with
	// advanced security needs, e.g., in an online banking scenario. This is
	// not the case here, so we switch it off.
	conveyancePref := protocol.PreferNoAttestation

	options, sessionData, err := authnInstance.BeginRegistration(
		authNUser,
		webauthn.WithAuthenticatorSelection(authSelect),
		webauthn.WithConveyancePreference(conveyancePref),
		webauthn.WithExclusions(excludeList),
	)
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
	x.SessMgr.Put(ctx, "webauthn-session:register", val)
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

	keyName := r.FormValue("key_name")
	if keyName == "" {
		log.Debug().Msg("missing param data")
		x.Error(w, "Missing param data", http.StatusBadRequest)
		return
	}

	var sessionData webauthn.SessionData
	sessionBytes := x.SessMgr.Pop(ctx, "webauthn-session:register").([]byte)
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
	if err := authNUser.AddCredential(keyName, credential); err != nil {
		log.Info().Err(err).Msg("error finishing webauthn registration")
		x.Error(w, "webauthn registration error", http.StatusInternalServerError)
		return
	}
	resp := map[string]any{"verified": true}
	x.Json(w, http.StatusOK, resp)
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
	x.SessMgr.Put(ctx, "webauthn-session:login", val)
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
	sessionBytes := x.SessMgr.Pop(ctx, "webauthn-session:login").([]byte)
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

func (x *XHandler) DeleteWebAuthnKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		log.Debug().Err(err).Msg("bad session data")
		x.Error(w, "bad session data", http.StatusBadRequest)
		return
	}

	credentialRefID, err := model.ParseCredentialRefID(chi.URLParam(r, "cRefID"))
	if err != nil {
		x.Error(w, "bad credential-ref-id", http.StatusNotFound)
		return
	}

	credential, err := model.GetUserCredentialByRefID(ctx, x.Db, credentialRefID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Info().Err(err).Msg("credential not found")
		x.Error(w, "credential not found", http.StatusNotFound)
		return
	case err != nil:
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if credential.UserID != user.ID {
		x.Error(w, "access denied", http.StatusForbidden)
		return
	}

	count, err := model.GetUserCredentialCountByUser(ctx, x.Db, user.ID)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	if count == 1 && user.WebAuthn {
		log.Debug().Msg("refusing to remove last passkey when password auth disabled")
		x.Error(w, "pre-condition failed", http.StatusBadRequest)
		return
	}

	err = model.DeleteUserCredential(ctx, x.Db, credential.ID)
	if err != nil {
		log.Info().Err(err).Msg("db error")
		x.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
