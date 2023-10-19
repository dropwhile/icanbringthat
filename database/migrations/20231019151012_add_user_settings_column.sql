-- +goose Up
ALTER TABLE user_ ADD COLUMN settings JSONB DEFAULT '{}'::jsonb;
UPDATE user_ SET settings = '{}'::jsonb;
ALTER TABLE user_ ALTER COLUMN settings SET NOT NULL;

-- +goose Down
ALTER TABLE user_ DROP COLUMN settings;
