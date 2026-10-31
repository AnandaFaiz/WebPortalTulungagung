<script setup>
// ==== Fetch OPD dari API ====
const { data: opdData, pending: opdLoading } = await useFetch('http://localhost:8080/api/opd')
const { data: tautanData, pending: tautanLoading } = await useFetch('http://localhost:8080/api/tautan-eksternal')
const { data: heroData } = await useFetch('http://localhost:8080/api/banner?tipe=pemerintahan')
const heroImages = computed(() => heroData.value?.banner?.map(b => b.gambar) || [])

// Data mock untuk bagian lain (PPID, JDIH, DPRD) yang belum di-CRUD-kan
const data = {
    ppid: {
        deskripsi: 'PPID Kabupaten Tulungagung menyediakan layanan informasi publik sesuai UU No. 14 Tahun 2008 tentang Keterbukaan Informasi Publik.',
        layanan: [
            'Permohonan Informasi Publik',
            'Pengajuan Keberatan Informasi',
            'Sengketa Informasi',
            'Daftar Informasi Publik'
        ]
    },
    jdih: [
        { id: 1, jenis: 'Perda', nomor: '1 Tahun 2025', tentang: 'Retribusi Pelayanan Persampahan', tanggal: '2025-03-15' },
        { id: 2, jenis: 'Perbup', nomor: '12 Tahun 2025', tentang: 'Standar Operasional Prosedur Pelayanan Publik', tanggal: '2025-05-20' },
        { id: 3, jenis: 'SK Bupati', nomor: '188/45/2025', tentang: 'Penetapan Penerima Bantuan Sosial', tanggal: '2025-06-10' }
    ],
    dprd: {
        jumlah_anggota: 50,
        periode: '2024-2029',
        komisi: [
            { nama: 'Komisi A', bidang: 'Pemerintahan & Hukum' },
            { nama: 'Komisi B', bidang: 'Ekonomi & Keuangan' },
            { nama: 'Komisi C', bidang: 'Pembangunan' },
            { nama: 'Komisi D', bidang: 'Kesejahteraan Rakyat' }
        ]
    }
}

// Helper untuk inisial OPD
const getInitials = (nama) => {
    if (!nama) return '?'
    return nama
        .split(' ')
        .filter(w => !['dinas', 'badan', 'sekretariat', 'kantor'].includes(w.toLowerCase()))
        .slice(0, 2)
        .map(w => w.charAt(0).toUpperCase())
        .join('')
}

// Helper untuk domain dari URL
const extractDomain = (url) => {
    if (!url) return null
    try {
        return new URL(url).hostname.replace('www.', '')
    } catch {
        return url
    }
}
</script>

