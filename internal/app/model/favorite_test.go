package model

import (
	"context"
	"testing"
	"time"

	"github.com/dropwhile/refid"
	"github.com/pashagolub/pgxmock/v3"
	"gotest.tools/v3/assert"
)

var favoriteColumns = []string{"id", "user_id", "event_id", "created"}

func TestFavoriteInsert(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	ts := tstTs
	rows := pgxmock.NewRows(favoriteColumns).
		AddRow(1, 1, 1, ts)

	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO favorite_ (.+)*").
		WithArgs(1, 1).
		WillReturnRows(rows)
	mock.ExpectCommit()
	// hidden rollback after commit due to beginfunc being used
	mock.ExpectRollback()

	favorite, err := NewFavorite(ctx, mock, 1, 1)
	assert.NilError(t, err)

	assert.DeepEqual(t, favorite, &Favorite{
		Id:      1,
		UserId:  1,
		EventId: 1,
		Created: ts,
	})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFavoriteDelete(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	ts := tstTs

	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM favorite_ (.+)*").
		WithArgs(1).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	// hidden rollback after commit due to beginfunc being used
	mock.ExpectRollback()

	favorite := &Favorite{
		Id:      1,
		UserId:  1,
		EventId: 1,
		Created: ts,
	}
	err = favorite.Delete(ctx, mock)
	assert.NilError(t, err)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFavoriteGetById(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	ts := tstTs
	rows := pgxmock.NewRows(favoriteColumns).
		AddRow(1, 1, 1, ts)

	mock.ExpectQuery("^SELECT (.+) FROM favorite_ *").
		WithArgs(1).
		WillReturnRows(rows)

	favorite, err := GetFavoriteById(ctx, mock, 1)
	assert.NilError(t, err)

	assert.DeepEqual(t, favorite, &Favorite{
		Id:      1,
		UserId:  1,
		EventId: 1,
		Created: ts,
	})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFavoriteGetByUser(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	ts := tstTs
	rows := pgxmock.NewRows(favoriteColumns).
		AddRow(1, 1, 1, ts)

	mock.ExpectQuery("^SELECT (.+) FROM favorite_ *").
		WithArgs(1).
		WillReturnRows(rows)

	user := &User{
		Id:     1,
		RefID:  tstUserRefID,
		Email:  "user1@example.com",
		Name:   "j rando",
		PWHash: []byte("000x000"),
	}
	favorites, err := GetFavoritesByUser(ctx, mock, user)
	assert.NilError(t, err)

	assert.DeepEqual(t, favorites, []*Favorite{{
		Id:      1,
		UserId:  user.Id,
		EventId: 1,
		Created: ts,
	}})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFavoriteGetByEvent(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refId := tstEventRefID
	ts := tstTs
	rows := pgxmock.NewRows(favoriteColumns).
		AddRow(1, 1, 1, ts)

	mock.ExpectQuery("^SELECT (.+) FROM favorite_ *").
		WithArgs(1).
		WillReturnRows(rows)

	event := &Event{
		Id:            1,
		RefID:         refId,
		UserId:        1,
		Name:          "some name",
		Description:   "some desc",
		StartTime:     time.Time{},
		StartTimeTZ:   "Etc/UTC",
		ItemSortOrder: []int{},
		Created:       ts,
		LastModified:  ts,
	}
	favorites, err := GetFavoritesByEvent(ctx, mock, event)
	assert.NilError(t, err)

	assert.DeepEqual(t, favorites, []*Favorite{{
		Id:      1,
		UserId:  1,
		EventId: event.Id,
		Created: ts,
	}})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFavoriteGetByUserEvent(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	ts := tstTs
	rows := pgxmock.NewRows(favoriteColumns).
		AddRow(1, 1, 1, ts)

	mock.ExpectQuery("^SELECT (.+) FROM favorite_ *").
		WithArgs(1, 1).
		WillReturnRows(rows)

	user := &User{
		Id:     1,
		RefID:  tstUserRefID,
		Email:  "user1@example.com",
		Name:   "j rando",
		PWHash: []byte("000x000"),
	}
	event := &Event{
		Id:           1,
		RefID:        refid.Must(EventRefIDT.New()),
		UserId:       user.Id,
		Name:         "event",
		Description:  "description",
		StartTime:    ts,
		StartTimeTZ:  "Etc/UTC",
		Created:      ts,
		LastModified: ts,
	}
	favorite, err := GetFavoriteByUserEvent(ctx, mock, user, event)
	assert.NilError(t, err)

	assert.DeepEqual(t, favorite, &Favorite{
		Id:      1,
		UserId:  1,
		EventId: 1,
		Created: ts,
	})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFavoriteCountByUser(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	rows := pgxmock.NewRows([]string{"count"}).
		AddRow(2)

	mock.ExpectQuery("^SELECT count(.+) FROM favorite_ ").
		WithArgs(1).
		WillReturnRows(rows)

	user := &User{
		Id:     1,
		RefID:  tstUserRefID,
		Email:  "user1@example.com",
		Name:   "j rando",
		PWHash: []byte("000x000"),
	}

	count, err := GetFavoriteCountByUser(ctx, mock, user)
	assert.NilError(t, err)
	assert.Equal(t, count, 2)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFavoriteGetByUserPaginated(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	ts := tstTs
	rows := pgxmock.NewRows(favoriteColumns).
		AddRow(1, 1, 1, ts).
		AddRow(2, 1, 2, ts).
		AddRow(3, 1, 3, ts)

	mock.ExpectQuery("^SELECT (.+) FROM favorite_ ").
		WithArgs(1, 10, 0).
		WillReturnRows(rows)

	user := &User{
		Id:     1,
		RefID:  tstUserRefID,
		Email:  "user1@example.com",
		Name:   "j rando",
		PWHash: []byte("000x000"),
	}

	favorites, err := GetFavoritesByUserPaginated(ctx, mock, user, 10, 0)
	assert.NilError(t, err)
	assert.DeepEqual(t, favorites, []*Favorite{
		{
			Id:      1,
			UserId:  1,
			EventId: 1,
			Created: ts,
		},
		{
			Id:      2,
			UserId:  1,
			EventId: 2,
			Created: ts,
		},
		{
			Id:      3,
			UserId:  1,
			EventId: 3,
			Created: ts,
		},
	})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
