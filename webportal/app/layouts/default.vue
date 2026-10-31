<script setup>
const route = useRoute()
const isOpen = ref(false)
// ==== Logika Deteksi Scroll ====
const isScrolled = ref(false)

const handleScroll = () => {
    // Jika di-scroll lebih dari 50px, ubah header
    isScrolled.value = window.scrollY > 700
}

const { data: identitasData } = await useFetch('http://localhost:8080/api/identitas')

const { data: kategoriData } = await useFetch('http://localhost:8080/api/kategori-berita')
const kategoriList = computed(() => kategoriData.value?.kategori || [])
const { data: navigasiData } = await useFetch('http://localhost:8080/api/menu-navigasi')
const navigasiList = computed(() => navigasiData.value?.menu || [])

onMounted(() => {
    window.addEventListener('scroll', handleScroll)
})

onUnmounted(() => {
    window.removeEventListener('scroll', handleScroll)
})

// State untuk Accordion Mobile Menu
const activeMobileMenu = ref(null)
const toggleMobileMenu = (label) => {
    if (activeMobileMenu.value === label) {
        activeMobileMenu.value = null
    } else {
        activeMobileMenu.value = label
    }
}

const { user, isLoggedIn, logout, initialized, fetchUser } = useAuth()

onMounted(async () => {
    if (!initialized.value) await fetchUser()
})

// Data Mega Menu dari panel admin
const menuItems = computed(() => {
    const groups = [
        { label: 'Berita', to: '/berita', group: 'Berita' },
        { label: 'Layanan Publik', to: '/layanan', group: 'Layanan Publik' },
        { label: 'Pemerintahan', to: '/pemerintahan', group: 'Pemerintahan' }
    ]
    return groups.map(group => ({
        label: group.label,
        to: group.to,
        children: navigasiList.value.filter(item => item.grup === group.group).map(item => ({
            label: item.label,
            desc: item.deskripsi || 'Informasi publik terbaru',
            icon: item.icon || 'i-tabler-link',
            to: item.url,
            color: item.color || 'text-sky-600'
        }))
    })).concat([{ label: 'Kontak', to: '/kontak' }])
})

const userMenuItems = [
    [{ label: user.value?.nama || 'Admin', type: 'label' }],
    [
        { label: 'Dashboard Admin', icon: 'i-tabler-layout-dashboard', to: '/admin' },
        { label: 'Profil Saya', icon: 'i-tabler-user', to: '/admin/profil' }
    ],
    [{ label: 'Keluar', icon: 'i-tabler-logout', color: 'error', onSelect: async () => { await logout() } }]
]

// Tutup mobile menu saat pindah halaman
watch(() => route.path, () => {
    isOpen.value = false
    activeMobileMenu.value = null
})
</script>

