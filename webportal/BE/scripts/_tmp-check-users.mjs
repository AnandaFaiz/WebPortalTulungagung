import "dotenv/config";
import pg from "pg";

const p = new pg.Pool({
  host: process.env.DB_HOST,
  port: process.env.DB_PORT,
  user: process.env.DB_USER,
  password: process.env.DB_PASSWORD,
  database: process.env.DB_NAME,
});

try {
  const c = await p.query(
    "select column_name, data_type from information_schema.columns where table_name='users' and column_name in ('created_at','last_login_at')",
  );
  console.log(c.rows);
  const r = await p.query(
    "select id, email, role, created_at::text as created_at, last_login_at::text as last_login_at from users order by id",
  );
  console.table(r.rows);
  console.log((await p.query("show timezone")).rows);
} catch (e) {
  console.error("ERR:", e.message);
} finally {
  await p.end();
}