<template>
    <div>
        <div class="sticky top-0 z-0">
            <AnimatedHero badge="Pemerintahan" title="Pemerintahan Kabupaten Tulungagung"
                description="Informasi OPD, PPID, JDIH, dan DPRD Kabupaten Tulungagung."
                :images="heroImages.length > 0 ? heroImages : ['/hero/pemerintahan-1.webp']" :interval="5000" />
        </div>

        <div class="relative z-10 bg-white dark:bg-slate-950">
            <!-- Section OPD -->
            <UPageSection title="Perangkat Daerah (OPD)"
                description="Direktori Dinas, Badan, dan Sekretariat di lingkungan Pemkab Tulungagung.">
                <!-- Loading -->
                <div v-if="opdLoading" class="text-center py-12">
                    <UIcon name="i-tabler-loader-2" class="w-10 h-10 text-sky-500 animate-spin mx-auto" />
                    <p class="text-sm text-slate-500 mt-3">Memuat data OPD...</p>
                </div>

                <!-- Empty State -->
                <div v-else-if="!opdData?.opds || opdData.opds.length === 0" class="text-center py-12">
                    <UIcon name="i-tabler-building-off" class="w-16 h-16 text-slate-300 mx-auto mb-3" />
                    <p class="text-slate-500 font-medium">Belum ada data OPD</p>
                    <p class="text-xs text-slate-400 mt-1">Data OPD akan muncul di sini setelah ditambahkan oleh admin.
                    </p>
                </div>

                <!-- Grid OPD -->
                <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
                    <component :is="item.website ? 'a' : 'div'" v-for="item in opdData.opds" :key="item.id"
                        :href="item.website || undefined" :target="item.website ? '_blank' : undefined"
                        :rel="item.website ? 'noopener noreferrer' : undefined"
                        class="group relative block h-44 bg-sky-50 dark:bg-slate-800/50 rounded-2xl border border-sky-100 dark:border-slate-700 overflow-hidden transition-all duration-300"
                        :class="item.website
                            ? 'cursor-pointer hover:bg-sky-100 hover:border-sky-300 hover:shadow-lg hover:-translate-y-1'
                            : 'cursor-default'">
                        <!-- Logo Kabupaten Tulungagung -->
                        <img src="/Logo-Kabupaten-Tulungagung.png" alt="Logo Kabupaten Tulungagung" class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 h-14 w-auto shrink-0
                   transition-all duration-500 ease-out
                   group-hover:top-5 group-hover:left-5 group-hover:h-12
                   group-hover:translate-x-0 group-hover:translate-y-0" />

                        <!-- Nama OPD -->
                        <div class="absolute top-[68%] left-0 right-0 px-4 text-center
                   transition-all duration-500 ease-out
                   group-hover:top-5 group-hover:left-21 group-hover:right-4 group-hover:px-0 group-hover:text-left">
                            <h3 class="font-bold text-slate-900 dark:text-white leading-tight line-clamp-2
                       text-sm
                       transition-all duration-500 ease-out
                       group-hover:text-base">
                                {{ item.nama }}
                            </h3>
                        </div>

                        <!-- Deskripsi (muncul saat hover, align dengan nama) -->
                        <div class="absolute top-16 left-21 right-4
                   opacity-0 translate-y-3 pointer-events-none
                   transition-all duration-500 ease-out delay-150
                   group-hover:opacity-100 group-hover:translate-y-0">
                            <p class="text-xs text-slate-500 dark:text-slate-400 leading-relaxed line-clamp-3">
                                {{ item.deskripsi || 'Deskripsi belum tersedia.' }}
                            </p>
                        </div>

                        <!-- Hint website (bottom right) -->
                        <div v-if="item.website" class="absolute bottom-3 right-4 flex items-center gap-1 text-xs text-sky-600 dark:text-sky-400 font-medium
                   opacity-0 translate-y-2
                   transition-all duration-500 ease-out delay-250
                   group-hover:opacity-100 group-hover:translate-y-0">
                            <UIcon name="i-tabler-external-link" class="w-3 h-3" />
                            <span>Kunjungi Website</span>
                        </div>
                    </component>
                </div>
            </UPageSection>

            <!-- Section Tautan Eksternal -->
            <UPageSection title="Portal Terkait" description="Akses cepat ke portal resmi pemerintah lainnya.">
                <div v-if="tautanLoading" class="text-center py-12">
                    <UIcon name="i-tabler-loader-2" class="w-10 h-10 text-sky-500 animate-spin mx-auto" />
                </div>

                <div v-else-if="!tautanData?.tautan || tautanData.tautan.length === 0" class="text-center py-12">
                    <UIcon name="i-tabler-link-off" class="w-16 h-16 text-slate-300 mx-auto mb-3" />
                    <p class="text-slate-500 font-medium">Belum ada tautan eksternal</p>
                </div>

                <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    <a v-for="item in tautanData.tautan" :key="item.id" :href="item.url" target="_blank" rel="noopener"
                        class="group block">
                        <UCard
                            class="h-full hover:ring-2 hover:ring-sky-400 hover:shadow-lg hover:-translate-y-1 transition-all">
                            <div class="flex items-start gap-4">
                                <div
                                    class="w-12 h-12 rounded-xl bg-sky-100 dark:bg-sky-950 flex items-center justify-center shrink-0 group-hover:bg-sky-500 transition-colors">
                                    <UIcon :name="item.icon || 'i-tabler-external-link'"
                                        class="w-6 h-6 text-sky-600 dark:text-sky-400 group-hover:text-white transition-colors" />
                                </div>
                                <div class="flex-1 min-w-0">
                                    <h3
                                        class="font-bold text-slate-900 dark:text-white group-hover:text-sky-600 transition-colors">
                                        {{ item.nama }}
                                    </h3>
                                    <p v-if="item.deskripsi" class="text-sm text-slate-500 mt-1 line-clamp-2">
                                        {{ item.deskripsi }}
                                    </p>
                                    <span class="inline-flex items-center gap-1 text-xs text-sky-600 font-medium mt-2">
                                        Kunjungi
                                        <UIcon name="i-tabler-external-link" class="w-3 h-3" />
                                    </span>
                                </div>
                            </div>
                        </UCard>
                    </a>
                </div>
            </UPageSection>
        </div>
    </div>
</template>