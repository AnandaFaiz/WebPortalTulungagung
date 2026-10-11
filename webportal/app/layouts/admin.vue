<script setup>
const user = {
    name: 'Admin Pemkab',
    email: 'admin@tulungagung.go.id',
    avatar: ''
}

const sidebarLinks = [
    {
        label: 'Dashboard',
        icon: 'i-tabler-layout-dashboard',
        to: '/admin'
    },
    {
        label: 'Konten (CMS)',
        icon: 'i-tabler-file-text',
        children: [
            { label: 'Halaman Statis', icon: 'i-tabler-file-description', to: '/admin/cms/halaman' },
            { label: 'Berita', icon: 'i-tabler-news', to: '/admin/cms/berita' },
            { label: 'Pengumuman', icon: 'i-tabler-speakerphone', to: '/admin/cms/pengumuman' },
            { label: 'Galeri', icon: 'i-tabler-photo', to: '/admin/cms/galeri' },
            { label: 'Sambutan Bupati', icon: 'i-tabler-quote', to: '/admin/cms/sambutan' },
            { label: 'Banner', icon: 'i-tabler-presentation', to: '/admin/cms/banner' }
        ]
    },
    {
        label: 'Profil & Pemerintahan',
        icon: 'i-tabler-building-bank',
        children: [
            { label: 'Profil Daerah', icon: 'i-tabler-map-2', to: '/admin/pemerintahan/profil' },
            { label: 'Struktur Organisasi', icon: 'i-tabler-sitemap', to: '/admin/pemerintahan/struktur' },
            { label: 'Profil Pejabat', icon: 'i-tabler-user-star', to: '/admin/pemerintahan/pejabat' },
            { label: 'OPD', icon: 'i-tabler-building-community', to: '/admin/pemerintahan/opd' },
            { label: 'DPRD', icon: 'i-tabler-gavel', to: '/admin/pemerintahan/dprd' },
            { label: 'PPID', icon: 'i-tabler-file-search', to: '/admin/pemerintahan/ppid' },
            { label: 'JDIH', icon: 'i-tabler-scale', to: '/admin/pemerintahan/jdih' }
        ]
    },
    {
        label: 'Layanan Publik',
        icon: 'i-tabler-briefcase',
        children: [
            { label: 'Shortcut Layanan', icon: 'i-tabler-bolt', to: '/admin/layanan/shortcut' },
            { label: 'Kategori Layanan', icon: 'i-tabler-category', to: '/admin/layanan/kategori' },
            { label: 'Daftar Layanan', icon: 'i-tabler-list-details', to: '/admin/layanan/daftar' },
            { label: 'Informasi MPP', icon: 'i-tabler-building-store', to: '/admin/layanan/mpp' },
            { label: 'Layanan Darurat', icon: 'i-tabler-alert-circle', to: '/admin/layanan/darurat' }
        ]
    },
    {
        label: 'Interaksi',
        icon: 'i-tabler-message-2',
        children: [
            { label: 'Media Sosial', icon: 'i-tabler-share-3', to: '/admin/interaksi/medsos' },
            { label: 'Tiket Pengaduan', icon: 'i-tabler-ticket', to: '/admin/interaksi/tiket' },
            { label: 'Kategori & Status', icon: 'i-tabler-tag', to: '/admin/interaksi/kategori' }
        ]
    },
    {
        label: 'Potensi Daerah',
        icon: 'i-tabler-mountain',
        children: [
            { label: 'Pariwisata', icon: 'i-tabler-beach', to: '/admin/potensi/pariwisata' },
            { label: 'UMKM', icon: 'i-tabler-shopping-bag', to: '/admin/potensi/umkm' },
            { label: 'Investasi', icon: 'i-tabler-trending-up', to: '/admin/potensi/investasi' }
        ]
    },
    {
        label: 'Aplikasi',
        icon: 'i-tabler-apps',
        to: '/admin/aplikasi'
    },
    {
        label: 'Media',
        icon: 'i-tabler-photo',
        to: '/admin/media'
    },
    {
        label: 'Pengaturan',
        icon: 'i-tabler-settings',
        children: [
            { label: 'Identitas Situs', icon: 'i-tabler-id-badge-2', to: '/admin/pengaturan/identitas' },
            { label: 'Menu Navigasi', icon: 'i-tabler-menu-2', to: '/admin/pengaturan/menu' },
            { label: 'Kontak', icon: 'i-tabler-phone', to: '/admin/pengaturan/kontak' }
        ]
    },
    {
        label: 'Pengguna & Hak Akses',
        icon: 'i-tabler-users',
        children: [
            { label: 'Daftar Pengguna', icon: 'i-tabler-users-group', to: '/admin/pengguna/daftar' },
            { label: 'Role', icon: 'i-tabler-shield-lock', to: '/admin/pengguna/role' },
            { label: 'Log Aktivitas', icon: 'i-tabler-history', to: '/admin/pengguna/log' }
        ]
    }
]

const route = useRoute()
const isSidebarOpen = ref(true)

const userMenuItems = [
    [
        { label: 'Profil', icon: 'i-tabler-user', to: '/admin/profil' },
        { label: 'Pengaturan', icon: 'i-tabler-settings', to: '/admin/pengaturan' }
    ],
    [
        { label: 'Keluar', icon: 'i-tabler-logout', to: '/' }
    ]
]
</script>

