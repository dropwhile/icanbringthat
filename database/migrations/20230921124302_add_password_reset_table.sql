-- +goose Up
CREATE TABLE IF NOT EXISTS user_pw_reset_ (
    ref_id refid_bytea NOT NULL,
    user_id integer NOT NULL,
    created timestamp NOT NULL DEFAULT timezone('utc', now()),
    CONSTRAINT user_fk FOREIGN KEY(user_id) REFERENCES user_(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX user_pw_reset_ref_idx ON user_pw_reset_(ref_id);

-- +goose Down
DELETE INDEX IF EXISTS user_pw_reset_ref_idx;
DROP TABLE IF EXISTS user_pw_reset_;
