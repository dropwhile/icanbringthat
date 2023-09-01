// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package service

import (
	"context"
	"encoding/json"
	"testing"

	"gotest.tools/v3/assert"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/errs"
	"github.com/dropwhile/icanbringthat/internal/util"
)

type MockCredentialer struct {
	getFunc func(context.Context, int) ([]*model.UserCredential, errs.Error)
	newFunc func(context.Context, int, string, []byte) (*model.UserCredential, errs.Error)
}

func (m *MockCredentialer) GetUserCredentialsByUser(
	ctx context.Context, userID int,
) ([]*model.UserCredential, errs.Error) {
	return m.getFunc(ctx, userID)
}

func (m *MockCredentialer) NewUserCredential(
	ctx context.Context, userID int, keyName string, credential []byte,
) (*model.UserCredential, errs.Error) {
	return m.newFunc(ctx, userID, keyName, credential)
}

func MockCredentialerFuncs(
	getFunc func(context.Context, int) ([]*model.UserCredential, errs.Error),
	newFunc func(context.Context, int, string, []byte) (*model.UserCredential, errs.Error),
) *MockCredentialer {
	return &MockCredentialer{
		getFunc: getFunc,
		newFunc: newFunc,
	}
}

func TestWebAuthnUser_WebAuthnID(t *testing.T) {
	t.Parallel()

	user := &model.User{
		ID:           1,
		RefID:        util.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      tstTs,
		LastModified: tstTs,
	}

	o := WebAuthnUser{User: user, svc: nil}
	assert.DeepEqual(t, o.WebAuthnID(), user.RefID.Bytes())
}

func TestWebAuthnUser_WebAuthnName(t *testing.T) {
	t.Parallel()

	user := &model.User{
		ID:           1,
		RefID:        util.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      tstTs,
		LastModified: tstTs,
	}

	o := WebAuthnUser{User: user, svc: nil}
	assert.DeepEqual(t, o.WebAuthnName(), user.Email)
}

func TestWebAuthnUser_WebAuthnDisplayName(t *testing.T) {
	t.Parallel()

	user := &model.User{
		ID:           1,
		RefID:        util.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      tstTs,
		LastModified: tstTs,
	}

	o := WebAuthnUser{User: user, svc: nil}
	assert.Equal(t, o.WebAuthnDisplayName(), user.Email)
}

func TestWebAuthnUser_WebAuthnIcon(t *testing.T) {
	t.Parallel()

	user := &model.User{
		ID:           1,
		RefID:        util.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      tstTs,
		LastModified: tstTs,
	}

	o := WebAuthnUser{User: user, svc: nil}
	assert.Equal(t, o.WebAuthnIcon(), "")
}

func TestWebAuthnUser_WebAuthnCredentials(t *testing.T) {
	t.Parallel()

	user := &model.User{
		ID:           1,
		RefID:        util.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      tstTs,
		LastModified: tstTs,
	}

	t.Run("get credentials should succeed", func(t *testing.T) {
		cred := webauthn.Credential{
			ID:              []byte("some-id"),
			PublicKey:       []byte("some-pubkey"),
			AttestationType: "none",
			Transport:       []protocol.AuthenticatorTransport{"internal", "hybrid"},
			Flags: webauthn.CredentialFlags{
				UserPresent:    true,
				UserVerified:   true,
				BackupEligible: true,
				BackupState:    true,
			},
			Authenticator: webauthn.Authenticator{
				AAGUID:       []byte("some-aaguid"),
				SignCount:    0,
				CloneWarning: false,
				Attachment:   "platform",
			},
		}

		jdata, err := json.Marshal(cred)
		assert.NilError(t, err)

		credentials := []*model.UserCredential{{
			ID:         3,
			RefID:      util.Must(model.NewCredentialRefID()),
			UserID:     user.ID,
			Credential: jdata,
			KeyName:    "key-a",
			Created:    tstTs,
		}}

		o := WebAuthnUser{User: user, svc: &MockCredentialer{
			getFunc: func(ctx context.Context, i int) ([]*model.UserCredential, errs.Error) {
				return credentials, nil
			},
		}}
		assert.DeepEqual(t, o.WebAuthnCredentials()[0], cred)
	})

	t.Run("get credentials with empty results should succeed", func(t *testing.T) {
		o := WebAuthnUser{User: user, svc: &MockCredentialer{
			getFunc: func(ctx context.Context, i int) ([]*model.UserCredential, errs.Error) {
				return []*model.UserCredential{}, nil
			},
		}}
		assert.Equal(t, len(o.WebAuthnCredentials()), 0)
	})
}

func TestWebAuthnUser_AddCredential(t *testing.T) {
	t.Parallel()

	user := &model.User{
		ID:           1,
		RefID:        util.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      tstTs,
		LastModified: tstTs,
	}

	t.Run("add credential should succeed", func(t *testing.T) {
		cred := webauthn.Credential{
			ID:              []byte("some-id"),
			PublicKey:       []byte("some-pubkey"),
			AttestationType: "none",
			Transport:       []protocol.AuthenticatorTransport{"internal", "hybrid"},
			Flags: webauthn.CredentialFlags{
				UserPresent:    true,
				UserVerified:   true,
				BackupEligible: true,
				BackupState:    true,
			},
			Authenticator: webauthn.Authenticator{
				AAGUID:       []byte("some-aaguid"),
				SignCount:    0,
				CloneWarning: false,
				Attachment:   "platform",
			},
		}

		jdata, err := json.Marshal(cred)
		assert.NilError(t, err)

		credential := &model.UserCredential{
			ID:         3,
			RefID:      util.Must(model.NewCredentialRefID()),
			UserID:     user.ID,
			Credential: jdata,
			KeyName:    "key-a",
			Created:    tstTs,
		}

		o := WebAuthnUser{User: user, svc: &MockCredentialer{
			newFunc: func(ctx context.Context, i int, s string, b []byte) (*model.UserCredential, errs.Error) {
				assert.Equal(t, i, user.ID)
				assert.Equal(t, s, credential.KeyName)
				assert.DeepEqual(t, b, credential.Credential)
				return credential, nil
			},
		}}
		err = o.AddCredential(credential.KeyName, &cred)
		assert.NilError(t, err)
	})
}
