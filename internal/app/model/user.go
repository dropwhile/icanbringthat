package model

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/dropwhile/refid/v2/reftag"
	"github.com/jackc/pgx/v5"
	"github.com/samber/mo"

	"github.com/dropwhile/icbt/internal/crypto"
	"github.com/dropwhile/icbt/internal/util"
)

type UserRefID struct {
	reftag.IDt1
}

var NewUserRefID = reftag.New[UserRefID]

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
	ApiAccess    bool `db:"api_access"`
	WebAuthn     bool
}

func (u User) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Time("created", u.Created),
		slog.Time("last_modified", u.LastModified),
		slog.String("email", "_OMITTED_"),
		slog.String("pwhash", "_OMITTED_"),
		slog.Int("id", u.ID),
		slog.String("refid", u.RefID.String()),
		slog.Bool("verified", u.Verified),
		slog.Bool("pwauth", u.PWAuth),
		slog.Bool("api_access", u.ApiAccess),
		slog.Bool("webauthn", u.WebAuthn),
		slog.Any("settings", u.Settings),
	)
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
	refID := util.Must(NewUserRefID())
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

type UserUpdateModelValues struct {
	Name      mo.Option[string]
	Email     mo.Option[string]
	PWHash    mo.Option[[]byte]
	Verified  mo.Option[bool]
	PWAuth    mo.Option[bool]
	ApiAccess mo.Option[bool]
	WebAuthn  mo.Option[bool]
}

func UpdateUser(ctx context.Context, db PgxHandle, userID int,
	vals *UserUpdateModelValues,
) error {
	q := `
		UPDATE user_
		SET
			email = COALESCE(@email, email),
			name = COALESCE(@name, name),
			pwhash = COALESCE(@pwHash, pwhash),
			verified = COALESCE(@verified, verified),
			pwauth = COALESCE(@pwAuth, pwauth),
			api_access = COALESCE(@apiAccess, api_access),
			webauthn = COALESCE(@webAuthn, webauthn)
		WHERE id = @userID`
	args := pgx.NamedArgs{
		"userID":    userID,
		"email":     vals.Email,
		"name":      vals.Name,
		"pwHash":    vals.PWHash,
		"verified":  vals.Verified,
		"pwAuth":    vals.PWAuth,
		"apiAccess": vals.ApiAccess,
		"webAuthn":  vals.WebAuthn,
	}
	return ExecTx[User](ctx, db, q, args)
}

func UpdateUserSettings(ctx context.Context, db PgxHandle,
	userID int, pm *UserSettings,
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
