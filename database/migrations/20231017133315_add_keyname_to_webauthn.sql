-- +goose Up
ALTER TABLE user_webauthn_ ADD COLUMN key_name varchar(255) NOT NULL;
ALTER TABLE user_webauthn_ ADD COLUMN ref_id refid_bytea NOT NULL;
CREATE UNIQUE INDEX user_webauthn_ref_idx ON user_webauthn_(ref_id);
ALTER TABLE user_webauthn_ ALTER COLUMN credential SET NOT NULL;

-- +goose Down
ALTER TABLE user_webauthn_ ALTER COLUMN credential DROP NOT NULL;
DROP INDEX IF EXISTS user_webauthn_ref_idx;
ALTER TABLE user_webauthn_ DROP COLUMN IF EXISTS ref_id;
ALTER TABLE user_webauthn_ DROP COLUMN IF EXISTS key_name;