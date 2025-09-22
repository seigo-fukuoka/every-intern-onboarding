
-- +migrate Up
ALTER TABLE events ADD COLUMN description TEXT;

-- +migrate Down
ALTER TABLE events DROP COLUMN description;
