-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_last_modified() RETURNS TRIGGER
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.last_modified := timezone('utc', CURRENT_TIMESTAMP);
    RETURN NEW;
END;
$$;
-- +goose StatementEnd
-- create domain/type
CREATE DOMAIN refid_bytea AS BYTEA
  CONSTRAINT check_length CHECK (octet_length(VALUE) = 16);

-- create user table
CREATE TABLE IF NOT EXISTS user_ (
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    ref_id refid_bytea NOT NULL,
    email varchar(255) NOT NULL,
    name varchar(255) NOT NULL,
    pwhash bytea NOT NULL,
    created timestamp NOT NULL DEFAULT timezone('utc', now()),
    last_modified timestamp NOT NULL DEFAULT timezone('utc', now())
);
CREATE UNIQUE INDEX user_ref_idx ON user_(ref_id);
CREATE UNIQUE INDEX user_email_idx ON user_(email);
CREATE TRIGGER last_mod_user
	BEFORE UPDATE ON user_
	FOR EACH ROW
    EXECUTE PROCEDURE update_last_modified();

-- create event table
CREATE TABLE IF NOT EXISTS event_ (
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    ref_id refid_bytea NOT NULL,
    user_id integer NOT NULL,
    name text NOT NULL,
    description text NOT NULL,
    start_time timestamptz,
    created timestamp NOT NULL DEFAULT timezone('utc', now()),
    last_modified timestamp NOT NULL DEFAULT timezone('utc', now()),
    CONSTRAINT user_fk FOREIGN KEY(user_id) REFERENCES user_(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX event_ref_idx ON event_(ref_id);
CREATE TRIGGER last_mod_event
	BEFORE UPDATE ON event_
	FOR EACH ROW
    EXECUTE PROCEDURE update_last_modified();

-- create event item
CREATE TABLE IF NOT EXISTS event_item_ (
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    ref_id refid_bytea NOT NULL,
    event_id integer NOT NULL,
    description text NOT NULL,
    created timestamp NOT NULL DEFAULT timezone('utc', now()),
    last_modified timestamp NOT NULL DEFAULT timezone('utc', now()),
    CONSTRAINT event_fk FOREIGN KEY(event_id) REFERENCES event_(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX event_item_ref_idx ON event_item_(ref_id);
CREATE TRIGGER last_mod_event_item_
	BEFORE UPDATE ON event_item_
	FOR EACH ROW
    EXECUTE PROCEDURE update_last_modified();

-- create earmark
CREATE TABLE IF NOT EXISTS earmark_ (
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    ref_id refid_bytea NOT NULL,
    event_item_id integer NOT NULL,
    user_id integer NOT NULL,
    notes text NOT NULL,
    created timestamp DEFAULT timezone('utc', now()),
    last_modified timestamp DEFAULT timezone('utc', now()),
    CONSTRAINT event_item_fk FOREIGN KEY(event_item_id) REFERENCES event_item_(id) ON DELETE CASCADE,
    CONSTRAINT user_fk FOREIGN KEY(user_id) REFERENCES user_(id) ON DELETE CASCADE,
    UNIQUE(event_item_id)
);
CREATE UNIQUE INDEX earmark_ref_idx ON earmark_(ref_id);
CREATE TRIGGER last_mod_earmark
	BEFORE UPDATE ON earmark_
	FOR EACH ROW
    EXECUTE PROCEDURE update_last_modified();

-- +goose Down
-- drop earmarks table/indexes/triggers
DROP INDEX IF EXISTS earmark_ref_idx;
DROP TRIGGER IF EXISTS last_mod_earmark ON earmark_;
DROP TABLE IF EXISTS earmark_;

-- drop events items table/indexes/triggers
DROP INDEX IF EXISTS event_item_ref_idx;
DROP TABLE IF EXISTS event_item_;

-- drop events table/indexes/triggers
DROP INDEX IF EXISTS event_ref_idx;
DROP TRIGGER IF EXISTS last_mod_event on event_;
DROP TABLE IF EXISTS event_;

-- drop user table/indexes/triggers
DROP INDEX IF EXISTS user_ref_idx;
CREATE TRIGGER last_mod_user on user_;
DROP TABLE IF EXISTS user_;

-- drop functions
DROP FUNCTION IF EXISTS update_last_modified;