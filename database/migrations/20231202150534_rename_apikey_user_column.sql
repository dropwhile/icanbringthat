-- +goose Up
ALTER TABLE user_
RENAME COLUMN apikey TO api_access;

-- +goose Down
ALTER TABLE user_
RENAME COLUMN api_access TO apikey;