CREATE TABLE IF NOT EXISTS menu_navigasi (
    id SERIAL PRIMARY KEY,
    label VARCHAR(120) NOT NULL,
    url VARCHAR(255) NOT NULL DEFAULT '#',
    deskripsi TEXT NOT NULL DEFAULT '',
    icon VARCHAR(120) NOT NULL DEFAULT 'i-tabler-link',
    color VARCHAR(80) NOT NULL DEFAULT 'text-sky-600',
    grup VARCHAR(120) NOT NULL DEFAULT 'Lainnya',
    urutan INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS pengaduan_kategori (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(120) NOT NULL UNIQUE,
    deskripsi TEXT NOT NULL DEFAULT '',
    urutan INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS pengaduan_status (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(120) NOT NULL UNIQUE,
    deskripsi TEXT NOT NULL DEFAULT '',
    urutan INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

INSERT INTO pengaduan_kategori (nama, deskripsi, urutan)
VALUES
    ('Layanan Publik', 'Pertanyaan atau kendala layanan publik', 1),
    ('Infrastruktur', 'Jalan, fasilitas umum, dan lingkungan', 2),
    ('Administrasi Pemerintahan', 'Administrasi dan pelayanan pemerintahan', 3)
ON CONFLICT (nama) DO NOTHING;

INSERT INTO pengaduan_status (nama, deskripsi, urutan)
VALUES
    ('Baru', 'Aduan baru dan belum diproses', 1),
    ('Diproses', 'Aduan sedang ditindaklanjuti', 2),
    ('Selesai', 'Aduan telah selesai ditindaklanjuti', 3)
ON CONFLICT (nama) DO NOTHING;

INSERT INTO menu_navigasi (label, url, deskripsi, icon, color, grup, urutan)
SELECT defaults.label, defaults.url, defaults.deskripsi, defaults.icon, defaults.color, defaults.grup, defaults.urutan
FROM (VALUES
    ('Semua Berita', '/berita', 'Kumpulan kabar terbaru', 'i-tabler-news', 'text-sky-600', 'Berita', 1),
    ('Kependudukan', '/layanan#kependudukan', 'KTP, KK, Akta', 'i-tabler-users', 'text-sky-600', 'Layanan Publik', 1),
    ('Perizinan', '/layanan#perizinan', 'Izin usaha dan lainnya', 'i-tabler-file-text', 'text-amber-500', 'Layanan Publik', 2),
    ('Kesehatan', '/layanan#kesehatan', 'Layanan RSUD dan Puskesmas', 'i-tabler-heartbeat', 'text-rose-400', 'Layanan Publik', 3),
    ('Pengaduan', '/layanan#pengaduan', 'Sampaikan keluhan warga', 'i-tabler-speakerphone', 'text-emerald-500', 'Layanan Publik', 4),
    ('Darurat', '/layanan#darurat', 'Nomor penting darurat', 'i-tabler-phone-call', 'text-red-600', 'Layanan Publik', 5),
    ('MPP', '/layanan#mpp', 'Mal Pelayanan Publik', 'i-tabler-building-store', 'text-purple-600', 'Layanan Publik', 6),
    ('Profil Daerah', '/profil', 'Sejarah dan visi misi', 'i-tabler-map-2', 'text-sky-600', 'Pemerintahan', 1),
    ('OPD', '/pemerintahan#opd', 'Dinas dan badan', 'i-tabler-building-community', 'text-amber-500', 'Pemerintahan', 2),
    ('Pejabat Daerah', '/profil#pejabat', 'Bupati dan jajarannya', 'i-tabler-user-star', 'text-emerald-500', 'Pemerintahan', 3),
    ('Tautan Eksternal', '/pemerintahan#tautan', 'PPID, JDIH, DPRD', 'i-tabler-external-link', 'text-indigo-600', 'Pemerintahan', 4)
) AS defaults(label, url, deskripsi, icon, color, grup, urutan)
WHERE NOT EXISTS (SELECT 1 FROM menu_navigasi);
