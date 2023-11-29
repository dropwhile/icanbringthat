-- +goose Up
CREATE TABLE IF NOT EXISTS api_key_ (
    user_id integer NOT NULL,
    token varchar(60) NOT NULL,
    created timestamp NOT NULL DEFAULT timezone('utc', now()),
    CONSTRAINT user_fk FOREIGN KEY(user_id) REFERENCES user_(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX api_key_token_idx ON api_key_(token);
ALTER TABLE user_ ADD COLUMN apikey BOOLEAN NOT NULL DEFAULT FALSE;
UPDATE user_ SET apikey = FALSE;

-- +goose Down
ALTER TABLE user_ DROP COLUMN apikey;
DROP INDEX IF EXISTS api_key_token_idx;
DROP TABLE IF EXISTS api_key_;
