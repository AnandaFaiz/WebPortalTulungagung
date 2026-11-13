<script setup>
// 1. HERO & SEARCH LOGIC
const { data: bannerData } = await useFetch('http://localhost:8080/api/banner?tipe=beranda')
const currentBannerIndex = ref(0)
let bannerTimer = null

onMounted(() => {
    // Cek secara berkala jika data banner sudah siap
    const initTimer = setInterval(() => {
        if (heroSlides.value.length > 1) {
            clearInterval(initTimer) // Hentikan pengecekan jika data sudah ada

            // Jalankan auto-slide setiap 5 detik
            bannerTimer = setInterval(() => {
                currentBannerIndex.value = (currentBannerIndex.value + 1) % heroSlides.value.length
            }, 5000)
        }
    }, 500)
})

onUnmounted(() => {
    if (bannerTimer) clearInterval(bannerTimer)
})

const heroSlides = computed(() => {
    if (!bannerData.value?.banner) return []
    return bannerData.value.banner.map(b => ({
        image: b.gambar || '/hero/berita-1.webp',
        title: b.judul,
        link: b.link || ''
    }))
})

const { data: populerData } = await useFetch('http://localhost:8080/api/pencarian-populer')
const keywords = computed(() => populerData.value?.keywords || [])

const searchQuery = ref('')
const handleSearch = () => {
    if (searchQuery.value.trim()) {
        navigateTo(`/pencarian?q=${encodeURIComponent(searchQuery.value)}`)
    }
}

const fallbackKeywords = ['KTP & KK', 'Izin Usaha', 'Beasiswa', 'Wisata Tulungagung', 'Layanan Darurat']
const searchPlaceholders = computed(() => keywords.value.length > 0 ? keywords.value : fallbackKeywords)
const placeholderIndex = ref(0)
const currentPlaceholder = computed(() => searchPlaceholders.value[placeholderIndex.value % searchPlaceholders.value.length])
let placeholderTimer = null
onMounted(() => {
    placeholderTimer = setInterval(() => { placeholderIndex.value++ }, 3000)
})
onUnmounted(() => {
    if (placeholderTimer) clearInterval(placeholderTimer)
})

// 2. BERITA & PENGUMUMAN
const { data: beritaTerbaruData, pending: loadingTerbaru } = await useFetch('http://localhost:8080/api/berita?sort=terbaru&limit=5')
const { data: beritaPopulerData, pending: loadingPopuler } = await useFetch('http://localhost:8080/api/berita?sort=populer&limit=5')
const beritaTerbaru = computed(() => beritaTerbaruData.value?.berita || [])
const beritaPopuler = computed(() => beritaPopulerData.value?.berita || [])
const activeTab = ref('terbaru')
const currentBeritaList = computed(() => activeTab.value === 'terbaru' ? beritaTerbaru.value : beritaPopuler.value)
const currentLoading = computed(() => activeTab.value === 'terbaru' ? loadingTerbaru.value : loadingPopuler.value)

const { data: pengumumanData, pending: pengumumanLoading } = await useFetch('http://localhost:8080/api/pengumuman?limit=3')
const pengumuman = computed(() => pengumumanData.value?.pengumuman || [])

//Layanan Terpopuler
const { data: layananData } = await useFetch('http://localhost:8080/api/layanan-terpopuler')
const layananWarga = computed(() => layananData.value?.layanan?.filter(l => l.kategori === 'Warga') || [])
const layananUsaha = computed(() => layananData.value?.layanan?.filter(l => l.kategori === 'Usaha') || [])
const layananPekerja = computed(() => layananData.value?.layanan?.filter(l => l.kategori === 'Pencari Kerja') || [])

// ==== FETCH LAYANAN DARURAT UNTUK MODAL ====
const showDaruratModal = ref(false)

// Ambil data layanan darurat dari API Golang
const { data: daruratData } = await useFetch('http://localhost:8080/api/layanan-darurat')
const daruratList = computed(() => daruratData.value?.darurat || [])
const { data: kontakData } = await useFetch('http://localhost:8080/api/kontak')
const { data: pengaduanConfigData } = await useFetch('http://localhost:8080/api/pengaduan-konfigurasi')
const pengaduanKategori = computed(() => pengaduanConfigData.value?.kategori || [])
const socialLinks = computed(() => {
    const kontak = kontakData.value?.kontak || {}
    return [
        { key: 'facebook', label: 'Facebook', icon: 'i-mdi-facebook', url: kontak.facebook, class: 'hover:bg-[#1877F2]' },
        { key: 'instagram', label: 'Instagram', icon: 'i-mdi-instagram', url: kontak.instagram, class: 'hover:bg-[#DD2A7B]' },
        { key: 'twitter', label: 'X', icon: 'i-tabler-brand-x', url: kontak.twitter, class: 'hover:bg-black' },
        { key: 'youtube', label: 'YouTube', icon: 'i-mdi-youtube', url: kontak.youtube, class: 'hover:bg-[#FF0000]' }
    ].filter(item => item.url)
})

// Fetch Agenda dari Golang
const { data: agendaData, pending: agendaLoading } = await useFetch('http://localhost:8080/api/agenda')
const agendas = computed(() => agendaData.value?.agenda || [])

// Countdown Timer
const now = ref(new Date())
let countdownTimer = null

onMounted(() => {
    countdownTimer = setInterval(() => {
        now.value = new Date()
    }, 1000) // Update tiap menit
})

onUnmounted(() => {
    if (countdownTimer) clearInterval(countdownTimer)
})

const getCountdown = (targetDateStr) => {
    const target = new Date(targetDateStr).getTime()
    const diff = target - now.value.getTime()
    if (diff <= 0) return { d: '00', h: '00', m: '00', s: '00' }
    const d = Math.floor(diff / (1000 * 60 * 60 * 24))
    const h = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
    const m = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
    const s = Math.floor((diff % (1000 * 60)) / 1000)
    return {
        d: String(d).padStart(2, '0'),
        h: String(h).padStart(2, '0'),
        m: String(m).padStart(2, '0'),
        s: String(s).padStart(2, '0')
    }
}

