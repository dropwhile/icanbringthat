package model

import (
	"context"
	"fmt"
	"time"

	"github.com/dropwhile/refid"
	"github.com/dropwhile/refid/reftag"
	"github.com/jackc/pgx/v5"

	"github.com/dropwhile/icbt/internal/crypto"
)

type (
	UserRefID     = reftag.IDt1
	UserRefIDNull = reftag.NullIDt1
)

var (
	NewUserRefID       = reftag.New[UserRefID]
	UserRefIDMatcher   = reftag.NewMatcher[UserRefID]()
	UserRefIDFromBytes = reftag.FromBytes[UserRefID]
	ParseUserRefID     = reftag.Parse[UserRefID]
)

type User struct {
	Created      time.Time
	LastModified time.Time `db:"last_modified"`
	Email        string
	Name         string
	PWHash       []byte
	ID           int
	RefID        UserRefID `db:"ref_id"`
	Settings     UserSettings
	Verified     bool
	PWAuth       bool
	WebAuthn     bool
}

func HashPass(ctx context.Context, rawPass []byte) ([]byte, error) {
	return crypto.HashPW([]byte(rawPass))
}

func CheckPass(ctx context.Context,
	pwHash []byte, rawPass []byte,
) (bool, error) {
	return crypto.CheckPWHash(pwHash, rawPass)
}

func NewUser(ctx context.Context, db PgxHandle,
	email, name string, rawPass []byte,
) (*User, error) {
	refID := refid.Must(NewUserRefID())
	pwHash, err := crypto.HashPW([]byte(rawPass))
	if err != nil {
		return nil, fmt.Errorf("error hashing pw: %w", err)
	}
	return CreateUser(ctx, db, refID, email, name, pwHash)
}

func CreateUser(ctx context.Context, db PgxHandle,
	refID UserRefID, email, name string, pwHash []byte,
) (*User, error) {
	q := `
		INSERT INTO user_ (
			ref_id, email, name, pwhash, pwauth, settings
		)
		VALUES (
			@refID, @email, @name, @pwHash, @pwAuth, @settings
		)
		RETURNING *`
	args := pgx.NamedArgs{
		"refID":    refID,
		"email":    email,
		"name":     name,
		"pwHash":   pwHash,
		"pwAuth":   true,
		"settings": NewUserPropertyMap(),
	}
	return QueryOneTx[User](ctx, db, q, args)
}

func UpdateUser(ctx context.Context, db PgxHandle,
	email, name string, pwHash []byte, verified bool,
	pwAuth bool, webAuthn bool, userID int,
) error {
	q := `
		UPDATE user_
		SET
			email = @email,
			name = @name,
			pwhash = @pwHash,
			verified = @verified,
			pwauth = @pwAuth,
			webauthn = @webAuthn
		WHERE id = @userID`
	args := pgx.NamedArgs{
		"email":    email,
		"name":     name,
		"pwHash":   pwHash,
		"verified": verified,
		"pwAuth":   pwAuth,
		"webAuthn": webAuthn,
		"userID":   userID,
	}
	return ExecTx[User](ctx, db, q, args)
}

func UpdateUserSettings(ctx context.Context, db PgxHandle,
	pm *UserSettings, userID int,
) error {
	q := `
		UPDATE user_
		SET
			settings = @settings
		WHERE id = @userID`
	args := pgx.NamedArgs{
		"settings": pm,
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
