ALTER TABLE boards
DROP CONSTRAINT IF EXISTS fk_boards_owner_public_id;

ALTER TABLE boards
DROP COLUMN IF EXISTS owner_public_id;