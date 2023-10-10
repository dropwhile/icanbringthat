// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: earmark.sql

package modelx

import (
	"context"

	"github.com/dropwhile/refid"
)

const createEarmark = `-- name: CreateEarmark :one
INSERT INTO earmark_ (
    ref_id, event_item_id, user_id, note
)
VALUES (
    $1, $2, $3, $4
)
RETURNING id, ref_id, event_item_id, user_id, note, created, last_modified
`

type CreateEarmarkParams struct {
	RefID       refid.RefID `db:"ref_id" json:"ref_id"`
	EventItemID int32       `db:"event_item_id" json:"event_item_id"`
	UserID      int32       `db:"user_id" json:"user_id"`
	Note        string      `db:"note" json:"note"`
}

func (q *Queries) CreateEarmark(ctx context.Context, arg CreateEarmarkParams) (Earmark, error) {
	row := q.db.QueryRow(ctx, createEarmark,
		arg.RefID,
		arg.EventItemID,
		arg.UserID,
		arg.Note,
	)
	var i Earmark
	err := row.Scan(
		&i.ID,
		&i.RefID,
		&i.EventItemID,
		&i.UserID,
		&i.Note,
		&i.Created,
		&i.LastModified,
	)
	return i, err
}

const deleteEarmark = `-- name: DeleteEarmark :exec
DELETE FROM earmark_
WHERE id = $1
`

func (q *Queries) DeleteEarmark(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteEarmark, id)
	return err
}

const getEarmarkByEventItem = `-- name: GetEarmarkByEventItem :many
SELECT id, ref_id, event_item_id, user_id, note, created, last_modified FROM earmark_
WHERE event_item_id = $1
`

func (q *Queries) GetEarmarkByEventItem(ctx context.Context, eventItemID int32) ([]Earmark, error) {
	rows, err := q.db.Query(ctx, getEarmarkByEventItem, eventItemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Earmark
	for rows.Next() {
		var i Earmark
		if err := rows.Scan(
			&i.ID,
			&i.RefID,
			&i.EventItemID,
			&i.UserID,
			&i.Note,
			&i.Created,
			&i.LastModified,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEarmarkById = `-- name: GetEarmarkById :one
SELECT id, ref_id, event_item_id, user_id, note, created, last_modified FROM earmark_
WHERE id = $1
`

func (q *Queries) GetEarmarkById(ctx context.Context, id int32) (Earmark, error) {
	row := q.db.QueryRow(ctx, getEarmarkById, id)
	var i Earmark
	err := row.Scan(
		&i.ID,
		&i.RefID,
		&i.EventItemID,
		&i.UserID,
		&i.Note,
		&i.Created,
		&i.LastModified,
	)
	return i, err
}

const getEarmarkByRefID = `-- name: GetEarmarkByRefID :one
SELECT id, ref_id, event_item_id, user_id, note, created, last_modified FROM earmark_
WHERE ref_id = $1
`

func (q *Queries) GetEarmarkByRefID(ctx context.Context, refID refid.RefID) (Earmark, error) {
	row := q.db.QueryRow(ctx, getEarmarkByRefID, refID)
	var i Earmark
	err := row.Scan(
		&i.ID,
		&i.RefID,
		&i.EventItemID,
		&i.UserID,
		&i.Note,
		&i.Created,
		&i.LastModified,
	)
	return i, err
}

const getEarmarkCountByUser = `-- name: GetEarmarkCountByUser :one
SELECT count(*) FROM earmark_
WHERE user_id = $1
`

func (q *Queries) GetEarmarkCountByUser(ctx context.Context, userID int32) (int64, error) {
	row := q.db.QueryRow(ctx, getEarmarkCountByUser, userID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getEarmarksByEvent = `-- name: GetEarmarksByEvent :many
SELECT earmark_.id, earmark_.ref_id, earmark_.event_item_id, earmark_.user_id, earmark_.note, earmark_.created, earmark_.last_modified
FROM earmark_
JOIN event_item_ ON 
    event_item_.id = earmark_.event_item_id
WHERE 
    event_item_.event_id = $1
`

func (q *Queries) GetEarmarksByEvent(ctx context.Context, eventID int32) ([]Earmark, error) {
	rows, err := q.db.Query(ctx, getEarmarksByEvent, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Earmark
	for rows.Next() {
		var i Earmark
		if err := rows.Scan(
			&i.ID,
			&i.RefID,
			&i.EventItemID,
			&i.UserID,
			&i.Note,
			&i.Created,
			&i.LastModified,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEarmarksByUserId = `-- name: GetEarmarksByUserId :many
SELECT id, ref_id, event_item_id, user_id, note, created, last_modified FROM earmark_
WHERE
    user_id = $1
ORDER BY
    created DESC,
    id DESC
`

func (q *Queries) GetEarmarksByUserId(ctx context.Context, userID int32) ([]Earmark, error) {
	rows, err := q.db.Query(ctx, getEarmarksByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Earmark
	for rows.Next() {
		var i Earmark
		if err := rows.Scan(
			&i.ID,
			&i.RefID,
			&i.EventItemID,
			&i.UserID,
			&i.Note,
			&i.Created,
			&i.LastModified,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEarmarksByUserPaginated = `-- name: GetEarmarksByUserPaginated :many
SELECT id, ref_id, event_item_id, user_id, note, created, last_modified FROM earmark_
WHERE
    earmark_.user_id = $1
ORDER BY
    created DESC,
    id DESC
LIMIT $2 OFFSET $3
`

type GetEarmarksByUserPaginatedParams struct {
	UserID int32 `db:"user_id" json:"user_id"`
	Limit  int32 `db:"limit" json:"limit"`
	Offset int32 `db:"offset" json:"offset"`
}

func (q *Queries) GetEarmarksByUserPaginated(ctx context.Context, arg GetEarmarksByUserPaginatedParams) ([]Earmark, error) {
	rows, err := q.db.Query(ctx, getEarmarksByUserPaginated, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Earmark
	for rows.Next() {
		var i Earmark
		if err := rows.Scan(
			&i.ID,
			&i.RefID,
			&i.EventItemID,
			&i.UserID,
			&i.Note,
			&i.Created,
			&i.LastModified,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateEarmark = `-- name: UpdateEarmark :exec
UPDATE earmark_
SET note = $1
WHERE id = $2
`

func (q *Queries) UpdateEarmark(ctx context.Context, note string, iD int32) error {
	_, err := q.db.Exec(ctx, updateEarmark, note, iD)
	return err
}
