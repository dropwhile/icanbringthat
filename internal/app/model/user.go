package model

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/dropwhile/refid/v2"
	"github.com/dropwhile/refid/v2/reftag"
	"github.com/jackc/pgx/v5"
	"github.com/samber/mo"

	"github.com/dropwhile/icbt/internal/crypto"
)

type UserRefID = reftag.IDt1

type UserRefIDNull struct {
	reftag.NullIDt1
}

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

func UpdateUser(ctx context.Context, db PgxHandle, userID int,
	email, name mo.Option[string], pwHash mo.Option[[]byte],
	verified, pwAuth, apiAccess, webAuthn mo.Option[bool],
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
		"email":     email,
		"name":      name,
		"pwHash":    pwHash,
		"verified":  verified,
		"pwAuth":    pwAuth,
		"apiAccess": apiAccess,
		"webAuthn":  webAuthn,
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
