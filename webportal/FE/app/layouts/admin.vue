<script setup>
import { adminNavigation } from '~/utils/admin-navigation'

const { user, logout, can } = useAuth()
const route = useRoute()
const isSidebarOpen = ref(true)
const openGroups = ref({})

const userProfile = computed(() => user.value || {
    nama: user.value?.nama || 'Admin',
    email: user.value?.email || '',
    role: user.value?.role || 'super_admin'
})

const sidebarLinks = adminNavigation

// ==== Filter menu berdasarkan permission dari DB ====
const filteredSidebarLinks = computed(() => {
    return sidebarLinks
        .map(item => {
            // ==== Item dengan children ====
            if (item.children) {
                // Filter children dulu
                const filteredChildren = item.children.filter(child => can(child.permission))
                // Kalau tidak ada child yang lolos, sembunyikan parent
                if (filteredChildren.length === 0) return null
                // Kalau ada, tampilkan parent + children yang lolos
                return { ...item, children: filteredChildren }
            }

            // ==== Item tanpa children ====
            // Cek permission langsung
            return can(item.permission) ? item : null
        })
        .filter(Boolean) // buang yang null
})

const isGroupOpen = (label) => openGroups.value[label] ?? false

const toggleGroup = (label) => {
    const shouldOpen = !isGroupOpen(label)
    openGroups.value = shouldOpen ? { [label]: true } : {}
}

watch(() => route.path, () => {
    const activeGroup = filteredSidebarLinks.value.find(item =>
        item.children?.some(child => child.to === route.path)
    )

    openGroups.value = activeGroup?.children ? { [activeGroup.label]: true } : {}
}, { immediate: true })

// ==== Label role untuk display ====
const roleLabel = computed(() => {
    const map = {
        super_admin: 'Super Admin',
        editor: 'Editor'
    }
    return map[user.value?.role] || 'Admin'
})

const userMenuItems = computed(() => [
    [
        { label: 'Profil Saya', icon: 'i-tabler-user', to: '/admin/profil' },
        ...(can('pengaturan')
            ? [{ label: 'Pengaturan', icon: 'i-tabler-settings', to: '/admin/pengaturan' }]
            : [])
    ],
    [
        {
            label: 'Keluar',
            icon: 'i-tabler-logout',
            color: 'error',
            onSelect: async () => { await logout() }
        }
    ]
])
</script>

