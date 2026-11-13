<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

type NotificationItem = {
    id: number
    title: string
    message: string
    link?: string
    is_read: boolean
    notif_type: string
    created_at: string
}

const isOpen = ref(false)
const notifications = ref<NotificationItem[]>([])
const unreadCount = ref(0)
const loading = ref(false)
const popoverRef = ref<HTMLElement | null>(null)

const API_BASE = '/api'

const fetchNotifications = async () => {
    try {
        loading.value = true
        const response = await $fetch<{ notifications: NotificationItem[], unread_count: number }>(
            `${API_BASE}/admin/notifications`,
            { credentials: 'include' }
        )
        if (response) {
            notifications.value = response.notifications || []
            unreadCount.value = response.unread_count || 0
        }
    } catch (err) {
        console.error('Gagal mengambil notifikasi:', err)
    } finally {
        loading.value = false
    }
}

const markAsRead = async (id: number, link?: string) => {
    try {
        await $fetch(`${API_BASE}/admin/notifications/${id}/read`, {
            method: 'PUT',
            credentials: 'include'
        })
        const item = notifications.value.find(n => n.id === id)
        if (item && !item.is_read) {
            item.is_read = true
            unreadCount.value = Math.max(0, unreadCount.value - 1)
        }
    } catch (err) {
        console.error('Gagal memperbarui status notifikasi:', err)
    }

    if (link) {
        isOpen.value = false
        navigateTo(link)
    }
}

const markAllAsRead = async () => {
    try {
        await $fetch(`${API_BASE}/admin/notifications/read-all`, {
            method: 'PUT',
            credentials: 'include'
        })
        notifications.value.forEach(n => {
            n.is_read = true
        })
        unreadCount.value = 0
    } catch (err) {
        console.error('Gagal menandai semua notifikasi:', err)
    }
}

const togglePopover = () => {
    isOpen.value = !isOpen.value
    if (isOpen.value) {
        fetchNotifications()
    }
}

const handleClickOutside = (event: MouseEvent) => {
    if (popoverRef.value && !popoverRef.value.contains(event.target as Node)) {
        isOpen.value = false
    }
}

const formatTimeAgo = (dateStr: string) => {
    if (!dateStr) return ''
    const date = new Date(dateStr)
    const now = new Date()
    const diffInSeconds = Math.floor((now.getTime() - date.getTime()) / 1000)

    if (diffInSeconds < 60) return 'Baru saja'
    if (diffInSeconds < 3600) return `${Math.floor(diffInSeconds / 60)} mtk lalu`
    if (diffInSeconds < 86400) return `${Math.floor(diffInSeconds / 3600)} jam lalu`
    if (diffInSeconds < 604800) return `${Math.floor(diffInSeconds / 86400)} hari lalu`

    return date.toLocaleDateString('id-ID', { day: 'numeric', month: 'short' })
}

const getNotificationIcon = (type: string) => {
    switch (type) {
        case 'pengaduan':
            return 'i-tabler-message-dots'
        case 'berita':
            return 'i-tabler-news'
        case 'warning':
            return 'i-tabler-alert-triangle'
        case 'success':
            return 'i-tabler-check'
        default:
            return 'i-tabler-bell'
    }
}

const getNotificationIconColor = (type: string) => {
    switch (type) {
        case 'pengaduan':
            return 'text-amber-500 bg-amber-50 dark:bg-amber-950/40 dark:text-amber-400'
        case 'berita':
            return 'text-sky-500 bg-sky-50 dark:bg-sky-950/40 dark:text-sky-400'
        case 'warning':
            return 'text-red-500 bg-red-50 dark:bg-red-950/40 dark:text-red-400'
        case 'success':
            return 'text-emerald-500 bg-emerald-50 dark:bg-emerald-950/40 dark:text-emerald-400'
        default:
            return 'text-indigo-500 bg-indigo-50 dark:bg-indigo-950/40 dark:text-indigo-400'
    }
}

let timer: any = null

onMounted(() => {
    fetchNotifications()
    document.addEventListener('click', handleClickOutside)
    timer = setInterval(fetchNotifications, 30000)
})

onUnmounted(() => {
    document.removeEventListener('click', handleClickOutside)
    if (timer) clearInterval(timer)
})
</script>

