<script setup>
const { data: heroData } = await useFetch('http://localhost:8080/api/banner?tipe=berita')
const heroImages = computed(() => heroData.value?.banner?.map(b => b.gambar) || [])
// ==== Fetch Berita (Golang API) ====
const { data: beritaData, pending: beritaLoading } = await useFetch('http://localhost:8080/api/berita')

// ✅ PERBAIKAN: Tambahkan || []
const berita = computed(() => beritaData.value?.berita || [])

// Fetch Kategori dari API
const { data: kategoriData } = await useFetch('http://localhost:8080/api/kategori-berita')

// Ambil nama kategorinya saja, dan tambahkan "Semua" di paling depan
const kategori = computed(() => {
    const list = kategoriData.value?.kategori?.map(k => k.nama) || []
    return ['Semua', ...list]
})

const kategoriAktif = ref('Semua')

const beritaFiltered = computed(() => {
    if (kategoriAktif.value === 'Semua') return berita.value
    return berita.value.filter(item => item.kategori === kategoriAktif.value)
})
</script>

<template>
    <div>
        <div class="sticky top-0 z-0">
            <!-- Hero -->
            <AnimatedHero badge="Kabar Terkini" title="Informasi & Berita"
                description="Berita, informasi resmi, dan dokumentasi kegiatan Pemkab Tulungagung."
                :images="heroImages.length > 0 ? heroImages : ['/hero/berita-1.webp']" :interval="5000" />
        </div>

        <div class="relative z-10 bg-white dark:bg-slate-950">
            <!-- Section Berita -->
            <UPageSection title="Berita Terbaru" description="Kabar terkini dari Pemerintah Kabupaten Tulungagung.">
                <!-- Filter Kategori -->
                <div class="flex flex-wrap gap-2 mb-8">
                    <UButton v-for="kat in kategori" :key="kat" :label="kat"
                        :color="kategoriAktif === kat ? 'primary' : 'neutral'"
                        :variant="kategoriAktif === kat ? 'solid' : 'outline'" size="sm" @click="kategoriAktif = kat" />
                </div>

                <!-- Grid Berita -->
                <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    <UCard v-for="item in beritaFiltered" :key="item.id">
                        <template #header>
                            <div class="h-48 bg-sky-100 dark:bg-sky-950 rounded-lg flex items-center justify-center">
                                <UIcon name="i-tabler-news" class="w-16 h-16 text-sky-500 dark:text-sky-400" />
                            </div>
                        </template>
                        <UBadge color="primary" variant="subtle" size="sm" class="mb-2">
                            {{ item.kategori }}
                        </UBadge>
                        <h3 class="font-bold text-lg text-slate-900 dark:text-white">{{ item.judul }}</h3>
                        <p class="text-sm text-slate-500 mt-2 line-clamp-2">{{ item.deskripsi }}</p>
                        <template #footer>
                            <div class="flex items-center justify-between">
                                <p class="text-xs text-slate-400 flex items-center gap-1">
                                    <UIcon name="i-tabler-calendar" />
                                    {{ item.tanggal }}
                                </p>
                                <UButton label="Baca" variant="link" size="sm" trailing-icon="i-tabler-arrow-right" />
                            </div>
                        </template>
                    </UCard>
                </div>

                <!-- Empty State -->
                <div v-if="beritaFiltered.length === 0" class="text-center py-12">
                    <UIcon name="i-tabler-news-off" class="w-16 h-16 text-slate-300 mx-auto mb-3" />
                    <p class="text-slate-500 font-medium">Tidak ada berita di kategori ini</p>
                </div>
            </UPageSection>
        </div>
    </div>
</template>