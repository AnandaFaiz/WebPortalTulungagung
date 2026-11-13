# Backend (BE) — Portal Kabupaten Tulungagung

REST API berbasis **Go + Gin** dan tooling PostgreSQL untuk Portal Kabupaten Tulungagung. Folder ini adalah root repository backend.

## Tech Stack

| Komponen | Versi / Keterangan |
|---|---|
| Go | 1.27.1 (lihat `go.mod`) |
| Framework | Gin 1.12.0 |
| Database | PostgreSQL 14+ (`lib/pq`) |
| Auth | JWT (`golang-jwt/jwt/v5`), password di-hash bcrypt |
| Email | `gomail.v2` (opsional, untuk notifikasi pengaduan) |
| Tooling migration/seed | Node.js 22+ (`pg`, `bcryptjs`, `dotenv`) |

## Prasyarat

- Go sesuai versi pada `go.mod`
- Node.js 22+ dan npm untuk migration/seed tooling
- PostgreSQL 14+ yang berjalan dan bisa diakses

## Setup

### 1. Environment

Salin `.env.example` menjadi `.env`, lalu sesuaikan koneksi database dan secret. Di PowerShell:

```powershell
Copy-Item .env.example .env
```

Variabel yang digunakan:

| Variabel | Wajib | Keterangan |
|---|---|---|
| `DB_HOST` | ya | Host PostgreSQL, mis. `localhost` |
| `DB_PORT` | ya | Port PostgreSQL, mis. `5432` |
| `DB_USER` | ya | User database |
| `DB_PASSWORD` | ya | Password database |
| `DB_NAME` | ya | Nama database, mis. `pemkab_tulungagung` |
| `JWT_SECRET` | ya | Secret token JWT — **gunakan nilai acak panjang, berbeda per environment** |
| `PORT` | tidak | Port server, default `8080` |
| `SMTP_HOST` `SMTP_PORT` `SMTP_USER` `SMTP_PASS` | tidak | Jika kosong, pengiriman email nonaktif |
| `SEED_DEFAULT_PASSWORD` | hanya untuk seed | Password akun demo; jangan disimpan di repository |

Buat database terlebih dahulu (sekali saja):

```sql
CREATE DATABASE pemkab_tulungagung;
```

### 2. Dependency tooling Node

```bash
npm install
```

### 3. Migration database

```bash
npm run db:migrate
```

Migration baru diterapkan berurutan dan dicatat pada tabel `schema_migrations`. **Jangan mengubah file migration yang sudah dipakai bersama** — tambahkan file SQL bernomor berikutnya (`006_...`). Database baru menjalankan baseline `000_`, sedangkan database lama yang sudah punya tabel `users` melewati baseline.

Daftar migration saat ini:

| File | Isi |
|---|---|
| `000_baseline_schema.sql` | Skema awal seluruh tabel portal |
| `001_navigation_and_complaint_config.sql` | Menu navigasi & konfigurasi pengaduan |
| `002_user_profile.sql` | Profil pengguna |
| `003_publish_dates_and_announcement_image.sql` | Tanggal publikasi & gambar pengumuman |
| `004_per_user_crud_permissions.sql` | Permission CRUD per user |
| `005_pengaturan_permissions.sql` | Permission untuk modul pengaturan (identitas, kontak, menu) |

### 4. Seed data (opsional, untuk development/review)

```bash
# Permission default untuk user super_admin & editor
npm run db:seed:permissions

# Akun demo — wajib set SEED_DEFAULT_PASSWORD di .env lebih dulu
npm run db:seed:users
```

Akun demo yang dibuat:

| Email | Role |
|---|---|
| `superadmin@tulungagung.go.id` | `super_admin` |
| `editor@tulungagung.go.id` | `editor` |

Password mengikuti nilai `SEED_DEFAULT_PASSWORD`. **Jangan gunakan akun demo di production.**

### 5. Jalankan server

```bash
go run .
```

