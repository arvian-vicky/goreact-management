ALTER TABLE boards
    ADD COLUMN owner_public_id UUID NOT NULL;

ALTER TABLE boards
    ADD CONSTRAINT fk_boards_owner_public_id
        FOREIGN KEY (owner_public_id)
            REFERENCES users(public_id)
            ON DELETE CASCADE;