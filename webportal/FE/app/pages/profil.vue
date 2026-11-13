<script setup>
// Fetch Profil Daerah
const { data: profilDaerahData, pending: profilDaerahLoading } = await useFetch('http://localhost:8080/api/profil-daerah')
const profil = computed(() => profilDaerahData.value?.profil || {})

// Fetch Pejabat (sudah ada)
const { data: pejabatData, pending: pejabatLoading } = await useFetch('http://localhost:8080/api/pejabat')

// Fetch Struktur (sudah ada)
const { data: strukturData, pending: strukturLoading } = await useFetch('http://localhost:8080/api/struktur')

// Helper inisial (sudah ada)
const getInitials = (nama) => {
    if (!nama) return '?'
    return nama.split(' ').filter(w => !['h', 'hj', 'drs', 'dr', 'ir', 'prof'].includes(w.toLowerCase())).slice(0, 2).map(w => w.charAt(0).toUpperCase()).join('')
}

// Helper kategori (sudah ada)
const getKategoriColor = (kategori) => {
    const map = { 'Kepala Daerah': 'primary', 'Sekda': 'warning', 'Asisten': 'info', 'Staf Ahli': 'success', 'Kepala OPD': 'neutral', 'Camat': 'neutral' }
    return map[kategori] || 'neutral'
}
</script>

