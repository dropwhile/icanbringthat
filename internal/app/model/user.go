package model

import (
	"context"
	"fmt"
	"time"

	"github.com/dropwhile/refid"
	"github.com/jackc/pgx/v5"

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
	WebAuthn     bool
	Created      time.Time
	LastModified time.Time `db:"last_modified"`
}

func HashPass(ctx context.Context, rawPass []byte) ([]byte, error) {
	pwHash, err := util.HashPW([]byte(rawPass))
	if err != nil {
		return nil, fmt.Errorf("error hashing pw: %w", err)
	}
	return pwHash, nil
}

func CheckPass(ctx context.Context,
	pwHash []byte, rawPass []byte,
) (bool, error) {
	ok, err := util.CheckPWHash(pwHash, rawPass)
	if err != nil {
		return false, fmt.Errorf("error when comparing pass")
	}
	return ok, nil
}

func NewUser(ctx context.Context, db PgxHandle,
	email, name string, rawPass []byte,
) (*User, error) {
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

func CreateUser(ctx context.Context, db PgxHandle,
	refID UserRefID, email, name string, pwHash []byte,
) (*User, error) {
	q := `
		INSERT INTO user_ (
			ref_id, email, name, pwhash
		)
		VALUES (
			@refID, @email, @name, @pwHash
		)
		RETURNING *`
	args := pgx.NamedArgs{
		"refID":  refID,
		"email":  email,
		"name":   name,
		"pwHash": pwHash,
	}
	res, err := QueryOneTx[User](ctx, db, q, args)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func UpdateUser(ctx context.Context, db PgxHandle,
	email, name string, pwHash []byte, verified bool,
	webAuthn bool, userID int,
) error {
	q := `
		UPDATE user_
		SET
			email = @email,
			name = @name,
			pwhash = @pwHash,
			verified = @verified,
			webauthn = @webAuthn
		WHERE id = @userID`
	args := pgx.NamedArgs{
		"email":    email,
		"name":     name,
		"pwHash":   pwHash,
		"verified": verified,
		"webAuthn": webAuthn,
		"userID":   userID,
	}
	return ExecTx[User](ctx, db, q, args)
}

func DeleteUser(ctx context.Context, db PgxHandle, userID int) error {
	q := `DELETE FROM user_ WHERE id = $1`
	return ExecTx[User](ctx, db, q, userID)
}

func GetUserByID(ctx context.Context, db PgxHandle,
	userID int,
) (*User, error) {
	q := `SELECT * FROM user_ WHERE id = $1`
	return QueryOne[User](ctx, db, q, userID)
}

func GetUserByRefID(ctx context.Context, db PgxHandle,
	refID UserRefID,
) (*User, error) {
	q := `SELECT * FROM user_ WHERE ref_id = $1`
	return QueryOne[User](ctx, db, q, refID)
}

func GetUserByEmail(ctx context.Context, db PgxHandle,
	email string,
) (*User, error) {
	q := `SELECT * FROM user_ WHERE email = $1`
	return QueryOne[User](ctx, db, q, email)
}

func GetUsersByIDs(ctx context.Context, db PgxHandle,
	userIDs []int,
) ([]*User, error) {
	q := `SELECT * FROM user_ WHERE id = ANY($1)`
	return Query[User](ctx, db, q, userIDs)
}
