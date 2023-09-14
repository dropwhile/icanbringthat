package model

import (
	"context"
	"testing"

	"github.com/dropwhile/icbt/internal/util/refid"
	"github.com/pashagolub/pgxmock/v2"
	"gotest.tools/v3/assert"
)

var tstEventRefId = refid.MustParse("080032mdz2b9pwbbj2wemswmh95qr")

func TestEventInsert(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refId := tstEventRefId
	ts := tstTs
	rows := pgxmock.NewRows(
		[]string{"id", "ref_id", "user_id", "name", "description", "start_time", "created", "last_modified"}).
		AddRow(1, refId, 1, "some name", "some desc", ts, ts, ts)

	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO event_ (.+)*").
		WithArgs(1, pgxmock.AnyArg(), "some name", "some desc", ts).
		WillReturnRows(rows)
	mock.ExpectCommit()
	// hidden rollback after commit due to beginfunc being used
	mock.ExpectRollback()

	event, err := NewEvent(ctx, mock, 1, "some name", "some desc", ts)
	assert.NilError(t, err)

	assert.Check(t, event.RefId.HasTag(2))
	assert.DeepEqual(t, event, &Event{
		Id:           1,
		RefId:        refId,
		UserId:       1,
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

func TestEventSave(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refId := tstEventRefId
	ts := tstTs

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE event_ (.+)*").
		WithArgs("some name", "some desc", ts, 1).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	// hidden rollback after commit due to beginfunc being used
	mock.ExpectRollback()

	event := &Event{
		Id:           1,
		RefId:        refId,
		UserId:       1,
		Name:         "some name",
		Description:  "some desc",
		StartTime:    ts,
		Created:      ts,
		LastModified: ts,
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

	refId := tstEventRefId
	ts := tstTs

	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM event_ (.+)*").
		WithArgs(1).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	// hidden rollback after commit due to beginfunc being used
	mock.ExpectRollback()

	event := &Event{
		Id:           1,
		RefId:        refId,
		UserId:       1,
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

func TestEventGetById(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refId := tstEventRefId
	ts := tstTs
	rows := pgxmock.NewRows(
		[]string{"id", "ref_id", "user_id", "name", "description", "start_time", "created", "last_modified"}).
		AddRow(1, refId, 1, "some name", "some desc", ts, ts, ts)

	mock.ExpectQuery("^SELECT (.+) FROM event_ *").
		WithArgs(1).
		WillReturnRows(rows)

	event, err := GetEventById(ctx, mock, 1)
	assert.NilError(t, err)

	assert.DeepEqual(t, event, &Event{
		Id:           1,
		RefId:        refId,
		UserId:       1,
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

func TestEventGetByRefId(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refId := tstEventRefId
	ts := tstTs
	rows := pgxmock.NewRows(
		[]string{"id", "ref_id", "user_id", "name", "description", "start_time", "created", "last_modified"}).
		AddRow(1, refId, 1, "some name", "some desc", ts, ts, ts)

	mock.ExpectQuery("^SELECT (.+) FROM event_ *").
		WithArgs(refId).
		WillReturnRows(rows)

	event, err := GetEventByRefId(ctx, mock, refId)
	assert.NilError(t, err)

	assert.DeepEqual(t, event, &Event{
		Id:           1,
		RefId:        refId,
		UserId:       1,
		Name:         "some name",
		Description:  "some desc",
		StartTime:    ts,
		Created:      ts,
		LastModified: ts,
	})

	assert.Equal(t, event.RefId.String(), "080032mdz2b9pwbbj2wemswmh95qr")

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

	refId := tstEventRefId
	ts := tstTs
	rows := pgxmock.NewRows(
		[]string{"id", "ref_id", "user_id", "name", "description", "start_time", "created", "last_modified"}).
		AddRow(1, refId, 1, "some name", "some desc", ts, ts, ts)

	mock.ExpectQuery("^SELECT (.+) FROM event_ *").
		WithArgs(1).
		WillReturnRows(rows)

	user := &User{
		Id:     1,
		RefId:  refId,
		Email:  "user1@example.com",
		Name:   "j rando",
		PWHash: []byte("000x000"),
	}
	events, err := GetEventsByUser(ctx, mock, user)
	assert.NilError(t, err)

	assert.DeepEqual(t, events, []*Event{{
		Id:           1,
		RefId:        refId,
		UserId:       1,
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
