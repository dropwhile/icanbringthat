package model

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dropwhile/icbt/internal/util/refid"
	"gotest.tools/v3/assert"
)

var tstEarmarkRefId = refid.MustParse("0g0032mdytkwexqrmezkjf02qrz5w")

func TestEarmarkInsert(t *testing.T) {
	t.Parallel()
	model, mock := setupDBMock(t)
	t.Cleanup(func() { model.Close() })

	refId := tstEarmarkRefId
	ts := tstTs
	rows := sqlmock.NewRows(
		[]string{"id", "ref_id", "event_item_id", "user_id", "notes", "created", "last_modified"}).
		AddRow(1, refId, 1, 1, "some note", ts, ts)

	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO earmark_ (.+)*").
		WithArgs(1, 1, "some note").
		WillReturnRows(rows)
	mock.ExpectCommit()

	ctx := context.TODO()
	earmark, err := NewEarmark(model, ctx, 1, 1, "some note")
	assert.NilError(t, err)
	assert.Check(t, earmark.RefId.HasTag(4))

	assert.DeepEqual(t, earmark, &Earmark{
		Id:           1,
		RefId:        refId,
		EventItemId:  1,
		UserId:       1,
		Notes:        "some note",
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
	model, mock := setupDBMock(t)
	t.Cleanup(func() { model.Close() })

	refId := tstEarmarkRefId
	ts := tstTs

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE earmark_ (.+)*").
		WithArgs("some note", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	ctx := context.TODO()
	earmark := &Earmark{
		Id:           1,
		RefId:        refId,
		EventItemId:  1,
		UserId:       1,
		Notes:        "some note",
		Created:      ts,
		LastModified: ts,
	}
	err := earmark.Save(model, ctx)
	assert.NilError(t, err)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEarmarkDelete(t *testing.T) {
	t.Parallel()
	model, mock := setupDBMock(t)
	t.Cleanup(func() { model.Close() })

	refId := tstEarmarkRefId
	ts := tstTs

	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM earmark_ (.+)*").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	ctx := context.TODO()
	earmark := &Earmark{
		Id:           1,
		RefId:        refId,
		EventItemId:  1,
		UserId:       1,
		Notes:        "some note",
		Created:      ts,
		LastModified: ts,
	}
	err := earmark.Delete(model, ctx)
	assert.NilError(t, err)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEarmarkGetById(t *testing.T) {
	t.Parallel()
	model, mock := setupDBMock(t)
	t.Cleanup(func() { model.Close() })

	refId := tstEarmarkRefId
	ts := tstTs
	rows := sqlmock.NewRows(
		[]string{"id", "ref_id", "event_item_id", "user_id", "notes", "created", "last_modified"}).
		AddRow(1, refId, 1, 1, "some note", ts, ts)

	mock.ExpectQuery("^SELECT (.+) FROM earmark_ *").
		WithArgs(1).
		WillReturnRows(rows)

	ctx := context.TODO()
	earmark, err := GetEarmarkById(model, ctx, 1)
	assert.NilError(t, err)

	assert.DeepEqual(t, earmark, &Earmark{
		Id:           1,
		RefId:        refId,
		EventItemId:  1,
		UserId:       1,
		Notes:        "some note",
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
	model, mock := setupDBMock(t)
	t.Cleanup(func() { model.Close() })

	refId := tstEarmarkRefId
	ts := tstTs
	rows := sqlmock.NewRows(
		[]string{"id", "ref_id", "event_item_id", "user_id", "notes", "created", "last_modified"}).
		AddRow(1, refId, 1, 1, "some note", ts, ts)

	mock.ExpectQuery("^SELECT (.+) FROM earmark_ *").
		WithArgs(refId).
		WillReturnRows(rows)

	ctx := context.TODO()
	earmark, err := GetEarmarkByRefId(model, ctx, refId)
	assert.NilError(t, err)

	assert.DeepEqual(t, earmark, &Earmark{
		Id:           1,
		RefId:        refId,
		EventItemId:  1,
		UserId:       1,
		Notes:        "some note",
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
	model, mock := setupDBMock(t)
	t.Cleanup(func() { model.Close() })

	refId := tstEarmarkRefId
	ts := tstTs
	rows := sqlmock.NewRows(
		[]string{"id", "ref_id", "event_item_id", "user_id", "notes", "created", "last_modified"}).
		AddRow(1, refId, 1, 1, "some note", ts, ts)

	mock.ExpectQuery("^SELECT (.+) FROM earmark_ *").
		WithArgs(1).
		WillReturnRows(rows)

	ctx := context.TODO()
	eventItem := &EventItem{
		Id:          1,
		RefId:       refId,
		EventId:     1,
		Description: "some desc",
	}
	earmark, err := GetEarmarkByEventItem(model, ctx, eventItem)
	assert.NilError(t, err)

	assert.DeepEqual(t, earmark, &Earmark{
		Id:           1,
		RefId:        refId,
		EventItemId:  1,
		UserId:       1,
		Notes:        "some note",
		Created:      ts,
		LastModified: ts,
	})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
