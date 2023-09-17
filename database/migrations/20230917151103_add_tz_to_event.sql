-- +goose Up
ALTER TABLE event_
    ADD COLUMN start_time_tz varchar(255) DEFAULT 'Etc/UTC';

UPDATE event_ 
    SET start_time_tz='Etc/UTC';

ALTER TABLE event_
    ALTER COLUMN start_time_tz SET NOT NULL;

-- +goose Down
ALTER TABLE event_ DROP COLUMN start_time_tz;