<template>
    <div ref="popoverRef" class="relative inline-block">
        <!-- Bell Icon Button (matches UButton ghost variant) -->
        <button
            type="button"
            @click="togglePopover"
            class="relative mt-2 p-2 rounded-lg text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors focus:outline-none focus-visible:ring-2 focus-visible:ring-sky-500 focus-visible:ring-offset-2"
            aria-label="Notifikasi"
            aria-expanded="false"
            :aria-haspopup="true"
        >
            <UIcon name="i-tabler-bell" class="w-5 h-5" />
            
            <!-- Badge Count -->
            <span
                v-if="unreadCount > 0"
                class="absolute -top-0.5 -right-0.5 flex items-center justify-center min-w-4.5 h-4.5 px-1 text-[10px] font-bold text-white bg-rose-500 rounded-full animate-pulse border-2 border-white dark:border-slate-900"
            >
                {{ unreadCount > 99 ? '99+' : unreadCount }}
            </span>
        </button>

        <!-- Dropdown Popover -->
        <Transition
            enter-active-class="transition duration-150 ease-out"
            enter-from-class="transform scale-95 opacity-0"
            enter-to-class="transform scale-100 opacity-100"
            leave-active-class="transition duration-100 ease-in"
            leave-from-class="transform scale-100 opacity-100"
            leave-to-class="transform scale-95 opacity-0"
        >
            <div
                v-if="isOpen"
                class="absolute right-0 mt-1.5 w-80 sm:w-96 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl shadow-xl z-50 overflow-hidden"
            >
                <!-- Popover Header -->
                <div class="px-4 py-3 border-b border-slate-100 dark:border-slate-800 flex items-center justify-between bg-slate-50/50 dark:bg-slate-800/50">
                    <div class="flex items-center gap-2">
                        <h3 class="font-semibold text-sm text-slate-800 dark:text-slate-200">Notifikasi</h3>
                        <span v-if="unreadCount > 0" class="px-2 py-0.5 text-xs font-medium bg-sky-100 text-sky-700 dark:bg-sky-950/60 dark:text-sky-300 rounded-full">
                            {{ unreadCount }} baru
                        </span>
                    </div>
                    <button
                        v-if="unreadCount > 0"
                        @click="markAllAsRead"
                        class="text-xs font-medium text-sky-600 dark:text-sky-400 hover:underline transition-colors"
                    >
                        Tandai semua dibaca
                    </button>
                </div>

                <!-- Notifications List -->
                <div class="max-h-95 overflow-y-auto divide-y divide-slate-100 dark:divide-slate-800">
                    <template v-if="loading && notifications.length === 0">
                        <div class="p-4 text-center text-xs text-slate-400">
                            Memuat notifikasi...
                        </div>
                    </template>

                    <template v-else-if="notifications.length > 0">
                        <div
                            v-for="item in notifications"
                            :key="item.id"
                            @click="markAsRead(item.id, item.link)"
                            class="p-3.5 flex gap-3 hover:bg-slate-50 dark:hover:bg-slate-800/60 transition-colors cursor-pointer relative"
                            :class="{ 'bg-sky-50/30 dark:bg-sky-950/20': !item.is_read }"
                        >
                            <!-- Icon -->
                            <div
                                class="w-8 h-8 rounded-full flex items-center justify-center shrink-0 mt-0.5"
                                :class="getNotificationIconColor(item.notif_type)"
                            >
                                <UIcon :name="getNotificationIcon(item.notif_type)" class="w-4 h-4" />
                            </div>

                            <!-- Content -->
                            <div class="flex-1 min-w-0">
                                <div class="flex items-start justify-between gap-1 mb-0.5">
                                    <p class="text-xs font-semibold text-slate-800 dark:text-slate-200 truncate" :class="{ 'font-bold': !item.is_read }">
                                        {{ item.title }}
                                    </p>
                                    <span class="text-[10px] text-slate-400 shrink-0">
                                        {{ formatTimeAgo(item.created_at) }}
                                    </span>
                                </div>
                                <p class="text-xs text-slate-600 dark:text-slate-400 line-clamp-2 leading-relaxed">
                                    {{ item.message }}
                                </p>
                            </div>

                            <!-- Unread Dot Indicator -->
                            <div v-if="!item.is_read" class="w-2 h-2 rounded-full bg-sky-500 shrink-0 self-center"></div>
                        </div>
                    </template>

                    <!-- Empty State -->
                    <template v-else>
                        <div class="p-8 text-center flex flex-col items-center justify-center">
                            <div class="w-10 h-10 rounded-full bg-slate-100 dark:bg-slate-800 flex items-center justify-center text-slate-400 mb-2">
                                <UIcon name="i-tabler-bell-off" class="w-5 h-5" />
                            </div>
                            <p class="text-xs font-medium text-slate-600 dark:text-slate-400">Belum ada notifikasi</p>
                            <p class="text-[11px] text-slate-400 mt-0.5">Notifikasi aktivitas akan muncul di sini</p>
                        </div>
                    </template>
                </div>
            </div>
        </Transition>
    </div>
</template>