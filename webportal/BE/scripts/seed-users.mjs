import pg from "pg";
import bcrypt from "bcryptjs";
import "dotenv/config";

const { Pool } = pg;

const defaultPassword = process.env.SEED_DEFAULT_PASSWORD;
if (!defaultPassword) {
  throw new Error(
    "SEED_DEFAULT_PASSWORD wajib diisi dan tidak boleh disimpan di repository.",
  );
}

const pool = new Pool({
  host: process.env.DB_HOST,
  port: Number(process.env.DB_PORT),
  user: process.env.DB_USER,
  password: process.env.DB_PASSWORD,
  database: process.env.DB_NAME,
});

// ==== User Default ====
const users = [
  {
    nama: "Super Administrator",
    email: "superadmin@tulungagung.go.id",
    password: defaultPassword,
    role: "super_admin",
    opd: null,
  },
  {
    nama: "Editor Portal",
    email: "editor@tulungagung.go.id",
    password: defaultPassword,
    role: "editor",
    opd: null,
  },
];

async function seed() {
  try {
    console.log("🌱 Memulai seeding user...\n");

    for (const user of users) {
      // Cek apakah email sudah ada
      const existing = await pool.query(
        "SELECT id FROM users WHERE email = $1",
        [user.email],
      );

      if (existing.rows.length > 0) {
        console.log(`⏭️  Skip: ${user.email} (sudah ada di database)`);
        continue;
      }

      // Hash password dengan bcrypt (10 rounds)
      const hashedPassword = await bcrypt.hash(user.password, 10);

      // Insert ke database
      await pool.query(
        "INSERT INTO users (nama, email, password, role, opd) VALUES ($1, $2, $3, $4, $5)",
        [user.nama, user.email, hashedPassword, user.role, user.opd],
      );

      console.log(`✅ Dibuat: ${user.email} (${user.role})`);
    }

    // Tampilkan semua user
    const all = await pool.query(
      "SELECT id, nama, email, role, opd, is_active FROM users ORDER BY id",
    );

    console.log("\n📋 Daftar user di database:");
    console.table(all.rows);

    console.log("\n🎉 Seeding selesai!");
    console.log(
      "\n📝 User demo dibuat menggunakan SEED_DEFAULT_PASSWORD dari environment.",
    );
    console.log("┌─────────────┬────────────────────────────────┐");
    console.log("│ Role        │ Email                          │");
    console.log("├─────────────┼────────────────────────────────┤");
    console.log("│ Super Admin │ superadmin@tulungagung.go.id   │");
    console.log("│ Editor      │ editor@tulungagung.go.id       │");
    console.log("└─────────────┴────────────────────────────────┘");

    await pool.end();
    process.exit(0);
  } catch (err) {
    console.error("\n❌ Error saat seeding:", err.message);
    await pool.end();
    process.exit(1);
  }
}

seed();
