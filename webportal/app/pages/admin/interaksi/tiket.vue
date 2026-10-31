<script setup>
definePageMeta({ layout: 'admin', title: 'Tiket Pengaduan' })
const toast = useToast()
const API_BASE = 'http://localhost:8080/api'

const list = ref([])
const loading = ref(true)
const saving = ref(false)
const search = ref('')
const selectedStatus = ref('Semua')
const page = ref(1)
const perPage = 10

const showModal = ref(false)
const showDeleteModal = ref(false)
const editingItem = ref(null)
const deletingItem = ref(null)

const form = ref({ id: null, status: 'Baru', balasan_admin: '' })
const selectedTemplate = ref(null)

const templateBalasan = [
    { label: 'Aduan sedang diproses', value: '<p>Terima kasih atas laporan Anda. Aduan sedang kami proses.</p>' },
    { label: 'Aduan telah selesai', value: '<p>Terima kasih atas laporan Anda. Aduan telah kami tindak lanjuti.</p>' },
    { label: 'Memerlukan informasi tambahan', value: '<p>Mohon melengkapi informasi tambahan agar aduan dapat kami proses.</p>' }
]

const statusOptions = ['Semua', 'Baru', 'Diproses', 'Selesai']
const statusConfig = ref([])

const fetchStatusConfig = async () => {
    try {
        const response = await $fetch(`${API_BASE}/pengaduan-konfigurasi`, { credentials: 'include' })
        statusConfig.value = response.status || []
        if (statusConfig.value.length) {
            statusOptions.splice(1, statusOptions.length - 1, ...statusConfig.value.map(item => item.nama))
        }
    } catch (err) {
        // Keep the built-in statuses as a fallback when configuration is unavailable.
    }
}

const fetchData = async () => {
    loading.value = true
    try {
        const res = await $fetch(`${API_BASE}/admin/pengaduan`, { credentials: 'include' })
        list.value = res.pengaduan || []
    } catch (err) {
        toast.add({ title: 'Gagal memuat', color: 'error' })
    } finally {
        loading.value = false
    }
}

const filteredList = computed(() => {
    return list.value.filter(item => {
        const matchSearch = item.nama.toLowerCase().includes(search.value.toLowerCase()) ||
            item.isi_aduan.toLowerCase().includes(search.value.toLowerCase())
        const matchStatus = selectedStatus.value === 'Semua' || item.status === selectedStatus.value
        return matchSearch && matchStatus
    })
})

const paginatedList = computed(() => {
    const start = (page.value - 1) * perPage
    return filteredList.value.slice(start, start + perPage)
})

watch([search, selectedStatus], () => { page.value = 1 })

const openDetailModal = (item) => {
    editingItem.value = item
    form.value = { id: item.id, status: item.status, balasan_admin: item.balasan_admin }
    selectedTemplate.value = null
    showModal.value = true
}

const applyTemplate = (value) => {
    const template = templateBalasan.find(item => item.value === value)
    if (template) form.value.balasan_admin = template.value
}

const openDeleteModal = (item) => {
    deletingItem.value = item
    showDeleteModal.value = true
}

const saveData = async () => {
    saving.value = true
    try {
        await $fetch(`${API_BASE}/admin/pengaduan/${form.value.id}`, {
            method: 'PUT',
            body: { status: form.value.status, balasan_admin: form.value.balasan_admin },
            credentials: 'include'
        })
        showModal.value = false
        await fetchData()
        toast.add({ title: 'Berhasil', description: 'Aduan telah diperbarui', color: 'success' })
    } catch (err) {
        toast.add({ title: 'Gagal menyimpan', color: 'error' })
    } finally { saving.value = false }
}

const confirmDelete = async () => {
    saving.value = true
    try {
        await $fetch(`${API_BASE}/admin/pengaduan/${deletingItem.value.id}`, { method: 'DELETE', credentials: 'include' })
        showDeleteModal.value = false
        await fetchData()
        toast.add({ title: 'Berhasil dihapus', color: 'success' })
    } catch (err) {
        toast.add({ title: 'Gagal menghapus', color: 'error' })
    } finally { saving.value = false }
}

const getStatusColor = (status) => {
    if (status === 'Baru') return 'error'
    if (status === 'Diproses') return 'warning'
    if (status === 'Selesai') return 'success'
    return 'neutral'
}

onMounted(() => { fetchData(); fetchStatusConfig() })

const sendingEmail = ref(false)

