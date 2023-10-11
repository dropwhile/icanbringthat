package xhandler

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/dropwhile/refid"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/modelx"
	"github.com/dropwhile/icbt/internal/util"
)

func TestHandler_ResetPassword(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &modelx.User{
		ID:           1,
		RefID:        refid.Must(modelx.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PwHash:       []byte("00x00"),
		Verified:     false,
		Created:      ts,
		LastModified: ts,
	}

	pwr := &modelx.UserPwReset{
		RefID:   refid.Must(modelx.NewUserPwResetRefID()),
		UserID:  user.ID,
		Created: ts,
	}

	pwColumns := []string{"ref_id", "user_id", "created"}
	userColumns := []string{"id", "ref_id", "email", "name", "pwhash"}

	t.Run("pwreset", func(t *testing.T) {
		t.Parallel()

		pwrRows := pgxmock.NewRows(pwColumns).
			AddRow(pwr.RefID, pwr.UserID, pwr.Created)
		userRows := pgxmock.NewRows(userColumns).
			AddRow(user.ID, user.RefID, user.Email, user.Name, user.PwHash)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("upwRefID", pwr.RefID.String())

		// generate hmac
		macBytes := handler.Hmac.Generate([]byte(pwr.RefID.String()))
		// base32 encode hmac
		macStr := util.Base32EncodeToString(macBytes)

		rctx.URLParams.Add("hmac", macStr)

		// refid as anyarg because new refid is created on call to create
		mock.ExpectQuery("^SELECT (.+) FROM user_pw_reset_ ").
			WithArgs(pwr.RefID).
			WillReturnRows(pwrRows)
		mock.ExpectQuery("^SELECT (.+) FROM user_ ").
			WithArgs(user.ID).
			WillReturnRows(userRows)
		// start outer tx
		mock.ExpectBegin()
		// begin first inner tx for user update
		mock.ExpectBegin()
		mock.ExpectExec("^UPDATE user_ (.+)").
			WithArgs(user.Email, user.Name, pgxmock.AnyArg(), user.Verified, user.ID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		// commit+rollback first inner tx
		mock.ExpectCommit()
		mock.ExpectRollback()
		// begin second inner tx for user_pw_reset_ delete
		mock.ExpectBegin()
		mock.ExpectExec("^DELETE FROM user_pw_reset_ (.+)").
			WithArgs(pwr.RefID).
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

		userId := handler.SessMgr.GetInt32(ctx, "user-id")
		assert.Assert(t, userId == user.ID)

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
		rctx.URLParams.Add("upwRefID", pwr.RefID.String())

		// generate hmac
		macBytes := handler.Hmac.Generate([]byte(pwr.RefID.String()))
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
		rctx.URLParams.Add("upwRefID", pwr.RefID.String())

		// generate hmac
		macBytes := handler.Hmac.Generate([]byte(pwr.RefID.String()))
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
		rctx.URLParams.Add("upwRefID", pwr.RefID.String())

		// generate hmac
		macBytes := handler.Hmac.Generate([]byte(pwr.RefID.String()))
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
		refId := refid.Must(modelx.NewEventItemRefID())
		rctx.URLParams.Add("upwRefID", refId.String())

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
		rctx.URLParams.Add("upwRefID", pwr.RefID.String())

		// generate hmac
		macBytes := handler.Hmac.Generate([]byte(pwr.RefID.String()))
		// base32 encode hmac
		macStr := util.Base32EncodeToString(macBytes)

		rctx.URLParams.Add("hmac", macStr)

		// refid as anyarg because new refid is created on call to create
		mock.ExpectQuery("^SELECT (.+) FROM user_pw_reset_ ").
			WithArgs(pwr.RefID).
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
			AddRow(pwr.RefID, pwr.UserID, pwr.Created)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("upwRefID", pwr.RefID.String())

		// generate hmac
		macBytes := handler.Hmac.Generate([]byte(pwr.RefID.String()))
		// base32 encode hmac
		macStr := util.Base32EncodeToString(macBytes)

		rctx.URLParams.Add("hmac", macStr)

		// refid as anyarg because new refid is created on call to create
		mock.ExpectQuery("^SELECT (.+) FROM user_pw_reset_ ").
			WithArgs(pwr.RefID).
			WillReturnRows(pwrRows)
		mock.ExpectQuery("^SELECT (.+) FROM user_ ").
			WithArgs(user.ID).
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

		refId := refid.Must(modelx.NewUserPwResetRefID())
		rfts, _ := time.Parse(time.RFC3339, "2023-01-14T18:29:00Z")
		refId.SetTime(rfts)
		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("upwRefID", refId.String())

		// generate hmac
		macBytes := handler.Hmac.Generate([]byte(refId.String()))
		// base32 encode hmac
		macStr := util.Base32EncodeToString(macBytes)

		rctx.URLParams.Add("hmac", macStr)

		pwrRows := pgxmock.NewRows(pwColumns).
			AddRow(refId, pwr.UserID, pwr.Created)

		// refid as anyarg because new refid is created on call to create
		mock.ExpectQuery("^SELECT (.+) FROM user_pw_reset_ ").
			WithArgs(modelx.UserPwResetRefIDMatcher{}).
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

	t.Run("pwreset with user upw delete failure", func(t *testing.T) {
		t.Parallel()

		pwrRows := pgxmock.NewRows(pwColumns).
			AddRow(pwr.RefID, pwr.UserID, pwr.Created)
		userRows := pgxmock.NewRows(userColumns).
			AddRow(user.ID, user.RefID, user.Email, user.Name, user.PwHash)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("upwRefID", pwr.RefID.String())

		// generate hmac
		macBytes := handler.Hmac.Generate([]byte(pwr.RefID.String()))
		// base32 encode hmac
		macStr := util.Base32EncodeToString(macBytes)

		rctx.URLParams.Add("hmac", macStr)

		// refid as anyarg because new refid is created on call to create
		mock.ExpectQuery("^SELECT (.+) FROM user_pw_reset_ ").
			WithArgs(pwr.RefID).
			WillReturnRows(pwrRows)
		mock.ExpectQuery("^SELECT (.+) FROM user_ ").
			WithArgs(user.ID).
			WillReturnRows(userRows)
		// start outer tx
		mock.ExpectBegin()
		// begin first inner tx for user update
		mock.ExpectBegin()
		mock.ExpectExec("^UPDATE user_ (.+)").
			WithArgs(user.Email, user.Name, pgxmock.AnyArg(), user.Verified, user.ID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		// commit+rollback first inner tx
		mock.ExpectCommit()
		mock.ExpectRollback()
		// begin second inner tx for user_pw_reset_ delete
		mock.ExpectBegin()
		mock.ExpectExec("^DELETE FROM user_pw_reset_ (.+)").
			WithArgs(pwr.RefID).
			WillReturnError(fmt.Errorf("honk honk"))
		// rollback second inner tx
		mock.ExpectRollback()
		mock.ExpectRollback()
		// rollback outer tx
		// rollback before putting conn back in pool
		mock.ExpectRollback()
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
		assert.Assert(t, userId == 0)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusInternalServerError)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestHandler_SendResetPasswordEmail(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &modelx.User{
		ID:           1,
		RefID:        refid.Must(modelx.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PwHash:       []byte("00x00"),
		Verified:     false,
		Created:      ts,
		LastModified: ts,
	}

	pwr := &modelx.UserPwReset{
		RefID:   refid.Must(modelx.NewUserPwResetRefID()),
		UserID:  user.ID,
		Created: ts,
	}

	userColumns := []string{"id", "ref_id", "email", "name", "pwhash"}
	pwColumns := []string{"ref_id", "user_id", "created"}
	passResetTpl := template.Must(template.New("").Parse(`{{.Subject}}: {{.PasswordResetUrl}}`))

	t.Run("send pw reset email", func(t *testing.T) {
		t.Parallel()

		userRows := pgxmock.NewRows(userColumns).
			AddRow(user.ID, user.RefID, user.Email, user.Name, user.PwHash)
		upwRows := pgxmock.NewRows(pwColumns).
			AddRow(pwr.RefID, pwr.UserID, ts)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		handler.Tpl["mail_password_reset.gohtml"] = passResetTpl

		// refid as anyarg because new refid is created on call to create
		mock.ExpectQuery("^SELECT (.+) FROM user_ ").
			WithArgs(user.Email).
			WillReturnRows(userRows)
		mock.ExpectBegin()
		mock.ExpectQuery("^INSERT INTO user_pw_reset_ (.+)").
			WithArgs(modelx.UserPwResetRefIDMatcher{}, user.ID).
			WillReturnRows(upwRows)
		mock.ExpectCommit()
		mock.ExpectRollback()

		data := url.Values{"email": {"user@example.com"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/send-password-reset", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.SendResetPasswordEmail(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		tm := handler.Mailer.(*TestMailer)
		assert.Equal(t, len(tm.Sent), 1)
		message := tm.Sent[0].BodyPlain
		after, found := strings.CutPrefix(message, "Password reset url: http://example.com/forgot-password/")
		assert.Assert(t, found)
		refParts := strings.Split(after, "-")
		rId := refid.Must(modelx.ParseUserPwResetRefID(refParts[0]))
		hmacBytes, err := util.Base32DecodeString(refParts[1])
		assert.NilError(t, err)
		assert.Assert(t, handler.Hmac.Validate([]byte(rId.String()), hmacBytes))

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/login",
			"handler returned wrong redirect")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("send pw reset email no user", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		handler.Tpl["mail_password_reset.gohtml"] = passResetTpl

		// refid as anyarg because new refid is created on call to create
		mock.ExpectQuery("^SELECT (.+) FROM user_ ").
			WithArgs(user.Email).
			WillReturnError(pgx.ErrNoRows)

		data := url.Values{"email": {"user@example.com"}}

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/send-password-reset", FormData(data))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.SendResetPasswordEmail(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		tm := handler.Mailer.(*TestMailer)
		assert.Equal(t, len(tm.Sent), 0)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/login",
			"handler returned wrong redirect")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}
