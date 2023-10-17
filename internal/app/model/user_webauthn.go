package model

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

//go:generate go run ../../../cmd/refidgen -t UserPWReset -v 5

type UserCredential struct {
	ID         int
	UserID     int    `db:"user_id"`
	KeyName    string `db:"key_name"`
	Created    time.Time
	Credential []byte
}

func NewUserCredential(ctx context.Context, db PgxHandle,
	userID int, keyName string, credential []byte,
) (*UserCredential, error) {
	return CreateUserCredential(ctx, db, userID, keyName, credential)
}

func CreateUserCredential(ctx context.Context, db PgxHandle,
	userID int, keyName string, credential []byte,
) (*UserCredential, error) {
	q := `
		INSERT INTO user_webauthn_ (
			user_id, key_name, credential
		)
		VALUES (@userID, @keyName, @credential)
		RETURNING *`
	args := pgx.NamedArgs{
		"userID":     userID,
		"credential": credential,
		"keyName":    keyName,
	}
	return QueryOneTx[UserCredential](ctx, db, q, args)
}

func DeleteUserCredential(ctx context.Context, db PgxHandle,
	ID int,
) error {
	q := `DELETE FROM user_webauthn_ WHERE id = $1`
	return ExecTx[UserCredential](ctx, db, q, ID)
}

func GetUserCredentialsByUser(ctx context.Context, db PgxHandle,
	userID int,
) ([]*UserCredential, error) {
	q := `SELECT * FROM user_webauthn_ WHERE user_id = $1`
	return Query[UserCredential](ctx, db, q, userID)
}
