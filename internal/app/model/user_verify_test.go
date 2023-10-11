package model

import (
	"context"
	"testing"

	"github.com/dropwhile/refid"
	"github.com/pashagolub/pgxmock/v3"
	"gotest.tools/v3/assert"
)

var tstRefIDUserVerify = refid.Must(refid.Parse("065h0vbe7450c8ks641me3ny9c"))

func TestUserVerifyUserInsert(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	columns := []string{"ref_id", "user_id"}
	rows := pgxmock.NewRows(columns).
		AddRow(tstRefIDUserVerify, 1)

	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO user_verify_ (.+)").
		WithArgs(VerifyRefIDT.AnyMatcher(), 1).
		WillReturnRows(rows)
	mock.ExpectCommit()
	// hidden rollback after commit due to beginfunc being used
	mock.ExpectRollback()

	user := &User{
		ID:     1,
		RefID:  tstUserRefID,
		Email:  "user1@example.com",
		Name:   "j rando",
		PWHash: []byte("000x000"),
	}

	upw, err := NewUserVerify(ctx, mock, user)
	assert.NilError(t, err)

	assert.Check(t, upw.RefID.HasTag(6))
	assert.DeepEqual(t, upw, &UserVerify{
		RefID:  tstRefIDUserVerify,
		UserID: 1,
	})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserVerifyUserDelete(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM user_verify_ (.+)").
		WithArgs(tstRefIDUserVerify).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	// hidden rollback after commit due to beginfunc being used
	mock.ExpectRollback()

	upw := &UserVerify{
		RefID:  tstRefIDUserVerify,
		UserID: 1,
	}
	err = upw.Delete(ctx, mock)
	assert.NilError(t, err)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserVerifyUserGetByRefID(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	columns := []string{"ref_id", "user_id"}
	rows := pgxmock.NewRows(columns).
		AddRow(tstRefIDUserVerify, 1)

	mock.ExpectQuery("^SELECT (.+) FROM user_verify_ ").
		WithArgs(VerifyRefIDT.AnyMatcher()).
		WillReturnRows(rows)

	upw, err := GetUserVerifyByRefID(ctx, mock, tstRefIDUserVerify)
	assert.NilError(t, err)

	assert.DeepEqual(t, upw, &UserVerify{
		RefID:  tstRefIDUserVerify,
		UserID: 1,
	})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