<template>
    <div class="min-h-screen flex flex-col font-sans">
        <!-- HEADER RESPONSIF -->
        <header
            class="fixed top-0 left-0 right-0 z-50 backdrop-blur-md border-b border-white/10 transition-all duration-300 h-16"
            :class="isScrolled ? 'bg-sky-600/95 shadow-lg' : 'bg-slate-900/40 hover:bg-sky-600/90'">
            <UContainer class="h-full">
                <div class="flex items-center justify-between h-full gap-4">

                    <!-- Logo -->
                    <NuxtLink to="/" class="flex items-center gap-2.5 shrink-0 group">
                        <img src="/Logo-Kabupaten-Tulungagung.png" alt="Logo"
                            class="h-9 md:h-10 w-auto drop-shadow-md" />
                        <div class="hidden sm:block leading-tight border-l-2 border-white/30 pl-2.5">
                            <p class="text-[10px] font-medium text-white/80">Pemerintah Kabupaten</p>
                            <p class="text-sm font-extrabold text-white tracking-tight uppercase">Tulungagung</p>
                        </div>
                    </NuxtLink>

                    <!-- Navigasi Desktop (Mega Menu) -->
                    <nav class="hidden lg:flex items-center h-full ml-auto">
                        <div v-for="item in menuItems" :key="item.label"
                            class="group relative h-full flex items-center px-1">
                            <div @click="navigateTo(item.to)"
                                class="relative z-20 flex items-center gap-1 px-3 py-2 text-sm font-semibold text-white/90 hover:text-white transition-colors cursor-pointer">
                                {{ item.label }}
                                <UIcon v-if="item.children" name="i-tabler-chevron-down"
                                    class="w-4 h-4 opacity-70 group-hover:rotate-180 transition-transform duration-300 pointer-events-none" />
                            </div>

                            <!-- Mega Menu Dropdown Desktop -->
                            <div v-if="item.children"
                                class="absolute top-full left-1/2 -translate-x-1/2 pt-4 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-300 translate-y-2 group-hover:translate-y-0 z-10 pointer-events-none group-hover:pointer-events-auto">
                                <div
                                    class="w-max min-w-70 max-w-125 bg-sky-800/95 backdrop-blur-xl rounded-2xl shadow-2xl border border-sky-400/30 overflow-hidden">
                                    <div class="grid grid-cols-1 sm:grid-cols-2 gap-2 p-3">
                                        <NuxtLink v-for="child in item.children" :key="child.label" :to="child.to"
                                            class="flex items-start gap-3 p-3 rounded-xl hover:bg-sky-700/80 transition-colors group/item">
                                            <!-- Background Putih, Shadow tipis, efek membesar saat hover -->
                                            <div
                                                class="w-10 h-10 rounded-lg bg-white shadow-sm flex items-center justify-center shrink-0 group-hover/item:scale-110 transition-transform">
                                                <!-- Warna ikon mengambil dari data 'child.color' -->
                                                <UIcon :name="child.icon" class="w-5 h-5"
                                                    :class="child.color || 'text-sky-600'" />
                                            </div>
                                            <div>
                                                <p class="text-sm font-bold text-white">{{ child.label }}</p>
                                                <p class="text-xs text-sky-200/80 mt-0.5 leading-snug">{{ child.desc }}
                                                </p>
                                            </div>
                                        </NuxtLink>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </nav>

                    <!-- Aksi Kanan (Desktop & Tablet) -->
                    <div class="flex items-center gap-2 shrink-0 ml-auto lg:ml-4">
                        <template v-if="!isLoggedIn">
                            <!-- ✅ PERBAIKAN: hidden lg:inline-flex agar hilang saat hamburger muncul -->
                            <UButton to="/login" label="Login" icon="i-tabler-login" color="primary" size="sm"
                                class="hidden lg:inline-flex" />
                        </template>
                        <template v-else>
                            <!-- ✅ PERBAIKAN: hidden lg:block pada Dropdown Profil -->
                            <div class="hidden lg:block">
                                <UDropdownMenu :items="userMenuItems">
                                    <UButton color="neutral" variant="ghost" class="p-1 hover:bg-transparent!">
                                        <div
                                            class="w-8 h-8 rounded-full bg-sky-500 flex items-center justify-center text-white text-xs font-bold ring-2 ring-white/50">
                                            {{ user?.nama?.charAt(0)?.toUpperCase() || 'A' }}
                                        </div>
                                    </UButton>
                                </UDropdownMenu>
                            </div>
                        </template>

                        <!-- Hamburger Mobile -->
                        <UButton :icon="isOpen ? 'i-tabler-x' : 'i-tabler-menu-2'" color="primary" variant="ghost"
                            class="lg:hidden text-white ml-1" @click="isOpen = !isOpen" />
                    </div>
                </div>
            </UContainer>

            <!-- Mobile & Tablet Menu -->
            <Transition enter-active-class="transition-all duration-300 ease-out"
                enter-from-class="opacity-0 -translate-y-4" leave-active-class="transition-all duration-200 ease-in"
                leave-to-class="opacity-0 -translate-y-4">
                <div v-if="isOpen"
                    class="lg:hidden absolute top-16 left-0 right-0 bg-sky-800/95 backdrop-blur-xl shadow-2xl border-t border-sky-600/50 max-h-[85vh] overflow-y-auto">
                    <UContainer class="py-4 flex flex-col gap-1">

                        <!-- Loop Menu Mobile -->
                        <div v-for="item in menuItems" :key="item.label">
                            <!-- Jika punya sub-menu, jadikan tombol accordion -->
                            <button v-if="item.children" @click="toggleMobileMenu(item.label)"
                                class="w-full flex items-center justify-between px-4 py-3 rounded-lg text-sm font-bold text-white hover:bg-sky-700 transition-colors">
                                {{ item.label }}
                                <UIcon name="i-tabler-chevron-down" class="w-4 h-4 transition-transform duration-300"
                                    :class="{ 'rotate-180': activeMobileMenu === item.label }" />
                            </button>

                            <!-- Jika tidak ada sub-menu, langsung link biasa -->
                            <NuxtLink v-else :to="item.to"
                                class="block px-4 py-3 rounded-lg text-sm font-bold text-white hover:bg-sky-700 transition-colors">
                                {{ item.label }}
                            </NuxtLink>

                            <!-- Sub-menu Accordion -->
                            <Transition enter-active-class="transition-all duration-300 ease-out"
                                enter-from-class="max-h-0 opacity-0" enter-to-class="max-h-[500px] opacity-100"
                                leave-active-class="transition-all duration-200 ease-in"
                                leave-from-class="max-h-[500px] opacity-100" leave-to-class="max-h-0 opacity-0">
                                <div v-if="item.children && activeMobileMenu === item.label" class="overflow-hidden">
                                    <div class="pl-4 pr-2 py-2 space-y-1 bg-sky-900/60 rounded-lg mt-1 mb-2 mx-2">
                                        <NuxtLink v-for="child in item.children" :key="child.label" :to="child.to"
                                            class="flex items-center gap-3 p-2 rounded-lg hover:bg-sky-700/50 transition-colors">
                                            <div
                                                class="w-8 h-8 rounded-md bg-sky-600/50 flex items-center justify-center shrink-0">
                                                <UIcon :name="child.icon" class="w-4 h-4 text-sky-200" />
                                            </div>
                                            <div>
                                                <p class="text-sm font-bold text-white">{{ child.label }}</p>
                                                <p class="text-xs text-sky-200/70">{{ child.desc }}</p>
                                            </div>
                                        </NuxtLink>
                                    </div>
                                </div>
                            </Transition>
                        </div>

                        <!-- Tombol Login Mobile (Posisi Kanan Bawah, Ukuran Kecil) -->
                        <div class="border-t border-sky-600/50 pt-4 mt-2 flex justify-end">
                            <UButton v-if="!isLoggedIn" to="/login" label="Masuk" icon="i-tabler-login" color="primary"
                                size="sm" @click="isOpen = false" />
                            <UButton v-else to="/admin" label="Dashboard Admin" icon="i-tabler-layout-dashboard"
                                color="primary" size="sm" @click="isOpen = false" />
                        </div>

                    </UContainer>
                </div>
            </Transition>
        </header>

        <!-- Konten Utama -->
        <main class="flex-1">
            <slot />
        </main>

        <!-- Footer -->
        <footer class="bg-slate-900 text-slate-300 mt-20">
            <UContainer>
                <div class="py-12 grid grid-cols-1 md:grid-cols-3 gap-8">
                    <div>
                        <div class="flex items-center gap-3 mb-4">
                            <img src="/Logo-Kabupaten-Tulungagung.png" alt="Logo" class="h-12 w-auto" />
                            <div class="leading-tight">
                                <p class="text-xs text-slate-400">Pemerintah Kabupaten</p>
                                <p class="font-extrabold text-white">TULUNGAGUNG</p>
                            </div>
                        </div>
                        <p class="text-sm text-slate-400 leading-relaxed">Portal resmi Pemerintah Kabupaten Tulungagung.
                            Sumber
                            informasi dan layanan publik terpercaya.</p>
                    </div>
                    <div>
                        <h4 class="font-bold text-white mb-4">Tautan Cepat</h4>
                        <ul class="space-y-2 text-sm">
                            <li>
                                <NuxtLink to="/profil" class="text-slate-400 hover:text-sky-400">Profil Daerah
                                </NuxtLink>
                            </li>
                            <li>
                                <NuxtLink to="/layanan" class="text-slate-400 hover:text-sky-400">Layanan Publik
                                </NuxtLink>
                            </li>
                            <li>
                                <NuxtLink to="/berita" class="text-slate-400 hover:text-sky-400">Informasi & Berita
                                </NuxtLink>
                            </li>
                            <li>
                                <NuxtLink to="/kontak" class="text-slate-400 hover:text-sky-400">Kontak</NuxtLink>
                            </li>
                        </ul>
                    </div>
                    <div>
                        <h4 class="font-bold text-white mb-4">Hubungi Kami</h4>
                        <ul class="space-y-3 text-sm">
                            <li class="flex items-start gap-2">
                                <UIcon name="i-tabler-map-pin" class="w-4 h-4 mt-0.5 text-sky-400 shrink-0" /><span
                                    class="text-slate-400">Jl. A. Yani Timur No. 37, Tulungagung 66216</span>
                            </li>
                            <li class="flex items-start gap-2">
                                <UIcon name="i-tabler-phone" class="w-4 h-4 mt-0.5 text-sky-400 shrink-0" /><span
                                    class="text-slate-400">(0355) 321-260</span>
                            </li>
                            <li class="flex items-start gap-2">
                                <UIcon name="i-tabler-mail" class="w-4 h-4 mt-0.5 text-sky-400 shrink-0" /><span
                                    class="text-slate-400">info@tulungagung.go.id</span>
                            </li>
                        </ul>
                    </div>
                </div>
                <div class="border-t border-slate-800 py-6 text-center text-xs text-slate-500">
                    © {{ new Date().getFullYear() }} Pemerintah Kabupaten Tulungagung. Hak Cipta Dilindungi.
                </div>
            </UContainer>
        </footer>
    </div>
</template>