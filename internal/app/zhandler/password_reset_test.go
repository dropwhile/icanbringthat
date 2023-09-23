package zhandler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/util"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"gotest.tools/v3/assert"
)

func TestHandler_ResetPassword(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		Id:           1,
		RefId:        model.UserRefIdT.MustNew(),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Created:      ts,
		LastModified: ts,
	}

	pwr := &model.UserPWReset{
		RefId:   model.UserPWResetRefIdT.MustNew(),
		UserId:  user.Id,
		Created: ts,
	}

	pwColumns := []string{"ref_id", "user_id", "created"}
	userColumns := []string{"id", "ref_id", "email", "name", "pwhash"}

	t.Run("pwreset", func(t *testing.T) {
		t.Parallel()

		pwrRows := pgxmock.NewRows(pwColumns).
			AddRow(pwr.RefId, pwr.UserId, pwr.Created)
		userRows := pgxmock.NewRows(userColumns).
			AddRow(user.Id, user.RefId, user.Email, user.Name, user.PWHash)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("upwRefId", pwr.RefId.String())

		// generate hmac
		macBytes := handler.Hmac.Generate([]byte(pwr.RefId.String()))
		// base32 encode hmac
		macStr := util.Base32EncodeToString(macBytes)

		rctx.URLParams.Add("hmac", macStr)

		// refid as anyarg because new refid is created on call to create
		mock.ExpectQuery("^SELECT (.+) FROM user_pw_reset_ ").
			WithArgs(pwr.RefId).
			WillReturnRows(pwrRows)
		mock.ExpectQuery("^SELECT (.+) FROM user_ ").
			WithArgs(user.Id).
			WillReturnRows(userRows)
		// start outer tx
		mock.ExpectBegin()
		// begin first inner tx for user update
		mock.ExpectBegin()
		mock.ExpectExec("^UPDATE user_ (.+)").
			WithArgs(user.Email, user.Name, pgxmock.AnyArg(), user.Id).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		// commit+rollback first inner tx
		mock.ExpectCommit()
		mock.ExpectRollback()
		// begin second inner tx for user_pw_reset_ delete
		mock.ExpectBegin()
		mock.ExpectExec("^DELETE FROM user_pw_reset_ (.+)").
			WithArgs(pwr.RefId).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))
		// commit+rollback second inner tx
		mock.ExpectCommit()
		mock.ExpectRollback()
		// commit+rollback outer tx
		mock.ExpectCommit()
		mock.ExpectRollback()

		data := url.Values{
			"password":         {"newpass"},
			"confirm_password": {"newpass"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/password-reset", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.ResetPassword(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		userId := handler.SessMgr.GetInt(ctx, "user-id")
		assert.Assert(t, userId == user.Id)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/dashboard",
			"handler returned wrong redirect")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("pwreset with user logged in", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("upwRefId", pwr.RefId.String())

		// generate hmac
		macBytes := handler.Hmac.Generate([]byte(pwr.RefId.String()))
		// base32 encode hmac
		macStr := util.Base32EncodeToString(macBytes)

		rctx.URLParams.Add("hmac", macStr)

		data := url.Values{
			"password":         {"newpass"},
			"confirm_password": {"newpass"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/password-reset", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.ResetPassword(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusForbidden)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("pwreset with password mismatch", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("upwRefId", pwr.RefId.String())

		// generate hmac
		macBytes := handler.Hmac.Generate([]byte(pwr.RefId.String()))
		// base32 encode hmac
		macStr := util.Base32EncodeToString(macBytes)

		rctx.URLParams.Add("hmac", macStr)

		data := url.Values{
			"password":         {"newpass"},
			"confirm_password": {"newpassx"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/password-reset", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.ResetPassword(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("pwreset with bad hmac", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("upwRefId", pwr.RefId.String())

		// generate hmac
		macBytes := handler.Hmac.Generate([]byte(pwr.RefId.String()))
		macBytes[0] += 1
		// base32 encode hmac
		macStr := util.Base32EncodeToString(macBytes)

		rctx.URLParams.Add("hmac", macStr)

		data := url.Values{
			"password":         {"newpass"},
			"confirm_password": {"newpass"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/password-reset", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.ResetPassword(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("pwreset with bad refid", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		refId := model.EventItemRefIdT.MustNew()
		rctx.URLParams.Add("upwRefId", refId.String())

		// generate hmac
		macBytes := handler.Hmac.Generate([]byte(refId.String()))
		// base32 encode hmac
		macStr := util.Base32EncodeToString(macBytes)

		rctx.URLParams.Add("hmac", macStr)

		data := url.Values{
			"password":         {"newpass"},
			"confirm_password": {"newpass"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/password-reset", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.ResetPassword(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("pwreset upw not in db", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("upwRefId", pwr.RefId.String())

		// generate hmac
		macBytes := handler.Hmac.Generate([]byte(pwr.RefId.String()))
		// base32 encode hmac
		macStr := util.Base32EncodeToString(macBytes)

		rctx.URLParams.Add("hmac", macStr)

		// refid as anyarg because new refid is created on call to create
		mock.ExpectQuery("^SELECT (.+) FROM user_pw_reset_ ").
			WithArgs(pwr.RefId).
			WillReturnError(pgx.ErrNoRows)

		data := url.Values{
			"password":         {"newpass"},
			"confirm_password": {"newpass"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/password-reset", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.ResetPassword(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("pwreset user not found", func(t *testing.T) {
		t.Parallel()

		pwrRows := pgxmock.NewRows(pwColumns).
			AddRow(pwr.RefId, pwr.UserId, pwr.Created)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("upwRefId", pwr.RefId.String())

		// generate hmac
		macBytes := handler.Hmac.Generate([]byte(pwr.RefId.String()))
		// base32 encode hmac
		macStr := util.Base32EncodeToString(macBytes)

		rctx.URLParams.Add("hmac", macStr)

		// refid as anyarg because new refid is created on call to create
		mock.ExpectQuery("^SELECT (.+) FROM user_pw_reset_ ").
			WithArgs(pwr.RefId).
			WillReturnRows(pwrRows)
		mock.ExpectQuery("^SELECT (.+) FROM user_ ").
			WithArgs(user.Id).
			WillReturnError(pgx.ErrNoRows)

		data := url.Values{
			"password":         {"newpass"},
			"confirm_password": {"newpass"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/password-reset", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.ResetPassword(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("pwreset upw is expired", func(t *testing.T) {
		t.Parallel()

		refId := model.UserPWResetRefIdT.MustNew()
		rfts, _ := time.Parse(time.RFC3339, "2023-01-14T18:29:00Z")
		refId.SetTime(rfts)
		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("upwRefId", refId.String())

		// generate hmac
		macBytes := handler.Hmac.Generate([]byte(refId.String()))
		// base32 encode hmac
		macStr := util.Base32EncodeToString(macBytes)

		rctx.URLParams.Add("hmac", macStr)

		pwrRows := pgxmock.NewRows(pwColumns).
			AddRow(refId, pwr.UserId, pwr.Created)

		// refid as anyarg because new refid is created on call to create
		mock.ExpectQuery("^SELECT (.+) FROM user_pw_reset_ ").
			WithArgs(model.UserPWResetRefIdT.AnyMatcher()).
			WillReturnRows(pwrRows)

		data := url.Values{
			"password":         {"newpass"},
			"confirm_password": {"newpass"},
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/password-reset", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.ResetPassword(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusBadRequest)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}
