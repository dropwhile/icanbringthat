-- name: CreateEventItem :one
INSERT INTO event_item_ (
    ref_id, event_id, description
)
VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: UpdateEventItem :exec
UPDATE event_item_
SET
    description = $1
WHERE id = $2;

-- name: DeleteEventItem :exec
DELETE FROM event_item_
WHERE id = $1;

-- name: GetEventItemById :one
SELECT * FROM event_item_
WHERE id = $1;

-- name: GetEventItemByRefId :one
SELECT * FROM event_item_
WHERE ref_id = $1;

-- name: GetEventItemsByEvent :many
SELECT * FROM event_item_
WHERE event_id = $1
ORDER BY
    created DESC,
    id DESC;
