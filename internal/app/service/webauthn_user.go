package service

import (
	"context"
	"encoding/json"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/app/model"
)

type WebAuthnUser struct {
	*model.User
	db model.PgxHandle
}

func WebAuthnUserFrom(db model.PgxHandle, user *model.User) *WebAuthnUser {
	return &WebAuthnUser{user, db}
}

// WebAuthnID provides the user handle of the user account. A user handle is an opaque byte sequence with a maximum
// size of 64 bytes, and is not meant to be displayed to the user.
//
// To ensure secure operation, authentication and authorization decisions MUST be made on the basis of this id
// member, not the displayName nor name members. See Section 6.1 of [RFC8266].
//
// It's recommended this value is completely random and uses the entire 64 bytes.
//
// Specification: §5.4.3. User Account Parameters for Credential Generation (https://w3c.github.io/webauthn/#dom-publickeycredentialuserentity-id)
func (u *WebAuthnUser) WebAuthnID() []byte {
	return u.RefID.Bytes()
}

// WebAuthnName provides the name attribute of the user account during registration and is a human-palatable name for the user
// account, intended only for display. For example, "Alex Müller" or "田中倫". The Relying Party SHOULD let the user
// choose this, and SHOULD NOT restrict the choice more than necessary.
//
// Specification: §5.4.3. User Account Parameters for Credential Generation (https://w3c.github.io/webauthn/#dictdef-publickeycredentialuserentity)
func (u *WebAuthnUser) WebAuthnName() string {
	return u.Email
}

// WebAuthnDisplayName provides the name attribute of the user account during registration and is a human-palatable
// name for the user account, intended only for display. For example, "Alex Müller" or "田中倫". The Relying Party
// SHOULD let the user choose this, and SHOULD NOT restrict the choice more than necessary.
//
// Specification: §5.4.3. User Account Parameters for Credential Generation (https://www.w3.org/TR/webauthn/#dom-publickeycredentialuserentity-displayname)
func (u *WebAuthnUser) WebAuthnDisplayName() string {
	return u.Email
}

// WebAuthnCredentials provides the list of Credential objects owned by the user.
func (u *WebAuthnUser) WebAuthnCredentials() []webauthn.Credential {
	res := make([]webauthn.Credential, 0)
	credentials, err := model.GetUserCredentialsByUser(context.Background(), u.db, u.ID)
	if err != nil {
		log.Info().Err(err).Msg("error retrieving credentials from db")
		return res
	}
	for _, c := range credentials {
		var cred webauthn.Credential
		err := json.Unmarshal(c.Credential, &cred)
		if err != nil {
			log.Info().Err(err).Msg("error unmarshalling webauthn credential")
			continue
		}
		res = append(res, cred)
		return res

	}
	return res
}

// WebAuthnIcon is a deprecated option.
// Deprecated: this has been removed from the specification recommendation. Suggest a blank string.
func (u *WebAuthnUser) WebAuthnIcon() string {
	return ""
}

func (u *WebAuthnUser) AddCredential(keyName string, credential *webauthn.Credential) error {
	credBytes, err := json.Marshal(credential)
	if err != nil {
		log.Info().Err(err).Msg("error marshalling webauthn credential")
		return err
	}
	_, err = model.NewUserCredential(context.Background(), u.db, u.ID, keyName, credBytes)
	if err != nil {
		log.Info().Err(err).Msg("db error creating credential")
		return err
	}
	return nil
}
