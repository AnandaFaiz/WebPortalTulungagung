<script setup>
// ==== Fetch OPD dari API ====
const { data: opdData, pending: opdLoading } = await useFetch('http://localhost:8080/api/opd')
const { data: heroData } = await useFetch('http://localhost:8080/api/banner?tipe=pemerintahan')
const heroImages = computed(() => heroData.value?.banner?.map(b => b.gambar) || [])

</script>

<template>
    <div>
        <div>
            <AnimatedHero badge="Pemerintahan" title="Pemerintahan Kabupaten Tulungagung"
                description="Informasi OPD, PPID, JDIH, dan DPRD Kabupaten Tulungagung."
                :images="heroImages.length > 0 ? heroImages : ['/hero/pemerintahan-1.webp']" :interval="5000" />
        </div>

        <div class="relative z-10 bg-transparent">
            <!-- Section OPD -->
            <UPageSection class="profil-text-white" title="Perangkat Daerah (OPD)"
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


        </div>
    </div>
</template>