-- +goose Up
CREATE TABLE IF NOT EXISTS user_event_notification_ (
    event_id integer NOT NULL,
    user_id integer NOT NULL,
    created timestamp NOT NULL DEFAULT timezone('utc', now()),
    CONSTRAINT event_fk FOREIGN KEY(event_id) REFERENCES event_(id) ON DELETE CASCADE,
    CONSTRAINT user_fk FOREIGN KEY(user_id) REFERENCES user_(id) ON DELETE CASCADE,
    UNIQUE(event_id, user_id)
);
ALTER TABLE user_ ADD COLUMN settings JSONB DEFAULT '{}'::jsonb;
UPDATE user_ SET settings = '{}'::jsonb;
ALTER TABLE user_ ALTER COLUMN settings SET NOT NULL;

-- +goose Down
DROP TABLE IF EXISTS user_event_notification_;
ALTER TABLE user_ DROP COLUMN settings;
