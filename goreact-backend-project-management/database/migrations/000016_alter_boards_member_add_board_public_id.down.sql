ALTER TABLE lists
DROP CONSTRAINT IF EXISTS fk_board_public_id;

ALTER TABLE lists
DROP COLUMN IF EXISTS board_public_id;

