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

	refID := refid.Must(refid.Parse("065f77rd5400b4dk0p20b37n7r"))
	columns := []string{"ref_id", "user_id"}
	rows := pgxmock.NewRows(columns).
		AddRow(refID, 1)

	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO user_pw_reset_ (.+)").
		WithArgs(UserPWResetRefIDT.AnyMatcher(), 1).
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

	upw, err := NewUserPWReset(ctx, mock, user)
	assert.NilError(t, err)

	assert.Check(t, upw.RefID.HasTag(5))
	assert.DeepEqual(t, upw, &UserPWReset{
		RefID:  refID,
		UserID: 1,
	})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserPWResetDelete(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refID := refid.Must(refid.Parse("065f77rd5400b4dk0p20b37n7r"))
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM user_pw_reset_ (.+)").
		WithArgs(refID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	// hidden rollback after commit due to beginfunc being used
	mock.ExpectRollback()

	upw := &UserPWReset{
		RefID:  refID,
		UserID: 1,
	}
	err = upw.Delete(ctx, mock)
	assert.NilError(t, err)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserPWResetGetByRefID(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refID := refid.Must(refid.Parse("065f77rd5400b4dk0p20b37n7r"))
	columns := []string{"ref_id", "user_id"}
	rows := pgxmock.NewRows(columns).
		AddRow(refID, 1)

	mock.ExpectQuery("^SELECT (.+) FROM user_pw_reset_ ").
		WithArgs(refID).
		WillReturnRows(rows)

	upw, err := GetUserPWResetByRefID(ctx, mock, refID)
	assert.NilError(t, err)

	assert.DeepEqual(t, upw, &UserPWReset{
		RefID:  refID,
		UserID: 1,
	})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
