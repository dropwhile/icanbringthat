package model

import (
	"context"
	"fmt"
	"time"

	"github.com/dropwhile/icbt/internal/util"
	"github.com/dropwhile/refid"
)

var UserRefIdT = refid.Tagger(1)

type User struct {
	Id           int
	RefId        refid.RefId `db:"ref_id"`
	Email        string
	Name         string `db:"name"`
	PWHash       []byte `db:"pwhash"`
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

func (user *User) Insert(ctx context.Context, db PgxHandle) error {
	if user.RefId.IsNil() {
		user.RefId = refid.Must(UserRefIdT.New())
	}
	q := `INSERT INTO user_ (ref_id, email, name, pwhash) VALUES ($1, $2, $3, $4) RETURNING *`
	res, err := QueryOneTx[User](ctx, db, q, user.RefId, user.Email, user.Name, user.PWHash)
	if err != nil {
		return err
	}
	user.Id = res.Id
	user.RefId = res.RefId
	user.Created = res.Created
	user.LastModified = res.LastModified
	return nil
}

func (user *User) Save(ctx context.Context, db PgxHandle) error {
	q := `UPDATE user_ SET email = $1, name = $2, pwhash = $3 WHERE id = $4`
	return ExecTx[User](ctx, db, q, user.Email, user.Name, user.PWHash, user.Id)
}

func (user *User) Delete(ctx context.Context, db PgxHandle) error {
	q := `DELETE FROM user_ WHERE id = $1`
	return ExecTx[User](ctx, db, q, user.Id)
}

func NewUser(ctx context.Context, db PgxHandle, email, name string, rawPass []byte) (*User, error) {
	user := &User{
		Email: email,
		Name:  name,
	}
	err := user.SetPass(ctx, rawPass)
	if err != nil {
		return nil, fmt.Errorf("error hashing pw: %w", err)
	}
	err = user.Insert(ctx, db)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserById(ctx context.Context, db PgxHandle, id int) (*User, error) {
	q := `SELECT * FROM user_ WHERE id = $1`
	return QueryOne[User](ctx, db, q, id)
}

func GetUserByRefId(ctx context.Context, db PgxHandle, refId refid.RefId) (*User, error) {
	if !UserRefIdT.HasCorrectTag(refId) {
		err := fmt.Errorf(
			"bad refid type: got %d expected %d",
			refId.Tag(), UserRefIdT.Tag(),
		)
		return nil, err
	}
	q := `SELECT * FROM user_ WHERE ref_id = $1`
	return QueryOne[User](ctx, db, q, refId)
}

func GetUserByEmail(ctx context.Context, db PgxHandle, email string) (*User, error) {
	q := `SELECT * FROM user_ WHERE email = $1`
	return QueryOne[User](ctx, db, q, email)
}

func GetUsersByIds(ctx context.Context, db PgxHandle, userIds []int) ([]*User, error) {
	q := `SELECT user_.* FROM user_ WHERE id = ANY($1)`
	return Query[User](ctx, db, q, userIds)
}
