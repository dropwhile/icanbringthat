package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pashagolub/pgxmock/v2"
)

func TestHandler_ShowIndex(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	ts := tstTs
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	rows := pgxmock.NewRows(
		[]string{"id", "ref_id", "event_item_id", "user_id", "note", "created", "last_modified"}).
		AddRow(1, "a", 1, 1, "some note", ts, ts)

	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO earmark_ (.+)*").
		WithArgs(pgxmock.AnyArg(), 1, 1, "some note").
		WillReturnRows(rows)
	mock.ExpectCommit()
	// hidden rollback after commit due to beginfunc being used
	mock.ExpectRollback()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	h := NewTestHandler(mock)
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.ShowIndex)
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := "/login"
	if rr.Header().Get("location") != expected {
		t.Errorf("handler returned unexpected redirect: got %v want %v",
			rr.Header().Get("location"),
			expected,
		)
	}
}