<template>
    <div class="min-h-screen bg-slate-50 dark:bg-slate-950">
        <!-- SIDEBAR (Terang dengan tint biru) -->
        <aside
            class="fixed top-0 left-0 z-40 h-screen transition-all duration-300 bg-sky-50 dark:bg-slate-900 border-r border-sky-100 dark:border-slate-800 flex flex-col"
            :class="isSidebarOpen ? 'w-64' : 'w-16'">
            <!-- Logo Area -->
            <div class="flex items-center gap-3 h-16 px-4 border-b border-sky-100 dark:border-slate-800 shrink-0">
                <img src="/Logo-Kabupaten-Tulungagung.png" alt="Logo" class="w-8 h-9 shrink-0" />
                <div v-if="isSidebarOpen" class="leading-tight overflow-hidden">
                    <p class="text-[10px] text-sky-700/70 dark:text-slate-400">Admin Panel</p>
                    <p class="text-sm font-bold text-sky-900 dark:text-white whitespace-nowrap">Tulungagung</p>
                </div>
            </div>

            <!-- Nav Menu -->
            <nav class="flex-1 overflow-y-auto py-4 px-2 sidebar-scroll">
                <template v-for="(item, i) in filteredSidebarLinks" :key="i">
                    <!-- Item dengan children -->
                    <div v-if="item.children">
                        <button type="button"
                            class="flex items-center gap-3 w-full px-3 py-2 rounded-lg text-xs font-bold text-sky-700/70 dark:text-slate-500 uppercase tracking-wider mb-1 mt-3 hover:bg-sky-100 dark:hover:bg-slate-800 transition-colors"
                            :aria-expanded="isGroupOpen(item.label)" @click="toggleGroup(item.label)">
                            <UIcon :name="item.icon" class="w-4 h-4 shrink-0" />
                            <span v-if="isSidebarOpen" class="flex-1 text-left">{{ item.label }}</span>
                            <UIcon v-if="isSidebarOpen"
                                :name="isGroupOpen(item.label) ? 'i-tabler-chevron-up' : 'i-tabler-chevron-down'"
                                class="w-4 h-4 shrink-0" />
                        </button>
                        <Transition name="sidebar-group">
                            <div v-if="isSidebarOpen && isGroupOpen(item.label)" class="space-y-0.5 mb-2">
                                <NuxtLink v-for="child in item.children" :key="child.to" :to="child.to"
                                    class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors"
                                    :class="route.path === child.to
                                        ? 'bg-sky-500 text-white font-semibold shadow-sm shadow-sky-500/30'
                                        : 'text-slate-600 dark:text-slate-300 hover:bg-sky-100 dark:hover:bg-slate-800 hover:text-sky-700 dark:hover:text-white'">
                                    <UIcon :name="child.icon" class="w-4 h-4 shrink-0" />
                                    <span class="truncate">{{ child.label }}</span>
                                </NuxtLink>
                            </div>
                        </Transition>
                    </div>

                    <!-- Item tanpa children -->
                    <NuxtLink v-else :to="item.to"
                        class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-xs font-bold text-sky-700/70 dark:text-slate-500 uppercase tracking-wider transition-colors mb-1"
                        :class="route.path === item.to
                            ? 'bg-sky-500 text-white font-semibold shadow-sm shadow-sky-500/30'
                            : 'hover:bg-sky-100 dark:hover:bg-slate-800 hover:text-sky-700 dark:hover:text-white'">
                        <UIcon :name="item.icon" class="w-4 h-4 shrink-0" />
                        <span v-if="isSidebarOpen">{{ item.label }}</span>
                    </NuxtLink>
                </template>
            </nav>

            <!-- User Info -->
            <div class="border-t border-sky-100 dark:border-slate-800 p-3 shrink-0 bg-sky-100/50 dark:bg-slate-900">
                <div class="flex items-center gap-3">
                    <div
                        class="w-9 h-9 rounded-full bg-sky-500 flex items-center justify-center text-white font-bold shrink-0">
                        {{ userProfile.nama.charAt(0) }}
                    </div>
                    <div v-if="isSidebarOpen" class="flex-1 min-w-0">
                        <p class="text-sm font-semibold text-sky-900 dark:text-white truncate">{{ userProfile.nama }}
                        </p>
                        <div class="flex items-center gap-1 mt-0.5">
                            <UBadge color="primary" variant="subtle" size="xs">{{ roleLabel }}</UBadge>
                        </div>
                    </div>
                </div>
            </div>
        </aside>

        <!-- MAIN AREA -->
        <div class="transition-all duration-300 flex flex-col min-h-screen" :class="isSidebarOpen ? 'ml-64' : 'ml-16'">
            <!-- TOP BAR -->
            <header
                class="sticky top-0 z-40 bg-white/95 dark:bg-slate-900/95 backdrop-blur-md border-b border-slate-200 dark:border-slate-800 h-16 flex items-center px-6 gap-4">
                <UButton
                    :icon="isSidebarOpen ? 'i-tabler-layout-sidebar-left-collapse' : 'i-tabler-layout-sidebar-left-expand'"
                    color="neutral" variant="ghost" @click="isSidebarOpen = !isSidebarOpen" />

                <div class="flex-1">
                    <h1 class="text-lg font-bold text-slate-900 dark:text-white">
                        {{ route.meta.title || 'Dashboard' }}
                    </h1>
                </div>

                <AdminNotifications />
                <UButton icon="i-tabler-home" color="neutral" variant="ghost" to="/" />
                <UDropdownMenu :items="userMenuItems" :modal="false">
                    <UButton color="neutral" variant="ghost">
                        <div
                            class="w-7 h-7 rounded-full bg-sky-500 flex items-center justify-center text-white text-xs font-bold">
                            {{ userProfile.nama.charAt(0) }}
                        </div>
                    </UButton>
                </UDropdownMenu>
            </header>

            <!-- CONTENT -->
            <main class="flex-1 p-6">
                <slot />
            </main>

            <!-- FOOTER -->
            <footer
                class="border-t border-slate-200 dark:border-slate-800 px-6 py-4 text-xs text-slate-400 text-center">
                © {{ new Date().getFullYear() }} Pemerintah Kabupaten Tulungagung — Admin Panel
            </footer>
        </div>
    </div>
</template>

<style scoped>
.sidebar-group-enter-active,
.sidebar-group-leave-active {
    overflow: hidden;
    transition: opacity 180ms ease, transform 180ms ease, max-height 180ms ease;
    max-height: 28rem;
}

.sidebar-group-enter-from,
.sidebar-group-leave-to {
    opacity: 0;
    transform: translateY(-4px);
    max-height: 0;
}

/* Scrollbar tipis untuk sidebar */
.sidebar-scroll::-webkit-scrollbar {
    width: 6px;
}

.sidebar-scroll::-webkit-scrollbar-track {
    background: transparent;
}

.sidebar-scroll::-webkit-scrollbar-thumb {
    background: rgba(14, 165, 233, 0.2);
    border-radius: 3px;
}

.sidebar-scroll::-webkit-scrollbar-thumb:hover {
    background: rgba(14, 165, 233, 0.4);
}
</style>