<template>
    <div class="min-h-screen bg-slate-50 dark:bg-slate-950">
        <!-- SIDEBAR (Terang dengan tint biru) -->
        <aside
            class="fixed top-0 left-0 z-40 h-screen transition-all duration-300 bg-sky-50 dark:bg-slate-900 border-r border-sky-100 dark:border-slate-800 flex flex-col"
            :class="isSidebarOpen ? 'w-64' : 'w-16'">
            <!-- Logo Area -->
            <div class="flex items-center gap-3 h-16 px-4 border-b border-sky-100 dark:border-slate-800 shrink-0">
                <img src="/Logo-Kabupaten-Tulungagung.png" alt="Logo" class="w-8 h-9 shrink-0" />
                <div v-if="isSidebarOpen" class="leading-tight overflow-hidden">
                    <p class="text-[10px] text-sky-700/70 dark:text-slate-400">Admin Panel</p>
                    <p class="text-sm font-bold text-sky-900 dark:text-white whitespace-nowrap">Tulungagung</p>
                </div>
            </div>

            <!-- Nav Menu -->
            <nav class="flex-1 overflow-y-auto py-4 px-2 sidebar-scroll">
                <template v-for="(item, i) in sidebarLinks" :key="i">
                    <!-- Item dengan children -->
                    <div v-if="item.children">
                        <div
                            class="flex items-center gap-3 px-3 py-2 rounded-lg text-xs font-bold text-sky-700/70 dark:text-slate-500 uppercase tracking-wider mb-1 mt-3">
                            <UIcon :name="item.icon" class="w-4 h-4 shrink-0" />
                            <span v-if="isSidebarOpen">{{ item.label }}</span>
                        </div>
                        <div v-if="isSidebarOpen" class="space-y-0.5 mb-2">
                            <NuxtLink v-for="child in item.children" :key="child.to" :to="child.to"
                                class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors"
                                :class="route.path === child.to
                                    ? 'bg-sky-500 text-white font-semibold shadow-sm shadow-sky-500/30'
                                    : 'text-slate-600 dark:text-slate-300 hover:bg-sky-100 dark:hover:bg-slate-800 hover:text-sky-700 dark:hover:text-white'">
                                <UIcon :name="child.icon" class="w-4 h-4 shrink-0" />
                                <span class="truncate">{{ child.label }}</span>
                            </NuxtLink>
                        </div>
                    </div>

                    <!-- Item tanpa children -->
                    <NuxtLink v-else :to="item.to"
                        class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors mb-1"
                        :class="route.path === item.to
                            ? 'bg-sky-500 text-white font-semibold shadow-sm shadow-sky-500/30'
                            : 'text-slate-700 dark:text-slate-300 hover:bg-sky-100 dark:hover:bg-slate-800 hover:text-sky-700 dark:hover:text-white'">
                        <UIcon :name="item.icon" class="w-5 h-5 shrink-0" />
                        <span v-if="isSidebarOpen">{{ item.label }}</span>
                    </NuxtLink>
                </template>
            </nav>

            <!-- User Info -->
            <div class="border-t border-sky-100 dark:border-slate-800 p-3 shrink-0 bg-sky-100/50 dark:bg-slate-900">
                <div class="flex items-center gap-3">
                    <div
                        class="w-9 h-9 rounded-full bg-sky-500 flex items-center justify-center text-white font-bold shrink-0">
                        {{ user.name.charAt(0) }}
                    </div>
                    <div v-if="isSidebarOpen" class="flex-1 min-w-0">
                        <p class="text-sm font-semibold text-sky-900 dark:text-white truncate">{{ user.name }}</p>
                        <p class="text-xs text-sky-700/70 dark:text-slate-400 truncate">{{ user.email }}</p>
                    </div>
                </div>
            </div>
        </aside>

        <!-- MAIN AREA -->
        <div class="transition-all duration-300 flex flex-col min-h-screen" :class="isSidebarOpen ? 'ml-64' : 'ml-16'">
            <!-- TOP BAR -->
            <header
                class="sticky top-0 z-30 bg-white dark:bg-slate-900 border-b border-slate-200 dark:border-slate-800 h-16 flex items-center px-6 gap-4">
                <UButton
                    :icon="isSidebarOpen ? 'i-tabler-layout-sidebar-left-collapse' : 'i-tabler-layout-sidebar-left-expand'"
                    color="neutral" variant="ghost" @click="isSidebarOpen = !isSidebarOpen" />

                <div class="flex-1">
                    <h1 class="text-lg font-bold text-slate-900 dark:text-white">
                        {{ route.meta.title || 'Dashboard' }}
                    </h1>
                </div>

                <UButton icon="i-tabler-bell" color="neutral" variant="ghost" />
                <UButton icon="i-tabler-home" color="neutral" variant="ghost" to="/" />
                <UDropdownMenu :items="userMenuItems">
                    <UButton color="neutral" variant="ghost">
                        <div
                            class="w-7 h-7 rounded-full bg-sky-500 flex items-center justify-center text-white text-xs font-bold">
                            {{ user.name.charAt(0) }}
                        </div>
                    </UButton>
                </UDropdownMenu>
            </header>

            <!-- CONTENT -->
            <main class="flex-1 p-6">
                <slot />
            </main>

            <!-- FOOTER -->
            <footer
                class="border-t border-slate-200 dark:border-slate-800 px-6 py-4 text-xs text-slate-400 text-center">
                © {{ new Date().getFullYear() }} Pemerintah Kabupaten Tulungagung — Admin Panel
            </footer>
        </div>
    </div>
</template>

<style scoped>
/* Scrollbar tipis untuk sidebar */
.sidebar-scroll::-webkit-scrollbar {
    width: 6px;
}

.sidebar-scroll::-webkit-scrollbar-track {
    background: transparent;
}

.sidebar-scroll::-webkit-scrollbar-thumb {
    background: rgba(14, 165, 233, 0.2);
    border-radius: 3px;
}

.sidebar-scroll::-webkit-scrollbar-thumb:hover {
    background: rgba(14, 165, 233, 0.4);
}
</style>