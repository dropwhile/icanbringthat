package model

import (
	"context"
	"testing"

	"github.com/dropwhile/refid"
	"github.com/pashagolub/pgxmock/v3"
	"gotest.tools/v3/assert"
)

var tstEventRefID = refid.Must(ParseEventRefID("065f77p96rp05pze29y74wx7q8"))

func TestEventInsert(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refID := tstEventRefID
	ts := tstTs
	rows := pgxmock.NewRows(
		[]string{
			"id", "ref_id", "user_id", "name", "description", "start_time",
			"start_time_tz", "created", "last_modified",
		}).
		AddRow(1, refID, 1, "some name", "some desc", ts, "Etc/UTC", ts, ts)

	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO event_ (.+)*").
		WithArgs(1, EventRefIDMatcher{}, "some name", "some desc", []int{}, ts, "Etc/UTC").
		WillReturnRows(rows)
	mock.ExpectCommit()
	// hidden rollback after commit due to beginfunc being used
	mock.ExpectRollback()

	event, err := NewEvent(ctx, mock, 1, "some name", "some desc", ts, "Etc/UTC")
	assert.NilError(t, err)

	assert.Check(t, event.RefID.HasTag(2))
	assert.DeepEqual(t, event, &Event{
		ID:           1,
		RefID:        refID,
		UserID:       1,
		Name:         "some name",
		Description:  "some desc",
		StartTime:    ts,
		StartTimeTZ:  "Etc/UTC",
		Created:      ts,
		LastModified: ts,
	})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEventSave(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refID := tstEventRefID
	ts := tstTs

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE event_ (.+)*").
		WithArgs("some name", "some desc", []int{1, 2, 3}, ts, "Etc/UTC", 1).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	// hidden rollback after commit due to beginfunc being used
	mock.ExpectRollback()

	event := &Event{
		ID:            1,
		RefID:         refID,
		UserID:        1,
		Name:          "some name",
		Description:   "some desc",
		ItemSortOrder: []int{1, 2, 3},
		StartTime:     ts,
		StartTimeTZ:   "Etc/UTC",
		Created:       ts,
		LastModified:  ts,
	}
	err = event.Save(ctx, mock)
	assert.NilError(t, err)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEventDelete(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refID := tstEventRefID
	ts := tstTs

	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM event_ (.+)*").
		WithArgs(1).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	// hidden rollback after commit due to beginfunc being used
	mock.ExpectRollback()

	event := &Event{
		ID:           1,
		RefID:        refID,
		UserID:       1,
		Name:         "some name",
		Description:  "some desc",
		StartTime:    ts,
		Created:      ts,
		LastModified: ts,
	}
	err = event.Delete(ctx, mock)
	assert.NilError(t, err)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEventGetByID(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refID := tstEventRefID
	ts := tstTs
	rows := pgxmock.NewRows(
		[]string{"id", "ref_id", "user_id", "name", "description", "start_time", "created", "last_modified"}).
		AddRow(1, refID, 1, "some name", "some desc", ts, ts, ts)

	mock.ExpectQuery("^SELECT (.+) FROM event_ *").
		WithArgs(1).
		WillReturnRows(rows)

	event, err := GetEventByID(ctx, mock, 1)
	assert.NilError(t, err)

	assert.DeepEqual(t, event, &Event{
		ID:           1,
		RefID:        refID,
		UserID:       1,
		Name:         "some name",
		Description:  "some desc",
		StartTime:    ts,
		Created:      ts,
		LastModified: ts,
	})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEventGetByRefID(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refID := tstEventRefID
	ts := tstTs
	rows := pgxmock.NewRows(
		[]string{"id", "ref_id", "user_id", "name", "description", "start_time", "created", "last_modified"}).
		AddRow(1, refID, 1, "some name", "some desc", ts, ts, ts)

	mock.ExpectQuery("^SELECT (.+) FROM event_ *").
		WithArgs(refID).
		WillReturnRows(rows)

	event, err := GetEventByRefID(ctx, mock, refID)
	assert.NilError(t, err)

	assert.DeepEqual(t, event, &Event{
		ID:           1,
		RefID:        refID,
		UserID:       1,
		Name:         "some name",
		Description:  "some desc",
		StartTime:    ts,
		Created:      ts,
		LastModified: ts,
	})

	assert.DeepEqual(t, event.RefID.String(), tstEventRefID.String())

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEventsGetByUser(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refID := tstEventRefID
	ts := tstTs
	rows := pgxmock.NewRows(
		[]string{"id", "ref_id", "user_id", "name", "description", "start_time", "created", "last_modified"}).
		AddRow(1, refID, 1, "some name", "some desc", ts, ts, ts)

	mock.ExpectQuery("^SELECT (.+) FROM event_ *").
		WithArgs(1).
		WillReturnRows(rows)

	user := &User{
		ID:     1,
		RefID:  refid.Must(NewUserRefID()),
		Email:  "user1@example.com",
		Name:   "j rando",
		PWHash: []byte("000x000"),
	}
	events, err := GetEventsByUser(ctx, mock, user)
	assert.NilError(t, err)

	assert.DeepEqual(t, events, []*Event{{
		ID:           1,
		RefID:        refID,
		UserID:       1,
		Name:         "some name",
		Description:  "some desc",
		StartTime:    ts,
		Created:      ts,
		LastModified: ts,
	}})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
