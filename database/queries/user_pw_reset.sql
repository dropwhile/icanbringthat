-- name: CreateUserPWReset :one
INSERT INTO user_pw_reset_ (
    ref_id, user_id
)
VALUES ($1, $2) RETURNING *;

-- name: DeleteUserPWReset :exec
DELETE FROM user_pw_reset_
WHERE ref_id = $1;

-- name: GetUserPWResetByRefID :one
SELECT * FROM user_pw_reset_
WHERE ref_id = $1;