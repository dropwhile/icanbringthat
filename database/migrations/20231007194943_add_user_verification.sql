-- +goose Up
CREATE TABLE IF NOT EXISTS user_verify_ (
    ref_id refid_bytea NOT NULL,
    user_id integer NOT NULL,
    created timestamp NOT NULL DEFAULT timezone('utc', now()),
    CONSTRAINT user_fk FOREIGN KEY(user_id) REFERENCES user_(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX IF NOT EXISTS user_verify_ref_idx ON user_verify_(ref_id);
ALTER TABLE user_ ADD COLUMN verified BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
DROP INDEX IF EXISTS user_verify_ref_idx;
DROP TABLE IF EXISTS user_verify_;
ALTER TABLE user_ DROP COLUMN verified;