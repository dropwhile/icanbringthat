package xhandler

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dropwhile/refid"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/app/model"
)

func TestHandler_PostmarkCallback(t *testing.T) {
	t.Parallel()

	user := &model.User{
		ID:           1,
		RefID:        refid.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      tstTs,
		LastModified: tstTs,
		Settings: model.UserSettings{
			EnableReminders: true,
		},
	}

	userCols := []string{
		"id", "ref_id", "email", "pwhash", "created", "last_modified",
		"settings",
	}

	t.Run("subscription change to disable should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		rows := pgxmock.NewRows(userCols).
			AddRow(
				user.ID, user.RefID, user.Email, user.PWHash, user.Created,
				user.LastModified, user.Settings,
			)

		expectedSettings := model.UserSettings{
			EnableReminders: false,
		}

		mock.ExpectQuery("^SELECT (.+) FROM user_").
			WithArgs(user.Email).
			WillReturnRows(rows)
		// start outer tx
		mock.ExpectBegin()
		// begin first inner tx for user update
		mock.ExpectBegin()
		mock.ExpectExec("^UPDATE user_ ").
			WithArgs(pgx.NamedArgs{
				"userID":   user.ID,
				"settings": model.NewUserSettingsMatcher(expectedSettings),
			}).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		// commit+rollback first inner tx
		mock.ExpectCommit()
		mock.ExpectRollback()
		// begin second inner tx
		mock.ExpectBegin()
		mock.ExpectQuery("^INSERT INTO notification_ ").
			WithArgs(pgx.NamedArgs{
				"refID":   model.NotificationRefIDMatcher{},
				"userID":  user.ID,
				"message": pgxmock.AnyArg(),
			}).
			WillReturnRows(
				pgxmock.NewRows([]string{"id", "ref_id", "user_id", "message"}).
					AddRow(1, refid.Must(model.NewNotificationRefID()), user.ID, "hodor"),
			)
		// commit+rollback second inner tx
		mock.ExpectCommit()
		mock.ExpectRollback()
		// commit+rollback outer tx
		mock.ExpectCommit()
		mock.ExpectRollback()

		jsonData := []byte(`{  
			"RecordType":"SubscriptionChange",
			"MessageID": "883953f4-6105-42a2-a16a-77a8eac79483",
			"ServerID":123456,
			"MessageStream": "outbound",
			"ChangedAt": "2020-02-01T10:53:34.416071Z",
			"Recipient": "user@example.com",
			"Origin": "Recipient",
			"SuppressSending": true,
			"SuppressionReason": "HardBounce",
			"Tag": "my-tag"
		}`)

		req, _ := http.NewRequestWithContext(
			ctx, "POST", "http://example.com/callback",
			bytes.NewBuffer(jsonData),
		)
		req.Header.Add("content-type", "application/json")
		rr := httptest.NewRecorder()
		handler.PostmarkCallback(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusOK)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("subscription change to disable if already disabled should succeed", func(t *testing.T) {
		t.Parallel()

		user := &model.User{
			ID:           1,
			RefID:        refid.Must(model.NewUserRefID()),
			Email:        "user@example.com",
			Name:         "user",
			PWHash:       []byte("00x00"),
			Verified:     true,
			Created:      tstTs,
			LastModified: tstTs,
			Settings: model.UserSettings{
				EnableReminders: false,
			},
		}

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		rows := pgxmock.NewRows(userCols).
			AddRow(
				user.ID, user.RefID, user.Email, user.PWHash, user.Created,
				user.LastModified, user.Settings,
			)

		mock.ExpectQuery("^SELECT (.+) FROM user_").
			WithArgs(user.Email).
			WillReturnRows(rows)

		jsonData := []byte(`{  
			"RecordType":"SubscriptionChange",
			"MessageID": "883953f4-6105-42a2-a16a-77a8eac79483",
			"ServerID":123456,
			"MessageStream": "outbound",
			"ChangedAt": "2020-02-01T10:53:34.416071Z",
			"Recipient": "user@example.com",
			"Origin": "Recipient",
			"SuppressSending": true,
			"SuppressionReason": "HardBounce",
			"Tag": "my-tag",
			"Metadata": {
				"example": "value",
				"example_2": "value"
			}
		}`)

		req, _ := http.NewRequestWithContext(
			ctx, "POST", "http://example.com/callback",
			bytes.NewBuffer(jsonData),
		)
		req.Header.Add("content-type", "application/json")
		rr := httptest.NewRecorder()
		handler.PostmarkCallback(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusOK)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
	t.Run("subscription change to enable should not update", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		rows := pgxmock.NewRows(userCols).
			AddRow(
				user.ID, user.RefID, user.Email, user.PWHash, user.Created,
				user.LastModified, user.Settings,
			)

		mock.ExpectQuery("^SELECT (.+) FROM user_").
			WithArgs(user.Email).
			WillReturnRows(rows)

		jsonData := []byte(`{  
			"RecordType":"SubscriptionChange",
			"MessageID": "883953f4-6105-42a2-a16a-77a8eac79483",
			"ServerID":123456,
			"MessageStream": "outbound",
			"ChangedAt": "2020-02-01T10:53:34.416071Z",
			"Recipient": "user@example.com",
			"Origin": "Recipient",
			"SuppressSending": false,
			"SuppressionReason": null,
			"Tag": null,
			"Metadata": {}
		}`)

		req, _ := http.NewRequestWithContext(
			ctx, "POST", "http://example.com/callback",
			bytes.NewBuffer(jsonData),
		)
		req.Header.Add("content-type", "application/json")
		rr := httptest.NewRecorder()
		handler.PostmarkCallback(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusOK)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("user not found", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		mock.ExpectQuery("^SELECT (.+) FROM user_").
			WithArgs(user.Email).
			WillReturnError(pgx.ErrNoRows)

		jsonData := []byte(`{  
			"RecordType":"SubscriptionChange",
			"MessageID": "883953f4-6105-42a2-a16a-77a8eac79483",
			"ServerID":123456,
			"MessageStream": "outbound",
			"ChangedAt": "2020-02-01T10:53:34.416071Z",
			"Recipient": "user@example.com",
			"Origin": "Recipient",
			"SuppressSending": true,
			"SuppressionReason": "HardBounce",
			"Tag": "my-tag",
			"Metadata": {
				"example": "value",
				"example_2": "value"
			}
		}`)

		req, _ := http.NewRequestWithContext(
			ctx, "POST", "http://example.com/callback",
			bytes.NewBuffer(jsonData),
		)
		req.Header.Add("content-type", "application/json")
		rr := httptest.NewRecorder()
		handler.PostmarkCallback(rr, req)

		response := rr.Result()
		_, err := io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusOK)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}
