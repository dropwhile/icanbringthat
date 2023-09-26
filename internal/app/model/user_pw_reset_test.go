package model

import (
	"context"
	"testing"

	"github.com/dropwhile/refid"
	"github.com/pashagolub/pgxmock/v3"
	"gotest.tools/v3/assert"
)

func TestUserPWResetInsert(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refId := refid.Must(refid.Parse("0r2y7pynskjgapab7cg0a87ny8"))
	columns := []string{"ref_id", "user_id"}
	rows := pgxmock.NewRows(columns).
		AddRow(refId, 1)

	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO user_pw_reset_ (.+)").
		WithArgs(UserPWResetRefIdT.AnyMatcher(), 1).
		WillReturnRows(rows)
	mock.ExpectCommit()
	// hidden rollback after commit due to beginfunc being used
	mock.ExpectRollback()

	user := &User{
		Id:     1,
		RefId:  tstUserRefId,
		Email:  "user1@example.com",
		Name:   "j rando",
		PWHash: []byte("000x000"),
	}

	upw, err := NewUserPWReset(ctx, mock, user)
	assert.NilError(t, err)

	assert.Check(t, upw.RefId.HasTag(5))
	assert.DeepEqual(t, upw, &UserPWReset{
		RefId:  refId,
		UserId: 1,
	})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserPWReserDelete(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refId := refid.Must(refid.Parse("0r2y7pynskjgapab7cg0a87ny8"))
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM user_pw_reset_ (.+)").
		WithArgs(refId).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	// hidden rollback after commit due to beginfunc being used
	mock.ExpectRollback()

	upw := &UserPWReset{
		RefId:  refId,
		UserId: 1,
	}
	err = upw.Delete(ctx, mock)
	assert.NilError(t, err)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserPWReserGetByRefId(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refId := refid.Must(refid.Parse("0r2y7pynskjgapab7cg0a87ny8"))
	columns := []string{"ref_id", "user_id"}
	rows := pgxmock.NewRows(columns).
		AddRow(refId, 1)

	mock.ExpectQuery("^SELECT (.+) FROM user_pw_reset_ ").
		WithArgs(refId).
		WillReturnRows(rows)

	upw, err := GetUserPWResetByRefId(ctx, mock, refId)
	assert.NilError(t, err)

	assert.DeepEqual(t, upw, &UserPWReset{
		RefId:  refId,
		UserId: 1,
	})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
