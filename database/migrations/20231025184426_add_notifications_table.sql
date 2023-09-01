-- +goose Up
CREATE TABLE IF NOT EXISTS notification_ (
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    ref_id refid_bytea NOT NULL,
    user_id integer NOT NULL,
    message text NOT NULL,
    read BOOLEAN NOT NULL DEFAULT FALSE,
    created timestamp NOT NULL DEFAULT timezone('utc', now()),
    last_modified timestamp NOT NULL DEFAULT timezone('utc', now()),
    CONSTRAINT user_fk FOREIGN KEY(user_id) REFERENCES user_(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX notification_ref_idx ON notification_(ref_id);
CREATE INDEX user_read_idx ON notification_(user_id, read);
CREATE TRIGGER last_mod_notification
	BEFORE UPDATE ON notification_
	FOR EACH ROW
    EXECUTE PROCEDURE update_last_modified();

-- +goose Down
DROP INDEX IF EXISTS notification_ref_idx;
DROP INDEX IF EXISTS user_read_idx;
DROP TRIGGER IF EXISTS last_mod_notification ON notification_;
DROP TABLE IF EXISTS notification_;