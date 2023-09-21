-- +goose Up

ALTER TABLE user_pw_reset_
    ADD COLUMN created timestamp NOT NULL DEFAULT timezone('utc', now());

-- +goose Down

ALTER TABLE user_pw_reset_ DROP COLUMN created;