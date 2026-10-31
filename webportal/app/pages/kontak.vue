<script setup>
// ==== Fetch Kontak dari API ====
const { data: kontakData, pending: kontakLoading } = await useFetch('http://localhost:8080/api/kontak')
const { data: heroData } = await useFetch('http://localhost:8080/api/banner?tipe=kontak')
const heroImages = computed(() => heroData.value?.banner?.map(b => b.gambar) || [])

// ==== Helper: Parse embed maps ====
const mapsSrc = computed(() => {
    const raw = kontakData.value?.kontak?.embed_maps || ''
    if (!raw) return ''
    const match = raw.match(/src=["']([^"']+)["']/)
    if (match) return match[1]
    return raw
})

// ==== Medsos links (filter yang kosong) ====
const socialLinks = computed(() => {
    const k = kontakData.value?.kontak || {}
    return [
        { name: 'Facebook', icon: 'i-mdi-facebook', url: k.facebook },
        { name: 'Instagram', icon: 'i-mdi-instagram', url: k.instagram },
        { name: 'Twitter', icon: 'i-mdi-twitter', url: k.twitter },
        { name: 'YouTube', icon: 'i-mdi-youtube', url: k.youtube }
    ].filter(s => s.url)
})
</script>

<template>
    <div>
        <!-- Hero Sticky -->
        <div class="sticky top-0 z-0">
            <AnimatedHero badge="Kontak" title="Hubungi Kami"
                description="Sampaikan pertanyaan, saran, atau aspirasi Anda kepada Pemkab Tulungagung."
                :images="heroImages.length > 0 ? heroImages : ['/hero/kontak-1.webp']" />
        </div>

        <div class="relative z-10 bg-white dark:bg-slate-950">

            <!-- Section Kontak -->
            <UPageSection title="Informasi Kontak" description="Berbagai kanal komunikasi resmi Pemkab Tulungagung.">
                <div v-if="kontakLoading" class="text-center py-12">
                    <UIcon name="i-tabler-loader-2" class="w-10 h-10 text-sky-500 animate-spin mx-auto" />
                </div>

                <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <UCard class="hover:ring-2 hover:ring-sky-400 transition-all">
                        <div class="flex items-start gap-4">
                            <div
                                class="w-12 h-12 rounded-xl bg-sky-100 dark:bg-sky-950 flex items-center justify-center shrink-0">
                                <UIcon name="i-tabler-map-pin" class="w-6 h-6 text-sky-500" />
                            </div>
                            <div class="flex-1">
                                <p class="text-xs uppercase text-slate-400 tracking-wide font-semibold">Alamat Kantor
                                </p>
                                <p class="mt-1 font-medium text-slate-900 dark:text-white">{{ kontakData?.kontak?.alamat
                                    || '-' }}</p>
                            </div>
                        </div>
                    </UCard>

                    <UCard class="hover:ring-2 hover:ring-sky-400 transition-all">
                        <div class="flex items-start gap-4">
                            <div
                                class="w-12 h-12 rounded-xl bg-sky-100 dark:bg-sky-950 flex items-center justify-center shrink-0">
                                <UIcon name="i-tabler-phone-call" class="w-6 h-6 text-sky-500" />
                            </div>
                            <div class="flex-1">
                                <p class="text-xs uppercase text-slate-400 tracking-wide font-semibold">Telepon</p>
                                <a :href="`tel:${kontakData?.kontak?.telepon}`"
                                    class="mt-1 font-medium text-slate-900 dark:text-white block hover:text-sky-600 transition-colors">
                                    {{ kontakData?.kontak?.telepon || '-' }}
                                </a>
                            </div>
                        </div>
                    </UCard>

                    <UCard class="hover:ring-2 hover:ring-sky-400 transition-all">
                        <div class="flex items-start gap-4">
                            <div
                                class="w-12 h-12 rounded-xl bg-sky-100 dark:bg-sky-950 flex items-center justify-center shrink-0">
                                <UIcon name="i-tabler-mail" class="w-6 h-6 text-sky-500" />
                            </div>
                            <div class="flex-1">
                                <p class="text-xs uppercase text-slate-400 tracking-wide font-semibold">Email</p>
                                <a :href="`mailto:${kontakData?.kontak?.email}`"
                                    class="mt-1 font-medium text-slate-900 dark:text-white block hover:text-sky-600 transition-colors">
                                    {{ kontakData?.kontak?.email || '-' }}
                                </a>
                            </div>
                        </div>
                    </UCard>

                    <UCard class="hover:ring-2 hover:ring-sky-400 transition-all">
                        <div class="flex items-start gap-4">
                            <div
                                class="w-12 h-12 rounded-xl bg-sky-100 dark:bg-sky-950 flex items-center justify-center shrink-0">
                                <UIcon name="i-tabler-clock" class="w-6 h-6 text-sky-500" />
                            </div>
                            <div class="flex-1">
                                <p class="text-xs uppercase text-slate-400 tracking-wide font-semibold">Jam Operasional
                                </p>
                                <p class="mt-1 font-medium text-slate-900 dark:text-white">{{
                                    kontakData?.kontak?.jam_operasional || '-' }}</p>
                            </div>
                        </div>
                    </UCard>
                </div>

                <!-- Medsos -->
                <div v-if="socialLinks.length > 0" class="mt-8 flex flex-col items-center">
                    <p class="text-sm font-semibold text-slate-500 mb-3">Ikuti Media Sosial Kami</p>
                    <div class="flex items-center gap-3">
                        <a v-for="sos in socialLinks" :key="sos.name" :href="sos.url" target="_blank" rel="noopener"
                            :aria-label="sos.name"
                            class="w-11 h-11 rounded-full bg-sky-100 dark:bg-sky-950 flex items-center justify-center text-sky-600 hover:bg-sky-500 hover:text-white transition-all hover:scale-110">
                            <UIcon :name="sos.icon" class="w-5 h-5" />
                        </a>
                    </div>
                </div>
            </UPageSection>

            <!-- Section Peta -->
            <UPageSection title="Lokasi Kantor" description="Kantor Bupati Tulungagung.">
                <UCard v-if="mapsSrc">
                    <div class="rounded-lg overflow-hidden">
                        <iframe :src="mapsSrc" width="100%" height="450" style="border:0" allowfullscreen loading="lazy"
                            referrerpolicy="no-referrer-when-downgrade" />
                    </div>
                </UCard>

                <UCard v-else>
                    <div class="h-96 bg-slate-100 dark:bg-slate-800 rounded-lg flex items-center justify-center">
                        <div class="text-center">
                            <UIcon name="i-tabler-map-pin" class="w-16 h-16 text-sky-400 mx-auto" />
                            <p class="mt-4 text-slate-500">Peta lokasi belum tersedia</p>
                            <p class="text-xs text-slate-400 mt-1">Admin belum mengatur embed Google Maps.</p>
                        </div>
                    </div>
                </UCard>
            </UPageSection>

            <!-- Section CTA Pengaduan -->
            <UPageCTA title="Punya Keluhan atau Aspirasi?"
                description="Sampaikan melalui kanal pengaduan resmi Pemkab Tulungagung." :links="[
                    { label: 'Buat Pengaduan', color: 'primary', size: 'lg', to: '/layanan' },
                    { label: 'Lihat FAQ', color: 'neutral', variant: 'outline', size: 'lg', to: '#' }
                ]" />
        </div>
    </div>
</template>