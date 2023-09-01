-- +goose Up
ALTER TABLE event_ ADD COLUMN item_sort_order INTEGER[] NOT NULL DEFAULT '{}';

-- +goose Down
ALTER TABLE event_ DROP COLUMN item_sort_order;
