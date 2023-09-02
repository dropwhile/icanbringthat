package model

import (
	"context"
	"fmt"
	"time"

	"github.com/dropwhile/icbt/internal/util"
	"github.com/dropwhile/icbt/internal/util/refid"
)

type User struct {
	Id           uint
	RefId        refid.RefId `db:"ref_id"`
	Email        string
	Name         string `db:"name"`
	PWHash       []byte
	Created      time.Time
	LastModified time.Time `db:"last_modified"`
}

func (user *User) SetPass(ctx context.Context, rawPass []byte) error {
	pwHash, err := util.HashPW([]byte(rawPass))
	if err != nil {
		return fmt.Errorf("error hashing pw: %w", err)
	}
	user.PWHash = pwHash
	return nil
}

func (user *User) CheckPass(ctx context.Context, rawPass []byte) (bool, error) {
	ok, err := util.CheckPWHash(user.PWHash, rawPass)
	if err != nil {
		return false, fmt.Errorf("error when comparing pass")
	}
	return ok, nil
}

func (user *User) Insert(db *DB, ctx context.Context) error {
	if user.RefId.IsNil() {
		user.RefId = UserRefIdT.MustNew()
	}
	q := `INSERT INTO user_ (ref_id, email, name, pwhash) VALUES ($1, $2, $3, $4) RETURNING *`
	res, err := QueryRowTx[User](db, ctx, q, user.RefId, user.Email, user.Name, user.PWHash)
	if err != nil {
		return err
	}
	user.Id = res.Id
	user.RefId = res.RefId
	user.Created = res.Created
	user.LastModified = res.LastModified
	return nil
}

func (user *User) Save(db *DB, ctx context.Context) error {
	q := `UPDATE user_ SET email = $1, name = $2, pwhash = $3 WHERE id = $4`
	return ExecTx[User](db, ctx, q, user.Email, user.Name, user.PWHash, user.Id)
}

func (user *User) Delete(db *DB, ctx context.Context) error {
	q := `DELETE FROM user_ WHERE id = $1`
	return ExecTx[User](db, ctx, q, user.Id)
}

func NewUser(db *DB, ctx context.Context, email, name string, rawPass []byte) (*User, error) {
	user := &User{
		Email: email,
		Name:  name,
	}
	err := user.SetPass(ctx, rawPass)
	if err != nil {
		return nil, fmt.Errorf("error hashing pw: %w", err)
	}
	err = user.Insert(db, ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserById(db *DB, ctx context.Context, id uint) (*User, error) {
	q := `SELECT * FROM user_ WHERE id = $1`
	return QueryRow[User](db, ctx, q, id)
}

func GetUserByRefId(db *DB, ctx context.Context, refId refid.RefId) (*User, error) {
	q := `SELECT * FROM user_ WHERE ref_id = $1`
	return QueryRow[User](db, ctx, q, refId)
}

func GetUserByEmail(db *DB, ctx context.Context, email string) (*User, error) {
	q := `SELECT * FROM user_ WHERE email = $1`
	return QueryRow[User](db, ctx, q, email)
}
