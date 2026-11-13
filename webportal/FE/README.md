# Frontend (FE) — Portal Kabupaten Tulungagung

Aplikasi **Nuxt 4** untuk portal publik dan panel admin Pemerintah Kabupaten Tulungagung. Folder ini adalah root repository frontend.

## Tech Stack

| Komponen | Versi / Keterangan |
|---|---|
| Nuxt | 4.6 |
| Vue | 3.5 |
| UI | `@nuxt/ui` 4 (Tailwind CSS 4) |
| Font | `@nuxt/fonts` — Plus Jakarta Sans |
| Ikon | Iconify (`mdi`, `tabler`) |
| Editor teks | Tiptap 3 (rich text berita/halaman) |
| Auth | `jose` (verifikasi token sisi klien) |
| Lain-lain | `cropperjs`, `vue-draggable-plus`, `vue-router` |

## Prasyarat

- Node.js 22+ dan npm
- Backend API berjalan di `http://localhost:8080` (lihat `../BE/README.md`)

## Setup

### 1. Environment

Salin `.env.example` menjadi `.env`:

```powershell
Copy-Item .env.example .env
```

| Variabel | Default | Keterangan |
|---|---|---|
| `NUXT_PUBLIC_API_BASE` | `http://localhost:8080/api` | Base URL API backend |

Jika backend berjalan di host/port lain, cukup ubah nilai di `.env` — tidak perlu menyentuh source code.

### 2. Install & jalankan

```bash
npm install
npm run dev
```

Buka `http://localhost:3000`.

> `npm install` menjalankan `nuxt prepare` (postinstall). Jalankan backend lebih dulu agar halaman yang butuh data bisa dimuat.

### 3. Build

```bash
npm run build     # build production
npm run preview   # pratinjau hasil build
npm run generate  # static generation
```

## Struktur Folder

```
FE/
├── nuxt.config.ts           # Konfigurasi Nuxt, runtimeConfig.apiBase, tema
├── app/
│   ├── app.vue
│   ├── assets/css/main.css  # Tailwind & style global
│   ├── components/          # Komponen bersama (Admin/*, HeroSlider, AutoIcon, …)
│   ├── composables/useAuth.ts
│   ├── layouts/             # default (publik) & admin
│   ├── middleware/          # auth.global.ts — guard rute admin
│   ├── plugins/             # auth-check.client.ts
│   ├── utils/               # admin-navigation.ts, permissions.ts
│   └── pages/
│       ├── index.vue, profil.vue, pemerintahan.vue,
│       │   layanan.vue, berita/, kontak.vue, pencarian.vue, login.vue
│       └── admin/           # panel admin (lihat tabel di bawah)
└── public/                  # Gambar hero, logo, favicon, uploads
```

## Halaman Publik

| Rute | Isi |
|---|---|
| `/` | Beranda: hero slider, berita terbaru, layanan populer |
| `/profil` | Profil daerah, sambutan, struktur organisasi |
| `/pemerintahan` | Pejabat, OPD, tautan eksternal |
| `/layanan` | Kategori layanan, layanan darurat, MPP |
| `/berita` dan `/berita/:id` | Daftar & detail berita |
| `/kontak` | Kontak, media sosial, form pengaduan |
| `/pencarian` | Pencarian global |
| `/login` | Login admin |

## Panel Admin

Seluruh rute `/admin/*` dijaga oleh `middleware/auth.global.ts` — user yang belum login diarahkan ke `/login`.

| Grup | Halaman |
|---|---|
| Dashboard | `/admin` |
| CMS | `cms/berita`, `cms/kategori-berita`, `cms/pengumuman`, `cms/event`, `cms/galeri`, `cms/banner`, `cms/halaman`, `cms/layanan-terpopuler`, `cms/pencarian-populer`, `cms/sambutan` |
| Pemerintahan | `pemerintahan/profil`, `pemerintahan/struktur`, `pemerintahan/pejabat`, `pemerintahan/opd`, `pemerintahan/tautan` |
| Layanan | `layanan/daftar`, `layanan/kategori`, `layanan/darurat`, `layanan/mpp` |
| Interaksi | `interaksi/tiket` (pengaduan), `interaksi/kategori`, `interaksi/medsos` |
| Pengaturan | `pengaturan/identitas`, `pengaturan/kontak`, `pengaturan/menu`, `pengaturan` |
| Pengguna | `pengguna/daftar`, `pengguna/role`, `pengguna/log` — **khusus super admin** |
| Lain-lain | `media`, `aplikasi`, `profil` |

### Hak akses di UI

Menu admin dibangun dinamis dari `app/utils/admin-navigation.ts` dan difilter oleh `app/utils/permissions.ts` berdasarkan `permission_key` yang dimiliki user (contoh: `cms.berita`, `layanan.daftar`, `pengaturan.menu`).

- `super_admin` melihat semua menu.
- `editor` melihat menu sesuai permission yang diberikan.
- Backend tetap menjadi penentu akhir — UI hanya menyembunyikan menu, otorisasi sebenarnya ada di middleware API.

## Login untuk Mencoba

Setelah backend di-seed (`npm run db:seed:users` di folder `BE`):

| Email | Role |
|---|---|
| `superadmin@tulungagung.go.id` | `super_admin` |
| `editor@tulungagung.go.id` | `editor` |

Password mengikuti `SEED_DEFAULT_PASSWORD` di `BE/.env`.

## Catatan Pengembangan

- `runtimeConfig.public.apiBase` dibaca dari `NUXT_PUBLIC_API_BASE`; jangan hardcode URL API di komponen.
- Gambar hero/statis berada di `public/hero/`; pakai `object-fit: cover` dengan aspect ratio eksplisit.
- Upload runtime (`public/uploads/`) tidak di-commit.
- Jangan commit `.env`, credential, token, atau data pengguna.
