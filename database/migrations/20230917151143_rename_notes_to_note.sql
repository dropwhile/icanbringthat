-- +goose Up
ALTER TABLE earmark_ RENAME COLUMN notes TO note;

-- +goose Down
ALTER TABLE earmark_ RENAME COLUMN note TO notes;