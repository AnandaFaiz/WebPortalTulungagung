<script setup>
const route = useRoute()
const q = computed(() => String(route.query.q || ''))
const searchInput = ref(q.value)

watch(q, (newVal) => {
    searchInput.value = newVal
})

const handleSearch = () => {
    if (searchInput.value.trim()) {
        navigateTo(`/pencarian?q=${encodeURIComponent(searchInput.value.trim())}`)
    }
}

// Fetch hasil pencarian dari API Golang
const { data: searchData, pending: loading, refresh } = await useFetch(
    () => `http://localhost:8080/api/search?q=${encodeURIComponent(q.value)}`
)

// Memisahkan hasil berdasarkan kategori
const hasilBerita = computed(() => searchData.value?.berita || [])
const hasilHalaman = computed(() => searchData.value?.halaman || [])
const hasilOPD = computed(() => searchData.value?.opd || [])
const hasilAplikasi = computed(() => searchData.value?.aplikasi || [])
const hasilDarurat = computed(() => searchData.value?.darurat || [])
const hasilLayanan = computed(() => [
    ...(searchData.value?.layanan || []),
    ...(searchData.value?.layanan_populer || []),
    ...(searchData.value?.kategori_layanan || [])
])
const hasilInformasi = computed(() => [
    ...(searchData.value?.pengumuman || []),
    ...(searchData.value?.agenda || []),
    ...(searchData.value?.galeri || []),
    ...(searchData.value?.pejabat || []),
    ...(searchData.value?.tautan || []),
    ...(searchData.value?.menu || []),
    ...(searchData.value?.profil_daerah || []),
    ...(searchData.value?.struktur || []),
    ...(searchData.value?.identitas || []),
    ...(searchData.value?.kontak || []),
    ...(searchData.value?.lainnya || [])
])
const totalHasil = computed(() => {
    return hasilBerita.value.length + hasilHalaman.value.length + hasilOPD.value.length + hasilAplikasi.value.length + hasilDarurat.value.length + hasilLayanan.value.length + hasilInformasi.value.length
})

const getRoutePath = (item) => {
    if (item.tipe === 'berita') return `/berita/${item.id}`
    if (item.tipe === 'halaman') return `/${item.slug}`
    if (item.tipe === 'opd') return `/pemerintahan#opd` // Sementara arahkan ke halaman pemerintahan
    if (item.tipe === 'darurat') return '/layanan#darurat'
    if (item.tipe === 'layanan' || item.tipe === 'kategori_layanan') return '/layanan#kategori-layanan'
    if (item.tipe === 'layanan_populer') return item.url && item.url !== '#' ? item.url : '/layanan#layanan-populer'
    if (item.tipe === 'pejabat') return '/pemerintahan'
    if (item.tipe === 'static') return item.url || '/'
    if (item.tipe === 'agenda' || item.tipe === 'pengumuman' || item.tipe === 'galeri') return '/berita'
    if (item.tipe === 'tautan' || item.tipe === 'menu') return item.url || '/'
    return '#'
}

const isExternalUrl = (url) => /^https?:\/\//i.test(url || '')

const resultTypeLabel = (type) => ({
    pengumuman: 'Pengumuman',
    agenda: 'Agenda',
    galeri: 'Galeri',
    pejabat: 'Profil Pejabat',
    tautan: 'Tautan Eksternal',
    menu: 'Navigasi Portal',
    layanan: 'Layanan Publik',
    layanan_populer: 'Layanan Populer',
    kategori_layanan: 'Kategori Layanan',
    profil_daerah: 'Profil Daerah',
    struktur: 'Struktur Organisasi',
    identitas: 'Identitas Portal',
    kontak: 'Kontak',
    static: 'Portal'
}[type] || 'Informasi Portal')

</script>

