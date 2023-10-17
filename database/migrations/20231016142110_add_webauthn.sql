-- +goose Up
CREATE TABLE IF NOT EXISTS user_webauthn_ (
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id integer NOT NULL,
    created timestamp NOT NULL DEFAULT timezone('utc', now()),
    credential bytea,
    CONSTRAINT user_fk FOREIGN KEY(user_id) REFERENCES user_(id) ON DELETE CASCADE
);
ALTER TABLE user_ ADD COLUMN webauthn BOOLEAN NOT NULL DEFAULT FALSE;
UPDATE user_ SET webauthn = FALSE;

-- +goose Down
DROP TABLE IF EXISTS user_webauthn_;
ALTER TABLE user_ DROP COLUMN webauthn;