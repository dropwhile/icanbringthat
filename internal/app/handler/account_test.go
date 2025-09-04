// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package handler

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/dropwhile/assert"
	"github.com/samber/mo"
	"go.uber.org/mock/gomock"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/app/service"
	"github.com/dropwhile/icanbringthat/internal/crypto"
	"github.com/dropwhile/icanbringthat/internal/errs"
	"github.com/dropwhile/icanbringthat/internal/middleware/auth"
	"github.com/dropwhile/icanbringthat/internal/util"
)

func TestHandler_Account_Update(t *testing.T) {
	t.Parallel()

	refID := util.Must(model.NewUserRefID())
	ts := tstTs
	pwhash, _ := crypto.HashPW([]byte("00x00"))
	user := &model.User{
		ID:           1,
		RefID:        refID,
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       pwhash,
		Verified:     false,
		PWAuth:       true,
		WebAuthn:     false,
		Created:      ts,
		LastModified: ts,
	}

	t.Run("update email with same as existing", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		// copy user to avoid context user being modified
		// impacting future tests
		u := *user
		user := &u
		ctx = auth.ContextSet(ctx, "user", user)

		data := url.Values{"email": {"user@example.com"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.SettingsUpdate(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		messages := handler.sessMgr.FlashPopAll(ctx)
		assert.Equal(t, messages,
			map[string][]string{
				"error": {"Same Email specified was already present"},
			},
		)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
	})

	t.Run("update email", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		// copy user to avoid context user being modified
		// impacting future tests
		u := *user
		user := &u
		ctx = auth.ContextSet(ctx, "user", user)

		email := "user2@example.com"
		euvs := &service.UserUpdateValues{
			Email:    mo.Some(email),
			Verified: mo.Some(false),
		}

		mock.EXPECT().
			UpdateUser(ctx, user, euvs).
			Return(nil)

		data := url.Values{"email": {email}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.SettingsUpdate(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		messages := handler.sessMgr.FlashPopAll(ctx)
		assert.Equal(t,
			messages,
			map[string][]string{
				"success": {"Email update successfull"},
			},
		)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/settings",
			"handler returned wrong redirect")
	})

	t.Run("update name with same as existing", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		// copy user to avoid context user being modified
		// impacting future tests
		user := &model.User{
			ID:           1,
			RefID:        refID,
			Email:        "user@example.com",
			Name:         "user",
			PWHash:       pwhash,
			Created:      ts,
			LastModified: ts,
		}
		ctx = auth.ContextSet(ctx, "user", user)

		data := url.Values{"name": {"user"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.SettingsUpdate(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		messages := handler.sessMgr.FlashPopAll(ctx)
		assert.Equal(t,
			messages,
			map[string][]string{
				"error": {"Same Name specified was already present"},
			},
		)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/settings",
			"handler returned wrong redirect")
	})

	t.Run("update name", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		// copy user to avoid context user being modified
		// impacting future tests
		u := *user
		user := &u
		ctx = auth.ContextSet(ctx, "user", user)

		euvs := &service.UserUpdateValues{
			Name: mo.Some("user2"),
		}

		mock.EXPECT().
			UpdateUser(ctx, user, euvs).
			Return(nil)

		data := url.Values{"name": {"user2"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.SettingsUpdate(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		messages := handler.sessMgr.FlashPopAll(ctx)
		assert.Equal(t,
			messages,
			map[string][]string{
				"success": {"Name update successfull"},
			},
		)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/settings",
			"handler returned wrong redirect")
	})

	t.Run("update passwd with missing confirm password", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		// copy user to avoid context user being modified
		// impacting future tests
		u := *user
		user := &u
		ctx = auth.ContextSet(ctx, "user", user)

		data := url.Values{"password": {"hodor"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.SettingsUpdate(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		messages := handler.sessMgr.FlashPopAll(ctx)
		assert.Equal(t,
			messages,
			map[string][]string{
				"error": {"New Password and Confirm Password do not match"},
			},
		)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/settings",
			"handler returned wrong redirect")
	})

	t.Run("update passwd with mismatched confirm password", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		// copy user to avoid context user being modified
		// impacting future tests
		u := *user
		user := &u
		ctx = auth.ContextSet(ctx, "user", user)

		data := url.Values{
			"password":         {"hodor"},
			"confirm_password": {"hodor2"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.SettingsUpdate(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		messages := handler.sessMgr.FlashPopAll(ctx)
		assert.Equal(t,
			messages,
			map[string][]string{
				"error": {"New Password and Confirm Password do not match"},
			},
		)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/settings",
			"handler returned wrong redirect")
	})

	t.Run("update passwd with invalid old password", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		// copy user to avoid context user being modified
		// impacting future tests
		u := *user
		user := &u
		ctx = auth.ContextSet(ctx, "user", user)

		mock.EXPECT().
			UpdateUser(
				ctx, user,
				gomock.AssignableToTypeOf(
					&service.UserUpdateValues{},
				)).
			Return(errs.ArgumentError("OldPass", "bad value"))

		data := url.Values{
			"password":         {"hodor"},
			"confirm_password": {"hodor"},
			"old_password":     {"invalid"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.SettingsUpdate(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		messages := handler.sessMgr.FlashPopAll(ctx)
		assert.Equal(t,
			messages,
			map[string][]string{
				"error": {"Old Password invalid"},
			},
		)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/settings",
			"handler returned wrong redirect")
	})

	t.Run("update passwd", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		// copy user to avoid context user being modified
		// impacting future tests
		u := *user
		user := &u
		ctx = auth.ContextSet(ctx, "user", user)

		mock.EXPECT().
			UpdateUser(ctx, user,
				gomock.AssignableToTypeOf(&service.UserUpdateValues{})).
			Return(nil)

		data := url.Values{
			"password":         {"hodor"},
			"confirm_password": {"hodor"},
			"old_password":     {"00x00"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.SettingsUpdate(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		messages := handler.sessMgr.FlashPopAll(ctx)
		assert.Equal(t,
			messages,
			map[string][]string{
				"success": {"Password update successfull"},
			},
		)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/settings",
			"handler returned wrong redirect")
		// we make sure that all expectations were met
	})
}

func TestHandler_Account_Update_Auth(t *testing.T) {
	t.Parallel()

	refID := util.Must(model.NewUserRefID())
	ts := tstTs
	pwhash, _ := crypto.HashPW([]byte("00x00"))
	user := &model.User{
		ID:           1,
		RefID:        refID,
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       pwhash,
		Verified:     true,
		PWAuth:       true,
		WebAuthn:     true,
		Created:      ts,
		LastModified: ts,
	}

	t.Run("disable passauth", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)

		euvs := &service.UserUpdateValues{
			PWAuth: mo.Some(false),
		}

		mock.EXPECT().
			GetUserCredentialCountByUser(ctx, user.ID).
			Return(1, nil)
		mock.EXPECT().
			UpdateUser(ctx, user, euvs).
			Return(nil)

		data := url.Values{"auth_passauth": {"off"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.SettingsAuthUpdate(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		messages := handler.sessMgr.FlashPopAll(ctx)
		assert.Equal(t, len(messages), 0)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/settings",
			"handler returned wrong redirect")
	})

	t.Run("disable passkeys", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)

		euvs := &service.UserUpdateValues{
			WebAuthn: mo.Some(false),
		}

		mock.EXPECT().
			GetUserCredentialCountByUser(ctx, user.ID).
			Return(1, nil)
		mock.EXPECT().
			UpdateUser(ctx, user, euvs).
			Return(nil)

		data := url.Values{"auth_passkeys": {"off"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.SettingsAuthUpdate(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		messages := handler.sessMgr.FlashPopAll(ctx)
		assert.Equal(t, len(messages), 0)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/settings",
			"handler returned wrong redirect")
	})

	t.Run("disable passkeys without pwauth should fail", func(t *testing.T) {
		t.Parallel()

		user := &model.User{
			ID:           1,
			RefID:        refID,
			Email:        "user@example.com",
			Name:         "user",
			PWHash:       pwhash,
			Verified:     true,
			PWAuth:       false,
			WebAuthn:     true,
			Created:      ts,
			LastModified: ts,
		}

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)

		mock.EXPECT().
			GetUserCredentialCountByUser(ctx, user.ID).
			Return(1, nil)

		data := url.Values{"auth_passkeys": {"off"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.SettingsAuthUpdate(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		messages := handler.sessMgr.FlashPopAll(ctx)
		assert.Equal(t, len(messages["error"]), 1)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/settings",
			"handler returned wrong redirect")
	})

	t.Run("enable passkeys without keys added should fail", func(t *testing.T) {
		t.Parallel()

		user := &model.User{
			ID:           1,
			RefID:        refID,
			Email:        "user@example.com",
			Name:         "user",
			PWHash:       pwhash,
			Verified:     true,
			PWAuth:       true,
			WebAuthn:     false,
			Created:      ts,
			LastModified: ts,
		}

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)

		mock.EXPECT().
			GetUserCredentialCountByUser(ctx, user.ID).
			Return(0, nil)

		data := url.Values{"auth_passkeys": {"on"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.SettingsAuthUpdate(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		messages := handler.sessMgr.FlashPopAll(ctx)
		assert.Equal(t, len(messages["error"]), 1)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/settings",
			"handler returned wrong redirect")
	})

	t.Run("enable api access with no existing key should succeed", func(t *testing.T) {
		t.Parallel()

		user := &model.User{
			ID:           1,
			RefID:        refID,
			Email:        "user@example.com",
			Name:         "user",
			PWHash:       pwhash,
			Verified:     true,
			PWAuth:       true,
			ApiAccess:    false,
			WebAuthn:     false,
			Created:      ts,
			LastModified: ts,
		}

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)

		euvs := &service.UserUpdateValues{
			ApiAccess: mo.Some(true),
		}

		mock.EXPECT().
			NewApiKeyIfNotExists(ctx, user.ID).
			Return(&model.ApiKey{
				UserID:  user.ID,
				Token:   "some-token",
				Created: tstTs,
			}, nil)
		mock.EXPECT().
			UpdateUser(ctx, user, euvs).
			Return(nil)

		data := url.Values{"api_access": {"on"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.SettingsAuthApiUpdate(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		messages := handler.sessMgr.FlashPopAll(ctx)
		assert.Equal(t, len(messages["error"]), 0)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/settings",
			"handler returned wrong redirect")
	})

	t.Run("enable api access with existing key should succeed", func(t *testing.T) {
		t.Parallel()

		user := &model.User{
			ID:           1,
			RefID:        refID,
			Email:        "user@example.com",
			Name:         "user",
			PWHash:       pwhash,
			Verified:     true,
			PWAuth:       true,
			ApiAccess:    false,
			WebAuthn:     false,
			Created:      ts,
			LastModified: ts,
		}

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)

		euvs := &service.UserUpdateValues{
			ApiAccess: mo.Some(true),
		}

		mock.EXPECT().
			NewApiKeyIfNotExists(ctx, user.ID).
			Return(&model.ApiKey{
				UserID:  user.ID,
				Token:   "some-token",
				Created: tstTs,
			}, nil)
		mock.EXPECT().
			UpdateUser(ctx, user, euvs).
			Return(nil)

		data := url.Values{"api_access": {"on"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.SettingsAuthApiUpdate(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		messages := handler.sessMgr.FlashPopAll(ctx)
		assert.Equal(t, len(messages["error"]), 0)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/settings",
			"handler returned wrong redirect")
	})

	t.Run("enable api access without a verified account should fail", func(t *testing.T) {
		t.Parallel()

		user := &model.User{
			ID:           1,
			RefID:        refID,
			Email:        "user@example.com",
			Name:         "user",
			PWHash:       pwhash,
			Verified:     false,
			PWAuth:       true,
			ApiAccess:    false,
			WebAuthn:     false,
			Created:      ts,
			LastModified: ts,
		}

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)

		data := url.Values{"api_access": {"on"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.SettingsAuthApiUpdate(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		messages := handler.sessMgr.FlashPopAll(ctx)
		assert.Equal(t, len(messages["error"]), 1)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/settings",
			"handler returned wrong redirect")
	})

	t.Run("rotate api key should succeed", func(t *testing.T) {
		t.Parallel()

		user := &model.User{
			ID:           1,
			RefID:        refID,
			Email:        "user@example.com",
			Name:         "user",
			PWHash:       pwhash,
			Verified:     true,
			PWAuth:       true,
			ApiAccess:    false,
			WebAuthn:     false,
			Created:      ts,
			LastModified: ts,
		}

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)

		mock.EXPECT().
			NewApiKey(ctx, user.ID).
			Return(&model.ApiKey{
				UserID:  user.ID,
				Token:   "new-api-key",
				Created: tstTs,
			}, nil)

		data := url.Values{"rotate_apikey": {"true"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.SettingsAuthApiUpdate(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		messages := handler.sessMgr.FlashPopAll(ctx)
		assert.Equal(t, len(messages["error"]), 0)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/settings",
			"handler returned wrong redirect")
	})
}

func TestHandler_Account_Delete(t *testing.T) {
	t.Parallel()

	refID := util.Must(model.NewUserRefID())
	ts := tstTs
	user := &model.User{
		ID:           1,
		RefID:        refID,
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Created:      ts,
		LastModified: ts,
	}

	ctx := context.TODO()
	mock, _, handler := SetupHandler(t, ctx)
	ctx, _ = handler.sessMgr.Load(ctx, "")
	ctx = auth.ContextSet(ctx, "user", user)

	mock.EXPECT().
		DeleteUser(ctx, user.ID).
		Return(nil)

	req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/account", nil)
	rr := httptest.NewRecorder()
	handler.AccountDelete(rr, req)

	response := rr.Result()
	_, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	// Check the status code is what we expect.
	AssertStatusEqual(t, rr, http.StatusOK)
	// we make sure that all expectations were met
}

func TestHandler_Account_Create(t *testing.T) {
	t.Parallel()

	pwhash, _ := crypto.HashPW([]byte("00x00"))
	user := &model.User{
		ID:           1,
		RefID:        util.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       pwhash,
		Verified:     false,
		PWAuth:       true,
		WebAuthn:     false,
		Created:      tstTs,
		LastModified: tstTs,
	}

	t.Run("create happy path", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")

		msg := `Account is not currently verified. Please verify account in link:/settings.`

		mock.EXPECT().
			NewUser(ctx, user.Email, user.Name, []byte("00x00")).
			Return(user, nil)
		mock.EXPECT().
			NewNotification(ctx, user.ID, msg).
			Return(&model.Notification{
				ID:      1,
				RefID:   model.NotificationRefID{},
				UserID:  user.ID,
				Message: msg,
				Read:    false,
			}, nil)

		data := url.Values{
			"email":            {user.Email},
			"name":             {user.Name},
			"password":         {"00x00"},
			"confirm_password": {"00x00"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.AccountCreate(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		messages := handler.sessMgr.FlashPopAll(ctx)
		assert.Equal(t, messages,
			map[string][]string{
				"success": {"Account created. You are now logged in."},
			},
		)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/dashboard",
			"handler returned wrong redirect")
	})

	t.Run("create missing form data", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")

		data := url.Values{
			"email":            {"user@example.com"},
			"password":         {"00x00"},
			"confirm_password": {"00x00"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.AccountCreate(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
	})

	t.Run("create password mismatch", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")

		data := url.Values{
			"email":            {"user@example.com"},
			"name":             {"user"},
			"password":         {"00x00"},
			"confirm_password": {"00x01"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.AccountCreate(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
	})

	t.Run("create user already exists", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")

		data := url.Values{
			"email":            {"user@example.com"},
			"name":             {"user"},
			"password":         {"00x00"},
			"confirm_password": {"00x00"},
		}

		mock.EXPECT().
			NewUser(ctx, user.Email, user.Name, []byte("00x00")).
			Return(nil, errs.AlreadyExists.Error("user already exists"))

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.AccountCreate(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
	})

	t.Run("create user already logged in", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.sessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)

		data := url.Values{
			"email":            {"user@example.com"},
			"name":             {"user"},
			"password":         {"00x00"},
			"confirm_password": {"00x00"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.AccountCreate(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
	})
}
