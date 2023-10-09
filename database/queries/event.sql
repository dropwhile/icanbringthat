-- name: CreateEvent :one
INSERT INTO event_ (
    user_id, ref_id, name, description,
    item_sort_order, start_time, start_time_tz
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: UpdateEvent :exec
UPDATE event_
SET
    name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    item_sort_order = COALESCE(sqlc.narg('item_sort_order'), item_sort_order),
    start_time = COALESCE(sqlc.narg('start_time'), start_time),
    start_time_tz =  COALESCE(sqlc.narg('start_time_tz'), start_time_tz)
WHERE id = @id;

-- name: DeleteEvent :exec
DELETE FROM event_
WHERE id = $1;

-- name: GetEventById :one
-- @sqlc-vet-disable
SELECT * FROM event_
WHERE id = $1;

-- name: GetEventByRefId :one
SELECT * FROM event_
WHERE ref_id = $1;

-- name: GetEventCountByUser :one
SELECT count(*) FROM event_
WHERE user_id = $1;

-- name: GetEventsByIds :many
-- @sqlc-vet-disable
SELECT * FROM event_
WHERE id = ANY(@ids::int[]);

-- name: GetEventsByUserId :many
SELECT * FROM event_ 
WHERE
    event_.user_id = $1
ORDER BY start_time DESC, id DESC;

-- name: GetEventsByUserPaginated :many
SELECT *
FROM event_
WHERE
    event_.user_id = $1
ORDER BY
    start_time DESC,
    id DESC
LIMIT $2 OFFSET $3;

-- name: GetEventsComingSoonByUserPaginated :many
SELECT *
FROM event_
WHERE
    event_.user_id = $1 AND
    start_time > CURRENT_TIMESTAMP(0)
ORDER BY 
    start_time ASC,
    id ASC
LIMIT $2 OFFSET $3;