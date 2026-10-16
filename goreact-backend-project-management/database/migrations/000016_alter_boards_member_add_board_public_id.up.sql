ALTER TABLE lists
    ADD COLUMN board_public_id UUID NOT NULL;

ALTER TABLE lists
    ADD CONSTRAINT fk_board_public_id
        FOREIGN KEY (board_public_id)
            REFERENCES boards(public_id)
            ON DELETE CASCADE;

ALTER TABLE lists
DROP COLUMN IF EXISTS board_internal_id;