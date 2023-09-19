package model

import (
	"context"
	"testing"

	"github.com/dropwhile/icbt/internal/util/refid"
	"github.com/pashagolub/pgxmock/v3"
	"gotest.tools/v3/assert"
)

var (
	columns      = []string{"id", "ref_id", "email", "name", "pwhash"}
	tstUserRefId = refid.MustParse("0r2nck7r2z7g35fefj2tv65thw")
)

func TestUserSetPassword(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refId := tstUserRefId
	user := &User{
		Id:     1,
		RefId:  refId,
		Email:  "user1@example.com",
		Name:   "j rando",
		PWHash: []byte("000x000"),
	}

	err = user.SetPass(ctx, []byte("111x111"))
	assert.NilError(t, err)

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE user_ (.+)*").
		WithArgs("user1@example.com", "j rando", user.PWHash, 1).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectCommit()
	// hidden rollback after commit due to beginfunc being used
	mock.ExpectRollback()

	err = user.Save(ctx, mock)
	assert.NilError(t, err)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserCheckPassword(t *testing.T) {
	t.Parallel()

	refId := tstUserRefId
	ctx := context.TODO()
	user := &User{
		Id:     1,
		RefId:  refId,
		Email:  "user1@example.com",
		Name:   "j rando",
		PWHash: []byte("000x000"),
	}

	err := user.SetPass(ctx, []byte("111x111"))
	assert.NilError(t, err)

	passok, err := user.CheckPass(ctx, []byte("111x111"))
	assert.NilError(t, err)
	assert.Equal(t, passok, true)

	passok, err = user.CheckPass(ctx, []byte("000x000"))
	assert.NilError(t, err)
	assert.Equal(t, passok, false)
}

func TestUserInsert(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refId := tstUserRefId
	rows := pgxmock.NewRows(columns).
		AddRow(1, refId, "user1@example.com", "j rando", []byte("000x000"))

	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO user_ (.+)*").
		WithArgs(UserRefIdT.AnyMatcher(), "user1@example.com", "j rando", pgxmock.AnyArg()).
		WillReturnRows(rows)
	mock.ExpectCommit()
	// hidden rollback after commit due to beginfunc being used
	mock.ExpectRollback()

	user, err := NewUser(ctx, mock, "user1@example.com", "j rando", []byte("000x000"))
	assert.NilError(t, err)

	assert.Check(t, user.RefId.HasTag(1))
	passOk, err := user.CheckPass(ctx, []byte("000x000"))
	assert.NilError(t, err)
	assert.Check(t, passOk)
	assert.DeepEqual(t, user, &User{
		Id:     1,
		RefId:  user.RefId,
		Email:  "user1@example.com",
		Name:   "j rando",
		PWHash: user.PWHash,
	})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserSave(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refId := tstUserRefId
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE user_ (.+)*").
		WithArgs("user1@example.com", "j rando", []byte("000x000"), 1).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	// hidden rollback after commit due to beginfunc being used
	mock.ExpectRollback()

	user := &User{
		Id:     1,
		RefId:  refId,
		Email:  "user1@example.com",
		Name:   "j rando",
		PWHash: []byte("000x000"),
	}
	err = user.Save(ctx, mock)
	assert.NilError(t, err)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserDelete(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refId := tstUserRefId
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM user_ (.+)*").
		WithArgs(1).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	// hidden rollback after commit due to beginfunc being used
	mock.ExpectRollback()

	user := &User{
		Id:     1,
		RefId:  refId,
		Email:  "user1@example.com",
		Name:   "j rando",
		PWHash: []byte("000x000"),
	}
	err = user.Delete(ctx, mock)
	assert.NilError(t, err)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserGetById(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refId := tstUserRefId
	rows := pgxmock.NewRows(columns).
		AddRow(1, refId, "user1@example.com", "j rando", []byte("000x000"))

	mock.ExpectQuery("^SELECT (.+) FROM user_ *").
		WithArgs(1).
		WillReturnRows(rows)

	user, err := GetUserById(ctx, mock, 1)
	assert.NilError(t, err)

	assert.DeepEqual(t, user, &User{
		Id:     1,
		RefId:  refId,
		Email:  "user1@example.com",
		Name:   "j rando",
		PWHash: []byte("000x000"),
	})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserGetByRefId(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refId := tstUserRefId
	rows := pgxmock.NewRows(columns).
		AddRow(1, refId, "user1@example.com", "j rando", []byte("000x000"))

	mock.ExpectQuery("^SELECT (.+) FROM user_ *").
		WithArgs(refId).
		WillReturnRows(rows)

	user, err := GetUserByRefId(ctx, mock, refId)
	assert.NilError(t, err)

	assert.DeepEqual(t, user, &User{
		Id:     1,
		RefId:  refId,
		Email:  "user1@example.com",
		Name:   "j rando",
		PWHash: []byte("000x000"),
	})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserGetByEmail(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	t.Cleanup(func() { mock.Close(ctx) })

	refId := tstUserRefId
	rows := pgxmock.NewRows(columns).
		AddRow(1, refId, "user1@example.com", "j rando", []byte("000x000"))

	mock.ExpectQuery("^SELECT (.+) FROM user_ *").
		WithArgs("user1@example.com").
		WillReturnRows(rows)

	user, err := GetUserByEmail(ctx, mock, "user1@example.com")
	assert.NilError(t, err)

	assert.DeepEqual(t, user, &User{
		Id:     1,
		RefId:  refId,
		Email:  "user1@example.com",
		Name:   "j rando",
		PWHash: []byte("000x000"),
	})

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
