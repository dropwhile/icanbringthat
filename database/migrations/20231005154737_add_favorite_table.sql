-- +goose Up
CREATE TABLE IF NOT EXISTS favorite_ (
    id integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id integer NOT NULL,
    event_id integer NOT NULL,
    created timestamp NOT NULL DEFAULT timezone('utc', now()),
    CONSTRAINT user_fk FOREIGN KEY(user_id) REFERENCES user_(id) ON DELETE CASCADE,
    CONSTRAINT event_fk FOREIGN KEY(event_id) REFERENCES event_(id) ON DELETE CASCADE,
    UNIQUE (user_id, event_id)
);

-- +goose Down
DROP TABLE IF EXISTS favorite_;
