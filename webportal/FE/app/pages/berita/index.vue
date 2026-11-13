<script setup>
const { data: heroData } = await useFetch('http://localhost:8080/api/banner?tipe=berita')
const heroImages = computed(() => heroData.value?.banner?.map(b => b.gambar) || [])
const { data: beritaData } = await useFetch('http://localhost:8080/api/berita')
const berita = computed(() => beritaData.value?.berita || [])
const { data: kategoriData } = await useFetch('http://localhost:8080/api/kategori-berita')
const kategori = computed(() => ['Semua', ...(kategoriData.value?.kategori?.map(k => k.nama) || [])])
const { data: pengumumanData } = await useFetch('http://localhost:8080/api/pengumuman')
const pengumuman = computed(() => pengumumanData.value?.pengumuman || [])
const { data: galeriData } = await useFetch('http://localhost:8080/api/galeri')
const galeri = computed(() => galeriData.value?.galeri || [])
const kategoriAktif = ref('Semua')
const pencarian = ref('')
const pencarianPengumuman = ref('')
const halamanPengumuman = ref(1)
const halamanGaleri = ref(1)
const pengumumanPerPage = 3
const galeriPerPage = 4
const stripHtml = (value = '') => value.replace(/<[^>]*>/g, ' ').replace(/\s+/g, ' ').trim()
const beritaText = (item) => `${item.judul} ${item.kategori} ${item.deskripsi} ${stripHtml(item.konten)}`.toLowerCase()
const beritaFiltered = computed(() => kategoriAktif.value === 'Semua'
    ? berita.value.filter(item => beritaText(item).includes(pencarian.value.toLowerCase().trim()))
    : berita.value.filter(item => item.kategori === kategoriAktif.value && beritaText(item).includes(pencarian.value.toLowerCase().trim())))
const beritaTerbaru = computed(() => beritaFiltered.value[0])
const beritaLainnya = computed(() => beritaFiltered.value.slice(1))
const pengumumanFiltered = computed(() => pengumuman.value.filter(item => {
    const query = pencarianPengumuman.value.toLowerCase().trim()
    return `${item.judul} ${stripHtml(item.konten)}`.toLowerCase().includes(query)
}))
const pengumumanPaginated = computed(() => {
    const start = (halamanPengumuman.value - 1) * pengumumanPerPage
    return pengumumanFiltered.value.slice(start, start + pengumumanPerPage)
})
const galeriPaginated = computed(() => {
    const start = (halamanGaleri.value - 1) * galeriPerPage
    return galeri.value.slice(start, start + galeriPerPage)
})
watch(pencarianPengumuman, () => { halamanPengumuman.value = 1 })
</script>

