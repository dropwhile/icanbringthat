// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: user_verify.sql

package modelx

import (
	"context"

	"github.com/dropwhile/refid"
)

const createUserVerify = `-- name: CreateUserVerify :one
INSERT INTO user_verify_ (
    ref_id, user_id
)
VALUES ( $1, $2) RETURNING ref_id, user_id, created
`

func (q *Queries) CreateUserVerify(ctx context.Context, refID refid.RefID, userID int32) (UserVerify, error) {
	row := q.db.QueryRow(ctx, createUserVerify, refID, userID)
	var i UserVerify
	err := row.Scan(&i.RefID, &i.UserID, &i.Created)
	return i, err
}

const deleteUserVerify = `-- name: DeleteUserVerify :exec
DELETE FROM user_verify_
WHERE ref_id = $1
`

func (q *Queries) DeleteUserVerify(ctx context.Context, refID refid.RefID) error {
	_, err := q.db.Exec(ctx, deleteUserVerify, refID)
	return err
}

const getUserVerifyByRefID = `-- name: GetUserVerifyByRefID :one
SELECT ref_id, user_id, created FROM user_verify_
WHERE ref_id = $1
`

func (q *Queries) GetUserVerifyByRefID(ctx context.Context, refID refid.RefID) (UserVerify, error) {
	row := q.db.QueryRow(ctx, getUserVerifyByRefID, refID)
	var i UserVerify
	err := row.Scan(&i.RefID, &i.UserID, &i.Created)
	return i, err
}