<template>
    <div>
        <div>
            <AnimatedHero badge="Profil Daerah" title="Profil Kabupaten Tulungagung"
                description="Mengenal lebih dekat sejarah, visi misi, dan jajaran pemerintahan Kabupaten Tulungagung."
                :images="[
                    '/hero/profil-1.webp',
                    '/hero/profil-2.webp',
                    '/hero/profil-3.webp'
                ]" :interval="5000" />
        </div>
        <div class="relative z-10 bg-transparent">

            <!-- Section Sejarah & Lambang -->
            <UPageSection
                class="profil-text-white [&>div]:py-8! [&>div]:sm:py-10! [&>div]:lg:py-12! [&>div]:gap-5! [&>div]:sm:gap-6!"
                title="Sejarah & Lambang" description="Jejak panjang perjalanan Kabupaten Tulungagung.">
                <div v-if="profilDaerahLoading" class="text-center py-12">
                    <UIcon name="i-tabler-loader-2" class="w-10 h-10 text-sky-500 animate-spin mx-auto" />
                </div>
                <div v-else class="grid grid-cols-1 gap-5 md:grid-cols-2 md:gap-6">
                    <UCard>
                        <h3 class="font-bold text-lg mb-3">Sejarah Singkat</h3>
                        <p class="text-slate-600 dark:text-slate-400 leading-relaxed whitespace-pre-line">{{
                            profil.sejarah || '-' }}</p>
                    </UCard>
                    <UCard>
                        <h3 class="font-bold text-lg mb-3">Makna Lambang</h3>
                        <p class="text-slate-600 dark:text-slate-400 leading-relaxed whitespace-pre-line">{{
                            profil.lambang || '-' }}</p>
                    </UCard>
                </div>
            </UPageSection>

            <!-- Section Visi & Misi -->
            <UPageSection
                class="profil-text-white [&>div]:py-8! [&>div]:sm:py-10! [&>div]:lg:py-12! [&>div]:gap-5! [&>div]:sm:gap-6!"
                title="Visi & Misi" description="Arah pembangunan Kabupaten Tulungagung.">
                <div class="grid grid-cols-1 items-start gap-5 md:grid-cols-2 md:gap-6">
                    <UCard class="h-fit">
                        <h3 class="mb-3 text-lg font-bold">Visi</h3>
                        <p class="text-xl font-medium text-primary">{{ profil.visi || '-' }}</p>
                    </UCard>

                    <div class="rounded-xl border border-slate-200 bg-white p-5 shadow-md sm:p-6">
                        <h3 class="mb-4 text-lg font-bold">Misi</h3>
                        <div v-if="profil.misi && profil.misi.length > 0" class="space-y-3">
                            <div v-for="(misi, index) in profil.misi" :key="index"
                                class="flex items-start gap-4 rounded-lg border border-slate-200 bg-slate-50 p-4">
                                <div
                                    class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-primary font-bold text-white">
                                    {{ index + 1 }}</div>
                                <p class="text-slate-700">{{ misi }}</p>
                            </div>
                        </div>
                        <p v-else class="italic text-slate-500">Belum ada data misi.</p>
                    </div>
                </div>
            </UPageSection>

            <!-- Section Struktur Organisasi -->
            <UPageSection
                class="profil-text-white [&>div]:py-8! [&>div]:sm:py-10! [&>div]:lg:py-12! [&>div]:gap-5! [&>div]:sm:gap-6!"
                title="Struktur Organisasi"
                description="Bagan resmi struktur organisasi Pemerintah Kabupaten Tulungagung.">
                <!-- Loading -->
                <div v-if="strukturLoading" class="text-center py-12">
                    <UIcon name="i-tabler-loader-2" class="w-10 h-10 text-sky-500 animate-spin mx-auto" />
                </div>

                <!-- Empty -->
                <div v-else-if="!strukturData?.struktur?.gambar"
                    class="text-center py-12 bg-slate-50 dark:bg-slate-800/50 rounded-2xl">
                    <UIcon name="i-tabler-sitemap" class="w-16 h-16 text-slate-300 mx-auto mb-3" />
                    <p class="text-slate-500 font-medium">Bagan struktur belum tersedia</p>
                    <p class="text-xs text-slate-400 mt-1">Bagan akan ditampilkan setelah admin mengunggahnya.</p>
                </div>

                <!-- Gambar Bagan -->
                <div v-else>
                    <p v-if="strukturData.struktur.deskripsi" class="text-center text-slate-500 max-w-2xl mx-auto mb-8">
                        {{ strukturData.struktur.deskripsi }}
                    </p>
                    <div
                        class="bg-white dark:bg-slate-800 rounded-2xl p-4 shadow-lg border border-slate-200 dark:border-slate-700 overflow-hidden">
                        <img :src="strukturData.struktur.gambar" :alt="strukturData.struktur.judul"
                            class="w-full h-auto rounded-lg" />
                    </div>
                </div>
            </UPageSection>

            <!-- Section Pejabat Daerah -->
            <UPageSection
                class="profil-text-white [&>div]:py-8! [&>div]:sm:py-10! [&>div]:lg:py-12! [&>div]:gap-5! [&>div]:sm:gap-6!"
                title="Pejabat Daerah"
                description="Bupati, Wakil Bupati, dan jajaran Pemerintah Kabupaten Tulungagung.">
                <!-- Loading -->
                <div v-if="pejabatLoading" class="text-center py-12">
                    <UIcon name="i-tabler-loader-2" class="w-10 h-10 text-sky-500 animate-spin mx-auto" />
                    <p class="text-sm text-slate-500 mt-3">Memuat data pejabat...</p>
                </div>

                <!-- Empty State -->
                <div v-else-if="!pejabatData?.pejabat || pejabatData.pejabat.length === 0" class="text-center py-12">
                    <UIcon name="i-tabler-user-off" class="w-16 h-16 text-slate-300 mx-auto mb-3" />
                    <p class="text-slate-500 font-medium">Belum ada data pejabat</p>
                    <p class="text-xs text-slate-400 mt-1">Data pejabat akan muncul di sini setelah ditambahkan oleh
                        admin.
                    </p>
                </div>

                <!-- Grid Pejabat -->
                <div v-else class="flex flex-wrap justify-center gap-6">
                    <UCard v-for="item in pejabatData.pejabat" :key="item.id"
                        class="text-center hover:ring-2 hover:ring-sky-400 hover:shadow-lg transition-all group w-40 sm:w-50">
                        <!-- Foto / Inisial -->
                        <div
                            class="w-32 h-32 mx-auto rounded-full bg-sky-100 dark:bg-sky-950 flex items-center justify-center overflow-hidden mb-4 ring-4 ring-white dark:ring-slate-800 shadow-md group-hover:ring-sky-200 dark:group-hover:ring-sky-900 transition-all">
                            <img v-if="item.foto" :src="item.foto" :alt="item.nama"
                                class="w-full h-full object-cover" />
                            <span v-else class="text-sky-600 dark:text-sky-400 font-bold text-3xl">
                                {{ getInitials(item.nama) }}
                            </span>
                        </div>

                        <!-- Kategori Badge -->
                        <div v-if="item.kategori" class="mb-2">
                            <UBadge :color="getKategoriColor(item.kategori)" variant="subtle" size="xs">
                                {{ item.kategori }}
                            </UBadge>
                        </div>

                        <!-- Nama & Jabatan -->
                        <h3 class="font-bold text-slate-900 dark:text-white leading-tight line-clamp-2">
                            {{ item.nama }}
                        </h3>
                        <p class="text-sky-600 dark:text-sky-400 text-sm mt-1 line-clamp-2">
                            {{ item.jabatan }}
                        </p>

                        <!-- Periode -->
                        <p v-if="item.periode" class="text-xs text-slate-400 mt-2">
                            {{ item.periode }}
                        </p>
                    </UCard>
                </div>
            </UPageSection>

            <!-- Section Geografis & Demografi -->
            <UPageSection
                class="profil-text-white [&>div]:py-8! [&>div]:sm:py-10! [&>div]:lg:py-12! [&>div]:gap-5! [&>div]:sm:gap-6!"
                title="Geografis & Demografi" description="Data wilayah dan kependudukan.">
                <div v-if="profilDaerahLoading" class="text-center py-12">
                    <UIcon name="i-tabler-loader-2" class="w-10 h-10 text-sky-500 animate-spin mx-auto" />
                </div>
                <div v-else class="space-y-8">
                    <UCard>
                        <p class="text-slate-600 dark:text-slate-400 leading-relaxed whitespace-pre-line">{{
                            profil.geografis || '-' }}</p>
                    </UCard>

                    <div v-if="profil.demografi" class="grid grid-cols-2 md:grid-cols-4 gap-4">
                        <UCard class="text-center">
                            <p class="text-2xl font-bold text-primary">{{ profil.demografi.luas || '-' }}</p>
                            <p class="text-xs text-slate-500 mt-1">Luas Wilayah</p>
                        </UCard>
                        <UCard class="text-center">
                            <p class="text-2xl font-bold text-primary">{{ profil.demografi.penduduk || '-' }}</p>
                            <p class="text-xs text-slate-500 mt-1">Penduduk</p>
                        </UCard>
                        <UCard class="text-center">
                            <p class="text-2xl font-bold text-primary">{{ profil.demografi.kecamatan || '-' }}</p>
                            <p class="text-xs text-slate-500 mt-1">Kecamatan</p>
                        </UCard>
                        <UCard class="text-center">
                            <p class="text-2xl font-bold text-primary">{{ profil.demografi.desa || '-' }}</p>
                            <p class="text-xs text-slate-500 mt-1">Desa/Kelurahan</p>
                        </UCard>
                    </div>
                </div>
            </UPageSection>
        </div>
    </div>
</template>