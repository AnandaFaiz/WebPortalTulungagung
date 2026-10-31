import fs from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";
import pg from "pg";
import "dotenv/config";

const { Pool } = pg;
const root = path.resolve(
  path.dirname(fileURLToPath(import.meta.url)),
  "..",
  "backend",
  "migrations",
);
const pool = new Pool({
  host: process.env.DB_HOST,
  port: Number(process.env.DB_PORT),
  user: process.env.DB_USER,
  password: process.env.DB_PASSWORD,
  database: process.env.DB_NAME,
});

async function migrate() {
  const client = await pool.connect();
  try {
    await client.query("BEGIN");
    await client.query(`
      CREATE TABLE IF NOT EXISTS schema_migrations (
        version VARCHAR(255) PRIMARY KEY,
        applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
      )
    `);

    const files = (await fs.readdir(root))
      .filter((file) => /^\d+_.+\.sql$/.test(file))
      .sort();

    for (const file of files) {
      const existing = await client.query(
        "SELECT 1 FROM schema_migrations WHERE version = $1",
        [file],
      );
      if (existing.rowCount > 0) {
        console.log(`skip ${file}`);
        continue;
      }

      if (
        file.startsWith("000_") &&
        (await client.query("SELECT to_regclass('public.users') AS table_name"))
          .rows[0].table_name
      ) {
        await client.query(
          "INSERT INTO schema_migrations (version) VALUES ($1)",
          [file],
        );
        console.log(
          `baseline ${file} sudah tersedia di database, ditandai selesai`,
        );
        continue;
      }

      const sql = await fs.readFile(path.join(root, file), "utf8");
      await client.query(sql);
      await client.query(
        "INSERT INTO schema_migrations (version) VALUES ($1)",
        [file],
      );
      console.log(`applied ${file}`);
    }

    await client.query("COMMIT");
    console.log("Database migration selesai.");
  } catch (error) {
    await client.query("ROLLBACK");
    console.error("Migration gagal:", error.message);
    process.exitCode = 1;
  } finally {
    client.release();
    await pool.end();
  }
}

migrate();
