ALTER TABLE pengumuman ADD COLUMN IF NOT EXISTS gambar character varying(255);

UPDATE berita
SET tanggal = NULL
WHERE status = 'Draft';