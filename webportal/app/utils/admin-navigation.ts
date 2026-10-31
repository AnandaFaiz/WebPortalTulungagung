export interface AdminNavigationItem {
  label: string;
  icon: string;
  to: string;
  permission: string;
}

export interface AdminNavigationGroup {
  label: string;
  icon: string;
  permission: string;
  children: AdminNavigationItem[];
}

export type AdminNavigationEntry = AdminNavigationGroup | AdminNavigationItem;

export const adminNavigation: AdminNavigationEntry[] = [
  {
    label: "Dashboard",
    icon: "i-tabler-layout-dashboard",
    to: "/admin",
    permission: "dashboard",
  },
  {
    label: "Konten (CMS)",
    icon: "i-tabler-file-text",
    permission: "cms",
    children: [
      {
        label: "Halaman Statis",
        icon: "i-tabler-file-description",
        to: "/admin/cms/halaman",
        permission: "cms.halaman",
      },
      {
        label: "Berita",
        icon: "i-tabler-news",
        to: "/admin/cms/berita",
        permission: "cms.berita",
      },
      {
        label: "Kategori Berita",
        icon: "i-tabler-category",
        to: "/admin/cms/kategori-berita",
        permission: "cms.kategori",
      },
      {
        label: "Pengumuman",
        icon: "i-tabler-speakerphone",
        to: "/admin/cms/pengumuman",
        permission: "cms.pengumuman",
      },
      {
        label: "Galeri",
        icon: "i-tabler-photo",
        to: "/admin/cms/galeri",
        permission: "cms.galeri",
      },
      {
        label: "Sambutan Bupati",
        icon: "i-tabler-quote",
        to: "/admin/cms/sambutan",
        permission: "cms.sambutan",
      },
      {
        label: "Wallpaper Hero",
        icon: "i-tabler-photo",
        to: "/admin/cms/banner",
        permission: "cms.banner",
      },
      {
        label: "Event/Agenda",
        icon: "i-tabler-calendar-event",
        to: "/admin/cms/event",
        permission: "cms.event",
      },
      {
        label: "Layanan Populer",
        icon: "i-tabler-star",
        to: "/admin/cms/layanan-terpopuler",
        permission: "cms.layanan",
      },
      {
        label: "Pencarian Populer",
        icon: "i-tabler-search",
        to: "/admin/cms/pencarian-populer",
        permission: "cms.pencarian",
      },
    ],
  },
  {
    label: "Profil & Pemerintahan",
    icon: "i-tabler-building-bank",
    permission: "pemerintahan",
    children: [
      {
        label: "Profil Daerah",
        icon: "i-tabler-map-2",
        to: "/admin/pemerintahan/profil",
        permission: "pemerintahan.profil",
      },
      {
        label: "Struktur Organisasi",
        icon: "i-tabler-sitemap",
        to: "/admin/pemerintahan/struktur",
        permission: "pemerintahan.struktur",
      },
      {
        label: "Profil Pejabat",
        icon: "i-tabler-user-star",
        to: "/admin/pemerintahan/pejabat",
        permission: "pemerintahan.pejabat",
      },
      {
        label: "OPD",
        icon: "i-tabler-building-community",
        to: "/admin/pemerintahan/opd",
        permission: "pemerintahan.opd",
      },
      {
        label: "Tautan Eksternal",
        icon: "i-tabler-external-link",
        to: "/admin/pemerintahan/tautan",
        permission: "pemerintahan.tautan",
      },
    ],
  },
  {
    label: "Layanan Publik",
    icon: "i-tabler-briefcase",
    permission: "layanan",
    children: [
      {
        label: "Layanan Darurat",
        icon: "i-tabler-alert-circle",
        to: "/admin/layanan/darurat",
        permission: "layanan.darurat",
      },
      {
        label: "Kategori Layanan",
        icon: "i-tabler-category",
        to: "/admin/layanan/kategori",
        permission: "layanan.kategori",
      },
      {
        label: "Daftar Layanan",
        icon: "i-tabler-list-details",
        to: "/admin/layanan/daftar",
        permission: "layanan.daftar",
      },
      {
        label: "Informasi MPP",
        icon: "i-tabler-building-store",
        to: "/admin/layanan/mpp",
        permission: "layanan.mpp",
      },
    ],
  },
  {
    label: "Interaksi",
    icon: "i-tabler-message-2",
    permission: "interaksi",
    children: [
      {
        label: "Media Sosial",
        icon: "i-tabler-share-3",
        to: "/admin/interaksi/medsos",
        permission: "interaksi.medsos",
      },
      {
        label: "Tiket Pengaduan",
        icon: "i-tabler-ticket",
        to: "/admin/interaksi/tiket",
        permission: "interaksi.tiket",
      },
      {
        label: "Kategori & Status",
        icon: "i-tabler-tag",
        to: "/admin/interaksi/kategori",
        permission: "interaksi.kategori",
      },
    ],
  },
  {
    label: "Aplikasi",
    icon: "i-tabler-apps",
    to: "/admin/aplikasi",
    permission: "aplikasi",
  },
  {
    label: "Media",
    icon: "i-tabler-photo",
    to: "/admin/media",
    permission: "media",
  },
  {
    label: "Pengguna & Hak Akses",
    icon: "i-tabler-users",
    permission: "pengguna",
    children: [
      {
        label: "Daftar Pengguna",
        icon: "i-tabler-users-group",
        to: "/admin/pengguna/daftar",
        permission: "pengguna.daftar",
      },
      {
        label: "Role",
        icon: "i-tabler-shield-lock",
        to: "/admin/pengguna/role",
        permission: "pengguna.role",
      },
      {
        label: "Log Aktivitas",
        icon: "i-tabler-history",
        to: "/admin/pengguna/log",
        permission: "pengguna.log",
      },
    ],
  },
];

export const adminPermissionGroups = adminNavigation.reduce(
  (groups, entry) => {
    if ("children" in entry) {
      groups.push({
        label: entry.label,
        items: entry.children.map((child) => ({
          label: child.label,
          key: child.permission,
          to: child.to,
        })),
      });
    } else {
      groups.push({
        label: entry.label,
        items: [{ label: entry.label, key: entry.permission, to: entry.to }],
      });
    }
    return groups;
  },
  [] as Array<{
    label: string;
    items: Array<{ label: string; key: string; to?: string }>;
  }>,
);

export const adminPermissionKeys = adminPermissionGroups.flatMap((group) =>
  group.items.map((item) => item.key),
);

adminPermissionGroups.push({
  label: "Pengaturan",
  items: [
    { label: "Pengaturan Portal", key: "pengaturan", to: "/admin/pengaturan" },
  ],
});

adminPermissionKeys.push("pengaturan");
