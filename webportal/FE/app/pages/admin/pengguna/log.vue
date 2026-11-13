<script setup>
definePageMeta({
    layout: 'admin',
    title: 'Log Aktivitas'
})

const toast = useToast()

// ==== State ====
const logs = ref([])
const total = ref(0)
const loading = ref(true)
const search = ref('')
const selectedAction = ref('all')
const page = ref(1)
const perPage = 15

// ==== Opsi Filter ====
const actionOptions = [
    { label: 'Semua Aksi', value: 'all' },
    { label: 'Login', value: 'login' },
    { label: 'Logout', value: 'logout' },
    { label: 'Tambah Data', value: 'create' },
    { label: 'Ubah Data', value: 'update' },
    { label: 'Hapus Data', value: 'delete' },
    { label: 'Ubah Permission', value: 'permission_change' }
]

// ==== Warna badge per action ====
const actionBadgeMap = {
    login: { color: 'success', label: 'Login', icon: 'i-tabler-login' },
    logout: { color: 'neutral', label: 'Logout', icon: 'i-tabler-logout' },
    create: { color: 'primary', label: 'Tambah', icon: 'i-tabler-plus' },
    update: { color: 'warning', label: 'Ubah', icon: 'i-tabler-edit' },
    delete: { color: 'error', label: 'Hapus', icon: 'i-tabler-trash' },
    permission_change: { color: 'info', label: 'Permission', icon: 'i-tabler-shield-lock' }
}

const getActionConfig = (action) => {
    return actionBadgeMap[action] || { color: 'neutral', label: action, icon: 'i-tabler-activity' }
}

// ==== Role badge ====
const getRoleLabel = (role) => {
    const map = {
        super_admin: 'Super Admin',
        editor: 'Editor'
    }
    return map[role] || role || 'System'
}

// ==== Format waktu ====
const formatDateTime = (timestamp) => {
    if (!timestamp) return '-'
    const date = new Date(timestamp)
    return date.toLocaleString('id-ID', {
        day: '2-digit',
        month: 'short',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    })
}

const timeAgo = (timestamp) => {
    if (!timestamp) return ''
    const now = new Date()
    const past = new Date(timestamp)
    const diff = Math.floor((now - past) / 1000)

    if (diff < 60) return `${diff} detik lalu`
    if (diff < 3600) return `${Math.floor(diff / 60)} menit lalu`
    if (diff < 86400) return `${Math.floor(diff / 3600)} jam lalu`
    if (diff < 2592000) return `${Math.floor(diff / 86400)} hari lalu`
    return ''
}

// ==== Fetch ====
const fetchLogs = async () => {
    loading.value = true
    try {
        const params = new URLSearchParams()
        if (selectedAction.value !== 'all') params.append('action', selectedAction.value)
        if (search.value) params.append('search', search.value)
        params.append('page', page.value)
        params.append('perPage', perPage)

        const response = await useRequestFetch()(`http://localhost:8080/api/admin/activity-logs?${params.toString()}`, { credentials: 'include' })
        logs.value = response.logs
        total.value = response.total
    } catch (err) {
        toast.add({
            title: 'Gagal memuat log',
            description: err?.data?.statusMessage || 'Terjadi kesalahan',
            icon: 'i-tabler-alert-circle',
            color: 'error'
        })
    } finally {
        loading.value = false
    }
}

// ==== Watch filter action → reset ke halaman 1 ====
watch(selectedAction, () => {
    page.value = 1
    fetchLogs()
})

// ==== Debounce search → reset ke halaman 1 ====
let searchTimer = null
watch(search, () => {
    if (searchTimer) clearTimeout(searchTimer)
    searchTimer = setTimeout(() => {
        page.value = 1
        fetchLogs()
    }, 400)
})

// ==== Watch page → refetch ====
watch(page, () => {
    fetchLogs()
})

// ==== Init ====
onMounted(() => {
    fetchLogs()
})
</script>

