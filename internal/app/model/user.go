package model

import (
	"context"
	"fmt"
	"time"

	"github.com/dropwhile/refid"

	"github.com/dropwhile/icbt/internal/util"
)

//go:generate go run ../../../cmd/refidgen -t User -v 1

type User struct {
	ID           int
	RefID        UserRefID `db:"ref_id"`
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

func (user *User) Save(ctx context.Context, db PgxHandle) error {
	return UpdateUser(ctx, db, user.Email, user.Name, user.PWHash, user.Verified, user.ID)
}

func NewUser(ctx context.Context, db PgxHandle, email, name string, rawPass []byte) (*User, error) {
	refID := refid.Must(NewUserRefID())
	pwHash, err := util.HashPW([]byte(rawPass))
	if err != nil {
		return nil, fmt.Errorf("error hashing pw: %w", err)
	}
	user, err := CreateUser(ctx, db, refID, email, name, pwHash)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func CreateUser(ctx context.Context, db PgxHandle, refID UserRefID, email, name string, pwHash []byte) (*User, error) {
	q := `INSERT INTO user_ (ref_id, email, name, pwhash) VALUES ($1, $2, $3, $4) RETURNING *`
	res, err := QueryOneTx[User](ctx, db, q, refID, email, name, pwHash)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func UpdateUser(ctx context.Context, db PgxHandle, email, name string, pwHash []byte, verified bool, userID int) error {
	q := `UPDATE user_ SET email = $1, name = $2, pwhash = $3, verified = $4 WHERE id = $5`
	return ExecTx[User](ctx, db, q, email, name, pwHash, verified, userID)
}

func DeleteUser(ctx context.Context, db PgxHandle, userID int) error {
	q := `DELETE FROM user_ WHERE id = $1`
	return ExecTx[User](ctx, db, q, userID)
}

func GetUserByID(ctx context.Context, db PgxHandle, id int) (*User, error) {
	q := `SELECT * FROM user_ WHERE id = $1`
	return QueryOne[User](ctx, db, q, id)
}

func GetUserByRefID(ctx context.Context, db PgxHandle, refID UserRefID) (*User, error) {
	q := `SELECT * FROM user_ WHERE ref_id = $1`
	return QueryOne[User](ctx, db, q, refID)
}

func GetUserByEmail(ctx context.Context, db PgxHandle, email string) (*User, error) {
	q := `SELECT * FROM user_ WHERE email = $1`
	return QueryOne[User](ctx, db, q, email)
}

func GetUsersByIDs(ctx context.Context, db PgxHandle, userIDs []int) ([]*User, error) {
	q := `SELECT * FROM user_ WHERE id = ANY($1)`
	return Query[User](ctx, db, q, userIDs)
}
