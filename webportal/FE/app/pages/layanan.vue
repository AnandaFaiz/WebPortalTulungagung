<script setup>

const { data: aplikasiData, pending: aplikasiLoading } = await useFetch('http://localhost:8080/api/aplikasi')
const { data: layananPopulerData, pending: layananPopulerLoading } = await useFetch('http://localhost:8080/api/layanan-terpopuler')
const { data: heroData } = await useFetch('http://localhost:8080/api/banner?tipe=layanan')
const { data: daruratData, pending: daruratLoading } = await useFetch('http://localhost:8080/api/layanan-darurat')
const daruratList = computed(() => daruratData.value?.darurat || [])
const heroImages = computed(() => heroData.value?.banner?.map(b => b.gambar) || [])
const { data: mppData, pending: mppLoading } = await useFetch('http://localhost:8080/api/mpp')
const mpp = computed(() => mppData.value?.mpp || {})

const appSearch = ref('')
const filteredApps = computed(() => {
    const apps = aplikasiData.value?.aplikasi || []
    if (!appSearch.value) return apps
    return apps.filter(app => app.nama.toLowerCase().includes(appSearch.value.toLowerCase()))
})
// Fetch Kategori & Sub-layanan dari Golang
const { data: layananApiData, pending: layananLoading } = await useFetch('http://localhost:8080/api/layanan')
const kategoriList = computed(() => layananApiData.value?.kategori || [])
const layananPopulerList = computed(() => layananPopulerData.value?.layanan || [])

const activeDropdown = ref(null)
const searchKategori = ref('')

const toggleDropdown = (id) => {
    activeDropdown.value = activeDropdown.value === id ? null : id
}

const filteredKategori = computed(() => {
    const list = kategoriList.value || []
    if (!searchKategori.value.trim) return list
    const keyword = searchKategori.value.toLowerCase()
    return list.filter(kat => {
        const matchNama = (kat.nama || '').toLowerCase().includes(keyword)
        const matchDeskripsi = (kat.deskripsi || '').toLowerCase().includes(keyword)
        const matchItems = kat.items && kat.items.some(item =>
            (item.nama || '').toLowerCase().includes(keyword)
        )
        return matchNama || matchDeskripsi || matchItems
    })
})
</script>

