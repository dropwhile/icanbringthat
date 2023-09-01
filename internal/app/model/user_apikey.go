// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package model

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dropwhile/refid/v2/reftag"
	"github.com/jackc/pgx/v5"

	"github.com/dropwhile/icanbringthat/internal/util"
)

type ApiKeyRefID struct {
	reftag.IDt9
}

var NewApiKeyRefID = reftag.NewRandom[ApiKeyRefID]

type ApiKey struct {
	Created time.Time
	Token   string
	UserID  int `db:"user_id"`
}

func NewApiKey(ctx context.Context, db PgxHandle,
	userID int,
) (*ApiKey, error) {
	if userID == 0 {
		return nil, fmt.Errorf("nil user supplied")
	}
	token := strings.Join(
		[]string{
			util.Must(NewApiKeyRefID()).String(),
			util.Must(NewApiKeyRefID()).String(),
		},
		":",
	)
	return CreateApiKey(ctx, db, userID, token)
}

func CreateApiKey(ctx context.Context, db PgxHandle,
	userID int, token string,
) (*ApiKey, error) {
	q := `
		INSERT INTO api_key_ (
			user_id, token
		)
		VALUES (@userID, @token)
		ON CONFLICT (user_id)
		DO
			UPDATE SET token = EXCLUDED.token
		RETURNING *`
	args := pgx.NamedArgs{"userID": userID, "token": token}
	return QueryOneTx[ApiKey](ctx, db, q, args)
}

func GetUserByApiKey(ctx context.Context, db PgxHandle,
	token string,
) (*User, error) {
	q := `
		SELECT user_.*
		FROM user_
		JOIN api_key_ ON
			api_key_.user_id = user_.id
		WHERE
			user_.api_access = TRUE AND
			api_key_.token = $1
		`
	return QueryOne[User](ctx, db, q, token)
}

func GetApiKeyByUser(ctx context.Context, db PgxHandle,
	userID int,
) (*ApiKey, error) {
	q := `
		SELECT *
		FROM api_key_
		WHERE
			user_id = $1
		`
	return QueryOne[ApiKey](ctx, db, q, userID)
}
