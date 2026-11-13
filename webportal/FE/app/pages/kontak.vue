<script setup>
// ==== Fetch Kontak dari API ====
const { data: kontakData, pending: kontakLoading } = await useFetch('http://localhost:8080/api/kontak')
const { data: heroData } = await useFetch('http://localhost:8080/api/banner?tipe=kontak')
const { data: pengaduanConfigData } = await useFetch('http://localhost:8080/api/pengaduan-konfigurasi')
const { data: daruratData } = await useFetch('http://localhost:8080/api/layanan-darurat')
const heroImages = computed(() => heroData.value?.banner?.map(b => b.gambar) || [])
const pengaduanKategori = computed(() => pengaduanConfigData.value?.kategori || [])
const daruratList = computed(() => daruratData.value?.darurat || [])
const socialLinks = computed(() => {
    const kontak = kontakData.value?.kontak || {}
    return [
        { key: 'facebook', label: 'Facebook', icon: 'i-mdi-facebook', url: kontak.facebook, class: 'bg-[#1877F2]' },
        { key: 'instagram', label: 'Instagram', icon: 'i-mdi-instagram', url: kontak.instagram, class: 'bg-gradient-to-br from-[#833AB4] via-[#E1306C] to-[#FCAF45]' },
        { key: 'twitter', label: 'X', icon: 'i-tabler-brand-x', url: kontak.twitter, class: 'bg-black' },
        { key: 'youtube', label: 'YouTube', icon: 'i-mdi-youtube', url: kontak.youtube, class: 'bg-[#FF0000]' }
    ].filter(link => link.url)
})

// ==== Helper: Parse embed maps ====
const mapsSrc = computed(() => {
    const raw = kontakData.value?.kontak?.embed_maps || ''
    if (!raw) return ''
    const match = raw.match(/src=["']([^"']+)["']/)
    if (match) return match[1]
    return raw
})

const formAduan = ref({ nama: '', kontak: '', kategori: 'Umum', isi_aduan: '' })
const submittingAduan = ref(false)
const tiketAduan = ref(null)
const errorAduan = ref('')

watch(pengaduanKategori, categories => {
    if (!categories.length) return
    if (!categories.some(category => category.nama === formAduan.value.kategori)) {
        formAduan.value.kategori = categories[0].nama
    }
}, { immediate: true })

const submitAduan = async () => {
    submittingAduan.value = true
    tiketAduan.value = null
    errorAduan.value = ''

    try {
        const response = await $fetch('http://localhost:8080/api/pengaduan', {
            method: 'POST',
            body: formAduan.value
        })
        tiketAduan.value = response.tiket_id
        formAduan.value = { nama: '', kontak: '', kategori: pengaduanKategori.value[0]?.nama || 'Umum', isi_aduan: '' }
    } catch (error) {
        errorAduan.value = error.data?.statusMessage || 'Aduan belum berhasil dikirim. Silakan coba kembali.'
    } finally {
        submittingAduan.value = false
    }
}

const kategoriAduanOptions = computed(() => pengaduanKategori.value.length
    ? pengaduanKategori.value.map(category => category.nama)
    : ['Umum'])
</script>