// Kalender Logic
const currentMonth = ref(new Date().getMonth())
const currentYear = ref(new Date().getFullYear())

const monthNames = ['Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni', 'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember']
const daysOfWeek = ['Min', 'Sen', 'Sel', 'Rab', 'Kam', 'Jum', 'Sab']

const daysInMonth = computed(() => {
    return new Date(currentYear.value, currentMonth.value + 1, 0).getDate()
})

const firstDayOfMonth = computed(() => {
    return new Date(currentYear.value, currentMonth.value, 1).getDay()
})

const calendarDays = computed(() => {
    const days = []
    for (let i = 0; i < firstDayOfMonth.value; i++) days.push(null)
    for (let i = 1; i <= daysInMonth.value; i++) days.push(i)
    return days
})

const prevMonth = () => {
    if (currentMonth.value === 0) {
        currentMonth.value = 11
        currentYear.value--
    } else {
        currentMonth.value--
    }
}

const nextMonth = () => {
    if (currentMonth.value === 11) {
        currentMonth.value = 0
        currentYear.value++
    } else {
        currentMonth.value++
    }
}

// ==== CAROUSEL REF & HANDLER KLIK KALENDER ====
const carouselRef = ref(null)
// ==== CAROUSEL LOGIC & AUTO SLIDE ====
let autoSlideTimer = null

const scrollCarousel = (direction) => {
    if (!carouselRef.value) return
    const scrollAmount = carouselRef.value.offsetWidth
    carouselRef.value.scrollBy({ left: direction * scrollAmount, behavior: 'smooth' })

    //reset timer
    if (autoSlideTimer) clearInterval(autoSlideTimer)
    startAutoSlide()
}

const startAutoSlide = () => {
    if (agendas.value.length <= 1) return
    autoSlideTimer = setInterval(() => {
        if (carouselRef.value) {
            const maxScroll = carouselRef.value.scrollWidth - carouselRef.value.offsetWidth
            // Jika sudah di ujung kanan, balik ke awal
            if (carouselRef.value.scrollLeft >= maxScroll - 10) {
                carouselRef.value.scrollTo({ left: 0, behavior: 'smooth' })
            } else {
                carouselRef.value.scrollBy({ left: carouselRef.value.offsetWidth, behavior: 'smooth' })
            }
        }
    }, 5000) // 5000ms = 5 detik
}

// Jalankan auto slide saat halaman dimuat
onMounted(() => {
    setTimeout(() => { startAutoSlide() }, 500)
})

onUnmounted(() => {
    if (autoSlideTimer) clearInterval(autoSlideTimer)
})

const handleDateClick = (day) => {
    if (!day) return

    // Format tanggal menjadi YYYY-MM-DD
    const dateStr = `${currentYear.value}-${String(currentMonth.value + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`

    // Cari agenda berdasarkan tanggal
    const agendaIndex = agendas.value.findIndex(a => a.tanggal.startsWith(dateStr))

    // Jika agenda ditemukan, geser slider ke agenda tersebut
    if (agendaIndex !== -1 && carouselRef.value) {
        const agendaElement = carouselRef.value.children[agendaIndex]
        if (agendaElement) {
            agendaElement.scrollIntoView({ behavior: 'smooth', block: 'nearest', inline: 'start' })
        }
    }
}

// Cek apakah tanggal punya agenda
const hasEventOnDate = (day) => {
    if (!day) return false
    const dateStr = `${currentYear.value}-${String(currentMonth.value + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`
    return agendas.value.some(a => a.tanggal.startsWith(dateStr))
}

// Fetch Galeri
const { data: galeriData, pending: galeriLoading } = await useFetch('http://localhost:8080/api/galeri')
const galeri = computed(() => galeriData.value?.galeri || [])

// State Modal
const showAksesCepat = ref(false)

// Fetch Aplikasi untuk Modal Akses Cepat
const { data: aplikasiModalData } = await useFetch('http://localhost:8080/api/aplikasi')
const listAplikasi = computed(() => aplikasiModalData.value?.aplikasi || [])

const showAduanModal = ref(false)
const submittingAduan = ref(false)

const formAduan = ref({
    nama: '',
    kontak: '',
    kategori: '',
    isi_aduan: ''
})

watch(pengaduanKategori, categories => {
    if (!formAduan.value.kategori && categories.length) formAduan.value.kategori = categories[0].nama
}, { immediate: true })

const submitAduan = async () => {
    if (!formAduan.value.nama || !formAduan.value.kontak || !formAduan.value.isi_aduan) {
        alert('Mohon lengkapi semua data!')
        return
    }
    submittingAduan.value = true
    try {
        const res = await $fetch('http://localhost:8080/api/pengaduan', {
            method: 'POST',
            body: formAduan.value
        })
        showAduanModal.value = false
        alert(`Terima kasih! Aduan Anda berhasil dikirim.\nNomor Tiket Anda: #${res.tiket_id}`)
        formAduan.value = { nama: '', kontak: '', kategori: pengaduanKategori.value[0]?.nama || '', isi_aduan: '' }
    } catch (err) {
        alert('Gagal mengirim aduan. Silakan coba lagi.')
    } finally {
        submittingAduan.value = false
    }
}
</script>

