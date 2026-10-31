<script setup>
const props = defineProps({
    title: { type: String, required: true },
    description: { type: String, default: '' },
    badge: { type: String, default: '' },
    images: { type: Array, default: () => [] },
    interval: { type: Number, default: 5000 },
    height: { type: String, default: 'h-[85vh] min-h-150' },  // default full screen
    links: { type: Array, default: () => [] }
})

const currentIndex = ref(0)
let timer = null

const isSlider = computed(() => props.images.length > 1)

const nextSlide = () => {
    if (!isSlider.value) return
    currentIndex.value = (currentIndex.value + 1) % props.images.length
}

onMounted(() => {
    if (isSlider.value) {
        timer = setInterval(nextSlide, props.interval)
    }
})

onUnmounted(() => {
    if (timer) clearInterval(timer)
})
</script>

<template>
    <section class="relative w-full overflow-hidden" :class="height">
        <!-- ==== Wallpaper ==== -->
        <template v-if="images.length > 0">
            <div v-for="(img, index) in images" :key="index"
                class="absolute inset-0 transition-opacity duration-1000 ease-in-out"
                :class="currentIndex === index ? 'opacity-100' : 'opacity-0'">
                <img :src="img" :alt="title" class="w-full h-full object-cover"
                    :loading="index === 0 ? 'eager' : 'lazy'" />
            </div>

            <!-- Overlay -->
            <div class="absolute inset-0 bg-linear-to-b from-slate-900/60 via-slate-900/50 to-slate-900/75" />
        </template>

        <div v-else class="absolute inset-0 bg-linear-to-br from-sky-500 to-sky-700" />

        <!-- ==== Konten dengan Animasi ==== -->
        <div class="relative z-10 h-full flex flex-col items-center justify-center text-center px-4 pt-16">
            <div v-if="badge" class="animate-fade-in-up mb-4">
                <UBadge variant="solid" size="lg"
                    class="bg-white/15 backdrop-blur-md text-white border border-white/30">
                    {{ badge }}
                </UBadge>
            </div>

            <h1
                class="animate-fade-in-up-delay-1 text-4xl md:text-5xl lg:text-6xl font-extrabold text-white max-w-5xl leading-tight drop-shadow-xl">
                {{ title }}
            </h1>

            <p v-if="description"
                class="animate-fade-in-up-delay-2 mt-6 text-base md:text-lg lg:text-xl text-white/90 max-w-2xl drop-shadow-md leading-relaxed">
                {{ description }}
            </p>

            <div v-if="links.length > 0" class="animate-fade-in-up-delay-3 mt-8 flex flex-wrap gap-3 justify-center">
                <UButton v-for="(link, i) in links" :key="i" v-bind="link" />
            </div>

            <slot />
        </div>

        <!-- ==== Dots (kalau slider) ==== -->
        <div v-if="isSlider" class="absolute bottom-8 left-1/2 -translate-x-1/2 flex gap-2 z-20">
            <button v-for="(img, index) in images" :key="index" :aria-label="`Slide ${index + 1}`"
                class="h-2 rounded-full transition-all duration-300" :class="currentIndex === index
                    ? 'w-8 bg-white'
                    : 'w-2 bg-white/50 hover:bg-white/80'" @click="currentIndex = index" />
        </div>
    </section>
</template>