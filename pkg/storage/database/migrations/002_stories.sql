-- +migrate Up
ALTER TABLE entries
ADD COLUMN "has_story" BOOLEAN NOT NULL DEFAULT FALSE;

-- +migrate Down
ALTER TABLE entries
DROP COLUMN "has_story";

