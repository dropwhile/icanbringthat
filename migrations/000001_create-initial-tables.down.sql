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