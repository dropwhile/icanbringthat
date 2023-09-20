package handler

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/util"
	"github.com/pashagolub/pgxmock/v3"
	"gotest.tools/v3/assert"
)

func TestHandler_Account_Update(t *testing.T) {
	t.Parallel()

	refId, _ := model.UserRefIdT.New()
	ts := tstTs
	pwhash, _ := util.HashPW([]byte("00x00"))
	user := &model.User{
		Id:           1,
		RefId:        refId,
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       pwhash,
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
				"errors": {"Same Email specified was already present"},
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
			WithArgs("user2@example.com", user.Name, user.PWHash, user.Id).
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
				"operations": {"Email update successfull"},
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
		u := *user
		user = &u
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
				"errors": {"Same Name specified was already present"},
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
		mock.ExpectExec("^UPDATE user_ SET (.+)").
			WithArgs(user.Email, "user2", user.PWHash, user.Id).
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
				"operations": {"Name update successfull"},
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
				"errors": {"New Password and Confirm Password do not match"},
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
				"errors": {"New Password and Confirm Password do not match"},
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
				"errors": {"Old Password invalid"},
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
			WithArgs(user.Email, user.Name, pgxmock.AnyArg(), user.Id).
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
				"operations": {"Password update successfull"},
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

func TestHandler_Account_Delete(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	mock, _, handler := SetupHandler(t, ctx)

	refId, _ := model.UserRefIdT.New()
	ts := tstTs
	user := &model.User{
		Id:           1,
		RefId:        refId,
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
		WithArgs(user.Id).
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
