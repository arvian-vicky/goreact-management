-- Menghapus kolom due_date dari tabel boards
ALTER TABLE boards
DROP COLUMN IF EXISTS due_date;

-- Hapus index jika dibuat
DROP INDEX IF EXISTS idx_boards_due_date;