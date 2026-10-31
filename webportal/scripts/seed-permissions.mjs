import pg from "pg";
import "dotenv/config";

const { Pool } = pg;

const pool = new Pool({
  host: process.env.DB_HOST,
  port: Number(process.env.DB_PORT),
  user: process.env.DB_USER,
  password: process.env.DB_PASSWORD,
  database: process.env.DB_NAME,
});

// ==== Salinan dari permissions.ts (nilai default) ====
const defaults = {
  super_admin: {
    dashboard: true,
    "cms.halaman": true,
    "cms.berita": true,
    "cms.pengumuman": true,
    "cms.galeri": true,
    "cms.sambutan": true,
    "cms.banner": true,
    "pemerintahan.profil": true,
    "pemerintahan.struktur": true,
    "pemerintahan.pejabat": true,
    "pemerintahan.opd": true,
    "pemerintahan.dprd": true,
    "pemerintahan.ppid": true,
    "pemerintahan.jdih": true,
    "layanan.shortcut": true,
    "layanan.kategori": true,
    "layanan.daftar": true,
    "layanan.mpp": true,
    "layanan.darurat": true,
    "interaksi.medsos": true,
    "interaksi.tiket": true,
    "interaksi.kategori": true,
    "potensi.pariwisata": true,
    "potensi.umkm": true,
    "potensi.investasi": true,
    aplikasi: true,
    media: true,
    pengaturan: true,
    "pengguna.daftar": true,
    "pengguna.role": true,
    "pengguna.log": true,
  },
  admin_opd: {
    dashboard: true,
    "cms.halaman": false,
    "cms.berita": true,
    "cms.pengumuman": true,
    "cms.galeri": true,
    "cms.sambutan": false,
    "cms.banner": false,
    "pemerintahan.profil": false,
    "pemerintahan.struktur": false,
    "pemerintahan.pejabat": false,
    "pemerintahan.opd": true,
    "pemerintahan.dprd": false,
    "pemerintahan.ppid": false,
    "pemerintahan.jdih": false,
    "layanan.shortcut": false,
    "layanan.kategori": false,
    "layanan.daftar": true,
    "layanan.mpp": false,
    "layanan.darurat": false,
    "interaksi.medsos": false,
    "interaksi.tiket": true,
    "interaksi.kategori": false,
    "potensi.pariwisata": true,
    "potensi.umkm": true,
    "potensi.investasi": true,
    aplikasi: false,
    media: true,
    pengaturan: false,
    "pengguna.daftar": false,
    "pengguna.role": false,
    "pengguna.log": false,
  },
  editor: {
    dashboard: true,
    "cms.halaman": true,
    "cms.berita": true,
    "cms.pengumuman": true,
    "cms.galeri": true,
    "cms.sambutan": true,
    "cms.banner": true,
    "pemerintahan.profil": false,
    "pemerintahan.struktur": false,
    "pemerintahan.pejabat": false,
    "pemerintahan.opd": false,
    "pemerintahan.dprd": false,
    "pemerintahan.ppid": false,
    "pemerintahan.jdih": false,
    "layanan.shortcut": false,
    "layanan.kategori": false,
    "layanan.daftar": false,
    "layanan.mpp": false,
    "layanan.darurat": false,
    "interaksi.medsos": false,
    "interaksi.tiket": false,
    "interaksi.kategori": false,
    "potensi.pariwisata": true,
    "potensi.umkm": true,
    "potensi.investasi": true,
    aplikasi: false,
    media: true,
    pengaturan: false,
    "pengguna.daftar": false,
    "pengguna.role": false,
    "pengguna.log": false,
  },
};

const sidebarPermissionKeys = [
  "dashboard",
  "cms.halaman",
  "cms.berita",
  "cms.kategori",
  "cms.pengumuman",
  "cms.galeri",
  "cms.sambutan",
  "cms.banner",
  "cms.event",
  "cms.layanan",
  "cms.pencarian",
  "pemerintahan.profil",
  "pemerintahan.struktur",
  "pemerintahan.pejabat",
  "pemerintahan.opd",
  "pemerintahan.tautan",
  "layanan.darurat",
  "layanan.kategori",
  "layanan.daftar",
  "layanan.mpp",
  "interaksi.medsos",
  "interaksi.tiket",
  "interaksi.kategori",
  "aplikasi",
  "media",
  "pengaturan",
  "pengguna.daftar",
  "pengguna.role",
  "pengguna.log",
];

defaults.super_admin["cms.kategori"] = true;
defaults.super_admin["cms.event"] = true;
defaults.super_admin["cms.layanan"] = true;
defaults.super_admin["cms.pencarian"] = true;
defaults.super_admin["pemerintahan.tautan"] = true;
defaults.admin_opd["cms.kategori"] = false;
defaults.admin_opd["cms.event"] = false;
defaults.admin_opd["cms.layanan"] = false;
defaults.admin_opd["cms.pencarian"] = false;
defaults.admin_opd["pemerintahan.tautan"] = false;
defaults.editor["cms.kategori"] = true;
defaults.editor["cms.event"] = true;
defaults.editor["cms.layanan"] = true;
defaults.editor["cms.pencarian"] = true;
defaults.editor["pemerintahan.tautan"] = false;

async function seed() {
  try {
    console.log("🌱 Seeding permissions...\n");

    let count = 0;
    for (const [role, perms] of Object.entries(defaults)) {
      for (const [key, allowed] of Object.entries(perms)) {
        if (!sidebarPermissionKeys.includes(key)) continue;
        await pool.query(
          `INSERT INTO role_permissions (role, permission_key, allowed)
           VALUES ($1, $2, $3)
           ON CONFLICT (role, permission_key) 
           DO UPDATE SET allowed = EXCLUDED.allowed, updated_at = NOW()`,
          [role, key, allowed],
        );
        count++;
      }
      await pool.query(
        `DELETE FROM role_permissions WHERE role = $1 AND permission_key <> ALL($2::text[])`,
        [role, sidebarPermissionKeys],
      );
    }

    console.log(`✅ Berhasil seed ${count} entri permission`);

    const result = await pool.query(
      "SELECT role, COUNT(*) as total FROM role_permissions GROUP BY role",
    );
    console.log("\n📋 Total per role:");
    console.table(result.rows);

    await pool.end();
    process.exit(0);
  } catch (err) {
    console.error("\n❌ Error:", err.message);
    await pool.end();
    process.exit(1);
  }
}

seed();
