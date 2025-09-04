// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package handler

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dropwhile/assert"
	"github.com/go-chi/chi/v5"

	"github.com/dropwhile/icanbringthat/internal/errs"
)

func TestHandler_PostmarkCallback(t *testing.T) {
	t.Parallel()

	t.Run("subscription change to disable should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		email := "user@example.com"
		reason := "HardBounce"

		mock.EXPECT().
			DisableRemindersWithNotification(ctx, email, reason).
			Return(nil)

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
		assert.Nil(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusOK)
		// we make sure that all expectations were met
	})

	t.Run("subscription change to disable if already disabled should succeed", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		mock.EXPECT().
			DisableRemindersWithNotification(ctx, "user@example.com", "HardBounce").
			Return(errs.FailedPrecondition.Error("reminders already disabled"))

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
		assert.Nil(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusOK)
		// we make sure that all expectations were met
	})

	t.Run("subscription change to enable should not update", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		_, _, handler := SetupHandler(t, ctx)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

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
		assert.Nil(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusOK)
		// we make sure that all expectations were met
	})

	t.Run("user not found", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, _, handler := SetupHandler(t, ctx)
		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

		mock.EXPECT().
			DisableRemindersWithNotification(ctx, "user@example.com", "HardBounce").
			Return(errs.NotFound.Error("user not found"))

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
		assert.Nil(t, err)

		// Check the status code is what we expect.
		AssertStatusEqual(t, rr, http.StatusOK)
		// we make sure that all expectations were met
	})
}
