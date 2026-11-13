<script setup>
const route = useRoute()
const router = useRouter()
const API_BASE = 'http://localhost:8080/api'

const { data: detailData, pending, error } = await useFetch(
    () => `${API_BASE}/berita/${route.params.id}`,
    { watch: [() => route.params.id] }
)
const { data: relatedData } = await useFetch(`${API_BASE}/berita?sort=terbaru`)
const { data: pengumumanData } = await useFetch(`${API_BASE}/pengumuman`)

const berita = computed(() => detailData.value?.berita || null)
const beritaLainnya = computed(() =>
    (relatedData.value?.berita || []).filter(item => String(item.id) !== String(route.params.id)).slice(0, 5)
)
const pengumuman = computed(() => (pengumumanData.value?.pengumuman || []).slice(0, 3))
const formatDate = (value) => {
    if (!value) return '-'
    return new Intl.DateTimeFormat('id-ID', { dateStyle: 'long' }).format(new Date(value))
}

useHead(() => ({
    title: berita.value?.judul ? `${berita.value.judul} | Berita` : 'Berita'
}))

const goBack = () => router.push('/berita')
</script>

<template>
    <div class="min-h-screen bg-transparent pt-30 pb-16">
        <UContainer>
            <div class="rounded-2xl border border-slate-200 bg-white p-5 shadow-xl md:p-8 lg:p-10">
                <div v-if="pending" class="py-24 text-center">
                    <UIcon name="i-tabler-loader-2" class="w-12 h-12 text-sky-500 animate-spin mx-auto" />
                    <p class="mt-4 text-slate-500">Memuat berita...</p>
                </div>

                <div v-else-if="error || !berita" class="py-24 text-center">
                    <UIcon name="i-tabler-news-off" class="w-16 h-16 text-slate-300 mx-auto" />
                    <h1 class="mt-4 text-2xl font-bold text-slate-800 dark:text-white">Berita tidak ditemukan</h1>
                    <p class="mt-2 text-slate-500">Berita mungkin sudah dihapus atau belum dipublikasikan.</p>
                    <UButton label="Kembali ke Berita" icon="i-tabler-arrow-left" class="mt-6" @click="goBack" />
                </div>

                <div v-else class="grid grid-cols-1 lg:grid-cols-[minmax(0,1.85fr)_minmax(280px,1fr)] gap-10 lg:gap-14">
                    <article>
                        <NuxtLink to="/berita"
                            class="inline-flex items-center gap-2 text-sm font-semibold text-sky-600 hover:text-sky-700 mb-8">
                            <UIcon name="i-tabler-arrow-left" /> Semua Berita
                        </NuxtLink>

                        <div class="flex flex-wrap items-center gap-3 mb-4">
                            <UBadge color="primary" variant="subtle">{{ berita.kategori }}</UBadge>
                            <span class="text-sm text-slate-500">{{ formatDate(berita.tanggal) }}</span>
                        </div>
                        <h1
                            class="text-3xl md:text-5xl font-extrabold leading-tight tracking-tight text-slate-900 dark:text-white">
                            {{ berita.judul }}
                        </h1>
                        <p v-if="berita.deskripsi"
                            class="mt-5 text-lg leading-relaxed text-slate-600 dark:text-slate-300">
                            {{ berita.deskripsi }}
                        </p>
                        <div
                            class="mt-5 flex items-center gap-2 text-sm text-slate-500 border-b border-slate-200 dark:border-slate-800 pb-6">
                            <UIcon name="i-tabler-user-circle" class="text-sky-500" />
                            <span>Ditulis oleh <strong class="text-slate-700 dark:text-slate-200">{{ berita.penulis
                                    }}</strong></span>
                            <span class="text-slate-300">|</span>
                            <span>{{ berita.views }} kali dibaca</span>
                        </div>
                        <div class="mt-8 aspect-video overflow-hidden rounded-2xl bg-sky-100 dark:bg-sky-950">
                            <img v-if="berita.thumbnail" :src="berita.thumbnail" :alt="berita.judul"
                                class="w-full h-full object-cover" />
                            <div v-else class="w-full h-full flex items-center justify-center">
                                <UIcon name="i-tabler-news" class="w-24 h-24 text-sky-300" />
                            </div>
                        </div>

                        <div
                            class="mt-10 max-w-none text-slate-700 dark:text-slate-300 leading-8 [&_a]:text-sky-600 [&_a]:underline [&_img]:rounded-xl [&_img]:max-w-full [&_p]:mb-5 [&_h2]:text-2xl [&_h2]:font-bold [&_h2]:text-slate-900 [&_h2]:dark:text-white [&_h3]:text-xl [&_h3]:font-bold [&_h3]:text-slate-900 [&_h3]:dark:text-white [&_ul]:list-disc [&_ul]:pl-6 [&_ol]:list-decimal [&_ol]:pl-6">
                            <div v-if="berita.konten" v-html="berita.konten" />
                            <p v-else class="italic text-slate-400">Konten berita belum tersedia.</p>
                        </div>
                    </article>

                    <aside class="lg:border-l lg:border-slate-200 lg:dark:border-slate-800 lg:pl-8">
                        <div class="sticky top-24">
                            <div
                                class="flex items-center justify-between border-b-2 border-slate-900 dark:border-white pb-3">
                                <h2 class="text-xl font-extrabold text-slate-900 dark:text-white">Berita Lainnya</h2>
                                <UIcon name="i-tabler-news" class="text-sky-500" />
                            </div>
                            <div class="divide-y divide-slate-200 dark:divide-slate-800">
                                <NuxtLink v-for="item in beritaLainnya" :key="item.id" :to="`/berita/${item.id}`"
                                    class="group flex gap-4 py-5">
                                    <div
                                        class="w-24 h-20 shrink-0 overflow-hidden rounded-lg bg-slate-200 dark:bg-slate-800">
                                        <img v-if="item.thumbnail" :src="item.thumbnail" :alt="item.judul"
                                            class="w-full h-full object-cover group-hover:scale-105 transition-transform" />
                                        <UIcon v-else name="i-tabler-news" class="w-full h-full p-6 text-slate-400" />
                                    </div>
                                    <div>
                                        <UBadge color="primary" variant="subtle" size="xs">{{ item.kategori }}</UBadge>
                                        <h3
                                            class="mt-1 font-bold leading-snug text-slate-800 dark:text-slate-200 group-hover:text-sky-600 line-clamp-3">
                                            {{ item.judul }}
                                        </h3>
                                        <p class="mt-2 text-xs text-slate-400">{{ formatDate(item.tanggal) }}</p>
                                    </div>
                                </NuxtLink>
                            </div>
                            <NuxtLink to="/berita"
                                class="mt-4 flex items-center justify-center gap-2 rounded-lg border border-sky-200 py-3 text-sm font-bold text-sky-600 hover:bg-sky-50 dark:border-sky-900 dark:hover:bg-sky-950">
                                Lihat semua berita
                                <UIcon name="i-tabler-arrow-right" />
                            </NuxtLink>

                            <div class="mt-10 border-t border-slate-200 dark:border-slate-800 pt-6">
                                <div
                                    class="flex items-center justify-between border-b-2 border-slate-900 dark:border-white pb-3">
                                    <h2 class="text-xl font-extrabold text-slate-900 dark:text-white">Pengumuman</h2>
                                    <UIcon name="i-tabler-speakerphone" class="text-sky-500" />
                                </div>
                                <div v-if="pengumuman.length" class="divide-y divide-slate-200 dark:divide-slate-800">
                                    <div v-for="item in pengumuman" :key="item.id" class="py-4">
                                        <div class="flex items-start gap-3">
                                            <div
                                                class="w-16 h-12 shrink-0 overflow-hidden rounded-lg bg-slate-200 dark:bg-slate-800 flex items-center justify-center">
                                                <img v-if="item.gambar" :src="item.gambar" :alt="item.judul"
                                                    class="w-full h-full object-cover" />
                                                <UIcon v-else name="i-tabler-speakerphone"
                                                    class="w-6 h-6 text-sky-400" />
                                            </div>
                                            <div class="min-w-0">
                                                <UBadge :color="item.prioritas === 'Tinggi' ? 'error' : 'primary'"
                                                    variant="subtle" size="xs">{{ item.prioritas }}</UBadge>
                                                <h3
                                                    class="mt-1 font-bold leading-snug text-slate-800 dark:text-slate-200 line-clamp-2">
                                                    {{ item.judul }}</h3>
                                                <p class="mt-1 text-xs text-slate-400">{{ formatDate(item.tanggal) }}
                                                </p>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                                <p v-else class="py-5 text-sm text-slate-400">Belum ada pengumuman.</p>
                                <NuxtLink to="/berita#pengumuman"
                                    class="mt-3 flex items-center justify-center gap-2 rounded-lg border border-sky-200 py-3 text-sm font-bold text-sky-600 hover:bg-sky-50 dark:border-sky-900 dark:hover:bg-sky-950">
                                    Lihat semua pengumuman
                                    <UIcon name="i-tabler-arrow-right" />
                                </NuxtLink>
                            </div>
                        </div>
                    </aside>
                </div>
            </div>
        </UContainer>
    </div>
</template>
