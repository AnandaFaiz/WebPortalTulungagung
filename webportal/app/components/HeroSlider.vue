<script setup>
const props = defineProps({
    slides: { type: Array, required: true },
    interval: { type: Number, default: 5000 }
})

const currentIndex = ref(0)
let timer = null

const nextSlide = () => {
    currentIndex.value = (currentIndex.value + 1) % props.slides.length
}

const goToSlide = (index) => {
    currentIndex.value = index
    resetTimer()
}

const startTimer = () => {
    timer = setInterval(nextSlide, props.interval)
}

const resetTimer = () => {
    if (timer) clearInterval(timer)
    startTimer()
}

onMounted(() => {
    startTimer()
})

onUnmounted(() => {
    if (timer) clearInterval(timer)
})
</script>

<template>
    <section class="relative w-full h-screen overflow-hidden">
        <!-- Slides Background -->
        <div v-for="(slide, index) in slides" :key="index"
            class="absolute inset-0 transition-opacity duration-1000 ease-in-out"
            :class="currentIndex === index ? 'opacity-100' : 'opacity-0'">
            <img :src="slide.image" :alt="slide.title" class="w-full h-full object-cover" loading="lazy" />
            <!-- Overlay gradient gelap -->
            <div class="absolute inset-0 bg-linear-to-b from-slate-900/50 via-slate-900/40 to-slate-900/70" />
        </div>

        <!-- Konten Utama -->
        <div class="relative z-10 h-full flex flex-col items-center justify-center text-center px-4">
            <UBadge variant="solid" size="lg"
                class="mb-6 bg-white/15 backdrop-blur-md text-white border border-white/30">
                Portal Resmi Pemerintah Daerah
            </UBadge>

            <h1
                class="text-4xl md:text-6xl lg:text-7xl font-extrabold text-white max-w-5xl leading-tight drop-shadow-xl">
                Selamat Datang di Portal
                <br class="hidden md:block" />
                Pemkab Tulungagung
            </h1>

            <p class="mt-6 text-base md:text-lg lg:text-xl text-white/90 max-w-2xl drop-shadow-md leading-relaxed">
                Portal resmi Pemerintah Kabupaten Tulungagung. Dapatkan informasi terbaru, layanan publik, dan berita
                kegiatan pemerintah daerah dalam satu pintu.
            </p>

            <div class="mt-8 flex flex-wrap gap-3 justify-center">
                <UButton label="Lihat Layanan" color="primary" size="xl" to="/layanan"
                    class="shadow-lg shadow-sky-500/30" />
                <UButton label="Baca Berita" variant="outline" size="xl" to="/berita"
                    class="bg-white/10 backdrop-blur-md text-white border-white/40 hover:bg-white/20" />
            </div>
        </div>

        <!-- Dots Indicator -->
        <div class="absolute bottom-8 left-1/2 -translate-x-1/2 flex gap-2 z-20">
            <button v-for="(slide, index) in slides" :key="index" :aria-label="`Slide ${index + 1}: ${slide.title}`"
                class="h-2 rounded-full transition-all duration-300" :class="currentIndex === index
                    ? 'w-8 bg-white'
                    : 'w-2 bg-white/50 hover:bg-white/80'" @click="goToSlide(index)" />
        </div>

        <!-- Scroll Hint -->
        <div class="absolute bottom-8 right-8 z-20 text-white/70 animate-bounce hidden md:flex flex-col items-center">
            <span class="text-xs mb-1 tracking-wider">SCROLL</span>
            <UIcon name="i-tabler-chevron-down" class="w-6 h-6" />
        </div>
    </section>
</template>