<template>
    <div>
        <div>
            <AnimatedHero badge="Kabar Terkini" title="Informasi & Berita"
                description="Berita, informasi resmi, dan dokumentasi kegiatan Pemkab Tulungagung."
                :images="heroImages.length > 0 ? heroImages : ['/hero/berita-1.webp']" :interval="5000" />
        </div>

        <div class="relative z-10 bg-transparent py-8 md:py-12">
            <UContainer>
                <section
                    class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl shadow-sm overflow-hidden">
                    <div
                        class="p-6 md:p-8 border-b border-slate-200 dark:border-slate-800 flex flex-col md:flex-row md:items-center justify-between gap-4">
                        <div class="text-left">
                            <h1 class="text-2xl font-bold text-slate-900 dark:text-white">Berita Terbaru</h1>
                            <p class="text-slate-500 mt-1">Kabar terkini dari Pemerintah Kabupaten Tulungagung.</p>
                        </div>
                        <div class="w-full md:w-80 shrink-0">
                            <UInput v-model="pencarian" icon="i-tabler-search"
                                placeholder="Cari berita atau isi berita..." size="md" class="w-full" />
                        </div>
                    </div>

                    <div class="p-6 md:p-8">
                        <div class="flex flex-wrap gap-2 mb-8">
                            <UButton v-for="kat in kategori" :key="kat" :label="kat"
                                :color="kategoriAktif === kat ? 'primary' : 'neutral'"
                                :variant="kategoriAktif === kat ? 'solid' : 'outline'" size="sm"
                                @click="kategoriAktif = kat" />
                        </div>

                        <div v-if="beritaTerbaru" class="grid grid-cols-1 lg:grid-cols-[13fr_7fr] gap-8">
                            <NuxtLink :to="`/berita/${beritaTerbaru.id}`" class="group block">
                                <UCard class="h-full">
                                    <div
                                        class="h-64 md:h-96 bg-sky-100 dark:bg-sky-950 rounded-lg overflow-hidden flex items-center justify-center">
                                        <img v-if="beritaTerbaru.thumbnail" :src="beritaTerbaru.thumbnail"
                                            :alt="beritaTerbaru.judul"
                                            class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" />
                                        <UIcon v-else name="i-tabler-news" class="w-20 h-20 text-sky-500" />
                                    </div>
                                    <UBadge color="primary" variant="subtle" size="sm" class="mt-5">{{
                                        beritaTerbaru.kategori }}
                                    </UBadge>
                                    <h3
                                        class="font-bold text-2xl text-slate-900 dark:text-white mt-2 group-hover:text-sky-600">
                                        {{ beritaTerbaru.judul }}</h3>
                                    <p class="text-slate-500 mt-3 line-clamp-3">{{ beritaTerbaru.deskripsi ||
                                        stripHtml(beritaTerbaru.konten) }}</p>
                                    <p class="text-xs text-slate-400 mt-4 flex items-center gap-1">
                                        <UIcon name="i-tabler-calendar" /> {{ beritaTerbaru.tanggal }}
                                    </p>
                                </UCard>
                            </NuxtLink>

                            <div class="space-y-4">
                                <h3 class="text-lg font-bold text-slate-900 dark:text-white">Berita Lainnya</h3>
                                <NuxtLink v-for="item in beritaLainnya" :key="item.id" :to="`/berita/${item.id}`"
                                    class="flex gap-4 p-3 rounded-xl border border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-900 group hover:border-sky-400 transition-colors">
                                    <div
                                        class="w-24 h-20 shrink-0 rounded-lg overflow-hidden bg-sky-100 flex items-center justify-center">
                                        <img v-if="item.thumbnail" :src="item.thumbnail" :alt="item.judul"
                                            class="w-full h-full object-cover" />
                                        <UIcon v-else name="i-tabler-news" class="w-8 h-8 text-sky-500" />
                                    </div>
                                    <div class="min-w-0">
                                        <UBadge color="primary" variant="subtle" size="xs">{{ item.kategori }}</UBadge>
                                        <h4
                                            class="font-semibold text-sm text-slate-900 dark:text-white mt-1 line-clamp-2 group-hover:text-sky-600">
                                            {{ item.judul }}</h4>
                                        <p class="text-xs text-slate-400 mt-2">{{ item.tanggal }}</p>
                                    </div>
                                </NuxtLink>
                            </div>
                        </div>

                        <div v-else class="text-center py-12">
                            <UIcon name="i-tabler-news-off" class="w-16 h-16 text-slate-300 mx-auto mb-3" />
                            <p class="text-slate-500 font-medium">Tidak ada berita di kategori ini</p>
                        </div>

                    </div>
                </section>

                <section id="pengumuman"
                    class="mt-8 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl shadow-sm overflow-hidden">
                    <div
                        class="p-6 md:p-8 border-b border-slate-200 dark:border-slate-800 flex flex-col md:flex-row md:items-center justify-between gap-4">
                        <div class="text-left">
                            <h2 class="text-2xl font-bold text-slate-900 dark:text-white">Pengumuman</h2>
                            <p class="text-slate-500 mt-1">Informasi resmi dari Pemerintah Kabupaten Tulungagung.</p>
                        </div>
                        <div class="w-full md:w-80 shrink-0">
                            <UInput v-model="pencarianPengumuman" icon="i-tabler-search"
                                placeholder="Cari pengumuman..." size="md" class="w-full" />
                        </div>
                    </div>
                    <div class="p-6 md:p-8">
                        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                            <UCard v-for="item in pengumumanPaginated" :key="item.id" class="group">
                                <div
                                    class="h-40 rounded-lg overflow-hidden bg-slate-100 dark:bg-slate-800 flex items-center justify-center mb-4">
                                    <img v-if="item.gambar" :src="item.gambar" :alt="item.judul"
                                        class="w-full h-full object-cover group-hover:scale-105 transition-transform" />
                                    <UIcon v-else name="i-tabler-speakerphone" class="w-12 h-12 text-sky-400" />
                                </div>
                                <UBadge :color="item.prioritas === 'Tinggi' ? 'error' : 'primary'" variant="subtle"
                                    size="sm">{{ item.prioritas }}</UBadge>
                                <h3 class="font-bold text-lg text-slate-900 dark:text-white mt-2">{{ item.judul }}</h3>
                                <p class="text-sm text-slate-500 mt-2 line-clamp-2">{{ stripHtml(item.konten) }}</p>
                                <p class="text-xs text-slate-400 mt-4">{{ item.tanggal }}</p>
                            </UCard>
                        </div>
                        <div v-if="pengumumanFiltered.length === 0" class="text-center py-10">
                            <UIcon name="i-tabler-speakerphone-off" class="w-12 h-12 text-slate-300 mx-auto mb-2" />
                            <p class="text-slate-500">Pengumuman tidak ditemukan.</p>
                        </div>
                        <div v-if="pengumumanFiltered.length > pengumumanPerPage" class="flex justify-end mt-8">
                            <UPagination v-model:page="halamanPengumuman" :total="pengumumanFiltered.length"
                                :items-per-page="pengumumanPerPage" />
                        </div>
                    </div>
                </section>

                <section id="galeri-kegiatan"
                    class="mt-8 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl shadow-sm overflow-hidden scroll-mt-6">
                    <div
                        class="p-6 md:p-8 border-b border-slate-200 dark:border-slate-800 flex flex-col md:flex-row md:items-center justify-between gap-4">
                        <div class="text-left">
                            <h2 class="text-2xl font-bold text-slate-900 dark:text-white">Galeri Kegiatan</h2>
                            <p class="text-slate-500 mt-1">Dokumentasi foto dan video kegiatan pemerintah daerah.</p>
                        </div>
                    </div>
                    <div class="p-6 md:p-8">
                        <div v-if="galeri.length" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
                            <NuxtLink v-for="item in galeriPaginated" :key="item.id" :to="item.file_url || '#'
                                "
                                class="group relative aspect-4/3 rounded-2xl overflow-hidden cursor-pointer bg-slate-100 dark:bg-slate-800 shadow-sm">
                                <img v-if="item.thumbnail" :src="item.thumbnail" :alt="item.judul"
                                    class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-700" />
                                <div v-else
                                    class="absolute inset-0 flex items-center justify-center bg-linear-to-br from-green-300 to-emerald-500">
                                    <UIcon :name="item.tipe === 'video' ? 'i-tabler-video' : 'i-tabler-photo'"
                                        class="w-16 h-16 text-white/40" />
                                </div>
                                <div
                                    class="absolute inset-0 bg-linear-to-t from-black/90 via-black/40 to-transparent opacity-80" />
                                <div class="absolute top-3 right-3 z-10">
                                    <UBadge :color="item.tipe === 'video' ? 'error' : 'neutral'" variant="solid"
                                        size="sm">
                                        {{ item.tipe.toUpperCase() }}
                                    </UBadge>
                                </div>
                                <div class="absolute bottom-0 left-0 right-0 p-5 z-10">
                                    <h3 class="text-white! font-bold text-lg leading-tight line-clamp-2">{{ item.judul
                                    }}</h3>
                                    <p class="text-white/80 text-xs mt-2 flex items-center gap-1.5">
                                        <UIcon name="i-tabler-calendar" class="w-3.5 h-3.5" /> {{ item.tanggal }}
                                    </p>
                                </div>
                            </NuxtLink>
                        </div>
                        <div v-else class="text-center py-10">
                            <UIcon name="i-tabler-photo-off" class="w-12 h-12 text-slate-300 mx-auto mb-2" />
                            <p class="text-slate-500">Belum ada galeri kegiatan.</p>
                        </div>
                        <div v-if="galeri.length > galeriPerPage" class="flex justify-end mt-8">
                            <UPagination v-model:page="halamanGaleri" :total="galeri.length"
                                :items-per-page="galeriPerPage" />
                        </div>
                    </div>
                </section>
            </UContainer>
        </div>
    </div>
</template>
