-- +goose Up
ALTER TABLE event_ ADD COLUMN archived BOOLEAN DEFAULT FALSE;
UPDATE event_ SET archived = FALSE;
ALTER TABLE event_ ALTER COLUMN archived SET NOT NULL;
CREATE INDEX event_archived_idx ON event_(archived);

-- +goose Down
ALTER TABLE event_ DROP COLUMN archived;