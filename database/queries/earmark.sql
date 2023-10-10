-- name: CreateEarmark :one
INSERT INTO earmark_ (
    ref_id, event_item_id, user_id, note
)
VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateEarmark :exec
UPDATE earmark_
SET note = $1
WHERE id = $2;

-- name: DeleteEarmark :exec
DELETE FROM earmark_
WHERE id = $1;

-- name: GetEarmarkById :one
SELECT * FROM earmark_
WHERE id = $1;

-- name: GetEarmarkByRefID :one
SELECT * FROM earmark_
WHERE ref_id = $1;

-- name: GetEarmarkByEventItem :one
SELECT * FROM earmark_
WHERE event_item_id = $1;

-- name: GetEarmarksByUserId :many
SELECT * FROM earmark_
WHERE
    user_id = $1
ORDER BY
    created DESC,
    id DESC;

-- name: GetEarmarksByEvent :many
SELECT earmark_.*
FROM earmark_
JOIN event_item_ ON 
    event_item_.id = earmark_.event_item_id
WHERE 
    event_item_.event_id = $1;

-- name: GetEarmarksByUserPaginated :many
SELECT * FROM earmark_
WHERE
    earmark_.user_id = $1
ORDER BY
    created DESC,
    id DESC
LIMIT $2 OFFSET $3;

-- name: GetEarmarkCountByUser :one
SELECT count(*) FROM earmark_
WHERE user_id = $1;