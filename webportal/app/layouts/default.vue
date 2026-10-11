<script setup>
const items = [
    { label: 'Beranda', to: '/' },
    { label: 'Profil', to: '/profil' },
    { label: 'Pemerintahan', to: '/pemerintahan' },
    { label: 'Layanan Publik', to: '/layanan' },
    { label: 'Informasi & Berita', to: '/berita' },
    { label: 'Kontak', to: '/kontak' },
]

const isOpen = ref(false)
const route = useRoute()

// Tutup mobile menu saat pindah halaman
watch(() => route.path, () => { isOpen.value = false })
</script>

<template>
    <div class="min-h-screen flex flex-col">
        <header class="sticky top-0 z-50 bg-linear-to-r from-sky-200 to-sky-300 shadow-sm">
            <UContainer>
                <div class="flex items-center justify-between h-20 gap-6">

                    <!-- Logo tanpa container putih -->
                    <NuxtLink to="/" class="flex items-center gap-3 shrink-0 group">
                        <img src="/Logo-Kabupaten-Tulungagung.png" alt="Logo Pemkab Tulungagung"
                            class="h-12 md:h-14 w-auto transition-transform duration-300 group-hover:scale-105 drop-shadow-md" />
                        <div class="hidden sm:block leading-tight">
                            <p class="text-[11px] md:text-xs font-medium text-sky-900/80">
                                Pemerintah Kabupaten
                            </p>
                            <p class="text-base md:text-lg font-extrabold text-sky-950 tracking-tight">
                                TULUNGAGUNG
                            </p>
                        </div>
                    </NuxtLink>

                    <!-- Navigasi Desktop -->
                    <nav class="hidden lg:flex items-center gap-1">
                        <NuxtLink v-for="item in items" :key="item.to" :to="item.to"
                            class="px-4 py-2 text-sm font-semibold rounded-lg transition-all duration-200" :class="route.path === item.to
                                ? 'bg-white text-sky-600 shadow-sm'
                                : 'text-sky-900/80 hover:text-sky-950 hover:bg-white/40'">
                            {{ item.label }}
                        </NuxtLink>
                    </nav>

                    <!-- Aksi Kanan -->
                    <div class="flex items-center gap-1">
                        <UButton icon="i-tabler-search" color="neutral" variant="ghost"
                            class="text-sky-900 hover:bg-white/40" aria-label="Cari" />
                        <UButton :icon="isOpen ? 'i-tabler-x' : 'i-tabler-menu-2'" color="neutral" variant="ghost"
                            class="lg:hidden text-sky-900 hover:bg-white/40" aria-label="Menu"
                            @click="isOpen = !isOpen" />
                    </div>
                </div>
            </UContainer>

            <!-- Mobile Menu -->
            <Transition enter-active-class="transition-all duration-300 ease-out"
                enter-from-class="opacity-0 -translate-y-2" leave-active-class="transition-all duration-200 ease-in"
                leave-to-class="opacity-0 -translate-y-2">
                <div v-if="isOpen" class="lg:hidden border-t border-white/50 bg-sky-100">
                    <UContainer>
                        <nav class="py-4 flex flex-col gap-1">
                            <NuxtLink v-for="item in items" :key="item.to" :to="item.to"
                                class="px-4 py-3 rounded-lg text-sm font-semibold transition-colors" :class="route.path === item.to
                                    ? 'bg-white text-sky-600 shadow-sm'
                                    : 'text-sky-900/80 hover:bg-white/50'">
                                {{ item.label }}
                            </NuxtLink>
                        </nav>
                    </UContainer>
                </div>
            </Transition>
        </header>
        <!-- Main Content -->
        <main class="flex-1">
            <slot />
        </main>

        <!-- Footer -->
        <footer class="bg-slate-900 text-slate-300 mt-20">
            <UContainer>
                <div class="py-12 grid grid-cols-1 md:grid-cols-3 gap-8">
                    <!-- Kolom 1: Identitas -->
                    <div>
                        <div class="flex items-center gap-3 mb-4">
                            <img src="/Logo-Kabupaten-Tulungagung.png" alt="Logo" class="h-12 w-auto" />
                            <div class="leading-tight">
                                <p class="text-xs text-slate-400">Pemerintah Kabupaten</p>
                                <p class="font-extrabold text-white">TULUNGAGUNG</p>
                            </div>
                        </div>
                        <p class="text-sm text-slate-400 leading-relaxed">
                            Portal resmi Pemerintah Kabupaten Tulungagung. Sumber informasi dan layanan publik
                            terpercaya.
                        </p>
                    </div>

                    <!-- Kolom 2: Tautan Cepat -->
                    <div>
                        <h4 class="font-bold text-white mb-4">Tautan Cepat</h4>
                        <ul class="space-y-2 text-sm">
                            <li v-for="item in items" :key="item.to">
                                <NuxtLink :to="item.to" class="text-slate-400 hover:text-sky-400 transition-colors">
                                    {{ item.label }}
                                </NuxtLink>
                            </li>
                        </ul>
                    </div>

                    <!-- Kolom 3: Kontak -->
                    <div>
                        <h4 class="font-bold text-white mb-4">Hubungi Kami</h4>
                        <ul class="space-y-3 text-sm">
                            <li class="flex items-start gap-2">
                                <UIcon name="i-tabler-map-pin" class="w-4 h-4 mt-0.5 text-sky-400 shrink-0" />
                                <span class="text-slate-400">Jl. A. Yani Timur No. 37, Tulungagung 66216</span>
                            </li>
                            <li class="flex items-start gap-2">
                                <UIcon name="i-tabler-phone" class="w-4 h-4 mt-0.5 text-sky-400 shrink-0" />
                                <span class="text-slate-400">(0355) 321-260</span>
                            </li>
                            <li class="flex items-start gap-2">
                                <UIcon name="i-tabler-mail" class="w-4 h-4 mt-0.5 text-sky-400 shrink-0" />
                                <span class="text-slate-400">info@tulungagung.go.id</span>
                            </li>
                        </ul>
                    </div>
                </div>

                <div
                    class="border-t border-slate-800 py-6 flex flex-col md:flex-row items-center justify-between gap-4">
                    <p class="text-xs text-slate-500">
                        © {{ new Date().getFullYear() }} Pemerintah Kabupaten Tulungagung. Hak Cipta Dilindungi.
                    </p>
                    <div class="flex items-center gap-3">
                        <a href="#" aria-label="Facebook" class="text-slate-400 hover:text-sky-400 transition-colors">
                            <UIcon name="i-mdi-facebook" class="w-5 h-5" />
                        </a>
                        <a href="#" aria-label="Instagram" class="text-slate-400 hover:text-sky-400 transition-colors">
                            <UIcon name="i-mdi-instagram" class="w-5 h-5" />
                        </a>
                        <a href="#" aria-label="Twitter" class="text-slate-400 hover:text-sky-400 transition-colors">
                            <UIcon name="i-mdi-twitter" class="w-5 h-5" />
                        </a>
                        <a href="#" aria-label="YouTube" class="text-slate-400 hover:text-sky-400 transition-colors">
                            <UIcon name="i-tabler-brand-youtube" class="w-5 h-5" />
                        </a>
                    </div>
                </div>
            </UContainer>
        </footer>
    </div>
</template>