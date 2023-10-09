-- name: CreateFavorite :one
INSERT INTO favorite_ (
    user_id, event_id
)
VALUES ($1, $2)
RETURNING *;

-- name: DeleteFavorite :exec
DELETE FROM favorite_
WHERE id = $1;

-- name: GetFavoriteById :one
-- @sqlc-vet-disable
SELECT * FROM favorite_
WHERE id = $1;

-- name: GetFavoritesByUserId :many
SELECT * FROM favorite_
WHERE user_id = $1
ORDER BY
    created DESC,
    id DESC;

-- name: GetFavoritesByEventId :many
SELECT * FROM favorite_
WHERE event_id = $1
ORDER BY
    created DESC,
    id DESC;

-- name: GetFavoriteByUserEvent :one
-- @sqlc-vet-disable
SELECT * FROM favorite_
WHERE
    user_id = $1 AND
    event_id = $2;

-- name: GetFavoriteCountByUser :one
SELECT count(*) FROM favorite_
WHERE user_id = $1;

-- name: GetFavoritesByUserPaginated :many
SELECT 
    favorite_.*
FROM favorite_ 
JOIN event_ ON
    favorite_.event_id = event_.id
WHERE favorite_.user_id = $1 
ORDER BY 
    event_.start_time DESC,
    event_.id DESC
LIMIT $2 OFFSET $3;