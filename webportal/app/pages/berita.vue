<script setup>
import beritaData from '~/data/berita-lengkap.json'
import galeriData from '~/data/galeri.json'

const berita = beritaData
const galeri = galeriData

const kategori = ['Semua', 'Pemerintahan', 'Ekonomi', 'Budaya', 'Pembangunan', 'Lingkungan']
const kategoriAktif = ref('Semua')

const beritaFiltered = computed(() => {
    if (kategoriAktif.value === 'Semua') return berita
    return berita.filter(item => item.kategori === kategoriAktif.value)
})
</script>

<template>
    <div>
        <UPageHero title="Informasi & Berita"
            description="Berita, informasi resmi, dan dokumentasi kegiatan Pemkab Tulungagung." />

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
                    <h3 class="font-bold text-lg">{{ item.judul }}</h3>
                    <p class="text-sm text-gray-500 mt-2 line-clamp-2">{{ item.deskripsi }}</p>
                    <template #footer>
                        <div class="flex items-center justify-between">
                            <p class="text-xs text-gray-400 flex items-center gap-1">
                                <UIcon name="i-mdi-calendar-blank" />
                                {{ item.tanggal }}
                            </p>
                            <UButton label="Baca" variant="link" size="sm" trailing-icon="i-tabler-arrow-right" />
                        </div>
                    </template>
                </UCard>
            </div>

            <!-- Pagination -->
            <div class="flex justify-center mt-10">
                <UPagination :total="24" :page-count="6" :page="1" />
            </div>
        </UPageSection>

        <!-- Section Galeri -->
        <UPageSection title="Galeri Kegiatan" description="Dokumentasi foto dan video kegiatan pemerintah daerah.">
            <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
                <UCard v-for="item in galeri" :key="item.id"
                    class="overflow-hidden group cursor-pointer hover:ring-2 hover:ring-primary hover:-translate-y-1 hover:shadow-xl transition-all duration-300"
                    :ui="{ header: 'p-0', body: 'p-4' }">
                    <template #header>
                        <!-- Thumbnail Container -->
                        <div
                            class="relative h-48 m-3 rounded-lg overflow-hidden bg-linear-to-br from-sky-400 to-sky-600 dark:from-sky-700 dark:to-sky-900">
                            <!-- Ikon besar di tengah -->
                            <div class="absolute inset-0 flex items-center justify-center">
                                <UIcon :name="item.tipe === 'video' ? 'i-tabler-video' : 'i-tabler-photo'"
                                    class="w-20 h-20 text-white/40 group-hover:scale-110 group-hover:text-white/60 transition-all duration-300" />
                            </div>

                            <!-- Badge tipe di pojok kanan atas -->
                            <div class="absolute top-3 right-3">
                                <UBadge :color="item.tipe === 'video' ? 'error' : 'neutral'" variant="solid" size="sm"
                                    class="font-semibold">
                                    {{ item.tipe.toUpperCase() }}
                                </UBadge>
                            </div>

                            <!-- Play button kalau video -->
                            <div v-if="item.tipe === 'video'"
                                class="absolute bottom-3 left-3 w-10 h-10 rounded-full bg-white/20 backdrop-blur-sm flex items-center justify-center group-hover:bg-white/30 transition-all">
                                <UIcon name="i-tabler-player-play-filled" class="w-5 h-5 text-white" />
                            </div>
                        </div>
                    </template>

                    <!-- Body: Judul + Tanggal -->
                    <div>
                        <h3 class="font-bold text-base line-clamp-2 group-hover:text-primary transition-colors">
                            {{ item.judul }}
                        </h3>
                        <div class="flex items-center gap-1.5 text-xs text-gray-400 mt-2">
                            <UIcon name="i-tabler-calendar" class="w-3.5 h-3.5" />
                            <span>{{ item.tanggal }}</span>
                        </div>
                    </div>
                </UCard>
            </div>

            <div class="flex justify-center mt-10">
                <UButton label="Lihat Semua Galeri" color="primary" variant="soft"
                    trailing-icon="i-tabler-arrow-right" />
            </div>
        </UPageSection>
    </div>
</template>