<template>
    <div>
        <!-- Hero Sticky -->
        <div>
            <AnimatedHero badge="Kontak" title="Hubungi Kami"
                description="Sampaikan pertanyaan, saran, atau aspirasi Anda kepada Pemkab Tulungagung."
                :images="heroImages.length > 0 ? heroImages : ['/hero/kontak-1.webp']" />
        </div>

        <div class="relative z-10 bg-transparent">
            <section class="px-4 pb-8 pt-6 sm:px-6 lg:pt-8">
                <UContainer>
                    <div class="mb-5 text-center md:mb-7">
                        <h2 class="text-3xl font-extrabold text-white sm:text-4xl">Informasi Kontak dan Lokasi</h2>
                        <p class="mt-2 text-sm text-white/85 sm:text-base">Berbagai kanal komunikasi dan lokasi resmi
                            Pemkab Tulungagung.</p>
                    </div>

                    <div class="space-y-6">
                        <div
                            class="rounded-2xl border border-white/70 bg-white/85 p-5 shadow-xl shadow-slate-950/10 backdrop-blur-md sm:p-7 lg:p-9">
                            <div
                                class="grid grid-cols-1 items-start gap-6 lg:grid-cols-[minmax(0,2fr)_minmax(0,3fr)] lg:gap-8">
                                <div class="space-y-3">
                                    <div v-if="kontakLoading" class="py-12 text-center">
                                        <UIcon name="i-tabler-loader-2"
                                            class="mx-auto h-10 w-10 animate-spin text-sky-500" />
                                    </div>
                                    <template v-else>
                                        <UCard class="transition-all hover:ring-2 hover:ring-sky-400">
                                            <div class="flex items-start gap-4">
                                                <div
                                                    class="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-sky-100">
                                                    <UIcon name="i-tabler-map-pin" class="h-5 w-5 text-sky-500" />
                                                </div>
                                                <div class="min-w-0 flex-1">
                                                    <p class="text-xs font-semibold uppercase text-slate-400">Alamat
                                                        Kantor
                                                    </p>
                                                    <p class="mt-1 font-medium text-slate-900">{{
                                                        kontakData?.kontak?.alamat
                                                        || '-' }}</p>
                                                </div>
                                            </div>
                                        </UCard>
                                        <UCard class="transition-all hover:ring-2 hover:ring-sky-400">
                                            <div class="flex items-start gap-4">
                                                <div
                                                    class="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-red-50">
                                                    <UIcon name="i-tabler-phone-call" class="h-5 w-5 text-red-500" />
                                                </div>
                                                <div class="min-w-0 flex-1">
                                                    <p class="text-xs font-semibold uppercase text-slate-400">Telepon
                                                    </p><a :href="`tel:${kontakData?.kontak?.telepon}`"
                                                        class="mt-1 block font-medium text-slate-900 transition-colors hover:text-sky-600">{{
                                                            kontakData?.kontak?.telepon || '-' }}</a>
                                                </div>
                                            </div>
                                        </UCard>
                                        <UCard class="transition-all hover:ring-2 hover:ring-sky-400">
                                            <div class="flex items-start gap-4">
                                                <div
                                                    class="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-amber-50">
                                                    <UIcon name="i-tabler-mail" class="h-5 w-5 text-amber-500" />
                                                </div>
                                                <div class="min-w-0 flex-1">
                                                    <p class="text-xs font-semibold uppercase text-slate-400">Email</p>
                                                    <a :href="`mailto:${kontakData?.kontak?.email}`"
                                                        class="mt-1 block break-all font-medium text-slate-900 transition-colors hover:text-sky-600">{{
                                                            kontakData?.kontak?.email || '-' }}</a>
                                                </div>
                                            </div>
                                        </UCard>
                                        <UCard class="transition-all hover:ring-2 hover:ring-sky-400">
                                            <div class="flex items-start gap-4">
                                                <div
                                                    class="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-teal-50">
                                                    <UIcon name="i-tabler-clock" class="h-5 w-5 text-teal-500" />
                                                </div>
                                                <div class="min-w-0 flex-1">
                                                    <p class="text-xs font-semibold uppercase text-slate-400">Jam
                                                        Operasional</p>
                                                    <p class="mt-1 font-medium text-slate-900">{{
                                                        kontakData?.kontak?.jam_operasional || '-' }}</p>
                                                </div>
                                            </div>
                                        </UCard>
                                    </template>
                                </div>

                                <div class="min-w-0">
                                    <UCard v-if="mapsSrc" class="overflow-hidden">
                                        <div class="h-72 overflow-hidden rounded-lg sm:h-80 lg:h-96">
                                            <iframe :src="mapsSrc" title="Peta lokasi Kantor Bupati Tulungagung"
                                                class="h-full w-full border-0" allowfullscreen loading="lazy"
                                                referrerpolicy="no-referrer-when-downgrade" />
                                        </div>
                                    </UCard>
                                    <UCard v-else>
                                        <div
                                            class="flex h-72 items-center justify-center rounded-lg bg-slate-100 sm:h-80 lg:h-96">
                                            <div class="text-center">
                                                <UIcon name="i-tabler-map-pin" class="mx-auto h-16 w-16 text-sky-400" />
                                                <p class="mt-4 text-slate-500">Peta lokasi belum tersedia</p>
                                                <p class="mt-1 text-xs text-slate-400">Admin belum mengatur embed Google
                                                    Maps.</p>
                                            </div>
                                        </div>
                                    </UCard>
                                </div>
                            </div>
                        </div>

                        <div
                            class="rounded-2xl border border-white/70 bg-white/85 p-5 shadow-xl shadow-slate-950/10 backdrop-blur-md sm:p-7 lg:p-9">
                            <div class="mb-6">
                                <h2 class="text-2xl font-extrabold text-slate-900 sm:text-3xl">Punya Keluhan atau
                                    Aspirasi?</h2>
                                <p class="mt-2 text-sm text-slate-600">Sampaikan melalui kanal pengaduan resmi Pemkab
                                    Tulungagung.</p>
                            </div>
                            <div
                                class="grid grid-cols-1 gap-8 lg:grid-cols-[minmax(0,0.8fr)_minmax(0,1.2fr)] lg:gap-10">
                                <aside class="space-y-6">
                                    <div>
                                        <div class="mb-3 flex items-center gap-2">
                                            <UIcon name="i-tabler-route" class="h-5 w-5 text-sky-600" />
                                            <h3 class="font-bold text-slate-900">Alur Pengaduan</h3>
                                        </div>
                                        <ol class="space-y-3 text-sm text-slate-600">
                                            <li class="flex gap-3"><span
                                                    class="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-sky-100 text-xs font-bold text-sky-700">1</span><span>Isi
                                                    identitas dan pilih kategori laporan.</span></li>
                                            <li class="flex gap-3"><span
                                                    class="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-sky-100 text-xs font-bold text-sky-700">2</span><span>Jelaskan
                                                    keluhan atau aspirasi secara jelas.</span></li>
                                            <li class="flex gap-3"><span
                                                    class="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-sky-100 text-xs font-bold text-sky-700">3</span><span>Simpan
                                                    nomor tiket untuk menanyakan tindak lanjut.</span></li>
                                        </ol>
                                    </div>
                                    <div class="rounded-xl border border-sky-100 bg-sky-50/80 p-4">
                                        <div class="flex items-start gap-3">
                                            <UIcon name="i-tabler-lock-check"
                                                class="mt-0.5 h-5 w-5 shrink-0 text-sky-700" />
                                            <div>
                                                <h3 class="font-bold text-slate-900">Privasi Pelapor</h3>
                                                <p class="mt-1 text-sm leading-relaxed text-slate-600">Informasi kontak
                                                    digunakan untuk
                                                    keperluan tindak lanjut pengaduan dan dijaga kerahasiaannya.</p>
                                            </div>
                                        </div>
                                    </div>
                                    <div>
                                        <div class="mb-3 flex items-center gap-2">
                                            <UIcon name="i-tabler-phone-call" class="h-5 w-5 text-red-600" />
                                            <h3 class="font-bold text-slate-900">Kontak Darurat</h3>
                                        </div>
                                        <div v-if="daruratList.length" class="space-y-2">
                                            <a v-for="item in daruratList" :key="item.id" :href="`tel:${item.nomor}`"
                                                class="flex items-center justify-between gap-3 rounded-lg border border-slate-200 bg-white px-3 py-2.5 transition hover:border-red-300 hover:bg-red-50">
                                                <span class="text-sm font-medium text-slate-700">{{ item.nama
                                                    }}</span><span
                                                    class="shrink-0 font-mono text-sm font-bold text-red-600">{{
                                                        item.nomor }}</span>
                                            </a>
                                        </div>
                                        <p v-else class="text-sm text-slate-500">Nomor kontak darurat belum tersedia.
                                        </p>
                                    </div>
                                </aside>

                                <form class="space-y-5" @submit.prevent="submitAduan">
                                    <div class="grid grid-cols-1 gap-5 sm:grid-cols-2">
                                        <label class="block space-y-2"><span
                                                class="text-sm font-semibold text-slate-700">Nama lengkap</span>
                                            <UInput v-model="formAduan.nama" placeholder="Nama Anda" size="lg"
                                                class="w-full" required />
                                        </label>
                                        <label class="block space-y-2"><span
                                                class="text-sm font-semibold text-slate-700">Email</span>
                                            <UInput v-model="formAduan.kontak"
                                                placeholder="email@contoh.com" size="lg" class="w-full"
                                                required />
                                        </label>
                                    </div>
                                    <label class="block space-y-2"><span
                                            class="text-sm font-semibold text-slate-700">Kategori</span>
                                        <USelect v-model="formAduan.kategori" :items="kategoriAduanOptions" size="lg"
                                            class="w-full" />
                                    </label>
                                    <label class="block space-y-2"><span
                                            class="text-sm font-semibold text-slate-700">Isi keluhan atau
                                            aspirasi</span>
                                        <UTextarea v-model="formAduan.isi_aduan"
                                            placeholder="Ceritakan keluhan atau aspirasi Anda..." :rows="5"
                                            class="w-full" required />
                                    </label>
                                    <div v-if="tiketAduan" role="status"
                                        class="rounded-lg border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-800">
                                        Aduan
                                        berhasil dikirim. Nomor tiket Anda: <strong>#{{ tiketAduan }}</strong></div>
                                    <div v-if="errorAduan" role="alert"
                                        class="rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
                                        {{ errorAduan }}</div>
                                    <div class="flex justify-end">
                                        <UButton type="submit" label="Kirim Aduan" icon="i-tabler-send" color="primary"
                                            size="lg" :loading="submittingAduan" />
                                    </div>
                                </form>
                            </div>
                        </div>
                    </div>

                    <div v-if="socialLinks.length" class="mt-7 flex flex-col items-center pb-2">
                        <p class="mb-3 text-sm font-semibold text-white">Ikuti Media Sosial Kami</p>
                        <div class="flex items-center gap-2.5">
                            <a v-for="social in socialLinks" :key="social.key" :href="social.url" target="_blank"
                                rel="noopener noreferrer" :aria-label="social.label" :title="social.label"
                                :class="['flex h-12 w-12 shrink-0 items-center justify-center rounded-2xl border border-white/25 text-white transition hover:-translate-y-0.5 hover:brightness-110 hover:border-white/50 sm:h-14 sm:w-14', social.class]">
                                <UIcon :name="social.icon" class="h-6 w-6" />
                            </a>
                        </div>
                    </div>
                </UContainer>
            </section>
        </div>
    </div>
</template>