package model

import (
	"context"
	"fmt"
	"time"

	"github.com/dropwhile/refid"

	"github.com/dropwhile/icbt/internal/app/modelx"
	"github.com/dropwhile/icbt/internal/util"
)

type User struct {
	Id           int
	RefID        modelx.UserRefID `db:"ref_id"`
	Email        string
	Name         string `db:"name"`
	PWHash       []byte `db:"pwhash"`
	Verified     bool
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
	if user.RefID.IsNil() {
		user.RefID = refid.Must(modelx.NewUserRefID())
	}
	q := `INSERT INTO user_ (ref_id, email, name, pwhash) VALUES ($1, $2, $3, $4) RETURNING *`
	res, err := QueryOneTx[User](ctx, db, q, user.RefID, user.Email, user.Name, user.PWHash)
	if err != nil {
		return err
	}
	user.Id = res.Id
	user.RefID = res.RefID
	user.Created = res.Created
	user.LastModified = res.LastModified
	return nil
}

func (user *User) Save(ctx context.Context, db PgxHandle) error {
	q := `UPDATE user_ SET email = $1, name = $2, pwhash = $3, verified = $4 WHERE id = $5`
	return ExecTx[User](ctx, db, q, user.Email, user.Name, user.PWHash, user.Verified, user.Id)
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

func GetUserByRefID(ctx context.Context, db PgxHandle, refId modelx.UserRefID) (*User, error) {
	q := `SELECT * FROM user_ WHERE ref_id = $1`
	return QueryOne[User](ctx, db, q, refId)
}

func GetUserByEmail(ctx context.Context, db PgxHandle, email string) (*User, error) {
	q := `SELECT * FROM user_ WHERE email = $1`
	return QueryOne[User](ctx, db, q, email)
}

func GetUsersByIds(ctx context.Context, db PgxHandle, userIds []int) ([]*User, error) {
	q := `SELECT * FROM user_ WHERE id = ANY($1)`
	return Query[User](ctx, db, q, userIds)
}
