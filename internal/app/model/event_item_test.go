package model

import (
	"context"
	"testing"

	"github.com/dropwhile/refid"
	"github.com/pashagolub/pgxmock/v3"
	"gotest.tools/v3/assert"
)

var tstEventItemRefID = refid.Must(ParseEventItemRefID("065f77kvj03069yhn864xhzq4r"))

func TestEventItemInsert(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refID := tstEventItemRefID
	ts := tstTs
	rows := pgxmock.NewRows(
		[]string{"id", "ref_id", "event_id", "description", "created", "last_modified"}).
		AddRow(1, refID, 1, "some desc", ts, ts)

	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO event_item_ (.+)*").
		WithArgs(EventItemRefIDMatcher{}, 1, "some desc").
		WillReturnRows(rows)
	mock.ExpectCommit()
	// hidden rollback after commit due to beginfunc being used
	mock.ExpectRollback()

	eventItem, err := NewEventItem(ctx, mock, 1, "some desc")
	assert.NilError(t, err)

	assert.Check(t, eventItem.RefID.HasTag(3))
	assert.DeepEqual(t, eventItem, &EventItem{
		ID:           1,
		RefID:        refID,
		EventID:      1,
		Description:  "some desc",
		Created:      ts,
		LastModified: ts,
	})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEventItemSave(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refID := tstEventItemRefID

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE event_item_ (.+)*").
		WithArgs("some desc", 1).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	// hidden rollback after commit due to beginfunc being used
	mock.ExpectRollback()

	eventItem := &EventItem{
		ID:          1,
		RefID:       refID,
		EventID:     1,
		Description: "some desc",
	}
	err = eventItem.Save(ctx, mock)
	assert.NilError(t, err)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEventItemDelete(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refID := tstEventItemRefID

	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM event_item_ (.+)*").
		WithArgs(1).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	// hidden rollback after commit due to beginfunc being used
	mock.ExpectRollback()

	eventItem := &EventItem{
		ID:          1,
		RefID:       refID,
		EventID:     1,
		Description: "some desc",
	}
	err = eventItem.Delete(ctx, mock)
	assert.NilError(t, err)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEventItemGetByID(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refID := tstEventItemRefID
	rows := pgxmock.NewRows(
		[]string{"id", "ref_id", "event_id", "description"}).
		AddRow(1, refID, 1, "some desc")

	mock.ExpectQuery("^SELECT (.+) FROM event_item_ *").
		WithArgs(1).
		WillReturnRows(rows)

	eventItem, err := GetEventItemByID(ctx, mock, 1)
	assert.NilError(t, err)

	assert.DeepEqual(t, eventItem, &EventItem{
		ID:          1,
		RefID:       refID,
		EventID:     1,
		Description: "some desc",
	})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEventItemGetByRefID(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refID := tstEventItemRefID
	rows := pgxmock.NewRows(
		[]string{"id", "ref_id", "event_id", "description"}).
		AddRow(1, refID, 1, "some desc")

	mock.ExpectQuery("^SELECT (.+) FROM event_item_ *").
		WithArgs(refID).
		WillReturnRows(rows)

	eventItem, err := GetEventItemByRefID(ctx, mock, refID)
	assert.NilError(t, err)

	assert.DeepEqual(t, eventItem, &EventItem{
		ID:          1,
		RefID:       refID,
		EventID:     1,
		Description: "some desc",
	})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEventItemGetByEvent(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refID := tstEventItemRefID
	ts := tstTs
	rows := pgxmock.NewRows(
		[]string{"id", "ref_id", "event_id", "description"}).
		AddRow(1, refID, 1, "some desc")

	mock.ExpectQuery("^SELECT (.+) FROM event_item_ *").
		WithArgs(1).
		WillReturnRows(rows)

	event := &Event{
		ID:           1,
		RefID:        refid.Must(NewEventRefID()),
		UserID:       1,
		Name:         "some name",
		Description:  "some desc",
		StartTime:    ts,
		Created:      ts,
		LastModified: ts,
	}
	eventItems, err := GetEventItemsByEvent(ctx, mock, event)
	assert.NilError(t, err)

	assert.DeepEqual(t, eventItems, []*EventItem{{
		ID:          1,
		RefID:       refID,
		EventID:     1,
		Description: "some desc",
	}})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}