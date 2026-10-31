# Web Portal Kabupaten Tulungagung

Portal publik dan admin panel Pemerintah Kabupaten Tulungagung.

## Stack

- Frontend: Nuxt 4, Vue 3, Nuxt UI
- Backend: Go, Gin
- Database: PostgreSQL
- Runtime data: PostgreSQL dan folder upload lokal

## Prasyarat

- Node.js 22+
- Go sesuai versi pada `backend/go.mod`
- PostgreSQL 14+
- `npm` dan `psql` tersedia di PATH

## Setup Lokal

1. Salin `.env.example` menjadi `.env`, lalu isi koneksi PostgreSQL dan secret aplikasi.
2. Install dependency dari folder `webportal`:

   ```bash
   npm install
   ```

3. Buat atau siapkan database PostgreSQL sesuai `DB_NAME`.
4. Jalankan migration:

   ```bash
   npm run db:migrate
   ```

   Database baru menjalankan `backend/migrations/000_baseline_schema.sql` lalu migration lanjutan. Database lama yang sudah memiliki tabel `users` melewati baseline dan tetap menjalankan migration baru.

5. Seed permission dan user demo jika diperlukan:

   ```bash
   npm run db:seed:permissions
   $env:SEED_DEFAULT_PASSWORD='ganti-dengan-password-kuat'; npm run db:seed:users
   ```

   Jangan gunakan password demo di production. `SEED_DEFAULT_PASSWORD` wajib diberikan melalui environment dan tidak boleh ditulis di repository.

## Menjalankan Aplikasi

Terminal 1, backend:

```bash
cd backend
go run main.go
```

Terminal 2, frontend:

```bash
npm run dev
```

- Frontend: `http://localhost:3000`
- Backend health check: `http://localhost:8080/api/health`

## Validasi Sebelum Review

```bash
cd backend
go test ./...
cd ..
npm run build
```

Periksa juga:

```bash
git status
git diff --stat
```

## Migration

Migration dijalankan oleh `scripts/migrate.mjs` dan dicatat pada tabel `schema_migrations`. File SQL menggunakan urutan numerik dan tidak boleh diubah setelah dipakai bersama oleh tim. Perubahan schema baru dibuat sebagai file migration berikutnya.

`000_baseline_schema.sql` adalah schema-only export dari PostgreSQL tanpa data aplikasi. `001_navigation_and_complaint_config.sql` membuat dan mengisi konfigurasi menu navigasi serta kategori/status pengaduan. `002_user_profile.sql` menambahkan waktu login terakhir secara aditif untuk halaman Profil Saya; migration ini tidak menghapus atau mengubah data akun yang sudah ada.

Halaman `/admin/profil` menyediakan pengubahan nama/email dan password akun sendiri. Halaman `/admin/pengaturan` menjadi pusat tautan ke identitas situs, kontak, dan menu navigasi. Jalankan `npm run db:migrate` setelah pull sebelum mencoba fitur tersebut.

## Upload

File runtime pada `public/uploads/` dan `backend/public/uploads/` tidak di-commit. Folder tetap dipertahankan dengan `.gitkeep`. Untuk production, gunakan object storage atau volume persisten dan backup terpisah dari source code.

## Keamanan

- Jangan commit `.env`, password, JWT secret, token, atau dump data.
- Ganti `JWT_SECRET` dan password database untuk setiap environment.
- Seed user hanya untuk development/review dan harus memakai password dari environment.
