package model

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dropwhile/icbt/internal/util/refid"
	"gotest.tools/v3/assert"
)

var tstEventItemRefId = refid.MustParse("0c0032mdyyanmzsgg4amq1at1aq2m")

func TestEventItemInsert(t *testing.T) {
	t.Parallel()
	model, mock := setupDBMock(t)
	t.Cleanup(func() { model.Close() })

	refId := tstEventItemRefId
	ts := tstTs
	rows := sqlmock.NewRows(
		[]string{"id", "ref_id", "event_id", "description", "created", "last_modified"}).
		AddRow(1, refId, 1, "some desc", ts, ts)

	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO event_item_ (.+)*").
		WithArgs(sqlmock.AnyArg(), 1, "some desc").
		WillReturnRows(rows)
	mock.ExpectCommit()

	ctx := context.TODO()
	eventItem, err := NewEventItem(model, ctx, 1, "some desc")
	assert.NilError(t, err)

	assert.Check(t, eventItem.RefId.HasTag(3))
	assert.DeepEqual(t, eventItem, &EventItem{
		Id:           1,
		RefId:        refId,
		EventId:      1,
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
	model, mock := setupDBMock(t)
	t.Cleanup(func() { model.Close() })

	refId := tstEventItemRefId

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE event_item_ (.+)*").
		WithArgs("some desc", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	ctx := context.TODO()
	eventItem := &EventItem{
		Id:          1,
		RefId:       refId,
		EventId:     1,
		Description: "some desc",
	}
	err := eventItem.Save(model, ctx)
	assert.NilError(t, err)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEventItemDelete(t *testing.T) {
	t.Parallel()
	model, mock := setupDBMock(t)
	t.Cleanup(func() { model.Close() })

	refId := tstEventItemRefId

	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM event_item_ (.+)*").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	ctx := context.TODO()
	eventItem := &EventItem{
		Id:          1,
		RefId:       refId,
		EventId:     1,
		Description: "some desc",
	}
	err := eventItem.Delete(model, ctx)
	assert.NilError(t, err)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEventItemGetById(t *testing.T) {
	t.Parallel()
	model, mock := setupDBMock(t)
	t.Cleanup(func() { model.Close() })

	refId := tstEventItemRefId
	rows := sqlmock.NewRows(
		[]string{"id", "ref_id", "event_id", "description"}).
		AddRow(1, refId, 1, "some desc")

	mock.ExpectQuery("^SELECT (.+) FROM event_item_ *").
		WithArgs(1).
		WillReturnRows(rows)

	ctx := context.TODO()
	eventItem, err := GetEventItemById(model, ctx, 1)
	assert.NilError(t, err)

	assert.DeepEqual(t, eventItem, &EventItem{
		Id:          1,
		RefId:       refId,
		EventId:     1,
		Description: "some desc",
	})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEventItemGetByRefId(t *testing.T) {
	t.Parallel()
	model, mock := setupDBMock(t)
	t.Cleanup(func() { model.Close() })

	refId := tstEventItemRefId
	rows := sqlmock.NewRows(
		[]string{"id", "ref_id", "event_id", "description"}).
		AddRow(1, refId, 1, "some desc")

	mock.ExpectQuery("^SELECT (.+) FROM event_item_ *").
		WithArgs(refId).
		WillReturnRows(rows)

	ctx := context.TODO()
	eventItem, err := GetEventItemByRefId(model, ctx, refId)
	assert.NilError(t, err)

	assert.DeepEqual(t, eventItem, &EventItem{
		Id:          1,
		RefId:       refId,
		EventId:     1,
		Description: "some desc",
	})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEventItemGetByEvent(t *testing.T) {
	t.Parallel()
	model, mock := setupDBMock(t)
	t.Cleanup(func() { model.Close() })

	refId := tstEventItemRefId
	ts := tstTs
	rows := sqlmock.NewRows(
		[]string{"id", "ref_id", "event_id", "description"}).
		AddRow(1, refId, 1, "some desc")

	mock.ExpectQuery("^SELECT (.+) FROM event_item_ *").
		WithArgs(1).
		WillReturnRows(rows)

	ctx := context.TODO()
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
	eventItems, err := GetEventItemsByEvent(model, ctx, event)
	assert.NilError(t, err)

	assert.DeepEqual(t, eventItems, []*EventItem{{
		Id:          1,
		RefId:       refId,
		EventId:     1,
		Description: "some desc",
	}})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
