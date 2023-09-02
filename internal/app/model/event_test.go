package model

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dropwhile/icbt/internal/util/refid"
	"gotest.tools/v3/assert"
)

var tstEventRefId = refid.MustParse("080032mdz2b9pwbbj2wemswmh95qr")

func TestEventInsert(t *testing.T) {
	t.Parallel()
	model, mock := setupDBMock(t)
	t.Cleanup(func() { model.Close() })

	refId := tstEventRefId
	ts := tstTs
	rows := sqlmock.NewRows(
		[]string{"id", "ref_id", "user_id", "name", "description", "start_time", "created", "last_modified"}).
		AddRow(1, refId, 1, "some name", "some desc", ts, ts, ts)

	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO event_ (.+)*").
		WithArgs(1, sqlmock.AnyArg(), "some name", "some desc", ts).
		WillReturnRows(rows)
	mock.ExpectCommit()

	ctx := context.TODO()
	event, err := NewEvent(model, ctx, 1, "some name", "some desc", ts)
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
	model, mock := setupDBMock(t)
	t.Cleanup(func() { model.Close() })

	refId := tstEventRefId
	ts := tstTs

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE event_ (.+)*").
		WithArgs("some name", "some desc", ts, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

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
	err := event.Save(model, ctx)
	assert.NilError(t, err)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEventDelete(t *testing.T) {
	t.Parallel()
	model, mock := setupDBMock(t)
	t.Cleanup(func() { model.Close() })

	refId := tstEventRefId
	ts := tstTs

	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM event_ (.+)*").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

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
	err := event.Delete(model, ctx)
	assert.NilError(t, err)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEventGetById(t *testing.T) {
	t.Parallel()
	model, mock := setupDBMock(t)
	t.Cleanup(func() { model.Close() })

	refId := tstEventRefId
	ts := tstTs
	rows := sqlmock.NewRows(
		[]string{"id", "ref_id", "user_id", "name", "description", "start_time", "created", "last_modified"}).
		AddRow(1, refId, 1, "some name", "some desc", ts, ts, ts)

	mock.ExpectQuery("^SELECT (.+) FROM event_ *").
		WithArgs(1).
		WillReturnRows(rows)

	ctx := context.TODO()
	event, err := GetEventById(model, ctx, 1)
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
	model, mock := setupDBMock(t)
	t.Cleanup(func() { model.Close() })

	refId := tstEventRefId
	ts := tstTs
	rows := sqlmock.NewRows(
		[]string{"id", "ref_id", "user_id", "name", "description", "start_time", "created", "last_modified"}).
		AddRow(1, refId, 1, "some name", "some desc", ts, ts, ts)

	mock.ExpectQuery("^SELECT (.+) FROM event_ *").
		WithArgs(refId).
		WillReturnRows(rows)

	ctx := context.TODO()
	event, err := GetEventByRefId(model, ctx, refId)
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
	model, mock := setupDBMock(t)
	t.Cleanup(func() { model.Close() })

	refId := tstEventRefId
	ts := tstTs
	rows := sqlmock.NewRows(
		[]string{"id", "ref_id", "user_id", "name", "description", "start_time", "created", "last_modified"}).
		AddRow(1, refId, 1, "some name", "some desc", ts, ts, ts)

	mock.ExpectQuery("^SELECT (.+) FROM event_ *").
		WithArgs(1).
		WillReturnRows(rows)

	ctx := context.TODO()
	user := &User{
		Id:     1,
		RefId:  refId,
		Email:  "user1@example.com",
		Name:   "j rando",
		PWHash: []byte("000x000"),
	}
	events, err := GetEventsByUser(model, ctx, user)
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