<template>
    <div>
        <!-- Header Hero Sederhana untuk Halaman Pencarian -->
        <div class="bg-sky-900 py-16 px-4">
            <div class="max-w-4xl mx-auto text-center">
                <h1 class="text-2xl md:text-4xl font-bold text-white mb-6 pt-8">
                    Hasil Pencarian
                </h1>

                <!-- Search Bar (Bisa dipakai untuk mencari ulang) -->
                <form @submit.prevent="handleSearch"
                    class="relative max-w-xl mx-auto flex items-center bg-white rounded-lg overflow-hidden shadow-xl">
                    <UIcon name="i-tabler-search" class="absolute left-4 w-6 h-6 text-slate-400" />
                    <input v-model="searchInput" type="text" placeholder="Cari layanan, berita, atau informasi..."
                        class="w-full pl-12 pr-4 py-3 text-slate-800 outline-none text-base" />
                    <button type="submit"
                        class="bg-sky-500 hover:bg-sky-600 text-white font-semibold px-6 py-3 transition-colors">
                        Cari
                    </button>
                </form>
            </div>
        </div>

        <!-- Konten Hasil Pencarian -->
        <div class="max-w-4xl mx-auto px-4 py-12 min-h-[60vh]">

            <!-- Loading State -->
            <div v-if="loading" class="text-center py-20">
                <UIcon name="i-tabler-loader-2" class="w-12 h-12 text-sky-500 animate-spin mx-auto" />
                <p class="text-slate-500 mt-4">Mencari "{{ q }}"...</p>
            </div>

            <!-- Empty State -->
            <div v-else-if="totalHasil === 0" class="text-center py-20 bg-slate-50 dark:bg-slate-900 rounded-2xl">
                <UIcon name="i-tabler-search-off" class="w-20 h-20 text-slate-300 mx-auto mb-4" />
                <h2 class="text-2xl font-bold text-slate-700 dark:text-slate-300 mb-2">Tidak ada hasil ditemukan</h2>
                <p class="text-slate-500">Maaf, kami tidak dapat menemukan hasil untuk kata kunci "{{ q }}".</p>
                <p class="text-sm text-slate-400 mt-2">Coba gunakan kata kunci lain atau periksa ejaan Anda.</p>
                <UButton to="/" class="mt-6" color="primary" variant="soft" icon="i-tabler-arrow-left">
                    Kembali ke Beranda
                </UButton>
            </div>

            <!-- Hasil Ditemukan -->
            <div v-else>
                <div class="mb-8 border-b border-slate-200 dark:border-slate-700 pb-4">
                    <p class="text-white">
                        Menampilkan <strong class="text-teal-300">{{ totalHasil }} hasil</strong> untuk kata kunci
                        <strong class="text-white">"{{ q }}"</strong>
                    </p>
                </div>

                <!-- Grup Berita -->
                <div v-if="hasilBerita.length > 0" class="mb-10">
                    <h3 class="text-xl font-bold text-white mb-4 flex items-center gap-2">
                        <UIcon name="i-tabler-news" class="text-amber-300" />
                        Berita ({{ hasilBerita.length }})
                    </h3>
                    <div class="space-y-3">
                        <NuxtLink v-for="item in hasilBerita" :key="item.id" :to="getRoutePath(item)"
                            class="block bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl p-4 hover:border-sky-400 hover:shadow-md transition-all group">
                            <div class="flex items-center justify-between gap-4">
                                <div>
                                    <h4
                                        class="font-semibold text-slate-900 dark:text-white group-hover:text-sky-600 transition-colors">
                                        {{ item.judul }}
                                    </h4>
                                    <p class="text-xs text-slate-400 mt-1">Kategori: Berita</p>
                                </div>
                                <UIcon name="i-tabler-arrow-right"
                                    class="w-5 h-5 text-slate-300 group-hover:text-amber-300 transition-colors shrink-0" />
                            </div>
                        </NuxtLink>
                    </div>
                </div>

                <!-- Grup Halaman Statis -->
                <div v-if="hasilHalaman.length > 0" class="mb-10">
                    <h3 class="text-xl font-bold text-white mb-4 flex items-center gap-2">
                        <UIcon name="i-tabler-file-description" class="text-amber-300" />
                        Halaman Informasi ({{ hasilHalaman.length }})
                    </h3>
                    <div class="space-y-3">
                        <NuxtLink v-for="item in hasilHalaman" :key="item.id" :to="getRoutePath(item)"
                            class="block bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl p-4 hover:border-sky-400 hover:shadow-md transition-all group">
                            <div class="flex items-center justify-between gap-4">
                                <div>
                                    <h4
                                        class="font-semibold text-slate-900 dark:text-white group-hover:text-sky-600 transition-colors">
                                        {{ item.judul }}
                                    </h4>
                                    <p class="text-xs text-slate-400 mt-1">Halaman /{{ item.slug }}</p>
                                </div>
                                <UIcon name="i-tabler-arrow-right"
                                    class="w-5 h-5 text-slate-300 group-hover:text-amber-300 transition-colors shrink-0" />
                            </div>
                        </NuxtLink>
                    </div>
                </div>

                <!-- Grup OPD -->
                <div v-if="hasilOPD.length > 0" class="mb-10">
                    <h3 class="text-xl font-bold text-white mb-4 flex items-center gap-2">
                        <UIcon name="i-tabler-building-community" class="text-amber-300" />
                        Organisasi Perangkat Daerah ({{ hasilOPD.length }})
                    </h3>
                    <div class="space-y-3">
                        <NuxtLink v-for="item in hasilOPD" :key="item.id" :to="getRoutePath(item)"
                            class="block bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl p-4 hover:border-sky-400 hover:shadow-md transition-all group">
                            <div class="flex items-center justify-between gap-4">
                                <div>
                                    <h4
                                        class="font-semibold text-slate-900 dark:text-white group-hover:text-sky-600 transition-colors">
                                        {{ item.judul }}
                                    </h4>
                                    <p class="text-xs text-slate-400 mt-1">Perangkat Daerah</p>
                                </div>
                                <UIcon name="i-tabler-arrow-right"
                                    class="w-5 h-5 text-slate-300 group-hover:text-amber-300 transition-colors shrink-0" />
                            </div>
                        </NuxtLink>
                    </div>
                </div>

                <!-- Grup Aplikasi Daerah -->
                <div v-if="hasilAplikasi.length > 0" class="mb-10">
                    <h3 class="text-xl font-bold text-white mb-4 flex items-center gap-2">
                        <UIcon name="i-tabler-apps" class="text-amber-300" />
                        Aplikasi Daerah ({{ hasilAplikasi.length }})
                    </h3>
                    <div class="space-y-3">
                        <a v-for="item in hasilAplikasi" :key="item.id" :href="item.url || '#'" target="_blank"
                            class="block bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl p-4 hover:border-sky-400 hover:shadow-md transition-all group">
                            <div class="flex items-center justify-between gap-4">
                                <div>
                                    <h4
                                        class="font-semibold text-slate-900 dark:text-white group-hover:text-sky-600 transition-colors">
                                        {{ item.judul }}
                                    </h4>
                                    <p class="text-xs text-slate-400 mt-1">Aplikasi Layanan</p>
                                </div>
                                <UIcon name="i-tabler-external-link"
                                    class="w-5 h-5 text-slate-300 group-hover:text-amber-300 transition-colors shrink-0" />
                            </div>
                        </a>
                    </div>
                </div>

                <!-- Grup Layanan Darurat -->
                <div v-if="hasilDarurat.length > 0" class="mb-10">
                    <h3 class="text-xl font-bold text-white mb-4 flex items-center gap-2">
                        <UIcon name="i-tabler-phone-call" class="text-rose-500" />
                        Layanan Darurat ({{ hasilDarurat.length }})
                    </h3>
                    <div class="space-y-3">
                        <NuxtLink v-for="item in hasilDarurat" :key="item.id" to="/layanan#darurat"
                            class="block bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl p-4 hover:border-rose-400 hover:shadow-md transition-all group">
                            <div class="flex items-center justify-between gap-4">
                                <div>
                                    <h4
                                        class="font-semibold text-slate-900 dark:text-white group-hover:text-rose-600 transition-colors">
                                        {{ item.judul }}
                                    </h4>
                                    <p class="text-xs text-slate-400 mt-1">Nomor: {{ item.deskripsi }}</p>
                                </div>
                                <UIcon name="i-tabler-arrow-right"
                                    class="w-5 h-5 text-slate-400 group-hover:text-rose-500 transition-colors shrink-0" />
                            </div>
                        </NuxtLink>
                    </div>
                </div>

                <!-- Grup Layanan Publik -->
                <div v-if="hasilLayanan.length > 0" class="mb-10">
                    <h3 class="text-xl font-bold text-white mb-4 flex items-center gap-2">
                        <UIcon name="i-tabler-briefcase" class="text-amber-300" />
                        Layanan Publik ({{ hasilLayanan.length }})
                    </h3>
                    <div class="space-y-3">
                        <template v-for="item in hasilLayanan" :key="`${item.tipe}-${item.id}`">
                            <a v-if="isExternalUrl(item.url)" :href="item.url" target="_blank" rel="noopener noreferrer"
                                class="block bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl p-4 hover:border-sky-400 hover:shadow-md transition-all group">
                                <div class="flex items-center justify-between gap-4">
                                    <div>
                                        <h4
                                            class="font-semibold text-slate-900 dark:text-white group-hover:text-sky-600 transition-colors">
                                            {{ item.judul }}</h4>
                                        <p class="text-xs text-slate-400 mt-1">{{ resultTypeLabel(item.tipe) }}</p>
                                    </div>
                                    <UIcon name="i-tabler-external-link"
                                        class="w-5 h-5 text-slate-300 group-hover:text-amber-300 transition-colors shrink-0" />
                                </div>
                            </a>
                            <NuxtLink v-else :to="getRoutePath(item)"
                                class="block bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl p-4 hover:border-sky-400 hover:shadow-md transition-all group">
                                <div class="flex items-center justify-between gap-4">
                                    <div>
                                        <h4
                                            class="font-semibold text-slate-900 dark:text-white group-hover:text-sky-600 transition-colors">
                                            {{ item.judul }}</h4>
                                        <p class="text-xs text-slate-400 mt-1">{{ resultTypeLabel(item.tipe) }}</p>
                                    </div>
                                    <UIcon name="i-tabler-arrow-right"
                                        class="w-5 h-5 text-slate-400 group-hover:text-amber-300 transition-colors shrink-0" />
                                </div>
                            </NuxtLink>
                        </template>
                    </div>
                </div>

                <!-- Grup Informasi Portal -->
                <div v-if="hasilInformasi.length > 0" class="mb-10">
                    <h3 class="text-xl font-bold text-white mb-4 flex items-center gap-2">
                        <UIcon name="i-tabler-file-search" class="text-amber-300" />
                        Informasi Portal ({{ hasilInformasi.length }})
                    </h3>
                    <div class="space-y-3">
                        <template v-for="item in hasilInformasi" :key="`${item.tipe}-${item.id}`">
                            <a v-if="isExternalUrl(item.url)" :href="item.url" target="_blank" rel="noopener noreferrer"
                                class="block bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl p-4 hover:border-sky-400 hover:shadow-md transition-all group">
                                <div class="flex items-center justify-between gap-4">
                                    <div>
                                        <h4
                                            class="font-semibold text-slate-900 dark:text-white group-hover:text-sky-600 transition-colors">
                                            {{ item.judul }}</h4>
                                        <p class="text-xs text-slate-400 mt-1">{{ resultTypeLabel(item.tipe) }}</p>
                                    </div>
                                    <UIcon name="i-tabler-external-link"
                                        class="w-5 h-5 text-slate-300 group-hover:text-amber-300 transition-colors shrink-0" />
                                </div>
                            </a>
                            <NuxtLink v-else :to="getRoutePath(item)"
                                class="block bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl p-4 hover:border-sky-400 hover:shadow-md transition-all group">
                                <div class="flex items-center justify-between gap-4">
                                    <div>
                                        <h4
                                            class="font-semibold text-slate-900 dark:text-white group-hover:text-sky-600 transition-colors">
                                            {{ item.judul }}</h4>
                                        <p class="text-xs text-slate-400 mt-1">{{ resultTypeLabel(item.tipe) }}</p>
                                    </div>
                                    <UIcon name="i-tabler-arrow-right"
                                        class="w-5 h-5 text-slate-400 group-hover:text-amber-300 transition-colors shrink-0" />
                                </div>
                            </NuxtLink>
                        </template>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>