<template>
    <div class="space-y-6">

        <!-- Header -->
        <UCard>
            <div class="flex items-center justify-between flex-wrap gap-4">
                <div>
                    <h2 class="text-2xl font-bold text-slate-900 dark:text-white">Log Aktivitas</h2>
                    <p class="text-slate-500 mt-1">Riwayat aktivitas admin di portal Pemkab Tulungagung.</p>
                </div>
                <UButton icon="i-tabler-refresh" label="Refresh" color="neutral" variant="outline" :loading="loading"
                    @click="fetchLogs" />
            </div>
        </UCard>

        <!-- Filter -->
        <UCard>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <UInput v-model="search" icon="i-tabler-search" placeholder="Cari nama user atau deskripsi..."
                    size="lg" />
                <USelect v-model="selectedAction" :items="actionOptions" icon="i-tabler-filter" size="lg" />
            </div>
        </UCard>

        <!-- Tabel -->
        <UCard>
            <template #header>
                <div class="flex items-center justify-between">
                    <h3 class="font-bold text-slate-900 dark:text-white flex items-center gap-2">
                        <UIcon name="i-tabler-history" class="text-sky-500" />
                        Riwayat Aktivitas
                    </h3>
                    <span class="text-sm text-slate-500">{{ total }} log</span>
                </div>
            </template>

            <!-- Loading -->
            <div v-if="loading" class="py-12 text-center">
                <UIcon name="i-tabler-loader-2" class="w-10 h-10 text-sky-500 animate-spin mx-auto" />
                <p class="text-sm text-slate-500 mt-3">Memuat log aktivitas...</p>
            </div>

            <!-- Empty -->
            <AdminEmptyState v-else-if="logs.length === 0" icon="i-tabler-history-off"
                title="Belum ada aktivitas tercatat"
                description="Log aktivitas akan muncul di sini setelah ada interaksi dari admin." />

            <!-- Table -->
            <div v-else class="overflow-x-auto">
                <table class="w-full">
                    <thead>
                        <tr class="border-b border-slate-200 dark:border-slate-700">
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Waktu</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Pengguna</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Aksi</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">
                                Deskripsi</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase tracking-wider">IP
                                Address
                            </th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="item in logs" :key="item.id"
                            class="border-b border-slate-100 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors">
                            <td class="py-3 px-4">
                                <div>
                                    <p class="text-sm text-slate-900 dark:text-white font-medium">
                                        {{ formatDateTime(item.created_at) }}
                                    </p>
                                    <p class="text-xs text-slate-400">{{ timeAgo(item.created_at) }}</p>
                                </div>
                            </td>
                            <td class="py-3 px-4">
                                <div class="flex items-center gap-2">
                                    <div
                                        class="w-8 h-8 rounded-full bg-sky-100 dark:bg-sky-950 flex items-center justify-center text-sky-600 dark:text-sky-400 text-xs font-bold shrink-0">
                                        {{ item.user_nama ? item.user_nama.charAt(0).toUpperCase() : '?' }}
                                    </div>
                                    <div>
                                        <p class="text-sm font-medium text-slate-900 dark:text-white">
                                            {{ item.user_nama || 'System' }}
                                        </p>
                                        <p class="text-xs text-slate-400">{{ getRoleLabel(item.user_role) }}</p>
                                    </div>
                                </div>
                            </td>
                            <td class="py-3 px-4">
                                <UBadge :color="getActionConfig(item.action).color" variant="subtle" size="sm">
                                    <UIcon :name="getActionConfig(item.action).icon" class="w-3.5 h-3.5 mr-1" />
                                    {{ getActionConfig(item.action).label }}
                                </UBadge>
                            </td>
                            <td class="py-3 px-4">
                                <p class="text-sm text-slate-700 dark:text-slate-300 max-w-md line-clamp-2">
                                    {{ item.description || '-' }}
                                </p>
                            </td>
                            <td class="py-3 px-4">
                                <code
                                    class="text-xs bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300 px-2 py-1 rounded">
                        {{ item.ip_address || '-' }}
                    </code>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>

            <template #footer>
                <div class="flex items-center justify-between flex-wrap gap-3">
                    <p class="text-xs text-slate-400">
                        Menampilkan {{ logs.length }} dari {{ total }} log
                        <span v-if="total > 0"> — Halaman {{ page }} dari {{ Math.ceil(total / perPage) }}</span>
                    </p>
                    <UPagination v-if="total > perPage" v-model:page="page" :total="total" :items-per-page="perPage" />
                </div>
            </template>
        </UCard>

    </div>
</template>