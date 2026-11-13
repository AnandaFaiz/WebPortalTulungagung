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

const resources = [
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
  "pengaturan.identitas",
  "pengaturan.kontak",
  "pengaturan.menu",
];

async function seed() {
  try {
    const users = await pool.query(
      "SELECT id, role FROM users WHERE role IN ('super_admin', 'editor')",
    );
    for (const user of users.rows) {
      for (const resource of resources) {
        const canManage = user.role === "super_admin" || user.role === "editor";
        await pool.query(
          `INSERT INTO user_permissions (user_id, permission_key, can_create, can_read, can_update, can_delete)
           VALUES ($1, $2, $3, TRUE, $3, $4)
           ON CONFLICT (user_id, permission_key) DO NOTHING`,
          [user.id, resource, canManage, user.role === "super_admin"],
        );
      }
    }
    await pool.query("DELETE FROM role_permissions");
    console.log(`Seed permission selesai untuk ${users.rowCount} admin.`);
  } catch (error) {
    console.error("Seed permission gagal:", error.message);
    process.exitCode = 1;
  } finally {
    await pool.end();
  }
}

seed();