<template>
    <div>
        <!-- Hero Slider dengan Pencarian -->
        <div class="relative h-[85vh] min-h-150 bg-slate-900 overflow-hidden">
            <!-- Background Slideshow -->
            <div v-if="bannerLoading" class="absolute inset-0 flex items-center justify-center">
                <UIcon name="i-tabler-loader-2" class="w-12 h-12 text-sky-500 animate-spin" />
            </div>
            <div v-else-if="heroSlides.length > 0" class="absolute inset-0">
                <div v-for="(slide, index) in heroSlides" :key="index"
                    class="absolute inset-0 transition-opacity duration-1000"
                    :class="currentBannerIndex === index ? 'opacity-100 z-0' : 'opacity-0 z-0'">
                    <img :src="slide.image" class="w-full h-full object-cover" />
                </div>
                <div class="absolute inset-0 bg-slate-900/70"></div>
            </div>
            <div v-else class="absolute inset-0 bg-slate-900"></div>

            <!-- Konten Hero -->
            <div class="relative z-10 h-full flex flex-col justify-center px-4 md:px-12 lg:px-24 max-w-7xl mx-auto">
                <div class="max-w-3xl">
                    <h1 class="text-3xl md:text-4xl lg:text-5xl font-bold text-white leading-tight mb-4 drop-shadow-lg">
                        Menjawab Kebutuhan<br />Informasi Warga Tulungagung
                    </h1>
                    <p class="text-sm md:text-base text-white/80 mb-6 drop-shadow-md max-w-xl leading-relaxed">
                        Temukan informasi publik terkini, layanan, dan berita dari Pemerintah Kabupaten Tulungagung.
                    </p>

                    <!-- Search Bar -->
                    <form @submit.prevent="handleSearch"
                        class="relative max-w-md flex items-center bg-black/30 backdrop-blur-md border border-white/20 rounded-xl overflow-hidden shadow-2xl transition-all focus-within:border-sky-400 focus-within:ring-1 focus-within:ring-sky-400">
                        <UIcon name="i-tabler-search" class="absolute left-4 w-5 h-5 text-white/70 z-20" />
                        <input v-model="searchQuery" type="text"
                            class="w-full pl-11 pr-4 py-3 bg-transparent text-white outline-none text-sm relative z-10" />

                        <div v-if="!searchQuery"
                            class="absolute left-11 top-0 bottom-0 flex items-center pointer-events-none overflow-hidden z-0 w-[calc(100%-140px)]">
                            <Transition enter-active-class="transition-all duration-300 ease-out"
                                enter-from-class="translate-y-4 opacity-0" enter-to-class="translate-y-0 opacity-100"
                                leave-active-class="transition-all duration-200 ease-in absolute inset-0"
                                leave-from-class="translate-y-0 opacity-100" leave-to-class="-translate-y-4 opacity-0"
                                mode="out-in">
                                <span :key="currentPlaceholder" class="text-white/60 text-sm whitespace-nowrap">{{
                                    currentPlaceholder }}</span>
                            </Transition>
                        </div>
                        <button type="submit"
                            class="bg-sky-500 hover:bg-sky-600 text-white font-semibold px-5 py-2 m-1.5 rounded-lg transition-colors text-sm relative z-20">Cari</button>
                    </form>

                    <!-- Action Buttons -->
                    <div class="mt-8 flex flex-wrap gap-3">
                        <button type="button" @click="showDaruratModal = true"
                            class="flex items-center gap-2 px-5 py-2.5 rounded-xl bg-red-500/20 hover:bg-red-500/30 backdrop-blur-md border border-red-400/30 text-white transition-all duration-300 hover:-translate-y-1 hover:shadow-lg hover:shadow-red-500/30">
                            <UIcon name="i-tabler-phone-call" class="w-5 h-5 text-red-400" />
                            <span class="text-sm font-semibold">Layanan Darurat</span>
                        </button>
                        <button @click="showAduanModal = true"
                            class="flex items-center gap-2 px-5 py-2.5 rounded-xl bg-amber-500/20 hover:bg-amber-500/30 backdrop-blur-md border border-amber-400/30 text-white transition-all duration-300 hover:-translate-y-1 hover:shadow-lg hover:shadow-amber-500/30">
                            <UIcon name="i-tabler-speakerphone" class="w-5 h-5 text-amber-400" />
                            <span class="text-sm font-semibold">Aduan Warga</span>
                        </button>
                        <button @click="showAksesCepat = true"
                            class="flex items-center gap-2 px-5 py-2.5 rounded-xl bg-emerald-500/20 hover:bg-emerald-500/30 backdrop-blur-md border border-emerald-400/20 text-white transition-all duration-300 hover:-translate-y-1 hover:shadow-lg hover:shadow-emerald-500/30">
                            <UIcon name="i-mdi-clock-fast" class="w-5 h-5 text-emerald-400" />
                            <span class="text-sm font-semibold">Akses Cepat</span>
                        </button>
                    </div>
                </div>
            </div>

            <!-- Tombol Media Sosial Melayang di Tengah Kanan -->
            <div class="absolute top-1/2 -translate-y-1/2 right-6 lg:right-8 z-20 flex flex-col gap-3">

                <a v-for="social in socialLinks" :key="social.key" :href="social.url" target="_blank"
                    :aria-label="social.label"
                    class="w-12 h-12 rounded-xl flex items-center justify-center bg-white/10 backdrop-blur-md border border-white/20 text-white transition-all duration-300 hover:-translate-x-1 hover:border-transparent hover:shadow-lg"
                    :class="social.class">
                    <UIcon :name="social.icon" class="w-6 h-6" />
                </a>

            </div>
        </div>

        <!-- Main Content -->
        <div class="relative bg-transparent">

            <!-- Section Berita & Pengumuman -->
            <section class="pt-8 pb-16 md:pb-24">
                <UContainer>
                    <div class="grid grid-cols-1 lg:grid-cols-12 gap-8">

                        <!-- Kolom Kiri: Berita -->
                        <div
                            class="lg:col-span-8 bg-white dark:bg-slate-900 rounded-2xl shadow-xl border border-slate-200 dark:border-slate-800 p-6 md:p-8">
                            <div
                                class="flex items-center justify-between border-b border-slate-200 dark:border-slate-700 pb-4 mb-6">
                                <h3 class="text-2xl font-bold text-slate-900 dark:text-white flex items-center gap-2">
                                    Berita Terkini</h3>
                                <div class="flex gap-4">
                                    <button @click="activeTab = 'terbaru'"
                                        class="text-sm font-semibold transition-colors pb-4 -mb-4.25"
                                        :class="activeTab === 'terbaru' ? 'text-sky-600 border-b-2 border-sky-600' : 'text-slate-400 hover:text-slate-600'">Terbaru</button>
                                    <button @click="activeTab = 'populer'"
                                        class="text-sm font-semibold transition-colors pb-4 -mb-4.25"
                                        :class="activeTab === 'populer' ? 'text-sky-600 border-b-2 border-sky-600' : 'text-slate-400 hover:text-slate-600'">Terpopuler</button>
                                </div>
                            </div>

                            <div v-if="currentLoading" class="py-12 text-center">
                                <UIcon name="i-tabler-loader-2" class="w-10 h-10 text-sky-500 animate-spin mx-auto" />
                            </div>
                            <div v-else-if="currentBeritaList.length > 0" class="space-y-6">
                                <NuxtLink :to="`/berita/${currentBeritaList[0].id}`" class="group block">
                                    <div class="relative h-64 md:h-80 w-full rounded-xl overflow-hidden mb-4">
                                        <img v-if="currentBeritaList[0].thumbnail" :src="currentBeritaList[0].thumbnail"
                                            class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" />
                                        <div v-else class="w-full h-full bg-sky-100 flex items-center justify-center">
                                            <UIcon name="i-tabler-news" class="w-16 h-16 text-sky-300" />
                                        </div>
                                        <div class="absolute top-4 left-4">
                                            <UBadge color="primary" variant="solid" size="sm">{{
                                                currentBeritaList[0].kategori }}</UBadge>
                                        </div>
                                    </div>
                                    <h2
                                        class="text-xl md:text-2xl font-bold text-slate-900 dark:text-white group-hover:text-sky-600 transition-colors line-clamp-2">
                                        {{ currentBeritaList[0].judul }}</h2>
                                    <p class="text-slate-500 mt-2 line-clamp-2">{{ currentBeritaList[0].deskripsi }}</p>
                                    <p class="text-sm text-slate-400 mt-3 flex items-center gap-1">
                                        <UIcon name="i-tabler-calendar" /> {{ currentBeritaList[0].tanggal }}
                                    </p>
                                </NuxtLink>
                                <hr class="border-slate-200 dark:border-slate-700" />
                                <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                                    <NuxtLink v-for="item in currentBeritaList.slice(1)" :key="item.id"
                                        :to="`/berita/${item.id}`" class="group flex gap-4">
                                        <div class="w-24 h-24 shrink-0 rounded-lg overflow-hidden bg-slate-100">
                                            <img v-if="item.thumbnail" :src="item.thumbnail"
                                                class="w-full h-full object-cover group-hover:scale-110 transition-transform" />
                                            <UIcon v-else name="i-tabler-news"
                                                class="w-full h-full p-6 text-slate-300" />
                                        </div>
                                        <div class="flex flex-col justify-center">
                                            <h4
                                                class="text-sm font-bold text-slate-900 dark:text-white group-hover:text-sky-600 transition-colors line-clamp-2">
                                                {{ item.judul }}</h4>
                                            <p class="text-xs text-slate-400 mt-2 flex items-center gap-1">
                                                <UIcon name="i-tabler-calendar" /> {{ item.tanggal }}
                                            </p>
                                        </div>
                                    </NuxtLink>
                                </div>
                            </div>
                            <div class="mt-8 text-center">
                                <UButton to="/berita" label="Lihat Semua Berita" color="primary" variant="outline"
                                    trailing-icon="i-tabler-arrow-right" />
                            </div>
                        </div>

                        <!-- Kolom Kanan: Pengumuman -->
                        <div class="lg:col-span-4">
                            <div
                                class="bg-slate-50 dark:bg-slate-900 rounded-2xl border border-slate-200 dark:border-slate-800 p-6 sticky top-24">
                                <h3
                                    class="text-xl font-bold text-slate-900 dark:text-white mb-6 flex items-center gap-2 border-b border-slate-200 dark:border-slate-700 pb-4">
                                    <UIcon name="i-tabler-speakerphone" class="text-sky-500" /> Pengumuman Resmi
                                </h3>
                                <div v-if="pengumumanLoading" class="text-center py-8">
                                    <UIcon name="i-tabler-loader-2" class="w-8 h-8 text-sky-500 animate-spin mx-auto" />
                                </div>
                                <div v-else class="space-y-4">
                                    <NuxtLink v-for="item in pengumuman" :key="item.id" to="/pengumuman"
                                        class="block group p-4 bg-white dark:bg-slate-800 rounded-xl border border-slate-100 dark:border-slate-700 hover:border-sky-400 hover:shadow-md transition-all">
                                        <div class="flex items-start gap-3">
                                            <div class="w-2 h-2 mt-2 rounded-full shrink-0"
                                                :class="item.prioritas === 'Tinggi' ? 'bg-red-500' : 'bg-sky-400'">
                                            </div>
                                            <div>
                                                <h4
                                                    class="text-sm font-semibold text-slate-900 dark:text-white group-hover:text-sky-600 transition-colors line-clamp-2">
                                                    {{ item.judul }}</h4>
                                                <p class="text-xs text-slate-400 mt-1">{{ item.tanggal?.split('T')[0] }}
                                                </p>
                                            </div>
                                        </div>
                                    </NuxtLink>
                                </div>
                            </div>
                        </div>

                    </div>
                </UContainer>
            </section>

            <!-- Section Layanan Terpopuler (Gaya Jabarprov) -->
            <section class="bg-slate-50 dark:bg-slate-900/50 py-16">
                <UContainer>
                    <h2 class="text-3xl font-bold text-slate-900 dark:text-white mb-10">Layanan Terpopuler</h2>

                    <div
                        class="grid grid-cols-1 md:grid-cols-3 gap-0 md:divide-x divide-slate-200 dark:divide-slate-700">

                        <!-- Kolom Warga (Hover Box: Amber) -->
                        <div class="px-0 md:px-8 pb-8 md:pb-0">
                            <div class="flex items-center gap-4 mb-8 p-4 rounded-2xl bg-slate-50 ">
                                <div
                                    class="w-14 h-14 rounded-2xl bg-amber-100 flex items-center justify-center shadow-sm">
                                    <UIcon name="i-tabler-users" class="w-8 h-8 text-amber-600" />
                                </div>
                                <div>
                                    <p class="text-xs font-bold text-black uppercase tracking-wider">Untuk</p>
                                    <h3 class="text-2xl font-bold text-slate-900 dark:text-white">Warga</h3>
                                </div>
                            </div>
                            <div class="space-y-2">
                                <NuxtLink v-for="item in layananWarga" :key="item.judul" :to="item.link"
                                    class="group flex gap-4 p-3 -mx-3 rounded-xl hover:bg-amber-100 dark:hover:bg-amber-950/30 transition-colors">
                                    <UIcon :name="item.icon"
                                        class="w-8 h-8 text-amber-500 group-hover:text-amber-500 transition-colors shrink-0" />
                                    <div>
                                        <h4
                                            class="font-bold text-slate-800 dark:text-slate-200 group-hover:text-amber-600 transition-colors">
                                            {{ item.judul }}</h4>
                                        <p class="text-sm text-slate-500 mt-1">{{ item.deskripsi }}</p>
                                    </div>
                                </NuxtLink>
                            </div>
                        </div>

                        <!-- Kolom Usaha (Hover Box: Emerald) -->
                        <div
                            class="px-0 md:px-8 py-8 md:py-0 border-t md:border-t-0 border-slate-200 dark:border-slate-700">
                            <div class="flex items-center gap-4 mb-8 p-4 rounded-2xl bg-slate-50">
                                <div
                                    class="w-14 h-14 rounded-2xl bg-emerald-100 flex items-center justify-center shadow-sm">
                                    <UIcon name="i-tabler-briefcase" class="w-8 h-8 text-emerald-600" />
                                </div>
                                <div>
                                    <p class="text-xs font-bold text-black uppercase tracking-wider">Untuk</p>
                                    <h3 class="text-2xl font-bold text-slate-900 dark:text-white">Usaha</h3>
                                </div>
                            </div>
                            <div class="space-y-2">
                                <NuxtLink v-for="item in layananUsaha" :key="item.judul" :to="item.link"
                                    class="group flex gap-4 p-3 -mx-3 rounded-xl hover:bg-emerald-100 dark:hover:bg-emerald-950/30 transition-colors">
                                    <UIcon :name="item.icon"
                                        class="w-8 h-8 text-emerald-500 group-hover:text-emerald-500 transition-colors shrink-0" />
                                    <div>
                                        <h4
                                            class="font-bold text-slate-800 dark:text-slate-200 group-hover:text-emerald-600 transition-colors">
                                            {{ item.judul }}</h4>
                                        <p class="text-sm text-slate-500 mt-1">{{ item.deskripsi }}</p>
                                    </div>
                                </NuxtLink>
                            </div>
                        </div>

                        <!-- Kolom Pencari Kerja (Hover Box: Sky) -->
                        <div
                            class="px-0 md:px-8 pt-8 md:pt-0 border-t md:border-t-0 border-slate-200 dark:border-slate-700">
                            <div class="flex items-center gap-4 mb-8 p-4 rounded-2xl bg-slate-50">
                                <div
                                    class="w-14 h-14 rounded-2xl bg-sky-100 flex items-center justify-center shadow-sm">
                                    <UIcon name="i-tabler-school" class="w-8 h-8 text-sky-600" />
                                </div>
                                <div>
                                    <p class="text-xs font-bold text-black uppercase tracking-wider">Untuk</p>
                                    <h3 class="text-2xl font-bold text-slate-900 dark:text-white">Pencari Kerja</h3>
                                </div>
                            </div>
                            <div class="space-y-2">
                                <NuxtLink v-for="item in layananPekerja" :key="item.judul" :to="item.link"
                                    class="group flex gap-4 p-3 -mx-3 rounded-xl hover:bg-sky-100 dark:hover:bg-sky-950/30 transition-colors">
                                    <UIcon :name="item.icon"
                                        class="w-8 h-8 text-sky-500 group-hover:text-sky-500 transition-colors shrink-0" />
                                    <div>
                                        <h4
                                            class="font-bold text-slate-800 dark:text-slate-200 group-hover:text-sky-600 transition-colors">
                                            {{ item.judul }}</h4>
                                        <p class="text-sm text-slate-500 mt-1">{{ item.deskripsi }}</p>
                                    </div>
                                </NuxtLink>
                            </div>
                        </div>

                    </div>
                </UContainer>
            </section>

            <!-- Section Agenda & Kalender -->
            <section class="py-16">
                <UContainer>
                    <div class="grid grid-cols-1 lg:grid-cols-12 gap-8">

                        <!-- Kolom Kiri: Agenda Slider -->
                        <div class="lg:col-span-8 overflow-hidden">
                            <div class="flex items-center justify-between mb-6">
                                <h2 class="text-2xl font-bold text-white">Event Tulungagung</h2>
                            </div>

                            <!-- Empty State jika tidak ada data -->
                            <div v-if="agendas.length === 0"
                                class="w-full h-75 flex flex-col items-center justify-center bg-slate-50 dark:bg-slate-800/50 rounded-2xl border border-slate-200 dark:border-slate-700">
                                <UIcon name="i-tabler-calendar-off" class="w-16 h-16 text-slate-300 mb-3" />
                                <p class="text-slate-500 font-medium">Belum ada event mendatang.</p>
                            </div>

                            <!-- Carousel Container dengan Tombol Navigasi -->
                            <div v-else class="relative group/carousel">
                                <!-- Tombol Kiri -->
                                <button v-if="agendas.length > 1" @click="scrollCarousel(-1)"
                                    class="absolute left-4 top-1/2 -translate-y-1/2 z-30 bg-black/50 hover:bg-black/75 text-white p-2.5 rounded-full backdrop-blur-md transition-all duration-300">
                                    <UIcon name="i-tabler-chevron-left" class="w-6 h-6" />
                                </button>

                                <!-- Tombol Kanan -->
                                <button v-if="agendas.length > 1" @click="scrollCarousel(1)"
                                    class="absolute right-4 top-1/2 -translate-y-1/2 z-30 bg-black/50 hover:bg-black/75 text-white p-2.5 rounded-full backdrop-blur-md transition-all duration-300">
                                    <UIcon name="i-tabler-chevron-right" class="w-6 h-6" />
                                </button>

                                <div ref="carouselRef" class="flex overflow-x-auto snap-x snap-mandatory gap-6 pb-6 px-2
                           [&::-webkit-scrollbar]:hidden [-ms-overflow-style:none] scrollbar-none">
                                    <div v-for="agenda in agendas" :key="agenda.id"
                                        class="min-w-full snap-start relative h-112.5 rounded-2xl overflow-hidden shadow-xl">
                                        <!-- Background Image -->
                                        <img :src="agenda.gambar" class="absolute inset-0 w-full h-full object-cover" />

                                        <!-- Gradient Overlay (Gelap di bawah, terang di atas) -->
                                        <div
                                            class="absolute inset-0 bg-linear-to-t from-black/95 via-black/60 to-black/20">
                                        </div>

                                        <!-- Konten di dalam Gambar -->
                                        <div class="absolute inset-0 p-6 md:p-8 flex flex-col justify-end text-white">

                                            <!-- Judul Agenda (Paksa Putih dengan !text-white) -->
                                            <h3 class="text-xl md:text-3xl font-bold mb-3 drop-shadow-lg text-white!">{{
                                                agenda.judul }}</h3>

                                            <!-- Tanggal & Lokasi -->
                                            <div
                                                class="flex flex-wrap items-center gap-4 text-sm md:text-base mb-4 font-medium drop-shadow-md">
                                                <span class="flex items-center gap-1.5">
                                                    <UIcon name="i-tabler-clock" class="w-5 h-5" />
                                                    {{ new Date(agenda.tanggal).toLocaleTimeString('id-ID', {
                                                        hour:
                                                            '2-digit', minute: '2-digit'
                                                    }) }} WIB
                                                </span>
                                                <span class="hidden md:inline text-white/50">|</span>
                                                <span class="flex items-center gap-1.5">
                                                    <UIcon name="i-tabler-calendar" class="w-5 h-5" />
                                                    {{ new Date(agenda.tanggal).toLocaleDateString('id-ID', {
                                                        weekday:
                                                            'long', day: 'numeric', month: 'short', year: 'numeric'
                                                    }) }}
                                                </span>
                                            </div>

                                            <!-- Deskripsi -->
                                            <p
                                                class="text-sm md:text-base text-white/90 line-clamp-2 md:line-clamp-3 max-w-4xl mb-6">
                                                {{ agenda.deskripsi }}
                                            </p>

                                            <!-- Area Bawah: Tags (Kiri) & Countdown + Button (Kanan) -->
                                            <div
                                                class="flex flex-col md:flex-row md:items-end justify-between gap-4 mt-auto w-full">

                                                <!-- Tags (Kiri) -->
                                                <div class="flex flex-wrap gap-2">
                                                    <UBadge color="primary" variant="solid" size="md"
                                                        class="bg-green-600 hover:bg-green-700">{{ agenda.kategori }}
                                                    </UBadge>
                                                    <UBadge color="neutral" variant="solid" size="md"
                                                        class="bg-white/20 backdrop-blur-sm text-white border border-white/30">
                                                        Agenda Offline</UBadge>
                                                    <UBadge color="neutral" variant="solid" size="md"
                                                        class="bg-white/20 backdrop-blur-sm text-white border border-white/30">
                                                        Umum</UBadge>
                                                    <UBadge color="neutral" variant="solid" size="md"
                                                        class="bg-white/20 backdrop-blur-sm text-white border border-white/30 flex items-center gap-1">
                                                        <UIcon name="i-tabler-user" class="w-4 h-4" /> Terbuka
                                                    </UBadge>
                                                </div>

                                                <!-- Countdown & Button (Kanan) -->
                                                <div class="flex flex-col items-end gap-3 shrink-0">

                                                    <!-- Countdown Timer -->
                                                    <div
                                                        class="bg-black/60 backdrop-blur-md border border-white/20 rounded-xl p-2 md:p-3 text-center shadow-lg">
                                                        <div class="flex items-center gap-2 md:gap-3">
                                                            <div class="flex flex-col items-center">
                                                                <span
                                                                    class="text-lg md:text-xl font-bold font-mono text-white">{{
                                                                        getCountdown(agenda.tanggal).d }}</span>
                                                                <span class="text-[9px] text-white/60">Hari</span>
                                                            </div>
                                                            <span class="text-white/50">:</span>
                                                            <div class="flex flex-col items-center">
                                                                <span
                                                                    class="text-lg md:text-xl font-bold font-mono text-white">{{
                                                                        getCountdown(agenda.tanggal).h }}</span>
                                                                <span class="text-[9px] text-white/60">Jam</span>
                                                            </div>
                                                            <span class="text-white/50">:</span>
                                                            <div class="flex flex-col items-center">
                                                                <span
                                                                    class="text-lg md:text-xl font-bold font-mono text-white">{{
                                                                        getCountdown(agenda.tanggal).m }}</span>
                                                                <span class="text-[9px] text-white/60">Menit</span>
                                                            </div>
                                                            <span class="text-white/50">:</span>
                                                            <div class="flex flex-col items-center">
                                                                <span
                                                                    class="text-lg md:text-xl font-bold font-mono text-white">{{
                                                                        getCountdown(agenda.tanggal).s }}</span>
                                                                <span class="text-[9px] text-white/60">Detik</span>
                                                            </div>
                                                        </div>
                                                    </div>

                                                    <!-- Button Selengkapnya (Dinamis) -->
                                                    <UButton :to="agenda.link_event ? agenda.link_event : '/agenda'"
                                                        :target="agenda.link_event ? '_blank' : '_self'"
                                                        label="Selengkapnya" variant="outline" size="lg"
                                                        class="text-white border-white bg-transparent! hover:bg-transparent! rounded-full px-8 shadow-lg" />
                                                </div>

                                            </div>

                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <!-- Kolom Kanan: Kalender -->
                        <div class="lg:col-span-4 mt-14">
                            <div
                                class="bg-white dark:bg-slate-900 rounded-2xl border border-slate-200 dark:border-slate-800 p-6 sticky top-24 shadow-sm">

                                <!-- Header Kalender -->
                                <div class="flex items-center justify-between mb-4">
                                    <button @click="prevMonth"
                                        class="p-2 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-colors">
                                        <UIcon name="i-tabler-chevron-left" class="w-5 h-5" />
                                    </button>
                                    <span class="font-bold text-slate-800 dark:text-white">{{ monthNames[currentMonth]
                                        }} {{ currentYear }}</span>
                                    <button @click="nextMonth"
                                        class="p-2 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-colors">
                                        <UIcon name="i-tabler-chevron-right" class="w-5 h-5" />
                                    </button>
                                </div>

                                <!-- Hari -->
                                <div class="grid grid-cols-7 gap-1 text-center mb-2">
                                    <span v-for="day in daysOfWeek" :key="day"
                                        class="text-xs font-bold text-slate-400">{{ day
                                        }}</span>
                                </div>

                                <!-- Tanggal -->
                                <div class="grid grid-cols-7 gap-1">
                                    <div v-for="(day, index) in calendarDays" :key="index"
                                        class="aspect-square flex items-center justify-center">
                                        <button v-if="day" @click="hasEventOnDate(day) ? handleDateClick(day) : null"
                                            class="w-full h-full rounded-lg text-sm transition-colors relative flex items-center justify-center"
                                            :class="hasEventOnDate(day) ? 'bg-sky-100 text-sky-700 font-bold hover:bg-sky-200 cursor-pointer' : 'text-slate-600 hover:bg-slate-100 cursor-default'">
                                            {{ day }}
                                            <!-- Dot indikator jika ada event -->
                                            <span v-if="hasEventOnDate(day)"
                                                class="absolute bottom-1 w-1.5 h-1.5 rounded-full bg-sky-500"></span>
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </UContainer>
            </section>

            <!-- Section Galeri Kegiatan -->
            <section class="py-16 bg-transparent">
                <UContainer>
                    <div class="flex items-center justify-between mb-8">
                        <div>
                            <h2 class="text-3xl font-bold text-white">Galeri Kegiatan</h2>
                            <p class="text-white mt-2">Dokumentasi foto dan video kegiatan pemerintah daerah.</p>
                        </div>
                        <UButton to="/berita#galeri-kegiatan" label="Lihat Semua" variant="link"
                            trailing-icon="i-tabler-arrow-right"
                            class="hidden sm:flex text-amber-300 hover:text-amber-100" />
                    </div>

                    <div v-if="galeriLoading" class="text-center py-12">
                        <UIcon name="i-tabler-loader-2" class="w-10 h-10 text-sky-500 animate-spin mx-auto" />
                    </div>

                    <div v-else-if="galeri.length === 0"
                        class="text-center py-12 bg-slate-50 dark:bg-slate-900 rounded-2xl">
                        <UIcon name="i-tabler-photo-off" class="w-16 h-16 text-slate-300 mx-auto mb-3" />
                        <p class="text-slate-500">Belum ada galeri kegiatan.</p>
                    </div>

                    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
                        <NuxtLink v-for="item in galeri.slice(0, 4)" :key="item.id" to="/berita#galeri-kegiatan"
                            class="group relative aspect-4/3 rounded-2xl overflow-hidden cursor-pointer bg-slate-100 dark:bg-slate-800 shadow-sm">
                            <!-- Gambar (Jika ada) -->
                            <img v-if="item.thumbnail" :src="item.thumbnail"
                                class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-700" />

                            <!-- Placeholder (Jika tidak ada gambar/thumbnail) -->
                            <div v-else
                                class="absolute inset-0 flex items-center justify-center bg-linear-to-br from-green-300 to-emerald-500">
                                <UIcon :name="item.tipe === 'video' ? 'i-tabler-video' : 'i-tabler-photo'"
                                    class="w-16 h-16 text-white/40 group-hover:scale-110 group-hover:text-white/60 transition-all duration-500" />
                            </div>

                            <!-- Overlay Gradient Gelap di Bawah -->
                            <div
                                class="absolute inset-0 bg-linear-to-t from-black/90 via-black/40 to-transparent opacity-80 group-hover:opacity-100 transition-opacity duration-300">
                            </div>

                            <!-- Badge Tipe (FOTO/VIDEO) -->
                            <div class="absolute top-3 right-3 z-10">
                                <UBadge :color="item.tipe === 'video' ? 'error' : 'neutral'" variant="solid" size="sm"
                                    class="font-semibold backdrop-blur-md">
                                    {{ item.tipe.toUpperCase() }}
                                </UBadge>
                            </div>

                            <!-- Info Text (Judul & Tanggal) -->
                            <div
                                class="absolute bottom-0 left-0 right-0 p-5 z-10 transform translate-y-2 group-hover:translate-y-0 transition-transform duration-300">
                                <h3
                                    class="text-white! font-bold text-lg leading-tight line-clamp-2 drop-shadow-md group-hover:text-teal-300! transition-colors">
                                    {{ item.judul }}
                                </h3>
                                <p class="text-white/80 text-xs mt-2 flex items-center gap-1.5">
                                    <UIcon name="i-tabler-calendar" class="w-3.5 h-3.5" /> {{ item.tanggal }}
                                </p>
                            </div>
                        </NuxtLink>
                    </div>

                    <!-- Tombol Lihat Semua untuk Mobile -->
                    <div class="mt-8 flex justify-center sm:hidden">
                        <UButton to="/berita#galeri-kegiatan" label="Lihat Semua Galeri" color="primary" variant="soft"
                            block />
                    </div>
                </UContainer>
            </section>
        </div>
    </div>

    <!-- Modal Akses Cepat -->
    <UModal v-model:open="showAksesCepat" title="Akses Cepat"
        description="Dapatkan kemudahan akses ke beberapa layanan Pemerintah Kabupaten Tulungagung untuk kebutuhan Anda."
        :ui="{ width: 'max-w-3xl' }">
        <template #body>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <a v-for="app in listAplikasi" :key="app.id" :href="app.url || '#'" target="_blank"
                    class="group border border-slate-200 dark:border-slate-700 rounded-xl p-4 hover:border-sky-400 hover:shadow-md transition-all bg-white dark:bg-slate-800 flex items-start gap-4">
                    <div
                        class="w-10 h-10 rounded-lg bg-sky-50 dark:bg-sky-950 flex items-center justify-center shrink-0 group-hover:bg-sky-500 transition-colors">
                        <UIcon :name="app.icon || 'i-tabler-apps'"
                            class="w-5 h-5 text-sky-600 dark:text-sky-400 group-hover:text-white transition-colors" />
                    </div>
                    <div>
                        <h4 class="font-bold text-slate-900 dark:text-white group-hover:text-sky-600 text-sm">{{
                            app.nama }}</h4>
                        <p class="text-xs text-slate-500 mt-1 line-clamp-2 leading-relaxed">{{ app.deskripsi }}</p>
                    </div>
                </a>
            </div>
        </template>

        <template #footer>
            <div class="flex justify-center w-full">
                <UButton label="Tutup" color="primary" size="lg" class="px-8" @click="showAksesCepat = false" />
            </div>
        </template>
    </UModal>

    <!-- MODAL LAYANAN DARURAT -->
    <UModal v-model:open="showDaruratModal" title="Layanan Darurat"
        description="Hubungi nomor berikut jika Anda membutuhkan bantuan darurat." :ui="{ width: 'max-w-3xl' }">
        <template #body>
            <!-- Grid Card untuk Nomor Darurat -->
            <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
                <a v-for="item in daruratList" :key="item.id" :href="`tel:${item.nomor}`"
                    class="group flex flex-row items-center gap-3 px-4 py-3 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl hover:border-red-400 hover:shadow-md transition-all duration-300 cursor-pointer">
                    <!-- Ikon  -->
                    <div
                        class="w-11 h-11 rounded-full bg-red-50 dark:bg-red-950/50 flex items-center justify-center shrink-0 group-hover:scale-110 transition-transform">
                        <AutoIcon :text="item.nama" class="w-5 h-5 text-red-500" />
                    </div>

                    <!-- Teks -->
                    <div class="flex flex-col text-left justify-center">
                        <!-- Nomor (Tebal & Berwarna) -->
                        <p
                            class="text-lg font-bold text-red-600 dark:text-red-400 leading-tight group-hover:text-red-700">
                            {{ item.nomor }}
                        </p>
                        <!-- Nama Layanan -->
                        <p
                            class="text-xs sm:text-sm font-medium text-slate-600 dark:text-slate-300 leading-tight mt-0.5">
                            {{ item.nama }}
                        </p>
                    </div>
                </a>
            </div>

            <!-- Empty State (Jika admin belum input data) -->
            <div v-if="daruratList.length === 0" class="text-center py-8">
                <UIcon name="i-tabler-phone-off" class="w-12 h-12 text-slate-300 mx-auto mb-2" />
                <p class="text-slate-500">Belum ada data layanan darurat.</p>
            </div>
        </template>

        <template #footer>
            <div class="flex justify-center w-full">
                <UButton label="Tutup" color="neutral" variant="soft" size="lg" class="px-8"
                    @click="showDaruratModal = false" />
            </div>
        </template>
    </UModal>

    <!-- Modal Aduan Warga -->
    <UModal v-model:open="showAduanModal" title="Aduan Warga" description="Sampaikan keluhan atau aspirasi Anda."
        :ui="{ width: 'max-w-xl' }">
        <template #body>
            <div class="space-y-4">
                <UFormField label="Nama Lengkap" required>
                    <UInput v-model="formAduan.nama" placeholder="Nama Anda" size="lg" class="w-full" />
                </UFormField>
                <UFormField label="Email / No. WhatsApp" required>
                    <UInput v-model="formAduan.kontak" placeholder="email@contoh.com atau 0812..." size="lg"
                        class="w-full" />
                </UFormField>
                <UFormField label="Kategori Aduan" required>
                    <USelect v-model="formAduan.kategori" :items="pengaduanKategori.map(item => item.nama)" size="lg"
                        class="w-full" />
                </UFormField>
                <UFormField label="Isi Aduan" required>
                    <UTextarea v-model="formAduan.isi_aduan" placeholder="Tulis detail aduan Anda di sini..." :rows="4"
                        size="lg" class="w-full" />
                </UFormField>
            </div>
        </template>
        <template #footer>
            <div class="flex justify-end gap-2 w-full">
                <UButton label="Batal" color="neutral" variant="outline" @click="showAduanModal = false" />
                <UButton label="Kirim Aduan" color="primary" :loading="submittingAduan" @click="submitAduan" />
            </div>
        </template>
    </UModal>
</template>