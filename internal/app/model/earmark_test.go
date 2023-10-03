package model

import (
	"context"
	"testing"

	"github.com/dropwhile/refid"
	"github.com/pashagolub/pgxmock/v3"
	"gotest.tools/v3/assert"
)

var tstEarmarkRefID = refid.Must(refid.Parse("065f77gsvzv092fr8x100fdx0c"))

func TestEarmarkInsert(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refId := tstEarmarkRefID
	ts := tstTs
	rows := pgxmock.NewRows(
		[]string{"id", "ref_id", "event_item_id", "user_id", "note", "created", "last_modified"}).
		AddRow(1, refId, 1, 1, "some note", ts, ts)

	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO earmark_ (.+)*").
		WithArgs(EarmarkRefIDT.AnyMatcher(), 1, 1, "some note").
		WillReturnRows(rows)
	mock.ExpectCommit()
	// hidden rollback after commit due to beginfunc being used
	mock.ExpectRollback()

	earmark, err := NewEarmark(ctx, mock, 1, 1, "some note")
	assert.NilError(t, err)
	assert.Check(t, earmark.RefID.HasTag(4))

	assert.DeepEqual(t, earmark, &Earmark{
		Id:           1,
		RefID:        refId,
		EventItemId:  1,
		UserId:       1,
		Note:         "some note",
		Created:      ts,
		LastModified: ts,
	})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEarmarkSave(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refId := tstEarmarkRefID
	ts := tstTs

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE earmark_ (.+)*").
		WithArgs("some note", 1).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	// hidden rollback after commit due to beginfunc being used
	mock.ExpectRollback()

	earmark := &Earmark{
		Id:           1,
		RefID:        refId,
		EventItemId:  1,
		UserId:       1,
		Note:         "some note",
		Created:      ts,
		LastModified: ts,
	}
	err = earmark.Save(ctx, mock)
	assert.NilError(t, err)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEarmarkDelete(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refId := tstEarmarkRefID
	ts := tstTs

	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM earmark_ (.+)*").
		WithArgs(1).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	// hidden rollback after commit due to beginfunc being used
	mock.ExpectRollback()

	earmark := &Earmark{
		Id:           1,
		RefID:        refId,
		EventItemId:  1,
		UserId:       1,
		Note:         "some note",
		Created:      ts,
		LastModified: ts,
	}
	err = earmark.Delete(ctx, mock)
	assert.NilError(t, err)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEarmarkGetById(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refId := tstEarmarkRefID
	ts := tstTs
	rows := pgxmock.NewRows(
		[]string{"id", "ref_id", "event_item_id", "user_id", "note", "created", "last_modified"}).
		AddRow(1, refId, 1, 1, "some note", ts, ts)

	mock.ExpectQuery("^SELECT (.+) FROM earmark_ *").
		WithArgs(1).
		WillReturnRows(rows)

	earmark, err := GetEarmarkById(ctx, mock, 1)
	assert.NilError(t, err)

	assert.DeepEqual(t, earmark, &Earmark{
		Id:           1,
		RefID:        refId,
		EventItemId:  1,
		UserId:       1,
		Note:         "some note",
		Created:      ts,
		LastModified: ts,
	})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEarmarkGetByRefID(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refId := tstEarmarkRefID
	ts := tstTs
	rows := pgxmock.NewRows(
		[]string{"id", "ref_id", "event_item_id", "user_id", "note", "created", "last_modified"}).
		AddRow(1, refId, 1, 1, "some note", ts, ts)

	mock.ExpectQuery("^SELECT (.+) FROM earmark_ *").
		WithArgs(refId).
		WillReturnRows(rows)

	earmark, err := GetEarmarkByRefID(ctx, mock, refId)
	assert.NilError(t, err)

	assert.DeepEqual(t, earmark, &Earmark{
		Id:           1,
		RefID:        refId,
		EventItemId:  1,
		UserId:       1,
		Note:         "some note",
		Created:      ts,
		LastModified: ts,
	})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEarmarkGetEventItem(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refId := tstEarmarkRefID
	ts := tstTs
	rows := pgxmock.NewRows(
		[]string{"id", "ref_id", "event_item_id", "user_id", "note", "created", "last_modified"}).
		AddRow(1, refId, 1, 1, "some note", ts, ts)

	mock.ExpectQuery("^SELECT (.+) FROM earmark_ *").
		WithArgs(1).
		WillReturnRows(rows)

	eventItem := &EventItem{
		Id:          1,
		RefID:       refId,
		EventId:     1,
		Description: "some desc",
	}
	earmark, err := GetEarmarkByEventItem(ctx, mock, eventItem)
	assert.NilError(t, err)

	assert.DeepEqual(t, earmark, &Earmark{
		Id:           1,
		RefID:        refId,
		EventItemId:  1,
		UserId:       1,
		Note:         "some note",
		Created:      ts,
		LastModified: ts,
	})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
