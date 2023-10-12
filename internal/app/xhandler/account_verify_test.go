package xhandler

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dropwhile/refid"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/util"
)

func TestHandler_SendVerificationEmail(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		ID:           1,
		RefID:        refid.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     false,
		Created:      ts,
		LastModified: ts,
	}

	uv := &model.UserVerify{
		RefID:   refid.Must(model.NewUserVerifyRefID()),
		UserID:  user.ID,
		Created: ts,
	}

	uvColumns := []string{"ref_id", "user_id", "created"}
	verifyTpl := template.Must(template.New("").Parse(`{{.Subject}}: {{.VerificationUrl}}`))

	t.Run("send verification email", func(t *testing.T) {
		t.Parallel()

		uvRows := pgxmock.NewRows(uvColumns).
			AddRow(uv.RefID, uv.UserID, ts)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		// copy user to avoid context user being modified
		// impacting future tests
		u := *user
		user = &u
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		handler.Tpl["mail_account_email_verify.gohtml"] = verifyTpl

		// refid as anyarg because new refid is created on call to create
		mock.ExpectBegin()
		mock.ExpectQuery("^INSERT INTO user_verify_ (.+)").
			WithArgs(util.NewPgxNamedArgsMatcher(pgx.NamedArgs{
				"refID":  model.UserVerifyRefIDMatcher{},
				"userID": user.ID,
			})).
			WillReturnRows(uvRows)
		mock.ExpectCommit()
		mock.ExpectRollback()

		req, _ := http.NewRequestWithContext(ctx, "POST", "http://example.com/send-verification", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.SendVerificationEmail(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		tm := handler.Mailer.(*TestMailer)
		assert.Equal(t, len(tm.Sent), 1)
		message := tm.Sent[0].BodyPlain
		after, found := strings.CutPrefix(message, "Account Verification url: http://example.com/verify/")
		assert.Assert(t, found)
		refParts := strings.Split(after, "-")
		rID := refid.Must(model.ParseUserVerifyRefID(refParts[0]))
		hmacBytes, err := util.Base32DecodeString(refParts[1])
		assert.NilError(t, err)
		assert.Assert(t, handler.Hmac.Validate([]byte(rID.String()), hmacBytes))

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/settings",
			"handler returned wrong redirect")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}

func TestHandler_VerifyEmail(t *testing.T) {
	t.Parallel()

	ts := tstTs
	user := &model.User{
		ID:           1,
		RefID:        refid.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     false,
		Created:      ts,
		LastModified: ts,
	}

	uv := &model.UserVerify{
		RefID:   refid.Must(model.NewUserVerifyRefID()),
		UserID:  user.ID,
		Created: ts,
	}

	uvColumns := []string{"ref_id", "user_id", "created"}

	t.Run("verify", func(t *testing.T) {
		t.Parallel()

		uvRows := pgxmock.NewRows(uvColumns).
			AddRow(uv.RefID, uv.UserID, uv.Created)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		// copy user to avoid context user being modified
		// impacting future tests
		u := *user
		user = &u
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("uvRefID", uv.RefID.String())

		// generate hmac
		macBytes := handler.Hmac.Generate([]byte(uv.RefID.String()))
		// base32 encode hmac
		macStr := util.Base32EncodeToString(macBytes)

		rctx.URLParams.Add("hmac", macStr)

		// refid as anyarg because new refid is created on call to create
		mock.ExpectQuery("^SELECT (.+) FROM user_verify_ ").
			WithArgs(uv.RefID).
			WillReturnRows(uvRows)
		// start outer tx
		mock.ExpectBegin()
		// begin first inner tx for user update
		mock.ExpectBegin()
		mock.ExpectExec("^UPDATE user_ (.+)").
			WithArgs(util.NewPgxNamedArgsMatcher(pgx.NamedArgs{
				"email":    user.Email,
				"name":     user.Name,
				"pwHash":   pgxmock.AnyArg(),
				"verified": true,
				"userID":   user.ID,
			})).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		// commit+rollback first inner tx
		mock.ExpectCommit()
		mock.ExpectRollback()
		// begin second inner tx for user_verify_ delete
		mock.ExpectBegin()
		mock.ExpectExec("^DELETE FROM user_verify_ (.+)").
			WithArgs(uv.RefID).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))
		// commit+rollback second inner tx
		mock.ExpectCommit()
		mock.ExpectRollback()
		// commit+rollback outer tx
		mock.ExpectCommit()
		mock.ExpectRollback()

		req, _ := http.NewRequestWithContext(ctx, "GET", "http://example.com/verify", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.VerifyEmail(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("location"), "/settings",
			"handler returned wrong redirect")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("verify with bad hmac", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		// copy user to avoid context user being modified
		// impacting future tests
		u := *user
		user = &u
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("uvRefID", uv.RefID.String())

		// generate hmac
		macBytes := handler.Hmac.Generate([]byte(uv.RefID.String()))
		macBytes[0] += 1
		// base32 encode hmac
		macStr := util.Base32EncodeToString(macBytes)

		rctx.URLParams.Add("hmac", macStr)

		req, _ := http.NewRequestWithContext(ctx, "GET", "http://example.com/verify", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.VerifyEmail(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("verify with bad refid", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		// copy user to avoid context user being modified
		// impacting future tests
		u := *user
		user = &u
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		refID := refid.Must(model.NewEventItemRefID())
		rctx.URLParams.Add("uvRefID", refID.String())

		// generate hmac
		macBytes := handler.Hmac.Generate([]byte(refID.String()))
		// base32 encode hmac
		macStr := util.Base32EncodeToString(macBytes)

		rctx.URLParams.Add("hmac", macStr)

		req, _ := http.NewRequestWithContext(ctx, "GET", "http://example.com/verify", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.VerifyEmail(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("verify not in db", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		// copy user to avoid context user being modified
		// impacting future tests
		u := *user
		user = &u
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("uvRefID", uv.RefID.String())

		// generate hmac
		macBytes := handler.Hmac.Generate([]byte(uv.RefID.String()))
		// base32 encode hmac
		macStr := util.Base32EncodeToString(macBytes)

		rctx.URLParams.Add("hmac", macStr)

		// refid as anyarg because new refid is created on call to create
		mock.ExpectQuery("^SELECT (.+) FROM user_verify_ ").
			WithArgs(uv.RefID).
			WillReturnError(pgx.ErrNoRows)

		req, _ := http.NewRequestWithContext(ctx, "GET", "http://example.com/verify", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.VerifyEmail(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("verify is expired", func(t *testing.T) {
		t.Parallel()

		refID := refid.Must(model.NewUserVerifyRefID())
		rfts, _ := time.Parse(time.RFC3339, "2023-01-14T18:29:00Z")
		refID.SetTime(rfts)
		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		// copy user to avoid context user being modified
		// impacting future tests
		u := *user
		user = &u
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("uvRefID", refID.String())

		// generate hmac
		macBytes := handler.Hmac.Generate([]byte(refID.String()))
		// base32 encode hmac
		macStr := util.Base32EncodeToString(macBytes)

		rctx.URLParams.Add("hmac", macStr)

		pwrRows := pgxmock.NewRows(uvColumns).
			AddRow(refID, uv.UserID, uv.Created)

		// refid as anyarg because new refid is created on call to create
		mock.ExpectQuery("^SELECT (.+) FROM user_verify_ ").
			WithArgs(model.UserVerifyRefIDMatcher{}).
			WillReturnRows(pwrRows)

		req, _ := http.NewRequestWithContext(ctx, "GET", "http://example.com/verify", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.VerifyEmail(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusNotFound)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("verify with user_verify delete failure", func(t *testing.T) {
		t.Parallel()

		pwrRows := pgxmock.NewRows(uvColumns).
			AddRow(uv.RefID, uv.UserID, uv.Created)

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		ctx, _ = handler.SessMgr.Load(ctx, "")
		// copy user to avoid context user being modified
		// impacting future tests
		u := *user
		user = &u
		ctx = auth.ContextSet(ctx, "user", user)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("uvRefID", uv.RefID.String())

		// generate hmac
		macBytes := handler.Hmac.Generate([]byte(uv.RefID.String()))
		// base32 encode hmac
		macStr := util.Base32EncodeToString(macBytes)

		rctx.URLParams.Add("hmac", macStr)

		// refid as anyarg because new refid is created on call to create
		mock.ExpectQuery("^SELECT (.+) FROM user_verify_ ").
			WithArgs(uv.RefID).
			WillReturnRows(pwrRows)
		// start outer tx
		mock.ExpectBegin()
		// begin first inner tx for user update
		mock.ExpectBegin()
		mock.ExpectExec("^UPDATE user_ (.+)").
			WithArgs(util.NewPgxNamedArgsMatcher(pgx.NamedArgs{
				"email":    user.Email,
				"name":     user.Name,
				"pwHash":   pgxmock.AnyArg(),
				"verified": true,
				"userID":   user.ID,
			})).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		// commit+rollback first inner tx
		mock.ExpectCommit()
		mock.ExpectRollback()
		// begin second inner tx for user_verify_ delete
		mock.ExpectBegin()
		mock.ExpectExec("^DELETE FROM user_verify_ (.+)").
			WithArgs(uv.RefID).
			WillReturnError(fmt.Errorf("honk honk"))
		// rollback second inner tx
		mock.ExpectRollback()
		mock.ExpectRollback()
		// rollback outer tx
		// rollback before putting conn back in pool
		mock.ExpectRollback()
		mock.ExpectRollback()

		req, _ := http.NewRequestWithContext(ctx, "GET", "http://example.com/verify", nil)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler.VerifyEmail(rr, req)

		response := rr.Result()
		util.MustReadAll(response.Body)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusInternalServerError)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}
