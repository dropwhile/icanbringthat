-- +goose Up
ALTER TABLE user_webauthn_ ADD COLUMN key_name varchar(255) NOT NULL;
ALTER TABLE user_webauthn_ ALTER COLUMN credential SET NOT NULL;

-- +goose Down
ALTER TABLE user_webauthn_ DROP COLUMN key_name;
ALTER TABLE user_webauthn_ ALTER COLUMN credential DROP NOT NULL;