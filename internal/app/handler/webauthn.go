package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"

	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/internal/errs"
	"github.com/dropwhile/icbt/internal/middleware/auth"
)

func getAuthnInstance(r *http.Request, isProd bool, baseURL string) (*webauthn.WebAuthn, error) {
	if !isProd {
		protocol := "http"
		if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
			protocol = "https"
		}
		host := r.Host
		if r.Header.Get("X-Forwarded-Host") != "" {
			host = r.Header.Get("X-Forwarded-Host")
		}
		baseURL = fmt.Sprintf("%s://%s", protocol, host)
	}

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

func (x *Handler) WebAuthnBeginRegistration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	authnInstance, err := getAuthnInstance(r, x.IsProd, x.BaseURL)
	if err != nil {
		slog.ErrorContext(ctx, "webauthn error", "error", err)
		x.InternalServerError(w, "webauthn error")
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
		slog.ErrorContext(ctx, "webauthn error", "error", err)
		x.InternalServerError(w, "webauthn error")
		return
	}

	val, err := json.Marshal(sessionData)
	if err != nil {
		slog.ErrorContext(ctx, "webauthn error", "error", err)
		x.InternalServerError(w, "webauthn error")
		return
	}
	x.SessMgr.Put(ctx, "webauthn-session:register", val)
	x.Json(w, http.StatusOK, options)
}

func (x *Handler) WebAuthnFinishRegistration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	authnInstance, err := getAuthnInstance(r, x.IsProd, x.BaseURL)
	if err != nil {
		slog.ErrorContext(ctx, "webauthn error", "error", err)
		x.InternalServerError(w, "webauthn error")
		return
	}

	authNUser := service.WebAuthnUserFrom(x.Db, user)

	keyName := r.FormValue("key_name")
	if keyName == "" {
		slog.DebugContext(ctx, "missing param data")
		x.BadFormDataError(w, err, "key_name")
		return
	}

	var sessionData webauthn.SessionData
	sessionBytes := x.SessMgr.Pop(ctx, "webauthn-session:register").([]byte)
	if err := json.Unmarshal(sessionBytes, &sessionData); err != nil {
		slog.InfoContext(ctx, "error decoding json webauthn session",
			"error", err)
		x.BadSessionDataError(w)
		return
	}

	credential, err := authnInstance.FinishRegistration(authNUser, sessionData, r)
	if err != nil {
		slog.InfoContext(ctx, "error finishing webauthn registration",
			"error", err)
		x.InternalServerError(w, "webauthn registration error")
		return
	}
	if err := authNUser.AddCredential(keyName, credential); err != nil {
		slog.InfoContext(ctx, "error finishing webauthn registration",
			"error", err)
		x.InternalServerError(w, "webauthn registration error")
		return
	}

	x.Json(w, http.StatusOK, MapSA{"verified": true})
}

func (x *Handler) WebAuthnBeginLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	_, err := auth.UserFromContext(ctx)
	// already a logged in user
	if err == nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	authnInstance, err := getAuthnInstance(r, x.IsProd, x.BaseURL)
	if err != nil {
		slog.ErrorContext(ctx, "webauthn error", "error", err)
		x.Json(w, http.StatusInternalServerError,
			MapSA{"error": "Passkey login failed"},
		)
		return
	}

	opts := []webauthn.LoginOption{
		webauthn.WithUserVerification(protocol.VerificationPreferred),
	}
	options, sessionData, err := authnInstance.BeginDiscoverableLogin(opts...)
	if err != nil {
		slog.ErrorContext(ctx, "webauthn error", "error", err)
		x.Json(w, http.StatusBadRequest,
			MapSA{"error": "Passkey login failed"},
		)
		return
	}

	val, err := json.Marshal(sessionData)
	if err != nil {
		slog.ErrorContext(ctx, "webauthn error", "error", err)
		x.Json(w, http.StatusInternalServerError,
			MapSA{"error": "Passkey login failed"},
		)
		return
	}
	x.SessMgr.Put(ctx, "webauthn-session:login", val)
	x.Json(w, http.StatusOK, options)
}

func (x *Handler) WebAuthnFinishLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	_, err := auth.UserFromContext(ctx)
	// already a logged in user
	if err == nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	authnInstance, err := getAuthnInstance(r, x.IsProd, x.BaseURL)
	if err != nil {
		slog.ErrorContext(ctx, "webauthn error", "error", err)
		x.Json(w, http.StatusInternalServerError,
			MapSA{"error": "Passkey login failed"},
		)
		return
	}

	var sessionData webauthn.SessionData
	sessionBytes := x.SessMgr.Pop(ctx, "webauthn-session:login").([]byte)
	if err = json.Unmarshal(sessionBytes, &sessionData); err != nil {
		slog.ErrorContext(ctx, "error decoding json webauthn session", "error", err)
		x.Json(w, http.StatusBadRequest,
			MapSA{"error": "bad session data"},
		)
		return
	}

	var userID int
	// needs to be inline here (as opposed to a defined function elsewhere)
	// so we can capture the discovered userID value
	handler := func(rawID, userHandle []byte) (webauthn.User, error) {
		// rawID is the credentialID
		// userHandler is user.WebauthnID
		refID, err := service.UserRefIDFromBytes(userHandle)
		if err != nil {
			return nil, fmt.Errorf("bad user id: %w", err)
		}
		user, errx := service.GetUser(ctx, x.Db, refID)
		if errx != nil || user == nil {
			return nil, fmt.Errorf("could not find user: %w", err)
		}
		if !user.WebAuthn {
			return nil, fmt.Errorf("user found but webauthn disabled")
		}
		// capture userID for session login operation after auth success
		userID = user.ID
		authNUser := service.WebAuthnUserFrom(x.Db, user)
		return authNUser, nil
	}
	_, err = authnInstance.FinishDiscoverableLogin(handler, sessionData, r)
	if err != nil {
		slog.InfoContext(ctx, "error finishing webauthn login", "error", err)
		x.Json(w, http.StatusForbidden, MapSA{"error": "Passkey login failed"})
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
	x.SessMgr.Put(r.Context(), "user-id", userID)
	x.SessMgr.FlashAppend(ctx, "success", "Login successful")
	x.Json(w, http.StatusOK, MapSA{"verified": true})
}

func (x *Handler) DeleteWebAuthnKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get user from session
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		x.BadSessionDataError(w)
		return
	}

	credentialRefID, err := service.ParseCredentialRefID(chi.URLParam(r, "cRefID"))
	if err != nil {
		x.BadRefIDError(w, "credential", err)
		return
	}

	credential, errx := service.GetUserCredentialByRefID(ctx, x.Db, credentialRefID)
	if errx != nil {
		switch errx.Code() {
		case errs.NotFound:
			slog.InfoContext(ctx, "credential not found", "error", err)
			x.NotFoundError(w)
		default:
			x.InternalServerError(w, errx.Msg())
		}
		return
	}

	if credential.UserID != user.ID {
		x.AccessDeniedError(w)
		return
	}

	count, errx := service.GetUserCredentialCountByUser(ctx, x.Db, user.ID)
	if errx != nil {
		x.InternalServerError(w, errx.Msg())
		return
	}

	if count == 1 && user.WebAuthn {
		slog.DebugContext(ctx, "refusing to remove last passkey when password auth disabled")
		x.BadRequestError(w, "pre-condition failed")
		return
	}

	errx = service.DeleteUserCredential(ctx, x.Db, credential.ID)
	if errx != nil {
		x.InternalServerError(w, errx.Msg())
		return
	}

	w.WriteHeader(http.StatusOK)
}