// Fungsi untuk mengirim balasan via Email
const sendReply = async () => {
    if (!form.value.balasan_admin) {
        toast.add({ title: 'Validasi gagal', description: 'Balasan tidak boleh kosong', color: 'error' })
        return
    }

    sendingEmail.value = true
    try {
        // Panggil API khusus untuk kirim email
        await $fetch(`${API_BASE}/admin/pengaduan/${form.value.id}/send`, {
            method: 'POST',
            body: {
                status: form.value.status,
                balasan_admin: form.value.balasan_admin
            },
            credentials: 'include'
        })

        showModal.value = false
        await fetchData()
        toast.add({ title: 'Terkirim!', description: 'Balasan telah dikirim ke email pengadu', color: 'success' })
    } catch (err) {
        toast.add({ title: 'Gagal mengirim', description: err?.data?.statusMessage || 'Terjadi kesalahan', color: 'error' })
    } finally {
        sendingEmail.value = false
    }
}
</script>

<template>
    <div class="space-y-6">
        <AdminPageHeader title="Tiket Pengaduan" description="Kelola dan balas aduan dari masyarakat." />

        <UCard>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <AdminSearch v-model="search" placeholder="Cari nama atau isi aduan..." />
                <USelect v-model="selectedStatus" :items="statusOptions" icon="i-tabler-filter" size="lg" />
            </div>
        </UCard>

        <UCard>
            <div v-if="loading" class="py-12 text-center">
                <UIcon name="i-tabler-loader-2" class="w-10 h-10 animate-spin mx-auto text-sky-500" />
            </div>
            <AdminEmptyState v-else-if="filteredList.length === 0" icon="i-tabler-inbox" title="Belum ada aduan"
                description="Aduan dari masyarakat akan muncul di sini." />

            <div v-else class="overflow-x-auto">
                <table class="w-full">
                    <thead>
                        <tr class="border-b border-slate-200 dark:border-slate-700">
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase">Tanggal</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase">Pengadu</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase">Kategori</th>
                            <th class="text-left py-3 px-4 text-xs font-bold text-slate-500 uppercase">Status</th>
                            <th class="text-right py-3 px-4 text-xs font-bold text-slate-500 uppercase">Aksi</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="item in paginatedList" :key="item.id"
                            class="border-b border-slate-100 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-800/50">
                            <td class="py-3 px-4 text-sm text-slate-600">{{ item.created_at }}</td>
                            <td class="py-3 px-4">
                                <p class="font-medium text-slate-900 dark:text-white">{{ item.nama }}</p>
                                <p class="text-xs text-slate-500">{{ item.kontak }}</p>
                            </td>
                            <td class="py-3 px-4">
                                <UBadge color="neutral" variant="subtle" size="sm">{{ item.kategori }}</UBadge>
                            </td>
                            <td class="py-3 px-4">
                                <UBadge :color="getStatusColor(item.status)" variant="subtle" size="sm">{{ item.status
                                    }}</UBadge>
                            </td>
                            <td class="py-3 px-4">
                                <div class="flex items-center justify-end gap-2">
                                    <UButton icon="i-tabler-eye" color="primary" variant="soft" size="xs"
                                        @click="openDetailModal(item)" label="Lihat & Balas" />
                                    <AdminButtonDelete size="xs" @click="openDeleteModal(item)" />
                                </div>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>

            <template #footer>
                <AdminPagination v-model:page="page" :total="filteredList.length" :items-per-page="perPage" />
            </template>
        </UCard>

        <!-- Modal Detail & Balasan -->
        <UModal v-model:open="showModal" :title="'Detail Aduan'" :ui="{ width: 'max-w-2xl' }">
            <template #body>
                <div class="space-y-4" v-if="editingItem">
                    <!-- ... (Info Pengadu & Isi Aduan tetap sama) ... -->

                    <hr class="border-slate-200 dark:border-slate-700" />

                    <!-- Form Balasan -->
                    <UFormField label="Status Aduan" required>
                        <USelect v-model="form.status" :items="statusOptions.slice(1)" size="lg" class="w-full" />
                    </UFormField>

                    <UFormField label="Balasan Admin">
                        <div class="space-y-3">
                            <USelect v-model="selectedTemplate" :items="templateBalasan" value-key="value"
                                placeholder="Pilih template..." size="lg" class="w-full"
                                @update:model-value="applyTemplate" />
                            <ClientOnly>
                                <AdminTiptapEditor v-model="form.balasan_admin" />
                            </ClientOnly>
                        </div>
                    </UFormField>
                </div>
            </template>
            <template #footer>
                <div class="flex items-center justify-between w-full">
                    <UButton label="Batal" color="neutral" variant="outline" @click="showModal = false" />
                    <div class="flex gap-2">
                        <UButton label="Simpan Draft" color="neutral" variant="soft" :loading="saving"
                            @click="saveData" />
                        <UButton label="Kirim via Email" color="primary" icon="i-tabler-send" :loading="sendingEmail"
                            @click="sendReply" />
                    </div>
                </div>
            </template>
        </UModal>

        <AdminModalConfirm v-model:open="showDeleteModal" :item-name="'Aduan dari ' + deletingItem?.nama"
            :loading="saving" @confirm="confirmDelete" />
    </div>
</template>