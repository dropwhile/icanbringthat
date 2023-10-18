-- +goose Up
ALTER TABLE user_ ADD COLUMN pwauth BOOLEAN DEFAULT TRUE;
UPDATE user_ SET pwauth = TRUE;
ALTER TABLE user_ ALTER COLUMN pwauth SET NOT NULL;

-- +goose Down
ALTER TABLE user_ DROP COLUMN pwauth;