API berjalan di `http://localhost:8080`; health check tersedia di `GET /api/health`. Frontend dev server (`http://localhost:3000`) diizinkan oleh konfigurasi CORS.

## Struktur Folder

```
BE/
├── main.go              # Entry point, CORS, static /uploads, health check
├── config/              # Koneksi database
├── routes/              # Definisi seluruh rute API
├── middleware/          # Auth JWT & permission CRUD per rute admin
├── controllers/         # Handler per domain (berita, layanan, pengguna, dst.)
├── models/              # Struct & query data
├── migrations/          # File SQL migration bernomor
├── scripts/             # Tooling Node (migrate, seed) + script setup Go
└── public/uploads/      # Upload runtime (tidak di-commit)
```

## Autentikasi & Permission

Rute admin dilindungi dua lapis: **JWT** (memastikan sudah login) lalu **permission CRUD** (`middleware/permissions.go`).

Setiap rute admin dipetakan ke sebuah `permission_key`, dan hak akses dicek per method HTTP:

| Method | Kolom yang dicek |
|---|---|
| `GET` | `can_read` |
| `POST` | `can_create` |
| `PUT` / `PATCH` | `can_update` |
| `DELETE` | `can_delete` |

Aturan khusus:

- Role `super_admin` **melewati** pengecekan permission (akses penuh).
- Rute `/api/admin/users`, `/api/admin/permissions`, dan `/api/admin/activity-logs` **hanya** untuk `super_admin`.
- Menu admin di frontend otomatis menyesuaikan permission yang dimiliki user.

Daftar `permission_key` yang dipakai:

`cms.berita`, `cms.kategori`, `cms.pengumuman`, `cms.galeri`, `cms.banner`, `cms.event`, `cms.layanan`, `cms.pencarian`, `cms.halaman`, `pemerintahan.profil`, `pemerintahan.struktur`, `pemerintahan.pejabat`, `pemerintahan.opd`, `pemerintahan.tautan`, `layanan.kategori`, `layanan.darurat`, `layanan.daftar`, `layanan.mpp`, `interaksi.kategori`, `interaksi.tiket`, `interaksi.medsos`, `aplikasi`, `media`, `pengaturan.identitas`, `pengaturan.kontak`, `pengaturan.menu`.

## Modul API

Seluruh rute berada di bawah prefix `/api`:

- `auth` — login, sesi, profil sendiri
- `berita`, `kategori-berita`, `pengumuman`, `agenda`, `galeri`, `banner`, `halaman` — konten publik & CMS
- `profil-daerah`, `struktur-organisasi`, `pejabat`, `opd`, `tautan-eksternal` — data pemerintahan
- `kategori-layanan`, `layanan`, `layanan-darurat`, `layanan-terpopuler`, `mpp` — layanan publik
- `pengaduan`, `pengaduan-kategori`, `pengaduan-status`, `notifications` — pengaduan masyarakat
- `aplikasi`, `media`, `upload`, `identitas`, `kontak`, `media-sosial`, `menu-navigasi` — pendukung situs
- `users`, `permissions`, `activity-logs` — administrasi (khusus super admin)
- `search`, `pencarian-populer`, `dashboard` — pencarian & ringkasan statistik
- `health` — health check

## Validasi

```bash
go run .     # jalankan server
go test ./...  # jalankan test
```

## Upload & Keamanan

- Upload runtime disimpan di `public/uploads/` dan **tidak di-commit** (lihat `.gitignore`).
- Gunakan object storage atau volume persisten untuk production — filesystem lokal tidak cocok untuk deployment multi-instance.
- Jangan commit `.env`, password, JWT secret, token, atau dump data.
- Gunakan secret dan credential berbeda untuk setiap environment.
- Konfigurasi CORS saat ini mengizinkan semua origin (`AllowOriginFunc` selalu `true`) demi kemudahan development. **Batasi ke domain yang diketahui sebelum production.**
