package model

import (
	"context"
	"time"

	"github.com/dropwhile/refid"
	"github.com/jackc/pgx/v5"
)

//go:generate go run ../../../cmd/refidgen -t Credential -v 7

type UserCredential struct {
	ID         int
	RefID      CredentialRefID `db:"ref_id"`
	UserID     int             `db:"user_id"`
	KeyName    string          `db:"key_name"`
	Created    time.Time
	Credential []byte
}

func NewUserCredential(ctx context.Context, db PgxHandle,
	userID int, keyName string, credential []byte,
) (*UserCredential, error) {
	refID := refid.Must(NewCredentialRefID())
	return CreateUserCredential(ctx, db, refID, userID, keyName, credential)
}

func CreateUserCredential(ctx context.Context, db PgxHandle,
	refID CredentialRefID, userID int, keyName string, credential []byte,
) (*UserCredential, error) {
	q := `
		INSERT INTO user_webauthn_ (
			ref_id, user_id, key_name, credential
		)
		VALUES (@refID, @userID, @keyName, @credential)
		RETURNING *`
	args := pgx.NamedArgs{
		"refID":      refID,
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

func GetUserCredentialByRefID(ctx context.Context, db PgxHandle,
	refID CredentialRefID,
) (*UserCredential, error) {
	q := `SELECT * FROM user_webauthn_ WHERE ref_id = $1`
	return QueryOne[UserCredential](ctx, db, q, refID)
}

func GetUserCredentialCountByUser(ctx context.Context, db PgxHandle,
	userID int,
) (int, error) {
	q := `SELECT count(*) FROM user_webauthn_ WHERE user_id = $1`
	return Get[int](ctx, db, q, userID)
}
