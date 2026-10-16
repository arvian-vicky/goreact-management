-- Menambahkan kolom due_date ke tabel boards
ALTER TABLE boards
    ADD COLUMN due_date TIMESTAMP WITH TIME ZONE;

-- Optional: Tambahkan index jika diperlukan
CREATE INDEX idx_boards_due_date ON boards(due_date);