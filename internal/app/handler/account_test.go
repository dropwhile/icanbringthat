package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/dropwhile/refid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/crypto"
)

func TestHandler_Account_Update(t *testing.T) {
	t.Parallel()

	refID := refid.Must(model.NewUserRefID())
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
		ctx, _ = handler.SessMgr.Load(ctx, "")
		// copy user to avoid context user being modified
		// impacting future tests
		u := *user
		user = &u
		ctx = auth.ContextSet(ctx, "user", user)

		data := url.Values{"email": {"user@example.com"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateSettings(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		messages := handler.SessMgr.FlashPopAll(ctx)
		assert.DeepEqual(t, messages,
			map[string][]string{
				"error": {"Same Email specified was already present"},
			},
		)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/settings",
			"handler returned wrong redirect")
	})

	t.Run("update email", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		// copy user to avoid context user being modified
		// impacting future tests
		u := *user
		user = &u
		ctx = auth.ContextSet(ctx, "user", user)
		mock.ExpectBegin()
		mock.ExpectExec("^UPDATE user_ SET (.+)").
			WithArgs(pgx.NamedArgs{
				"email":     "user2@example.com",
				"name":      user.Name,
				"pwHash":    user.PWHash,
				"verified":  user.Verified,
				"pwAuth":    user.PWAuth,
				"apiAccess": user.ApiAccess,
				"webAuthn":  user.WebAuthn,
				"userID":    user.ID,
			}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		// hidden rollback after commit due to beginfunc being used
		mock.ExpectRollback()

		data := url.Values{"email": {"user2@example.com"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateSettings(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		messages := handler.SessMgr.FlashPopAll(ctx)
		assert.DeepEqual(t,
			messages,
			map[string][]string{
				"success": {"Email update successfull"},
			},
		)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/settings",
			"handler returned wrong redirect")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update name with same as existing", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
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
		handler.UpdateSettings(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		messages := handler.SessMgr.FlashPopAll(ctx)
		assert.DeepEqual(t,
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
		ctx, _ = handler.SessMgr.Load(ctx, "")
		// copy user to avoid context user being modified
		// impacting future tests
		u := *user
		user = &u
		ctx = auth.ContextSet(ctx, "user", user)
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE user_ SET (.+)").
			WithArgs(pgx.NamedArgs{
				"email":     user.Email,
				"name":      "user2",
				"pwHash":    user.PWHash,
				"verified":  user.Verified,
				"pwAuth":    user.PWAuth,
				"apiAccess": user.ApiAccess,
				"webAuthn":  user.WebAuthn,
				"userID":    user.ID,
			}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		// hidden rollback after commit due to beginfunc being used
		mock.ExpectRollback()

		data := url.Values{"name": {"user2"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateSettings(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		messages := handler.SessMgr.FlashPopAll(ctx)
		assert.DeepEqual(t,
			messages,
			map[string][]string{
				"success": {"Name update successfull"},
			},
		)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/settings",
			"handler returned wrong redirect")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("update passwd with missing confirm password", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		// copy user to avoid context user being modified
		// impacting future tests
		u := *user
		user = &u
		ctx = auth.ContextSet(ctx, "user", user)

		data := url.Values{"password": {"hodor"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateSettings(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		messages := handler.SessMgr.FlashPopAll(ctx)
		assert.DeepEqual(t,
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
		ctx, _ = handler.SessMgr.Load(ctx, "")
		// copy user to avoid context user being modified
		// impacting future tests
		u := *user
		user = &u
		ctx = auth.ContextSet(ctx, "user", user)

		data := url.Values{
			"password":         {"hodor"},
			"confirm_password": {"hodor2"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateSettings(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		messages := handler.SessMgr.FlashPopAll(ctx)
		assert.DeepEqual(t,
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
		_, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		// copy user to avoid context user being modified
		// impacting future tests
		u := *user
		user = &u
		ctx = auth.ContextSet(ctx, "user", user)

		data := url.Values{
			"password":         {"hodor"},
			"confirm_password": {"hodor"},
			"old_password":     {"invalid"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateSettings(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		messages := handler.SessMgr.FlashPopAll(ctx)
		assert.DeepEqual(t,
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
		ctx, _ = handler.SessMgr.Load(ctx, "")
		// copy user to avoid context user being modified
		// impacting future tests
		u := *user
		user = &u
		ctx = auth.ContextSet(ctx, "user", user)

		// note: can't pregenerate an expected pwhash to fulfill the sql query,
		// due to argon2 salting in the user.SetPass call, so just use Any instead.
		mock.ExpectBegin()
		mock.ExpectExec("^UPDATE user_ SET (.+)").
			WithArgs(pgx.NamedArgs{
				"email":     user.Email,
				"name":      user.Name,
				"pwHash":    pgxmock.AnyArg(),
				"verified":  user.Verified,
				"pwAuth":    user.PWAuth,
				"apiAccess": user.ApiAccess,
				"webAuthn":  user.WebAuthn,
				"userID":    user.ID,
			}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		// hidden rollback after commit due to beginfunc being used
		mock.ExpectRollback()

		data := url.Values{
			"password":         {"hodor"},
			"confirm_password": {"hodor"},
			"old_password":     {"00x00"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateSettings(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		messages := handler.SessMgr.FlashPopAll(ctx)
		assert.DeepEqual(t,
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
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestHandler_Account_Update_Auth(t *testing.T) {
	t.Parallel()

	refID := refid.Must(model.NewUserRefID())
	ts := tstTs
	pwhash, _ := crypto.HashPW([]byte("00x00"))

	t.Run("disable passauth", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")

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

		ctx = auth.ContextSet(ctx, "user", user)
		mock.ExpectQuery("SELECT (.+) FROM user_webauthn_").
			WithArgs(user.ID).
			WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(1))
		mock.ExpectBegin()
		mock.ExpectExec("^UPDATE user_ SET (.+)").
			WithArgs(pgx.NamedArgs{
				"email":     user.Email,
				"name":      user.Name,
				"pwHash":    user.PWHash,
				"verified":  user.Verified,
				"pwAuth":    false,
				"apiAccess": user.ApiAccess,
				"webAuthn":  user.WebAuthn,
				"userID":    user.ID,
			}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		// hidden rollback after commit due to beginfunc being used
		mock.ExpectRollback()

		data := url.Values{"auth_passauth": {"off"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateAuthSettings(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		messages := handler.SessMgr.FlashPopAll(ctx)
		assert.Equal(t, len(messages), 0)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/settings",
			"handler returned wrong redirect")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("disable passkeys", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")

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

		ctx = auth.ContextSet(ctx, "user", user)
		mock.ExpectQuery("SELECT (.+) FROM user_webauthn_").
			WithArgs(user.ID).
			WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(1))
		mock.ExpectBegin()
		mock.ExpectExec("^UPDATE user_ SET (.+)").
			WithArgs(pgx.NamedArgs{
				"email":     user.Email,
				"name":      user.Name,
				"pwHash":    user.PWHash,
				"verified":  user.Verified,
				"pwAuth":    user.PWAuth,
				"apiAccess": user.ApiAccess,
				"webAuthn":  false,
				"userID":    user.ID,
			}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		// hidden rollback after commit due to beginfunc being used
		mock.ExpectRollback()

		data := url.Values{"auth_passkeys": {"off"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateAuthSettings(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		messages := handler.SessMgr.FlashPopAll(ctx)
		assert.Equal(t, len(messages), 0)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/settings",
			"handler returned wrong redirect")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("disable passkeys without pwauth should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")

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

		ctx = auth.ContextSet(ctx, "user", user)
		mock.ExpectQuery("SELECT (.+) FROM user_webauthn_").
			WithArgs(user.ID).
			WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(1))

		data := url.Values{"auth_passkeys": {"off"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateAuthSettings(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		messages := handler.SessMgr.FlashPopAll(ctx)
		assert.Equal(t, len(messages["error"]), 1)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/settings",
			"handler returned wrong redirect")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("enable passkeys without keys added should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")

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

		ctx = auth.ContextSet(ctx, "user", user)
		mock.ExpectQuery("SELECT (.+) FROM user_webauthn_").
			WithArgs(user.ID).
			WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(0))

		data := url.Values{"auth_passkeys": {"on"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateAuthSettings(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		messages := handler.SessMgr.FlashPopAll(ctx)
		assert.Equal(t, len(messages["error"]), 1)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/settings",
			"handler returned wrong redirect")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("enable api access should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")

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

		ctx = auth.ContextSet(ctx, "user", user)
		mock.ExpectBegin()
		mock.ExpectExec("^UPDATE user_ SET (.+)").
			WithArgs(pgx.NamedArgs{
				"email":     user.Email,
				"name":      user.Name,
				"pwHash":    user.PWHash,
				"verified":  user.Verified,
				"pwAuth":    user.PWAuth,
				"apiAccess": true,
				"webAuthn":  user.WebAuthn,
				"userID":    user.ID,
			}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		// hidden rollback after commit due to beginfunc being used
		mock.ExpectRollback()

		data := url.Values{"api_access": {"on"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateApiAuthSettings(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		messages := handler.SessMgr.FlashPopAll(ctx)
		assert.Equal(t, len(messages["error"]), 0)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/settings",
			"handler returned wrong redirect")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("enable api access without a verified account should fail", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")

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

		ctx = auth.ContextSet(ctx, "user", user)
		data := url.Values{"api_access": {"on"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateApiAuthSettings(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		messages := handler.SessMgr.FlashPopAll(ctx)
		assert.Equal(t, len(messages["error"]), 1)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/settings",
			"handler returned wrong redirect")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("rotate api key should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")

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

		mock.ExpectQuery("SELECT (.+) FROM api_key_").
			WithArgs(user.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{"user_id", "token", "created"},
			).AddRow(
				1, "00000000000000000000000000:11111111111111111111111111", ts,
			))
		mock.ExpectBegin()
		mock.ExpectExec("^UPDATE api_key_ SET (.+)").
			WithArgs(pgx.NamedArgs{
				"userID": user.ID,
				"token":  pgxmock.AnyArg(),
			}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mock.ExpectCommit()
		// hidden rollback after commit due to beginfunc being used
		mock.ExpectRollback()

		ctx = auth.ContextSet(ctx, "user", user)
		data := url.Values{"rotate_apikey": {"true"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.UpdateApiAuthSettings(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		messages := handler.SessMgr.FlashPopAll(ctx)
		assert.Equal(t, len(messages["error"]), 0)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/settings",
			"handler returned wrong redirect")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestHandler_Account_Delete(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	mock, _, handler := SetupHandler(t, ctx)

	refID := refid.Must(model.NewUserRefID())
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

	ctx, _ = handler.SessMgr.Load(ctx, "")
	ctx = auth.ContextSet(ctx, "user", user)

	// note: can't pregenerate an expected pwhash to fulfill the sql query,
	// due to argon2 salting in the user.SetPass call, so just use Any instead.
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM user_ (.+)").
		WithArgs(user.ID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	// hidden rollback after commit due to beginfunc being used
	mock.ExpectRollback()

	req, _ := http.NewRequestWithContext(ctx, "DELETE", "http://example.com/account", nil)
	rr := httptest.NewRecorder()
	handler.DeleteAccount(rr, req)

	response := rr.Result()
	_, err := io.ReadAll(response.Body)
	assert.NilError(t, err)

	// Check the status code is what we expect.
	AssertStatusEqual(t, rr, http.StatusOK)
	// we make sure that all expectations were met
	assert.Assert(t, mock.ExpectationsWereMet(),
		"there were unfulfilled expectations")
}

func TestHandler_Account_Create(t *testing.T) {
	t.Parallel()

	t.Run("create happy path", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")

		data := url.Values{
			"email":            {"user@example.com"},
			"name":             {"user"},
			"password":         {"00x00"},
			"confirm_password": {"00x00"},
		}

		pwhash, _ := crypto.HashPW([]byte("00x00"))
		rows := pgxmock.NewRows(
			[]string{
				"id", "ref_id", "email", "pwhash", "created", "last_modified",
			}).AddRow(
			1, refid.Must(model.NewUserRefID()), "user@example.com", pwhash, tstTs, tstTs,
		)

		mock.ExpectBegin()
		mock.ExpectQuery("^INSERT INTO user_").
			WithArgs(pgx.NamedArgs{
				"refID":    model.UserRefIDMatcher,
				"email":    "user@example.com",
				"name":     "user",
				"pwHash":   pgxmock.AnyArg(),
				"pwAuth":   true,
				"settings": pgxmock.AnyArg(),
			}).
			WillReturnRows(rows)
		mock.ExpectCommit()
		mock.ExpectRollback()

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.CreateAccount(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		messages := handler.SessMgr.FlashPopAll(ctx)
		assert.DeepEqual(t, messages,
			map[string][]string{
				"success": {"Account created. You are now logged in."},
			},
		)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/dashboard",
			"handler returned wrong redirect")
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create missing form data", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")

		data := url.Values{
			"email":            {"user@example.com"},
			"password":         {"00x00"},
			"confirm_password": {"00x00"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.CreateAccount(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create password mismatch", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")

		data := url.Values{
			"email":            {"user@example.com"},
			"name":             {"user"},
			"password":         {"00x00"},
			"confirm_password": {"00x01"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.CreateAccount(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create user already exists", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")

		data := url.Values{
			"email":            {"user@example.com"},
			"name":             {"user"},
			"password":         {"00x00"},
			"confirm_password": {"00x00"},
		}

		mock.ExpectBegin()
		mock.ExpectQuery("^INSERT INTO user_").
			WithArgs(pgx.NamedArgs{
				"refID":    model.UserRefIDMatcher,
				"email":    "user@example.com",
				"name":     "user",
				"pwHash":   pgxmock.AnyArg(),
				"pwAuth":   true,
				"settings": pgxmock.AnyArg(),
			}).
			WillReturnError(fmt.Errorf("duplicate row"))
		mock.ExpectRollback()
		mock.ExpectRollback()

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.CreateAccount(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("create user already logged in", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")

		pwhash, _ := crypto.HashPW([]byte("00x00"))
		user := &model.User{
			ID:           1,
			RefID:        refid.Must(model.NewUserRefID()),
			Email:        "user@example.com",
			Name:         "user",
			PWHash:       pwhash,
			Created:      tstTs,
			LastModified: tstTs,
		}
		ctx = auth.ContextSet(ctx, "user", user)
		ctx = auth.ContextSet(ctx, "auth", true)

		data := url.Values{
			"email":            {"user@example.com"},
			"name":             {"user"},
			"password":         {"00x00"},
			"confirm_password": {"00x00"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/account", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.CreateAccount(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}
