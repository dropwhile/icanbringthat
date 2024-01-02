package handler

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"gotest.tools/v3/assert"
)

func TestHandler_ShowIndex_LoggedOut(t *testing.T) {
	t.Parallel()

	t.Run("logged in", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, mux, handler := SetupHandlerOld(t, ctx)
		mux.Get("/", handler.ShowIndex)

		// create request
		req, err := http.NewRequest("GET", "/", nil)
		assert.NilError(t, err)

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		response := rr.Result()
		_, err = io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		assert.Assert(
			t, StatusEqual(rr, http.StatusSeeOther),
			"handler returned wrong status code")
		assert.Equal(t, rr.Header().Get("location"), "/login",
			"handler returned wrong redirect")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("logged out", func(t *testing.T) {
		t.Parallel()

		ctx := context.TODO()
		mock, mux, handler := SetupHandlerOld(t, ctx)
		cookie := SetupUserSession(t, mux, mock, handler)
		mux.Get("/", handler.ShowIndex)

		// create request
		req, err := http.NewRequest("GET", "/", nil)
		assert.NilError(t, err)
		setCookie(req, cookie)

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		response := rr.Result()
		_, err = io.ReadAll(response.Body)
		assert.NilError(t, err)

		// Check the status code is what we expect.
		assert.Assert(t, StatusEqual(rr, http.StatusSeeOther),
			"handler returned wrong status code")
		assert.Equal(t, rr.Header().Get("location"), "/dashboard",
			"handler returned wrong redirect")
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}
