package model

import (
	"context"
	"testing"

	"github.com/dropwhile/icbt/internal/util/refid"
	"github.com/pashagolub/pgxmock/v2"
	"gotest.tools/v3/assert"
)

var tstEarmarkRefId = refid.MustParse("0r2ncjgvqbr09f7c304v2a4rh4")

func TestEarmarkInsert(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refId := tstEarmarkRefId
	ts := tstTs
	rows := pgxmock.NewRows(
		[]string{"id", "ref_id", "event_item_id", "user_id", "note", "created", "last_modified"}).
		AddRow(1, refId, 1, 1, "some note", ts, ts)

	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO earmark_ (.+)*").
		WithArgs(pgxmock.AnyArg(), 1, 1, "some note").
		WillReturnRows(rows)
	mock.ExpectCommit()
	// hidden rollback after commit due to beginfunc being used
	mock.ExpectRollback()

	earmark, err := NewEarmark(ctx, mock, 1, 1, "some note")
	assert.NilError(t, err)
	assert.Check(t, earmark.RefId.HasTag(4))

	assert.DeepEqual(t, earmark, &Earmark{
		Id:           1,
		RefId:        refId,
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

	refId := tstEarmarkRefId
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
		RefId:        refId,
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

	refId := tstEarmarkRefId
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
		RefId:        refId,
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

	refId := tstEarmarkRefId
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
		RefId:        refId,
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

func TestEarmarkGetByRefId(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refId := tstEarmarkRefId
	ts := tstTs
	rows := pgxmock.NewRows(
		[]string{"id", "ref_id", "event_item_id", "user_id", "note", "created", "last_modified"}).
		AddRow(1, refId, 1, 1, "some note", ts, ts)

	mock.ExpectQuery("^SELECT (.+) FROM earmark_ *").
		WithArgs(refId).
		WillReturnRows(rows)

	earmark, err := GetEarmarkByRefId(ctx, mock, refId)
	assert.NilError(t, err)

	assert.DeepEqual(t, earmark, &Earmark{
		Id:           1,
		RefId:        refId,
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

	refId := tstEarmarkRefId
	ts := tstTs
	rows := pgxmock.NewRows(
		[]string{"id", "ref_id", "event_item_id", "user_id", "note", "created", "last_modified"}).
		AddRow(1, refId, 1, 1, "some note", ts, ts)

	mock.ExpectQuery("^SELECT (.+) FROM earmark_ *").
		WithArgs(1).
		WillReturnRows(rows)

	eventItem := &EventItem{
		Id:          1,
		RefId:       refId,
		EventId:     1,
		Description: "some desc",
	}
	earmark, err := GetEarmarkByEventItem(ctx, mock, eventItem)
	assert.NilError(t, err)

	assert.DeepEqual(t, earmark, &Earmark{
		Id:           1,
		RefId:        refId,
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
