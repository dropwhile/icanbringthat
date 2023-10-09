-- name: CreateUserVerify :one
INSERT INTO user_verify_ (
    ref_id, user_id
)
VALUES ( $1, $2) RETURNING *;

-- name: DeleteUserVerify :exec
DELETE FROM user_verify_
WHERE ref_id = $1;

-- name: GetUserVerifyByRefID :one
SELECT * FROM user_verify_
WHERE ref_id = $1;