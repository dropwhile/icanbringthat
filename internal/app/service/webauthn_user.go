// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package service

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/go-webauthn/webauthn/webauthn"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/errs"
	"github.com/dropwhile/icanbringthat/internal/logger"
)

type Credentialer interface {
	GetUserCredentialsByUser(ctx context.Context, userID int) ([]*model.UserCredential, errs.Error)
	NewUserCredential(ctx context.Context, userID int, keyName string, credential []byte) (*model.UserCredential, errs.Error)
}

type WebAuthnUser struct {
	*model.User
	svc Credentialer
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

// WebAuthnIcon is a deprecated option.
// Deprecated: this has been removed from the specification recommendation. Suggest a blank string.
func (u *WebAuthnUser) WebAuthnIcon() string {
	return ""
}

// WebAuthnCredentials provides the list of Credential objects owned by the user.
func (u *WebAuthnUser) WebAuthnCredentials() []webauthn.Credential {
	ctx := context.Background()
	res := make([]webauthn.Credential, 0)
	credentials, err := u.svc.GetUserCredentialsByUser(ctx, u.ID)
	if err != nil {
		slog.InfoContext(ctx, "error retrieving credentials from db",
			logger.Err(err))
		return res
	}
	for _, c := range credentials {
		var cred webauthn.Credential
		err := json.Unmarshal(c.Credential, &cred)
		if err != nil {
			slog.InfoContext(ctx, "error unmarshalling webauthn credential",
				logger.Err(err))
			continue
		}
		res = append(res, cred)
		return res

	}
	return res
}

func (u *WebAuthnUser) AddCredential(keyName string, credential *webauthn.Credential) error {
	ctx := context.Background()
	credBytes, err := json.Marshal(credential)
	if err != nil {
		slog.InfoContext(ctx, "error marshalling webauthn credential",
			logger.Err(err))
		return err
	}
	_, err = u.svc.NewUserCredential(ctx, u.ID, keyName, credBytes)
	if err != nil {
		slog.InfoContext(ctx, "db error creating credential",
			logger.Err(err))
		return err
	}
	return nil
}

func (s *Service) WebAuthnUserFrom(user *model.User) *WebAuthnUser {
	return &WebAuthnUser{user, s}
}
