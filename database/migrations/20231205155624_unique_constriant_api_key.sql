-- +goose Up
CREATE UNIQUE INDEX api_key_user_id_idx ON api_key_(user_id);

-- +goose Down
DROP INDEX IF EXISTS api_key_user_id_idx;
