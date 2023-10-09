-- name: CreateUser :one
INSERT INTO user_ (
    ref_id, email, name, pwhash
)
VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateUser :exec
UPDATE user_
SET 
    email = COALESCE(sqlc.narg('email'), email),
    name = COALESCE(sqlc.narg('name'), name),
    pwhash = COALESCE(sqlc.narg('pwhash'), pwhash),
    verified = COALESCE(sqlc.narg('verified'), verified)
WHERE id = @id;

-- name: DeleteUser :exec
DELETE FROM user_
WHERE id = $1;

-- name: GetUserById :one
SELECT * FROM user_
WHERE id = $1;

-- name: GetUserByRefID :one
SELECT * FROM user_
WHERE ref_id = $1;

-- name: GetUserByEmail :one
SELECT * FROM user_
WHERE email = $1;

-- name: GetUsersByIds :many
SELECT * FROM user_
WHERE id = ANY(@ids::int[]);