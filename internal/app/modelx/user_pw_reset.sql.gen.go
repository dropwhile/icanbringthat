// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: user_pw_reset.sql

package modelx

import (
	"context"
)

const createUserPWReset = `-- name: CreateUserPWReset :one
INSERT INTO user_pw_reset_ (
    ref_id, user_id
)
VALUES ($1, $2) RETURNING ref_id, user_id, created
`

func (q *Queries) CreateUserPWReset(ctx context.Context, refID UserPwResetRefID, userID int32) (UserPwReset, error) {
	row := q.db.QueryRow(ctx, createUserPWReset, refID, userID)
	var i UserPwReset
	err := row.Scan(&i.RefID, &i.UserID, &i.Created)
	return i, err
}

const deleteUserPWReset = `-- name: DeleteUserPWReset :exec
DELETE FROM user_pw_reset_
WHERE ref_id = $1
`

func (q *Queries) DeleteUserPWReset(ctx context.Context, refID UserPwResetRefID) error {
	_, err := q.db.Exec(ctx, deleteUserPWReset, refID)
	return err
}

const getUserPWResetByRefID = `-- name: GetUserPWResetByRefID :one
SELECT ref_id, user_id, created FROM user_pw_reset_
WHERE ref_id = $1
`

func (q *Queries) GetUserPWResetByRefID(ctx context.Context, refID UserPwResetRefID) (UserPwReset, error) {
	row := q.db.QueryRow(ctx, getUserPWResetByRefID, refID)
	var i UserPwReset
	err := row.Scan(&i.RefID, &i.UserID, &i.Created)
	return i, err
}
