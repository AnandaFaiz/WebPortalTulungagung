import { adminNavigation } from "./admin-navigation";

// ==== Mapping path URL → permission ====
// Dipakai middleware untuk cek akses URL langsung
// Catatan: data permission (siapa boleh apa) sekarang tersimpan di PostgreSQL
// File ini hanya mapping URL → permission key, yang sifatnya statis.

const legacyPathPermissions: Record<string, string> = {
  "/admin": "dashboard",
  "/admin/cms/halaman": "cms.halaman",
  "/admin/cms/berita": "cms.berita",
  "/admin/cms/pengumuman": "cms.pengumuman",
  "/admin/cms/galeri": "cms.galeri",
  "/admin/cms/sambutan": "cms.sambutan",
  "/admin/cms/banner": "cms.banner",
  "/admin/pemerintahan/profil": "pemerintahan.profil",
  "/admin/pemerintahan/struktur": "pemerintahan.struktur",
  "/admin/pemerintahan/pejabat": "pemerintahan.pejabat",
  "/admin/pemerintahan/opd": "pemerintahan.opd",
  "/admin/pemerintahan/tautan": "pemerintahan.tautan",
  "/admin/layanan/shortcut": "layanan.shortcut",
  "/admin/layanan/kategori": "layanan.kategori",
  "/admin/layanan/daftar": "layanan.daftar",
  "/admin/layanan/darurat": "layanan.darurat",
  "/admin/layanan/mpp": "layanan.mpp",
  "/admin/interaksi/medsos": "interaksi.medsos",
  "/admin/interaksi/tiket": "interaksi.tiket",
  "/admin/interaksi/kategori": "interaksi.kategori",
  "/admin/potensi/pariwisata": "potensi.pariwisata",
  "/admin/potensi/umkm": "potensi.umkm",
  "/admin/potensi/investasi": "potensi.investasi",
  "/admin/aplikasi": "aplikasi",
  "/admin/media": "media",
  "/admin/pengaturan": "pengaturan",
  "/admin/pengaturan/identitas": "pengaturan",
  "/admin/pengaturan/menu": "pengaturan",
  "/admin/pengaturan/kontak": "pengaturan",
  "/admin/pengguna/daftar": "pengguna.daftar",
  "/admin/pengguna/role": "pengguna.role",
  "/admin/pengguna/log": "pengguna.log",
};

export const pathPermissions: Record<string, string> = {
  ...legacyPathPermissions,
  ...Object.fromEntries(
    adminNavigation.flatMap((entry) =>
      "children" in entry
        ? entry.children.map((child) => [child.to, child.permission])
        : [[entry.to, entry.permission]],
    ),
  ),
};

// ==== Helper: dapatkan permission yang dibutuhkan untuk suatu path ====
export function getRequiredPermission(path: string): string | null {
  if (pathPermissions[path]) return pathPermissions[path];

  // Fallback: cari parent path
  // Contoh: /admin/cms/berita/edit → cek /admin/cms/berita/edit dulu,
  // lalu /admin/cms/berita, lalu /admin/cms, dst.
  const segments = path.split("/").filter(Boolean);
  while (segments.length > 0) {
    const parentPath = "/" + segments.join("/");
    if (pathPermissions[parentPath]) return pathPermissions[parentPath];
    segments.pop();
  }
  return null;
}
