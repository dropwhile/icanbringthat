package model

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dropwhile/icbt/internal/util/refid"
	"gotest.tools/v3/assert"
)

var columns = []string{"id", "ref_id", "email", "name", "pwhash"}
var tstUserRefId = refid.MustParse("040032mdz53myygwgqj86s5dfjmvc")

func TestUserSetPassword(t *testing.T) {
	t.Parallel()
	model, mock := setupDBMock(t)
	t.Cleanup(func() { model.Close() })

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

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE user_ (.+)*").
		WithArgs("user1@example.com", "j rando", user.PWHash, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = user.Save(model, ctx)
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
	model, mock := setupDBMock(t)
	t.Cleanup(func() { model.Close() })

	refId := tstUserRefId
	rows := sqlmock.NewRows(columns).
		AddRow(1, refId, "user1@example.com", "j rando", []byte("000x000"))

	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO user_ (.+)*").
		WithArgs(sqlmock.AnyArg(), "user1@example.com", "j rando", sqlmock.AnyArg()).
		WillReturnRows(rows)
	mock.ExpectCommit()

	ctx := context.TODO()
	user, err := NewUser(model, ctx, "user1@example.com", "j rando", []byte("000x000"))
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
	model, mock := setupDBMock(t)
	t.Cleanup(func() { model.Close() })

	refId := tstUserRefId
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE user_ (.+)*").
		WithArgs("user1@example.com", "j rando", []byte("000x000"), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	ctx := context.TODO()
	user := &User{
		Id:     1,
		RefId:  refId,
		Email:  "user1@example.com",
		Name:   "j rando",
		PWHash: []byte("000x000"),
	}
	err := user.Save(model, ctx)
	assert.NilError(t, err)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserDelete(t *testing.T) {
	t.Parallel()
	model, mock := setupDBMock(t)
	t.Cleanup(func() { model.Close() })

	refId := tstUserRefId
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM user_ (.+)*").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	ctx := context.TODO()
	user := &User{
		Id:     1,
		RefId:  refId,
		Email:  "user1@example.com",
		Name:   "j rando",
		PWHash: []byte("000x000"),
	}
	err := user.Delete(model, ctx)
	assert.NilError(t, err)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserGetById(t *testing.T) {
	t.Parallel()
	model, mock := setupDBMock(t)
	t.Cleanup(func() { model.Close() })

	refId := tstUserRefId
	rows := sqlmock.NewRows(columns).
		AddRow(1, refId, "user1@example.com", "j rando", []byte("000x000"))

	mock.ExpectQuery("^SELECT (.+) FROM user_ *").
		WithArgs(1).
		WillReturnRows(rows)

	ctx := context.TODO()
	user, err := GetUserById(model, ctx, 1)
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
	model, mock := setupDBMock(t)
	t.Cleanup(func() { model.Close() })

	refId := tstUserRefId
	rows := sqlmock.NewRows(columns).
		AddRow(1, refId, "user1@example.com", "j rando", []byte("000x000"))

	mock.ExpectQuery("^SELECT (.+) FROM user_ *").
		WithArgs(refId).
		WillReturnRows(rows)

	ctx := context.TODO()
	user, err := GetUserByRefId(model, ctx, refId)
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
	model, mock := setupDBMock(t)
	t.Cleanup(func() { model.Close() })

	refId := tstUserRefId
	rows := sqlmock.NewRows(columns).
		AddRow(1, refId, "user1@example.com", "j rando", []byte("000x000"))

	mock.ExpectQuery("^SELECT (.+) FROM user_ *").
		WithArgs("user1@example.com").
		WillReturnRows(rows)

	ctx := context.TODO()
	user, err := GetUserByEmail(model, ctx, "user1@example.com")
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
