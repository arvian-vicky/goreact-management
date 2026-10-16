-- Aktifkan ekstensi pgcrypto jika belum ada
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Tambahkan ini sebelum migrasi boards
ALTER TABLE users RENAME COLUMN id TO internal_id;
-- Baru kemudian alter table
ALTER TABLE users
    ADD COLUMN public_id UUID NOT NULL DEFAULT gen_random_uuid(),
ADD CONSTRAINT users_public_id_unique UNIQUE (public_id);