<template>
    <div>
        <div>
            <AnimatedHero badge="Pelayanan Warga" title="Layanan Publik"
                description="Akses berbagai layanan publik Pemkab Tulungagung dalam satu halaman."
                :images="heroImages.length > 0 ? heroImages : ['/hero/layanan-1.webp']" :interval="5000" />
        </div>

        <div class="relative z-10 bg-transparent">
            <!-- Section Darurat -->
            <section id="darurat" class="py-16">
                <UContainer>
                    <div class="mb-12 text-center">
                        <h2 class="text-4xl font-bold text-white md:text-5xl">Layanan Darurat</h2>
                        <p class="mt-4 text-lg text-white">Nomor penting yang dapat dihubungi dalam situasi darurat.</p>
                    </div>
                    <div>
                        <div v-if="daruratLoading" class="text-center py-8">Memuat...</div>
                        <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                            <UCard v-for="(item, i) in daruratList" :key="i"
                                class="bg-white hover:ring-2 hover:ring-error transition-all">
                                <div class="flex items-center gap-4">
                                    <div
                                        class="w-12 h-12 rounded-full bg-error/10 flex items-center justify-center shrink-0">
                                        <AutoIcon :text="item.nama" class="w-6 h-6 text-error" />
                                    </div>
                                    <div>
                                        <p
                                            class="font-mono text-lg font-bold tracking-wide text-red-600 dark:text-red-400">
                                            {{ item.nomor }}</p>
                                        <p class="text-sm text-slate-500">{{ item.nama }}</p>
                                    </div>
                                </div>
                            </UCard>
                        </div>
                    </div>
                </UContainer>
            </section>

            <!-- Section Kategori Layanan -->
            <section id="kategori-layanan" class="bg-transparent py-12">
                <UContainer>
                    <div
                        class="bg-white/75 dark:bg-slate-900/75 backdrop-blur-sm border border-white/60 dark:border-slate-800 rounded-2xl shadow-sm overflow-hidden">

                        <div
                            class="p-6 md:p-8 border-b border-slate-200 dark:border-slate-800 flex flex-col md:flex-row md:items-center justify-between gap-4">
                            <div class="text-left">
                                <h2 class="text-2xl font-bold text-slate-900 dark:text-white">Kategori Layanan</h2>
                                <p class="text-slate-500 mt-1">Pilih kategori layanan yang Anda butuhkan.</p>
                            </div>
                            <div class="w-full md:w-64 shrink-0">
                                <UInput v-model="searchKategori" icon="i-tabler-search" placeholder="Cari kategori..."
                                    size="md" class="w-full" />
                            </div>
                        </div>

                        <div class="p-6 md:p-8">
                            <div v-if="layananLoading" class="text-center py-8">
                                <UIcon name="i-tabler-loader-2" class="w-8 h-8 text-sky-500 animate-spin mx-auto" />
                            </div>

                            <TransitionGroup v-else-if="filteredKategori.length > 0" tag="div"
                                move-class="transition-transform duration-500 ease-[cubic-bezier(0.22,1,0.36,1)]"
                                class="columns-1 md:columns-2 gap-4">
                                <div v-for="kat in filteredKategori" :key="kat.id"
                                    class="mb-4 break-inside-avoid rounded-xl border border-white/60 dark:border-slate-700 bg-white/30 dark:bg-slate-950/25 backdrop-blur-sm overflow-hidden transition-[border-color,box-shadow] duration-300"
                                    :class="activeDropdown === kat.id ? 'border-sky-400 shadow-md ring-1 ring-sky-400' : 'hover:border-sky-300 hover:shadow-sm'">
                                    <button @click="toggleDropdown(kat.id)" :aria-expanded="activeDropdown === kat.id"
                                        class="w-full flex items-center gap-4 p-4 text-left transition-colors"
                                        :class="activeDropdown === kat.id ? 'bg-sky-50 dark:bg-sky-950/30' : 'bg-transparent'">
                                        <div class="w-12 h-12 rounded-lg flex items-center justify-center shrink-0 transition-colors"
                                            :class="activeDropdown === kat.id ? 'bg-sky-500 text-white' : 'bg-sky-100 dark:bg-sky-900 text-sky-600 dark:text-sky-400'">
                                            <AutoIcon :text="kat.nama" size="w-6 h-6" />
                                        </div>

                                        <div class="flex-1 min-w-0">
                                            <h3 class="font-bold text-slate-900 dark:text-white text-sm md:text-base">{{
                                                kat.nama }}</h3>
                                            <p class="text-xs text-slate-500 mt-1 line-clamp-1">{{ kat.deskripsi }}</p>
                                        </div>

                                        <UIcon name="i-tabler-chevron-down"
                                            class="w-5 h-5 text-slate-400 transition-[transform,color] duration-500 ease-[cubic-bezier(0.22,1,0.36,1)]"
                                            :class="activeDropdown === kat.id ? 'rotate-180 text-sky-500' : ''" />
                                    </button>

                                    <div class="grid transition-[grid-template-rows,opacity] duration-500 ease-[cubic-bezier(0.22,1,0.36,1)]"
                                        :class="activeDropdown === kat.id ? 'grid-rows-[1fr] opacity-100' : 'grid-rows-[0fr] opacity-0'"
                                        :inert="activeDropdown !== kat.id">
                                        <div class="min-h-0 overflow-hidden">
                                            <div
                                                class="bg-white/30 dark:bg-slate-900/30 backdrop-blur-sm border-t border-white/50 dark:border-slate-700">
                                                <div class="p-4 space-y-1">
                                                    <!-- Loop data asli dari API Golang -->
                                                    <NuxtLink v-for="(item, index) in kat.items" :key="index"
                                                        :to="item.link || '#'"
                                                        class="flex items-center justify-between px-4 py-2.5 text-sm text-slate-700 dark:text-slate-300 hover:text-sky-600 hover:bg-white/70 dark:hover:bg-slate-800/70 rounded-lg transition-colors group">
                                                        <span class="flex items-center gap-3">
                                                            <span
                                                                class="w-5 text-right text-xs font-bold text-slate-400 group-hover:text-sky-500">
                                                                {{ index + 1 }}.
                                                            </span>
                                                            {{ item.nama }}
                                                        </span>
                                                        <UIcon name="i-tabler-arrow-right"
                                                            class="w-4 h-4 opacity-0 group-hover:opacity-100 transition-opacity" />
                                                    </NuxtLink>

                                                    <!-- Jika Kategori tidak punya sub-layanan -->
                                                    <p v-if="!kat.items || kat.items.length === 0"
                                                        class="text-xs text-slate-400 italic px-4 py-2">Belum ada daftar
                                                        layanan.</p>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </TransitionGroup>

                            <div v-else class="text-center py-8">
                                <UIcon name="i-tabler-search-off" class="w-12 h-12 text-slate-300 mx-auto mb-2" />
                                <p class="text-slate-500">Kategori layanan tidak ditemukan.</p>
                            </div>
                        </div>
                    </div>
                </UContainer>
            </section>

            <section id="layanan-populer" class="py-16">
                <UContainer>
                    <div>
                        <div class="mb-8 text-center">
                            <h2 class="text-4xl font-bold text-white md:text-5xl ">Layanan Populer
                            </h2>
                            <p class="mt-4 text-lg text-white dark:text-slate-300">Layanan yang dikelola oleh
                                pemerintah daerah.</p>
                        </div>
                        <div v-if="layananPopulerLoading" class="py-8 text-center">
                            <UIcon name="i-tabler-loader-2" class="mx-auto h-8 w-8 animate-spin text-sky-500" />
                        </div>
                        <div v-else-if="layananPopulerList.length"
                            class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
                            <a v-for="item in layananPopulerList" :key="item.id" :href="item.link || '#'"
                                class="group flex items-start gap-4 rounded-xl border border-emerald-500 bg-emerald-100 p-4 transition-colors hover:bg-emerald-300 ">
                                <UIcon :name="item.icon || 'i-tabler-star'" class="h-6 w-6 shrink-0 text-emerald-700" />
                                <span>
                                    <span class="block font-bold text-slate-900">{{ item.judul }}</span>
                                    <span class="mt-1 block text-xs text-slate-700">{{ item.deskripsi }}</span>
                                </span>
                            </a>
                        </div>
                        <p v-else class="py-8 text-center text-sm text-slate-600 dark:text-slate-300">Belum ada layanan
                            populer yang tersedia.</p>
                    </div>
                </UContainer>
            </section>

            <!-- Section Aplikasi Daerah -->
            <section id="aplikasi-daerah" class="bg-transparent py-16">
                <UContainer>
                    <!-- Card Bungkusan Utama -->
                    <div
                        class="bg-white/75 dark:bg-slate-900/75 backdrop-blur-sm border border-white/60 dark:border-slate-800 rounded-2xl shadow-sm overflow-hidden">

                        <!-- Header: Judul Kiri, Pencarian Kanan -->
                        <div
                            class="p-6 md:p-8 border-b border-slate-200 dark:border-slate-800 flex flex-col md:flex-row md:items-center justify-between gap-4">
                            <div class="text-left">
                                <h2 class="text-2xl font-bold text-slate-900 dark:text-white">Aplikasi Layanan Digital
                                </h2>
                                <p class="text-slate-500 mt-1">Akses cepat layanan mandiri Pemkab Tulungagung melalui
                                    aplikasi berikut.</p>
                            </div>

                            <!-- Input Pencarian Khusus Aplikasi -->
                            <div class="w-full md:w-64 shrink-0">
                                <UInput v-model="appSearch" icon="i-tabler-search" placeholder="Cari aplikasi..."
                                    size="md" class="w-full" />
                            </div>
                        </div>

                        <!-- Grid Aplikasi -->
                        <div class="p-6 md:p-8">
                            <div v-if="aplikasiLoading" class="text-center py-8">
                                <UIcon name="i-tabler-loader-2" class="w-8 h-8 text-sky-500 animate-spin mx-auto" />
                            </div>

                            <!-- Menggunakan 2 Kolom (Kiri & Kanan) dengan batas maksimal 8 item -->
                            <div v-else-if="filteredApps.length > 0" class="grid grid-cols-1 md:grid-cols-2 gap-4">
                                <a v-for="item in filteredApps.slice(0, 8)" :key="item.id" :href="item.url || '#'"
                                    target="_blank"
                                    class="group flex items-center gap-4 p-4 rounded-xl border border-white/60 dark:border-slate-700 bg-white/35 dark:bg-slate-950/25 backdrop-blur-sm hover:border-sky-400 hover:bg-sky-100 hover:shadow-md transition-all cursor-pointer">
                                    <!-- Ikon di Kiri -->
                                    <div
                                        class="w-12 h-12 rounded-lg bg-sky-100 dark:bg-sky-900 flex items-center justify-center shrink-0 group-hover:bg-sky-500 transition-colors">
                                        <UIcon :name="item.icon || 'i-tabler-apps'"
                                            class="w-6 h-6 text-sky-600 dark:text-sky-400 group-hover:text-white transition-colors" />
                                    </div>

                                    <!-- Teks di Kanan -->
                                    <div class="flex-1 min-w-0 text-left">
                                        <h3
                                            class="font-bold text-slate-900 dark:text-white group-hover:text-sky-600 text-sm md:text-base transition-colors truncate">
                                            {{ item.nama }}
                                        </h3>
                                        <p class="text-xs text-slate-500 mt-1 line-clamp-2 leading-relaxed">
                                            {{ item.deskripsi }}
                                        </p>
                                    </div>
                                </a>
                            </div>

                            <!-- Jika Pencarian Tidak Menemukan Hasil -->
                            <div v-else class="text-center py-8">
                                <UIcon name="i-tabler-apps-off" class="w-12 h-12 text-slate-300 mx-auto mb-2" />
                                <p class="text-slate-500">Aplikasi yang Anda cari tidak ditemukan.</p>
                            </div>
                        </div>

                    </div>
                </UContainer>
            </section>

            <!-- Section MPP -->
            <section class="py-16 overflow-hidden">
                <UContainer>
                    <div class="mb-12 text-center">
                        <h2 class="text-4xl font-bold text-white md:text-5xl">Mal Pelayanan Publik (MPP)</h2>
                        <p class="mt-4 text-lg text-white">Layanan terpadu dari berbagai instansi dalam satu gedung.</p>
                    </div>
                    <div v-if="mppLoading" class="text-center py-8">Memuat...</div>
                    <UCard v-else class="bg-white/75 dark:bg-slate-900/75 backdrop-blur-sm">
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-8 items-center">
                            <div>
                                <p class="text-slate-600 leading-relaxed mb-6">{{ mpp.deskripsi }}</p>
                                <div class="space-y-3 text-sm">
                                    <div class="flex items-start gap-3">
                                        <UIcon name="i-tabler-map-pin" class="text-sky-500 mt-0.5 shrink-0" />
                                        <span>{{ mpp.alamat }}</span>
                                    </div>
                                    <div class="flex items-start gap-3">
                                        <UIcon name="i-tabler-clock" class="text-sky-500 mt-0.5 shrink-0" />
                                        <span>{{ mpp.jam_operasional }}</span>
                                    </div>
                                </div>
                            </div>
                            <div class="grid grid-cols-2 gap-4">
                                <UCard variant="subtle"
                                    class="bg-white/35 dark:bg-slate-950/25 backdrop-blur-sm text-center">
                                    <p class="text-3xl font-bold text-sky-600">{{ mpp.jumlah_layanan }}+</p>
                                    <p class="text-xs text-slate-500 mt-1">Jenis Layanan</p>
                                </UCard>
                                <UCard variant="subtle"
                                    class="bg-white/35 dark:bg-slate-950/25 backdrop-blur-sm text-center">
                                    <p class="text-3xl font-bold text-sky-600">{{ mpp.jumlah_instansi }}</p>
                                    <p class="text-xs text-slate-500 mt-1">Instansi</p>
                                </UCard>
                            </div>
                        </div>
                    </UCard>
                </UContainer>
            </section>
        </div>
    </div>
</template>