package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/dropwhile/refid"
	"github.com/dropwhile/refid/reftag"
	"github.com/jackc/pgx/v5"
)

type ApiKeyRefID = reftag.IDt9

var (
	NewApiKeyRefID       = reftag.NewRandom[ApiKeyRefID]
	ApiKeyRefIDMatcher   = reftag.NewMatcher[ApiKeyRefID]()
	ApiKeyRefIDFromBytes = reftag.FromBytes[ApiKeyRefID]
	ParseApiKeyRefID     = reftag.Parse[ApiKeyRefID]
)

type ApiKey struct {
	Token  string
	UserID int `db:"user_id"`
}

func NewApiKey(ctx context.Context, db PgxHandle,
	user *User,
) (*ApiKey, error) {
	if user == nil {
		return nil, fmt.Errorf("nil user supplied")
	}
	token := strings.Join(
		[]string{
			refid.Must(NewApiKeyRefID()).String(),
			refid.Must(NewApiKeyRefID()).String(),
		},
		":",
	)
	return CreateApiKey(ctx, db, user.ID, token)
}

func CreateApiKey(ctx context.Context, db PgxHandle,
	userID int, token string,
) (*ApiKey, error) {
	q := `
		INSERT INTO api_key_ (
			user_id, token
		)
		VALUES (@userID, @token)
		RETURNING *`
	args := pgx.NamedArgs{"userID": userID, "token": token}
	return QueryOneTx[ApiKey](ctx, db, q, args)
}

func GetUserByApiKey(ctx context.Context, db PgxHandle,
	token string,
) (*User, error) {
	q := `
		SELECT user_.*
		FROM api_key_
		JOIN user_ ON
			api_key_.user_id = user_.id
		WHERE
			user_.apkey = TRUE AND
			api_key_.token = $1
		`
	return QueryOne[User](ctx, db, q, token)
}
